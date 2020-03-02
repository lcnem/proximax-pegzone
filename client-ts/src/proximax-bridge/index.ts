import * as proximax_bridge from "./module";
import { codec } from "cosmos-client";
import {
  MsgInvitationNotCosignedClaim,
  MsgPegClaim,
  MsgRequestInvitation,
  MsgUnpegNotCosignedClaim,
  MsgUnpeg
} from "./types";

export { proximax_bridge };
export * from "./types";

codec.registerCodec(
  "proximax_bridge/MsgInvitationNotCosignedClaim",
  MsgInvitationNotCosignedClaim,
  MsgInvitationNotCosignedClaim.fromJSON
);
codec.registerCodec(
  "proximax_bridge/MsgPegClaim",
  MsgPegClaim,
  MsgPegClaim.fromJSON
);
codec.registerCodec(
  "proximax_bridge/MsgRequestInvitation",
  MsgRequestInvitation,
  MsgRequestInvitation.fromJSON
);
codec.registerCodec(
  "proximax_bridge/MsgUnpegNotCosignedClaim",
  MsgUnpegNotCosignedClaim,
  MsgUnpegNotCosignedClaim.fromJSON
);
codec.registerCodec("proximax_bridge/MsgUnpeg", MsgUnpeg, MsgUnpeg.fromJSON);
