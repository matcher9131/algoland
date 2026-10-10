# sessions - ログインセッション

ログインセッション

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|created_at|TIMESTAMPZ|T|now()|-|-|
|expires_at|TIMESTAMPZ|T|-|-|1日？要検討|
|token_hash|bytea|T|-|PK|-|
|user_id|uuid|T|-|FK (app_users.id)|-|
