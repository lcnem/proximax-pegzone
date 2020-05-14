# ProximaX Peg Zone

- 基本設計は cosmos/peggy を踏襲。
- d,cli,relayer の 3 部設計。
- `1000000000` `stake`のネイティブトークン。

## Install

Sample of Ubuntu 20.04

```shell
apt update
apt install build-essential
cd ~
wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.14.linux-amd64.tar.gz
echo export PATH='$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
mkdir -p /usr/local/src/lcnem/
cd /usr/local/src/github.com/lcnem
git clone https://github.com/lcnem/proximax-pegzone.git
cd proximax-pegzone
git checkout vX.X.X
make install
```

## Config

```shell
vi $HOME/.pxbd/config/config.toml
```

## Start multiple nodes by docker-compose

```
make build-linux

make build-docker-pxbdnode

make localnet-start
```

## Commands

### CLI

Commands via CLI enable you to create a transaction and broadcast it with your signature made with your private key.

#### Peg

```shell
pxbcli tx proximaxbridge peg [key_or_address] [mainchain_tx_hash] [to_address] [amount]
```

#### Unpeg

```shell
pxbcli tx proximaxbridge unpeg [key_or_address] [amount] [mainchain_address] [first_cosigner_address]
```

#### Request Invitation

```shell
pxbcli tx proximaxbridge request-invitation [from_key_or_address] [multisig_account_address] [new_cosigner_public_key] [first_cosigner_address]
```

### Relayer

#### Init

```shell
pxbrelayer init
```

#### ProximaX relayer

```shell
pxbrelayer proximax [validator_from_name] [proximax_node] [proximax_private_key] [proximax_multisig_address] --chain-id [chain-id]
```

#### Cosmos Relayer

```shell
pxbrelayer cosmos [tendermint_node] [proximax_node] [validator_moniker] [proximax_cosigner_private_key] [multisig_account_public_key] --chain-id [chain-id]
```
