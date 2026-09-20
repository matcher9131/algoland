# app_users - ユーザ情報

登録されているユーザとロール

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|id|uuid|T|uuidv7()|PK|自動採番|
|is_admin|bool|F|false|-|-|
|mail_hash|text|T|-|-|メールアドレスのソルト付きハッシュ|

## 運用
- レコードの追加・変更・削除はしない
