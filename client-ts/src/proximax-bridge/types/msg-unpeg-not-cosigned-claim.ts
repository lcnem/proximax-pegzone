import { ValAddress, Msg } from "cosmos-client";

export class MsgUnpegNotCosignedClaim extends Msg {
  constructor(
    public address: ValAddress,
    public tx_hash: string,
    public not_cosigned_address: ValAddress[]
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
      value.tx_hash,
      value.not_cosigned_address
    );
  }
}
