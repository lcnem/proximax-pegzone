import { CosmosSDK } from "cosmos-client";
import { StdTx } from "cosmos-client/x/auth";
import { UnpegReq, RequestInvitationReq, Params } from "./types";
import Axios from "axios";

/**
 * Getting the multisig address on mainchain for collateral.
 * @param sdk
 */
export function getParams(sdk: CosmosSDK) {
  return Axios.get<Params>(`${sdk.url}/proximax_bridge/parameters`);
}

export function unpeg(sdk: CosmosSDK, req: UnpegReq) {
  return sdk.instancifyObjectWithoutAminoJSON<StdTx>(
    StdTx,
    Axios.post("/proximax_bridge/unpeg", req)
  );
}

export function requestInvitation(sdk: CosmosSDK, req: RequestInvitationReq) {
  return sdk.instancifyObjectWithoutAminoJSON<StdTx>(
    StdTx,
    Axios.post("/proximax_bridge/request_invitation", req)
  );
}
