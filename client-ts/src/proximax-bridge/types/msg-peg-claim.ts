import { AccAddress, Coin } from "cosmos-client";

export type MsgPegClaim = {
  mainchain_tx_hash: string;
  address: AccAddress;
  amount: Coin[];
};
