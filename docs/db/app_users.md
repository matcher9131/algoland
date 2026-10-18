# app_users - ユーザ情報

登録されているユーザとロール

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|email_hmac|bytea|F|-|-|メールアドレスのHMAC-SHA256|
|id|uuid|T|uuidv7()|PK|自動採番|
|is_admin|bool|F|false|-|-|
|sub|text|F|-|-|Google SSOのsub|
|updated_at|TIMESTAMPZ|T|now()|-|-|

## 運用
- レコードの追加・削除はしない
  - 更新は後述の手順によるもののみ許可
- ホワイトリストとし、事前にログインしうる全ユーザのデータを保持する
- 各ユーザの初回ログイン時に以下を行う
  - `email_hmac` を照合し、成功すれば使用済みとして `NULL` に書き換える
  - `sub` に値を設定する
- 各ユーザ2度目以降のログイン時は `sub` のみで照合する
