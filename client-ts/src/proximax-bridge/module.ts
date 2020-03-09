import { CosmosSDK } from "cosmos-client";
import { StdTx } from "cosmos-client/x/auth";
import { UnpegReq, RequestInvitationReq, Params } from "./types";

/**
 * Getting the multisig address on mainchain for collateral.
 * @param sdk
 */
export async function getParams(sdk: CosmosSDK) {
  return await sdk.get<Params>("/proximax_bridge/parameters");
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
