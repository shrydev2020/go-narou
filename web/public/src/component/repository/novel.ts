import { NovelListClient } from "../../../generated/list_grpc_pb";
import { Req } from "../../../generated/list_pb";
import { credentials, ServiceError } from "@grpc/grpc-js";
import { Novels } from "../../../generated/list_pb";

function getList() {
  const client = new NovelListClient(
    "localhost:18080",
    credentials.createInsecure()
  );
  const req: Req = new Req();
  client.get(req, (err: ServiceError | null, response: Novels) => {
    const lst = response.getNovelsList();
    console.log(lst);
  });
}
