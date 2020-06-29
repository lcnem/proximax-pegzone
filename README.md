# ProximaX Peg Zone

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
vi ~/.pxbd/config/config.toml
```

## Test Locally with single node

### Initialize and Start Deamon node

```shell
pxbd init test --chain-id testing
pxbcli config keyring-backend test

pxbcli keys add validator
pxbcli keys add account

pxbd add-genesis-account $(pxbcli keys show validator -a) 1000token,100000000stake

pxbd gentx --name validator --keyring-backend test
pxbd collect-gentxs
pxbd start
```

### Start Relayer

```shell
# Relayer for Cosmos
pxbrelayer start [URL for node by RPC] [URL for ProximaX node] [Validator Name] [ProximaX Cosigner Private Key] [ProximaX Multisig Account Public Key] --chain-id=[ChainID]
```

Example

```shell
pxbrelayer start http://127.0.0.1:26657 http://bctestnet1.brimstone.xpxsirius.io:3000 validator1  8611AF477E001C9D033216F94328BD22F91E782FD2D104FAE3F5B66997579154 8007692AB57547661CD0721FBE18AA1DB27E0CC55921D4C0C9A3BEBC96221AC7 --chain-id=testing --rpc-url=http://127.0.0.1:26657
```

## Test Locally with Multiple nodes by docker-compose

```shell
make build-linux

make build-docker-pxbdnode

make localnet-start
```

## Commands

### CLI

Commands via CLI enable you to create a transaction and broadcast it with your signature made with your private key.

#### Peg

Mint and send tokens to given account in cosmos by hash of transaction in ProximaX

```shell
pxbcli tx proximaxbridge peg [Validator's key or address] [Transaction Hash on ProximaX] [Recipient Account Address in Cosmos] [Amount]
```

#### Unpeg

Burn tokens from Cosmos Account and send to account in ProximaX

```shell
pxbcli tx proximaxbridge unpeg [Validator's key or address] [Amount] [Recipient Account Address in ProximaX] [First Cosigner Address in Cosmos]
```

#### Request Invitation

Invite new ProximaX account to Multisig Account

```shell
pxbcli tx proximaxbridge request-invitation [from_key_or_address] [multisig_account_address] [new_cosigner_public_key] [first_cosigner_address]
```
