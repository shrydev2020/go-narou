import { NovelListClient } from "./list_grpc_pb";
import { Req, Novel, Novels } from "./list_pb";
import { credentials, ServiceError } from "grpc";

export const GetList = (): Array<Novel> => {
  const client = new NovelListClient(
    "localhost:18080",
    credentials.createInsecure()
  );
  const req: Req = new Req();
  let ret = new Array<Novel>();
  client.get(req, (err: ServiceError | null, response: Novels) => {
    const lst = response.getNovelsList();
    console.log(lst);
    ret = lst;
  });
  return ret;
};
