# p3_inputs - 問題「ダンジョン」のユーザ入力

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|created_at|TIMESTAMPZ|T|now()|-|-|
|game_id|uuid|T|-|FK (games.id)|-|
|cell|int8|T|-|-|-|
|hp|int8|T|-|-|-|

## インデックス
|インデックス名|対象カラム|ユニーク|オプション|
|---|---|:-:|---|
|idx_p3_inputs_game_id|game_id|F|-|

## 運用
- レコードの追加のみ許可、変更・削除はしない
