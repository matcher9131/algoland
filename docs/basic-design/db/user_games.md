# user_games - プレイヤーとゲームセッションの中間テーブル

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|game_id|uuid|T|-|複合PK, FK (games.id)|-|
|user_id|uuid|T|-|複合PK, FK (app_users.id)|-|

## 運用
- レコードの追加のみ許可、変更・削除はしない
