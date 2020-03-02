import {
  TransferTransaction,
  Deadline,
  Address,
  Mosaic,
  PlainMessage,
  NetworkType
} from "tsjs-xpx-chain-sdk";
import { PegMessage } from "./peg-message";

/**
 *
 * @param multisigAddress `getMultisigAddress`
 * @param mosaics
 * @param cosmosAddress
 * @param networkType
 */
export function createPegTransaction(
  multisigAddress: string,
  mosaics: Mosaic[],
  cosmosAddress: string,
  networkType: NetworkType
) {
  const message: PegMessage = {
    cosmos_address: cosmosAddress
  };
  const json = JSON.stringify(message);

  return TransferTransaction.create(
    Deadline.create(),
    Address.createFromRawAddress(multisigAddress),
    mosaics,
    PlainMessage.create(json),
    networkType
  );
}
