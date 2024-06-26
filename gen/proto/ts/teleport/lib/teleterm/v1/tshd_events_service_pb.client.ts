/* eslint-disable */
// @generated by protobuf-ts 2.9.3 with parameter long_type_number,eslint_disable,add_pb_suffix,server_grpc1,ts_nocheck
// @generated from protobuf file "teleport/lib/teleterm/v1/tshd_events_service.proto" (package "teleport.lib.teleterm.v1", syntax proto3)
// tslint:disable
// @ts-nocheck
//
//
// Teleport
// Copyright (C) 2023  Gravitational, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { TshdEventsService } from "./tshd_events_service_pb";
import type { PromptMFAResponse } from "./tshd_events_service_pb";
import type { PromptMFARequest } from "./tshd_events_service_pb";
import type { SendPendingHeadlessAuthenticationResponse } from "./tshd_events_service_pb";
import type { SendPendingHeadlessAuthenticationRequest } from "./tshd_events_service_pb";
import type { SendNotificationResponse } from "./tshd_events_service_pb";
import type { SendNotificationRequest } from "./tshd_events_service_pb";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { ReloginResponse } from "./tshd_events_service_pb";
import type { ReloginRequest } from "./tshd_events_service_pb";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * TshdEventsService is served by the Electron app. The tsh daemon calls this service to notify the
 * app about actions that happen outside of the app itself.
 *
 * @generated from protobuf service teleport.lib.teleterm.v1.TshdEventsService
 */
export interface ITshdEventsServiceClient {
    /**
     * Relogin makes the Electron app display a login modal for the specific root cluster. The request
     * returns a response after the relogin procedure has been successfully finished.
     *
     * @generated from protobuf rpc: Relogin(teleport.lib.teleterm.v1.ReloginRequest) returns (teleport.lib.teleterm.v1.ReloginResponse);
     */
    relogin(input: ReloginRequest, options?: RpcOptions): UnaryCall<ReloginRequest, ReloginResponse>;
    /**
     * SendNotification causes the Electron app to display a notification in the UI. The request
     * accepts a specific message rather than a generic string so that the Electron is in control as
     * to what message is displayed and how exactly it looks.
     *
     * @generated from protobuf rpc: SendNotification(teleport.lib.teleterm.v1.SendNotificationRequest) returns (teleport.lib.teleterm.v1.SendNotificationResponse);
     */
    sendNotification(input: SendNotificationRequest, options?: RpcOptions): UnaryCall<SendNotificationRequest, SendNotificationResponse>;
    /**
     * SendPendingHeadlessAuthentication notifies the Electron app of a pending headless authentication,
     * which it can use to initiate headless authentication resolution in the UI.
     *
     * @generated from protobuf rpc: SendPendingHeadlessAuthentication(teleport.lib.teleterm.v1.SendPendingHeadlessAuthenticationRequest) returns (teleport.lib.teleterm.v1.SendPendingHeadlessAuthenticationResponse);
     */
    sendPendingHeadlessAuthentication(input: SendPendingHeadlessAuthenticationRequest, options?: RpcOptions): UnaryCall<SendPendingHeadlessAuthenticationRequest, SendPendingHeadlessAuthenticationResponse>;
    /**
     * PromptMFA notifies the Electron app that the daemon is waiting for the user to answer an MFA prompt.
     *
     * @generated from protobuf rpc: PromptMFA(teleport.lib.teleterm.v1.PromptMFARequest) returns (teleport.lib.teleterm.v1.PromptMFAResponse);
     */
    promptMFA(input: PromptMFARequest, options?: RpcOptions): UnaryCall<PromptMFARequest, PromptMFAResponse>;
}
/**
 * TshdEventsService is served by the Electron app. The tsh daemon calls this service to notify the
 * app about actions that happen outside of the app itself.
 *
 * @generated from protobuf service teleport.lib.teleterm.v1.TshdEventsService
 */
export class TshdEventsServiceClient implements ITshdEventsServiceClient, ServiceInfo {
    typeName = TshdEventsService.typeName;
    methods = TshdEventsService.methods;
    options = TshdEventsService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * Relogin makes the Electron app display a login modal for the specific root cluster. The request
     * returns a response after the relogin procedure has been successfully finished.
     *
     * @generated from protobuf rpc: Relogin(teleport.lib.teleterm.v1.ReloginRequest) returns (teleport.lib.teleterm.v1.ReloginResponse);
     */
    relogin(input: ReloginRequest, options?: RpcOptions): UnaryCall<ReloginRequest, ReloginResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<ReloginRequest, ReloginResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * SendNotification causes the Electron app to display a notification in the UI. The request
     * accepts a specific message rather than a generic string so that the Electron is in control as
     * to what message is displayed and how exactly it looks.
     *
     * @generated from protobuf rpc: SendNotification(teleport.lib.teleterm.v1.SendNotificationRequest) returns (teleport.lib.teleterm.v1.SendNotificationResponse);
     */
    sendNotification(input: SendNotificationRequest, options?: RpcOptions): UnaryCall<SendNotificationRequest, SendNotificationResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<SendNotificationRequest, SendNotificationResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * SendPendingHeadlessAuthentication notifies the Electron app of a pending headless authentication,
     * which it can use to initiate headless authentication resolution in the UI.
     *
     * @generated from protobuf rpc: SendPendingHeadlessAuthentication(teleport.lib.teleterm.v1.SendPendingHeadlessAuthenticationRequest) returns (teleport.lib.teleterm.v1.SendPendingHeadlessAuthenticationResponse);
     */
    sendPendingHeadlessAuthentication(input: SendPendingHeadlessAuthenticationRequest, options?: RpcOptions): UnaryCall<SendPendingHeadlessAuthenticationRequest, SendPendingHeadlessAuthenticationResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<SendPendingHeadlessAuthenticationRequest, SendPendingHeadlessAuthenticationResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * PromptMFA notifies the Electron app that the daemon is waiting for the user to answer an MFA prompt.
     *
     * @generated from protobuf rpc: PromptMFA(teleport.lib.teleterm.v1.PromptMFARequest) returns (teleport.lib.teleterm.v1.PromptMFAResponse);
     */
    promptMFA(input: PromptMFARequest, options?: RpcOptions): UnaryCall<PromptMFARequest, PromptMFAResponse> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<PromptMFARequest, PromptMFAResponse>("unary", this._transport, method, opt, input);
    }
}
