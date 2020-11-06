// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var list_pb = require('./list_pb.js');
var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js');

function serialize_Novels(arg) {
  if (!(arg instanceof list_pb.Novels)) {
    throw new Error('Expected argument of type Novels');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_Novels(buffer_arg) {
  return list_pb.Novels.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_Req(arg) {
  if (!(arg instanceof list_pb.Req)) {
    throw new Error('Expected argument of type Req');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_Req(buffer_arg) {
  return list_pb.Req.deserializeBinary(new Uint8Array(buffer_arg));
}


var NovelListService = exports.NovelListService = {
  get: {
    path: '/NovelList/Get',
    requestStream: false,
    responseStream: false,
    requestType: list_pb.Req,
    responseType: list_pb.Novels,
    requestSerialize: serialize_Req,
    requestDeserialize: deserialize_Req,
    responseSerialize: serialize_Novels,
    responseDeserialize: deserialize_Novels,
  },
};

exports.NovelListClient = grpc.makeGenericClientConstructor(NovelListService);
