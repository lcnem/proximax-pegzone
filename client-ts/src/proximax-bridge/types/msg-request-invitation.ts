import { ValAddress } from "cosmos-client";

export type MsgRequestInvitation = {
  validator_address: ValAddress;
  mainchain_address: string;
  first_cosigner_address: ValAddress;
};
