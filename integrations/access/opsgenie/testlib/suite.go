/*
 * Teleport
 * Copyright (C) 2024  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package testlib

import (
	"context"
	"runtime"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/integrations/access/opsgenie"
	"github.com/gravitational/teleport/integrations/access/pagerduty"
	"github.com/gravitational/teleport/integrations/lib/logger"
	"github.com/gravitational/teleport/integrations/lib/testing/integration"
)

const (
	NotifyServiceName       = "Teleport Notifications"
	NotifyServiceAnnotation = types.TeleportNamespace + types.ReqAnnotationNotifyServicesLabel
	ResponderName1          = "Responder 1"
	ResponderName2          = "Responder 2"
	ResponderName3          = "Responder 3"
)

// OpsgenieSuite is the OpsGenie access plugin test suite.
// It implements the testify.TestingSuite interface.
type OpsgenieSuite struct {
	*integration.AccessRequestSuite
	appConfig    opsgenie.Config
	raceNumber   int
	fakeOpsgenie *FakeOpsgenie

	ogNotifyResponder opsgenie.Responder
	ogResponder1      opsgenie.Responder
	ogResponder2      opsgenie.Responder
	ogResponder3      opsgenie.Responder
}

// SetupTest starts a fake OpsGenie and generates the plugin configuration.
// It also configures the role notifications for OpsGenie notifications and
// automatic approval.
// It is run for each test.
func (s *OpsgenieSuite) SetupTest() {
	t := s.T()
	ctx := context.Background()

	err := logger.Setup(logger.Config{Severity: "debug"})
	require.NoError(t, err)
	s.raceNumber = 2 * runtime.GOMAXPROCS(0)

	s.fakeOpsgenie = NewFakeOpsgenie(s.raceNumber)
	t.Cleanup(s.fakeOpsgenie.Close)

	// This service should be notified for every access request.
	s.ogNotifyResponder = s.fakeOpsgenie.StoreResponder(opsgenie.Responder{
		Name: NotifyServiceName,
	})
	s.AnnotateRequesterRoleAccessRequests(
		ctx,
		NotifyServiceAnnotation,
		[]string{NotifyServiceName},
	)

	// Responder 1 and 2 are on-call and should be automatically approved.
	// Responder 3 is not.
	s.ogResponder1 = s.fakeOpsgenie.StoreResponder(opsgenie.Responder{
		Name: ResponderName1,
	})
	s.ogResponder2 = s.fakeOpsgenie.StoreResponder(opsgenie.Responder{
		Name: ResponderName2,
	})
	s.ogResponder3 = s.fakeOpsgenie.StoreResponder(opsgenie.Responder{
		Name: ResponderName3,
	})
	s.AnnotateRequesterRoleAccessRequests(
		ctx,
		pagerduty.ServicesDefaultAnnotation,
		[]string{ResponderName1, ResponderName2},
	)

	var conf opsgenie.Config
	conf.Teleport = s.TeleportConfig()
	conf.ClientConfig.APIEndpoint = s.fakeOpsgenie.URL()

	s.appConfig = conf
}

// startApp starts the OpsGenie plugin, waits for it to become ready and returns.
func (s *OpsgenieSuite) startApp() {
	t := s.T()
	t.Helper()

	app, err := opsgenie.NewOpsgenieApp(context.Background(), &s.appConfig)
	require.NoError(t, err)
	s.RunAndWaitReady(t, app)
}

// TestAlertCreation validates that an alert is created to the service
// specified in the role's annotation.
func (s *OpsgenieSuite) TestAlertCreation() {
	t := s.T()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	t.Cleanup(cancel)

	s.startApp()

	// Test execution: create an access request
	req := s.CreateAccessRequest(ctx, integration.RequesterOSSUserName, nil)

	// Validate the alert has been created in OpsGenie and its ID is stored in
	// the plugin_data.
	pluginData := s.checkPluginData(ctx, req.GetName(), func(data opsgenie.PluginData) bool {
		return data.AlertID != ""
	})

	alert, err := s.fakeOpsgenie.CheckNewAlert(ctx)
	require.NoError(t, err, "no new alerts stored")

	assert.Equal(t, alert.ID, pluginData.AlertID)
}

// TestApproval tests that when a request is approved, its corresponding alert
// is updated to reflect the new request state and a note is added to the alert.
func (s *OpsgenieSuite) TestApproval() {
	t := s.T()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	t.Cleanup(cancel)

	s.startApp()

	// Test setup: we create an access request and wait for its alert.
	req := s.CreateAccessRequest(ctx, integration.RequesterOSSUserName, nil)

	alert, err := s.fakeOpsgenie.CheckNewAlert(ctx)
	require.NoError(t, err, "no new alerts stored")

	// Test execution: we approve the request
	err = s.Ruler().ApproveAccessRequest(ctx, req.GetName(), "okay")
	require.NoError(t, err)

	// Validating the plugin added a note to the alert describing the review.
	note, err := s.fakeOpsgenie.CheckNewAlertNote(ctx)
	require.NoError(t, err)
	assert.Equal(t, alert.ID, note.AlertID)
	assert.Contains(t, note.Note, "Access request has been approved")
	assert.Contains(t, note.Note, "Reason: okay")

	// Validating the plugin resolved the alert.
	alertUpdate, err := s.fakeOpsgenie.CheckAlertUpdate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "resolved", alertUpdate.Status)
}

// TestDenial tests that when a request is denied, its corresponding alert
// is updated to reflect the new request state and a note is added to the alert.
func (s *OpsgenieSuite) TestDenial() {
	t := s.T()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	t.Cleanup(cancel)

	s.startApp()

	// Test setup: we create an access request and wait for its alert.
	req := s.CreateAccessRequest(ctx, integration.RequesterOSSUserName, nil)

	alert, err := s.fakeOpsgenie.CheckNewAlert(ctx)
	require.NoError(t, err, "no new alerts stored")

	// Test execution: we deny the request
	err = s.Ruler().DenyAccessRequest(ctx, req.GetName(), "not okay")
	require.NoError(t, err)

	// Validating the plugin added a note to the alert describing the review.
	note, err := s.fakeOpsgenie.CheckNewAlertNote(ctx)
	require.NoError(t, err)
	assert.Equal(t, alert.ID, note.AlertID)
	assert.Contains(t, note.Note, "Access request has been denied")
	assert.Contains(t, note.Note, "Reason: not okay")

	// Validating the plugin resolved the alert.
	alertUpdate, err := s.fakeOpsgenie.CheckAlertUpdate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "resolved", alertUpdate.Status)
}

// TestReviewNotes tests that alert notes are sent after the access request
// is reviewed. Each review should create a new note.
func (s *OpsgenieSuite) TestReviewNotes() {
	t := s.T()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	t.Cleanup(cancel)

	if !s.TeleportFeatures().AdvancedAccessWorkflows {
		t.Skip("Doesn't work in OSS version")
	}

	s.startApp()

	// Test setup: we create an access request and wait for its alert.
	req := s.CreateAccessRequest(ctx, integration.Requester1UserName, nil)

	// Test execution: we submit two reviews
	err := s.Reviewer1().SubmitAccessRequestReview(ctx, req.GetName(), types.AccessReview{
		Author:        integration.Reviewer1UserName,
		ProposedState: types.RequestState_APPROVED,
		Created:       time.Now(),
		Reason:        "okay",
	})
	require.NoError(t, err)

	err = s.Reviewer2().SubmitAccessRequestReview(ctx, req.GetName(), types.AccessReview{
		Author:        integration.Reviewer2UserName,
		ProposedState: types.RequestState_APPROVED,
		Created:       time.Now(),
		Reason:        "not okay",
	})
	require.NoError(t, err)

	// Validate alert notes were sent with the correct content.
	pluginData := s.checkPluginData(ctx, req.GetName(), func(data opsgenie.PluginData) bool {
		return data.AlertID != "" && data.ReviewsCount == 2
	})

	note, err := s.fakeOpsgenie.CheckNewAlertNote(ctx)
	require.NoError(t, err)
	assert.Equal(t, pluginData.AlertID, note.AlertID)
	assert.Contains(t, note.Note, integration.Reviewer1UserName+" reviewed the request", "note must contain a review author")
	assert.Contains(t, note.Note, "Resolution: APPROVED", "note must contain an approval resolution")
	assert.Contains(t, note.Note, "Reason: okay", "note must contain an approval reason")

	note, err = s.fakeOpsgenie.CheckNewAlertNote(ctx)
	require.NoError(t, err)
	assert.Equal(t, pluginData.AlertID, note.AlertID)
	assert.Contains(t, note.Note, integration.Reviewer2UserName+" reviewed the request", "note must contain a review author")
	assert.Contains(t, note.Note, "Resolution: APPROVED", "note must contain a approval resolution")
	assert.Contains(t, note.Note, "Reason: not okay", "note must contain a denial reason")
}

// TestApprovalByReview tests that the alert is annotated and resolved after the
// access request approval threshold is reached.
func (s *OpsgenieSuite) TestApprovalByReview() {
	t := s.T()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	t.Cleanup(cancel)

	if !s.TeleportFeatures().AdvancedAccessWorkflows {
		t.Skip("Doesn't work in OSS version")
	}

	s.startApp()

	// Test setup: we create an access request and wait for its alert.
	req := s.CreateAccessRequest(ctx, integration.Requester1UserName, nil)

	alert, err := s.fakeOpsgenie.CheckNewAlert(ctx)
	require.NoError(t, err, "no new alerts stored")

	// Test execution: we submit two reviews
	err = s.Reviewer1().SubmitAccessRequestReview(ctx, req.GetName(), types.AccessReview{
		Author:        integration.Reviewer1UserName,
		ProposedState: types.RequestState_APPROVED,
		Created:       time.Now(),
		Reason:        "okay",
	})
	require.NoError(t, err)

	err = s.Reviewer2().SubmitAccessRequestReview(ctx, req.GetName(), types.AccessReview{
		Author:        integration.Reviewer2UserName,
		ProposedState: types.RequestState_APPROVED,
		Created:       time.Now(),
		Reason:        "finally okay",
	})
	require.NoError(t, err)

	// Validate alert notes were sent with the correct content.
	note, err := s.fakeOpsgenie.CheckNewAlertNote(ctx)
	require.NoError(t, err)
	assert.Equal(t, alert.ID, note.AlertID)
	assert.Contains(t, note.Note, integration.Reviewer1UserName+" reviewed the request", "note must contain a review author")

	note, err = s.fakeOpsgenie.CheckNewAlertNote(ctx)
	require.NoError(t, err)
	assert.Equal(t, alert.ID, note.AlertID)
	assert.Contains(t, note.Note, integration.Reviewer2UserName+" reviewed the request", "note must contain a review author")

	// Validate the alert got resolved.
	data := s.checkPluginData(ctx, req.GetName(), func(data opsgenie.PluginData) bool {
		return data.ReviewsCount == 2 && data.Resolution.Tag != opsgenie.Unresolved
	})
	assert.Equal(t, opsgenie.Resolution{Tag: opsgenie.ResolvedApproved, Reason: "finally okay"}, data.Resolution)

	note, err = s.fakeOpsgenie.CheckNewAlertNote(ctx)
	require.NoError(t, err)
	assert.Equal(t, alert.ID, note.AlertID)
	assert.Contains(t, note.Note, "Access request has been approved")
	assert.Contains(t, note.Note, "Reason: finally okay")

	alertUpdate, err := s.fakeOpsgenie.CheckAlertUpdate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "resolved", alertUpdate.Status)
}

// TestDenialByReview tests that the alert is annotated and resolved after the
// access request denial threshold is reached.
func (s *OpsgenieSuite) TestDenialByReview() {
	t := s.T()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	t.Cleanup(cancel)

	if !s.TeleportFeatures().AdvancedAccessWorkflows {
		t.Skip("Doesn't work in OSS version")
	}

	s.startApp()

	// Test setup: we create an access request and wait for its alert.
	req := s.CreateAccessRequest(ctx, integration.Requester1UserName, nil)

	alert, err := s.fakeOpsgenie.CheckNewAlert(ctx)
	require.NoError(t, err, "no new alerts stored")

	// Test execution: we submit two reviews
	err = s.Reviewer1().SubmitAccessRequestReview(ctx, req.GetName(), types.AccessReview{
		Author:        integration.Reviewer1UserName,
		ProposedState: types.RequestState_DENIED,
		Created:       time.Now(),
		Reason:        "not okay",
	})
	require.NoError(t, err)

	err = s.Reviewer2().SubmitAccessRequestReview(ctx, req.GetName(), types.AccessReview{
		Author:        integration.Reviewer2UserName,
		ProposedState: types.RequestState_DENIED,
		Created:       time.Now(),
		Reason:        "finally not okay",
	})
	require.NoError(t, err)

	// Validate alert notes were sent with the correct content.
	note, err := s.fakeOpsgenie.CheckNewAlertNote(ctx)
	require.NoError(t, err)
	assert.Equal(t, alert.ID, note.AlertID)
	assert.Contains(t, note.Note, integration.Reviewer1UserName+" reviewed the request", "note must contain a review author")

	note, err = s.fakeOpsgenie.CheckNewAlertNote(ctx)
	require.NoError(t, err)
	assert.Equal(t, alert.ID, note.AlertID)
	assert.Contains(t, note.Note, integration.Reviewer2UserName+" reviewed the request", "note must contain a review author")

	// Validate the alert got resolved.
	data := s.checkPluginData(ctx, req.GetName(), func(data opsgenie.PluginData) bool {
		return data.ReviewsCount == 2 && data.Resolution.Tag != opsgenie.Unresolved
	})
	assert.Equal(t, opsgenie.Resolution{Tag: opsgenie.ResolvedDenied, Reason: "finally not okay"}, data.Resolution)

	note, err = s.fakeOpsgenie.CheckNewAlertNote(ctx)
	require.NoError(t, err)
	assert.Equal(t, alert.ID, note.AlertID)
	assert.Contains(t, note.Note, "Access request has been denied")
	assert.Contains(t, note.Note, "Reason: finally not okay")

	alertUpdate, err := s.fakeOpsgenie.CheckAlertUpdate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "resolved", alertUpdate.Status)
}
