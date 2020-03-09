export type Params = {
  mainchain_multisig_address: string;
  cosigners: Cosigner[]
}

export type Cosigner = {
  validator_address: string;
  mainchain_address: string;
};
