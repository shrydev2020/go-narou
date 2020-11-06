// package: 
// file: list.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as list_pb from "./list_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

interface INovelListService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    get: INovelListService_IGet;
}

interface INovelListService_IGet extends grpc.MethodDefinition<list_pb.Req, list_pb.Novels> {
    path: "/NovelList/Get";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<list_pb.Req>;
    requestDeserialize: grpc.deserialize<list_pb.Req>;
    responseSerialize: grpc.serialize<list_pb.Novels>;
    responseDeserialize: grpc.deserialize<list_pb.Novels>;
}

export const NovelListService: INovelListService;

export interface INovelListServer {
    get: grpc.handleUnaryCall<list_pb.Req, list_pb.Novels>;
}

export interface INovelListClient {
    get(request: list_pb.Req, callback: (error: grpc.ServiceError | null, response: list_pb.Novels) => void): grpc.ClientUnaryCall;
    get(request: list_pb.Req, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: list_pb.Novels) => void): grpc.ClientUnaryCall;
    get(request: list_pb.Req, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: list_pb.Novels) => void): grpc.ClientUnaryCall;
}

export class NovelListClient extends grpc.Client implements INovelListClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public get(request: list_pb.Req, callback: (error: grpc.ServiceError | null, response: list_pb.Novels) => void): grpc.ClientUnaryCall;
    public get(request: list_pb.Req, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: list_pb.Novels) => void): grpc.ClientUnaryCall;
    public get(request: list_pb.Req, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: list_pb.Novels) => void): grpc.ClientUnaryCall;
}
