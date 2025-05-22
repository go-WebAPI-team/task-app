# Task App

## 概要
Task Appは、ユーザー認証付きのタスク管理Webアプリケーションです。タスクの作成・編集・削除、タグ付け、完了管理、フィルタリングなどができます。

## デモ動画

[demo](demo/demo-taskapp.mp4)

## 通信構成図

![sequence](assets/Golang-WebAPI-sequence_diagram(authentication).png)


## 主な機能
- ユーザー登録・ログイン・ログアウト（セッション認証）
- タスクの作成・編集・削除
- タスクの完了状態管理
- タグの作成・削除、タスクへのタグ付け
- タスク・タグのフィルタリング
- SwaggerによるAPIドキュメント

## 技術スタック
- バックエンド: Go, chi, SQL（MySQL等）
- フロントエンド: HTML, Bootstrap, JavaScript（Fetch API）
- テスト: Go標準testing, go-sqlmock

## セットアップ

1. リポジトリをクローン
    ```sh
    git clone <このリポジトリのURL>
    cd task-app
    ```

2. **環境変数の設定**  
   `.env` ファイルを作成し、以下のように必要な環境変数を設定してください。

    ```env
    DB_HOST=localhost
    DB_PORT=3306
    DB_USER=user
    DB_PASSWORD=password
    DB_NAME=taskdb
    ```

3. **データベースのセットアップ**  
   MySQL例:
    ```sql
    CREATE DATABASE `task-app`;
    USE `task-app`;
    ```



5. サーバ起動
    ```sh
    go run main.go
    ```

6. ブラウザで `http://localhost:8080` にアクセス

# 主要なAPIエンドポイント（Swaggerドキュメント準拠）

| メソッド | パス | summary（概要） | description（詳細） |
|----------|-------------------------------|----------------------|----------------------|
| GET      | /tags                         | タグ一覧取得         | 登録済みタグの一覧を取得します |
| POST     | /tags                         | タグを新規作成       | タグ情報を新規作成します |
| DELETE   | /tags/{id}                    | タグを削除           | 指定IDのタグを削除します |
| GET      | /tasks                        | タスク一覧取得       | ログインユーザのタスクをフィルタ付きで取得します（クエリ: is_done, tag_id） |
| POST     | /tasks                        | タスクを新規作成     | JSON で受け取ったタスク情報を保存します |
| GET      | /tasks/{id}                   | タスク詳細取得       | 指定IDのタスク詳細を取得します |
| PUT      | /tasks/{id}                   | タスクを更新         | 指定IDのタスク情報を更新します |
| DELETE   | /tasks/{id}                   | タスクを削除         | 指定IDのタスクを削除します |
| PATCH    | /tasks/{id}/complete          | タスクを完了状態に更新| タスクの is_done を true に変更します |
| POST     | /tasks/{id}/tags/{tag_id}     | タスクにタグを紐付け | タスク ID とタグ ID を指定して関連づけます |
| DELETE   | /tasks/{id}/tags/{tag_id}     | タスクからタグを解除 | タスク ID からタグ ID の紐付けを解除します |

---

- summary/descriptionはSwagger（OpenAPI）定義の内容を反映しています。
- より詳細なパラメータやレスポンス例は `/swagger/index.html` でAPIドキュメントを参照してください。

# Swaggerドキュメント確認手順

1. サーバ起動
    ```sh
    go run .
    ```
2. Swagger-UIをブラウザから確認  
   [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## テスト

```sh
go test ./...
```

## ライセンス
MIT

---