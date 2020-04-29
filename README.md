# ProximaX Peg Zone

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

## Relayer

Cosmos監視

- MsgUnpegがきて、自身がFirstCosignerならマルチシグ提案するだけ。

ProximaX監視

- マルチシグ提案がされたら連署名するだけ。

## Start multiple nodes by docker-compose

```
make build-linux

make build-docker-pxbdnode

make localnet-start
```