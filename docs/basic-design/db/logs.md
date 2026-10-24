# logs - ユーザのログイン・ログアウトのログ

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|created_at|TIMESTAMPZ|T|now()|-|-|
|failure_reason|enum (logs_failure_reason)|F|NULL|-|-|
|type|enum (logs_type)|T|-|-|-|
|user_id|uuid|F|-|FK (app_users.id)|登録外ユーザのとき`NULL`|

## 列挙体

|logs_failure_reason|備考|
|---|---|
|auth_state_mismatch|CSRF or タイムアウト|
|google_token_expired|Googleのコード有効期限切れ|
|google_api_fail|Google側の障害 or ネットワークエラー|
|app_user_not_allowed|許可されていないユーザ|

|logs_type|
|---|
|login_success|
|login_failure|
|logout|

## 運用
- レコードの追加のみ許可、変更・削除はしない
