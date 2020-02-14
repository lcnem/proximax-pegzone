import { ValAddress } from "cosmos-client";

export type MsgUnpegNotProposedClaim = {
  validator_address: ValAddress;
  tx_hash: string;
  first_cosigner_address: ValAddress;
};
