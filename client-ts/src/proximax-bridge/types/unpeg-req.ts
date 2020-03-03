import { BaseReq, Coin } from "cosmos-client";

export type UnpegReq = {
  base_req: BaseReq;
  address: string;
  mainchain_address: string;
  amount: Coin[];
  first_cosigner_address: string;
};
