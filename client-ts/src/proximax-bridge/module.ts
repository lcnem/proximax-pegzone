import { CosmosSDK } from "cosmos-client";
import { StdTx } from "cosmos-client/x/auth";
import { UnpegReq, RequestInvitationReq, Cosigner } from "./types";

export async function cosigners(sdk: CosmosSDK) {
  return await sdk.get<Cosigner[]>("/proximax_bridge/cosigners");
}

export async function unpeg(sdk: CosmosSDK, req: UnpegReq) {
  return await sdk.post<StdTx>("/proximax_bridge/unpeg", req);
}

export async function requestInvitation(
  sdk: CosmosSDK,
  req: RequestInvitationReq
) {
  return await sdk.post<StdTx>("/proximax_bridge/request_invitation", req);
}
