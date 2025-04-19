## 開発環境セットアップ
```
git clone <リポジトリのURL>
go run .                         # HTTPサーバー起動
curl localhost:8080/health       # 別ターミナルからリクエスト
go test -v ./...                 # テスト
```
