# task-app
## 仕様
https://docs.google.com/document/d/1boS3nGZQVaGXhDF39Ag8aGvQumKO0-l1MTGjaVrcGBY/edit?tab=t.0
# Task App

## プロジェクトの概要
Task Appは、タスク管理を簡単にするためのWeb APIベースのアプリケーションです。このプロジェクトはタスクの作成、管理、認証機能を提供します。また、session認証の実装をしています。

## 仕様書

## 使用技術や依存関係

## プログラミング言語
- **Go**
- **HTML**
- **javascript**

## データベース
- **MySQL**

## その他ライブラリ
- `gorilla/mux`: ルーティング用
- `dotenv`: 環境変数管理

## セットアップ手順
1. **リポジトリのクローン**
   ```bash
   git clone https://github.com/go-WebAPI-team/task-app.git
   cd task-app
   ```
2. **依存関係のインストール Go言語がインストールされていることを確認してください。その後、以下のコマンドを実行します。**

    ```bash
    go mod tidy
    ```
3.**環境変数の設定 .env ファイルを作成し、以下のように必要な環境変数を設定してください。**
```.env
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=taskdb
```

4.**データベースのセットアップ**
mysql
```sql
CREATE DATABASE `task-app`;
USE `task-app`;
```
