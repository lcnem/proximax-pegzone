# ProximaX Peg Zone

- 基本設計は cosmos/peggy を踏襲。
- d,cli,relayer の 3 部設計。
- `1000000000` `stake`のネイティブトークン。

## Start multiple nodes by docker-compose

```
make build-linux

make build-docker-pxbdnode

make localnet-start
```
