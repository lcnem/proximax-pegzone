import {AccAddress, Coin, ValAddress}from "cosmos-client"

export type MsgUnpeg = {
  address: AccAddress;
  mainchain_address: string;
  amount: Coin[];
  first_cosigner_address: ValAddress;
}