import { BaseReq } from "cosmos-client";

export type RequestInvitationReq = {
  base_req: BaseReq;
  address: string;
  mainchain_address: string;
  first_cosigner_address: string;
};
