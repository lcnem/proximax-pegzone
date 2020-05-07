import { AccAddress, ValAddress, Msg } from "cosmos-client";
import { Coin } from "cosmos-client/api";

export class MsgUnpeg extends Msg {
  constructor(
    public address: AccAddress,
    public mainchain_address: string,
    public amount: Coin[],
    public first_cosigner_address: ValAddress
  ) {
    super();
  }

  /**
   *
   * @param value
   */
  public static fromJSON(value: any) {
    return new this(
      AccAddress.fromBech32(value.address),
      value.mainchain_address,
      value.amount,
      ValAddress.fromBech32(value.first_cosigner_address)
    );
  }
}
