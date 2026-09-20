# games - ゲームセッション

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|created_at|TIMESTAMPZ|T|now()|-|-|
|id|uuid|T|uuidv7()|PK|自動採番|
|pattern_id|int8|F|NULL|-|問題「ダンジョン」におけるパターン割り振りに使用|
|problem_id|int8|T|-|FK (problems.id)|-|
|result|enum (game_result)|T|ongoing|-|-|
|updated_at|TIMESTAMPZ|T|now()|-|-|
|user_id|uuid|T|-|FK (app_users.id)|-|

## 列挙体

|game_result|
|---|
|ongoing|
|correct|
|incorrect|

## 運用
- レコードの追加許可、削除禁止。
- 更新は以下の通り
  - `result = "ongoing"` のとき、`result = "correct"`または`result = "incorrect"`に更新可能
    - 更新時は`updated_at`を現在時刻で更新する
  - 上記条件に当てはまらないケースは更新禁止
