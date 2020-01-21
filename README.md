# ProximaX Peg Zone

# Peg Zone

- 基本設計はcosmos/peggyを踏襲。
- d,cli,relayerの3部設計。
- `1000000000` `stake`のネイティブトークン。

## Ethereum Peggyとの相違点

- Etheruem Peggyは初期化する際にコントラクトアドレスを指定する、1ERC20トークン1Peggyチェーン仕様。
- Catapult/ProximaXはネイティブトークンと非ネイティブトークンが並列に扱われるため、1Peggyチェーンで全トークンに対応できる。
- Catapult/ProximaXにはスマートコントラクトによるトークンロックがないため、マルチシグアドレスにトークンを封じ込めることで再現する。
- 全ステーク量のうち、10%以上を占めるバリデータが、マルチシグアドレスの連署名者になれる。
  - Catapult/ProximaXではマルチシグアドレスの連署名者は10アドレスまでという制約と整合できる。
  - Catapult/ProixmaXではマルチシグに必要な連署名を過半数とする。

## モジュール

- Oracle
- Bridge

## Daemon

- OracleはEthereum Peggyのコードパクる。
- Bridge設計が肝になる。

### Bridge

#### MsgReceiveInvitation

`first_cosigner_address`に、Catapult/ProximaXにて連署名に招待するよう要求するメッセージ。

```TypeScript
{
  validator_address: ValAddress;
  mainchain_address: string;
  first_cosigner_address: ValAddress;
}
```

#### MsgInvitationNotProposedProphecy

ステークが10%以上あり、`MsgReceiveInvitation`を行ったにもかかわらず連署名招待されていないことを訴訟するメッセージ。
訴訟が正しければ`first_cosigner_address`はペナルティを与えられるが、訴訟が正しくなければ訴訟者が罰せられる。
ペナルティ分は、正しい方に分配される。

```TypeScript
{
  validator_address: ValAddress;
  mainchain_address: string;
  first_cosigner_address: ValAddress;
}
```

#### MsgInvitationNotCosignedProphecy

ステークが10%以上あり、`MsgReceiveInvitation`を行ったにもかかわらず連署名招待のマルチシグトランザクションに連署名が集まっていないことを訴訟するメッセージ。
訴訟が正しければ`invitee_address`以外の連署名者はペナルティを与えられるが、訴訟が正しくなければ訴訟者が罰せられる。
ペナルティ分は、正しい方に分配される。

```TypeScript
{
  validator_address: ValAddress;
  mainchain_address: string;
  first_cosigner_address: ValAddress;
}
```

#### MsgPegProphecy

メインチェーンにて、ペッグ用マルチシグアドレスに対してアセットを送信したことを訴訟する。
訴訟が正しければ本Peg ZoneにCoinがミントされる。

```TypeScript
{
  address: AccAddress;
  mainchain_tx_hash: string;
  amount: Coin[];
}
```

#### MsgUnPeg

```TypeScript
{
  address: AccAddress;
  mainchain_address: string;
  amount: Coin[];
  first_cosigner_address: ValAddress;
}
```

#### MsgUnpegNotProposedProphecy

`MsgUnPeg`したのにされてねえよという訴訟。

```TypeScript
{
  validator_address: ValAddress;
  tx_hash: string;
  first_cosigner_address: ValAddress;
}
```


#### MsgUnpegNotCosignedPropechy

`MsgUnPeg`したのに連署名されてねえよという訴訟。

```TypeScript
{
  validator_address: ValAddress;
  tx_hash: string;
  first_cosigner_address: ValAddress;
}
```

## CLI

- Oracleに関しては実装なし。
- BridgeモジュールにはCLIとRESTを実装する。
- 基本的にコマンドを使うことはなく、RelayerがフルオートでRESTを叩く

## Relayer

- WebSocketでPeg元チェーンを常駐監視。
- Ethereumの場合、ペッグ用コントラクトアドレスへのトークン振り込みを確認するとイベント発火するようになっている。
- Catapult/ProximaXの場合、ペッグ用マルチシグアドレスへのトークン振り込みを常駐監視するといい。
- イベント発火するとどうするかというと、対応するMsgのトランザクションをアナウンスすればよいだけ。
- 訴訟もフルオートで行う。
