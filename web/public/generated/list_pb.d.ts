// package: 
// file: list.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class Novels extends jspb.Message { 
    clearNovelsList(): void;
    getNovelsList(): Array<Novel>;
    setNovelsList(value: Array<Novel>): Novels;
    addNovels(value?: Novel, index?: number): Novel;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Novels.AsObject;
    static toObject(includeInstance: boolean, msg: Novels): Novels.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Novels, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Novels;
    static deserializeBinaryFromReader(message: Novels, reader: jspb.BinaryReader): Novels;
}

export namespace Novels {
    export type AsObject = {
        novelsList: Array<Novel.AsObject>,
    }
}

export class Novel extends jspb.Message { 
    getId(): number;
    setId(value: number): Novel;

    getAuthor(): string;
    setAuthor(value: string): Novel;

    getTitle(): string;
    setTitle(value: string): Novel;

    getFileTitle(): string;
    setFileTitle(value: string): Novel;

    getTopUrl(): string;
    setTopUrl(value: string): Novel;

    getSiteName(): string;
    setSiteName(value: string): Novel;

    getStory(): string;
    setStory(value: string): Novel;

    getNovelType(): Novel.Kind;
    setNovelType(value: Novel.Kind): Novel;

    getEnd(): boolean;
    setEnd(value: boolean): Novel;


    hasLastUpdate(): boolean;
    clearLastUpdate(): void;
    getLastUpdate(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setLastUpdate(value?: google_protobuf_timestamp_pb.Timestamp): Novel;


    hasNewArrivalsDate(): boolean;
    clearNewArrivalsDate(): void;
    getNewArrivalsDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setNewArrivalsDate(value?: google_protobuf_timestamp_pb.Timestamp): Novel;

    getUseSubdirectory(): boolean;
    setUseSubdirectory(value: boolean): Novel;


    hasGeneralFirstUp(): boolean;
    clearGeneralFirstUp(): void;
    getGeneralFirstUp(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setGeneralFirstUp(value?: google_protobuf_timestamp_pb.Timestamp): Novel;


    hasNovelUpdatedAt(): boolean;
    clearNovelUpdatedAt(): void;
    getNovelUpdatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setNovelUpdatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Novel;


    hasGeneralKastUp(): boolean;
    clearGeneralKastUp(): void;
    getGeneralKastUp(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setGeneralKastUp(value?: google_protobuf_timestamp_pb.Timestamp): Novel;

    getLength(): number;
    setLength(value: number): Novel;

    getSuspend(): boolean;
    setSuspend(value: boolean): Novel;

    getGeneralAllNo(): number;
    setGeneralAllNo(value: number): Novel;


    hasLastCheckAt(): boolean;
    clearLastCheckAt(): void;
    getLastCheckAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setLastCheckAt(value?: google_protobuf_timestamp_pb.Timestamp): Novel;

    clearSubsList(): void;
    getSubsList(): Array<Sub>;
    setSubsList(value: Array<Sub>): Novel;
    addSubs(value?: Sub, index?: number): Sub;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Novel.AsObject;
    static toObject(includeInstance: boolean, msg: Novel): Novel.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Novel, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Novel;
    static deserializeBinaryFromReader(message: Novel, reader: jspb.BinaryReader): Novel;
}

export namespace Novel {
    export type AsObject = {
        id: number,
        author: string,
        title: string,
        fileTitle: string,
        topUrl: string,
        siteName: string,
        story: string,
        novelType: Novel.Kind,
        end: boolean,
        lastUpdate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        newArrivalsDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        useSubdirectory: boolean,
        generalFirstUp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        novelUpdatedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        generalKastUp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        length: number,
        suspend: boolean,
        generalAllNo: number,
        lastCheckAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        subsList: Array<Sub.AsObject>,
    }

    export enum Kind {
    UNSPECIFIED = 0,
    SERIES = 1,
    SS = 2,
    }

}

export class Sub extends jspb.Message { 
    getNovelId(): number;
    setNovelId(value: number): Sub;

    getIndexId(): number;
    setIndexId(value: number): Sub;

    getHref(): string;
    setHref(value: string): Sub;

    getChapter(): string;
    setChapter(value: string): Sub;

    getSubtitle(): string;
    setSubtitle(value: string): Sub;


    hasSubDate(): boolean;
    clearSubDate(): void;
    getSubDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setSubDate(value?: google_protobuf_timestamp_pb.Timestamp): Sub;


    hasSubUpdatedAt(): boolean;
    clearSubUpdatedAt(): void;
    getSubUpdatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setSubUpdatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Sub;


    hasDownloadAt(): boolean;
    clearDownloadAt(): void;
    getDownloadAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setDownloadAt(value?: google_protobuf_timestamp_pb.Timestamp): Sub;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Sub.AsObject;
    static toObject(includeInstance: boolean, msg: Sub): Sub.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Sub, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Sub;
    static deserializeBinaryFromReader(message: Sub, reader: jspb.BinaryReader): Sub;
}

export namespace Sub {
    export type AsObject = {
        novelId: number,
        indexId: number,
        href: string,
        chapter: string,
        subtitle: string,
        subDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        subUpdatedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        downloadAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
}

export class Req extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Req.AsObject;
    static toObject(includeInstance: boolean, msg: Req): Req.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Req, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Req;
    static deserializeBinaryFromReader(message: Req, reader: jspb.BinaryReader): Req;
}

export namespace Req {
    export type AsObject = {
    }
}
