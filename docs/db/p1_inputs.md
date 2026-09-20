# p1_inputs - 問題「呑兵衛」のユーザ入力

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|created_at|TIMESTAMPZ|T|now()|-|-|
|game_id|uuid|T|-|FK (games.id)|-|
|value|text|T|-|-|-|

## インデックス
|インデックス名|対象カラム|ユニーク|オプション|
|---|---|:-:|---|
|idx_p1_inputs_game_id|game_id|F|-|

## 運用
- レコードの追加のみ許可、変更・削除はしない
