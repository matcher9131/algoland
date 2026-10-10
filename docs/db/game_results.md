# game_results - ゲームの結果

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|created_at|TIMESTAMPZ|T|now()|-|-|
|game_id|uuid|T|-|FK (games.id)|-|
|result|enum (game_result)|T|-|-|-|

## 制約
game_id と result は複合UNIQUE

## 列挙体

|game_result|
|---|
|ongoing|
|correct|
|incorrect|

## 運用
- レコードの追加のみ許可、変更・削除はしない
