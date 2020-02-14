import { ValAddress } from "cosmos-client";

export type MsgUnpegNotCosignedClaim = {
  validator_address: ValAddress;
  tx_hash: string;
  first_cosigner_address: ValAddress;
};
