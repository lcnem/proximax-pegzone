# ProximaX Peg Zone

# Peg Zone

- 基本設計は cosmos/peggy を踏襲。
- d,cli,relayer の 3 部設計。
- `1000000000` `stake`のネイティブトークン。

## Ethereum Peggy との相違点

- Etheruem Peggy は初期化する際にコントラクトアドレスを指定する、1ERC20 トークン 1Peggy チェーン仕様。
- Catapult/ProximaX はネイティブトークンと非ネイティブトークンが並列に扱われるため、1Peggy チェーンで全トークンに対応できる。
- Catapult/ProximaX にはスマートコントラクトによるトークンロックがないため、マルチシグアドレスにトークンを封じ込めることで再現する。
- 全ステーク量のうち、10%以上を占めるバリデータが、マルチシグアドレスの連署名者になれる。
  - Catapult/ProximaX ではマルチシグアドレスの連署名者は 10 アドレスまでという制約と整合できる。
  - Catapult/ProixmaX ではマルチシグに必要な連署名を過半数とする。

## モジュール

- Oracle
- Bridge

## Daemon

- Oracle は Ethereum Peggy のコードパクる。
- Bridge 設計が肝になる。

### Bridge

#### MsgPegClaim

メインチェーンにて、ペッグ用マルチシグアドレスに対してアセットを送信したことを訴訟する。
訴訟が正しければ本 Peg Zone に Coin がミントされる。

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

#### MsgUnpegNotProposedClaim

`MsgUnPeg`したのにされてねえよという訴訟。

```TypeScript
{
  validator_address: ValAddress;
  tx_hash: string;
  first_cosigner_address: ValAddress;
}
```

#### MsgRequestInvitation

`first_cosigner_address`に、Catapult/ProximaX にて連署名に招待するよう要求するメッセージ。

```TypeScript
{
  validator_address: ValAddress;
  mainchain_address: string;
  first_cosigner_address: ValAddress;
}
```

#### MsgInvitationNotCosignedClaim

ステークが 10%以上あり、`MsgReceiveInvitation`を行ったにもかかわらず連署名招待のマルチシグトランザクションに連署名が集まっていないことを訴訟するメッセージ。
訴訟が正しければ`invitee_address`以外の連署名者はペナルティを与えられるが、訴訟が正しくなければ訴訟者が罰せられる。
ペナルティ分は、正しい方に分配される。

```TypeScript
{
  validator_address: ValAddress;
  mainchain_address: string;
  first_cosigner_address: ValAddress;
}
```

## CLI

- Oracle に関しては実装なし。
- Bridge モジュールには CLI と REST を実装する。
- 基本的にコマンドを使うことはなく、Relayer がフルオートで REST を叩く

## Relayer

- WebSocket で Peg 元チェーンを常駐監視。
- Ethereum の場合、ペッグ用コントラクトアドレスへのトークン振り込みを確認するとイベント発火するようになっている。
- Catapult/ProximaX の場合、ペッグ用マルチシグアドレスへのトークン振り込みを常駐監視するといい。
- イベント発火するとどうするかというと、対応する Msg のトランザクションをアナウンスすればよいだけ。
- 訴訟もフルオートで行う。
