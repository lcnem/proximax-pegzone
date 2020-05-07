import { AccAddress, Msg, ValAddress } from "cosmos-client";
import { Coin } from "cosmos-client/api";

export class MsgPegClaim extends Msg {
  constructor(
    public address: ValAddress,
    public mainchain_tx_hash: string,
    public to_address: AccAddress,
    public amount: Coin[]
  ) {
    super();
  }

  /**
   *
   * @param value
   */
  public static fromJSON(value: any) {
    return new this(
      ValAddress.fromBech32(value.address),
      value.mainchain_tx_hash,
      AccAddress.fromBech32(value.to_address),
      value.amount
    );
  }
}
