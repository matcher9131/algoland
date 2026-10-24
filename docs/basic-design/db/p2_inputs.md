# p2_inputs - 問題「ポンプ」のユーザ入力

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|created_at|TIMESTAMPZ|T|now()|-|-|
|game_id|uuid|T|-|FK (games.id), UNIQUE|-|
|value|text|T|-|-|左上マスから順にポンプを`P`、水場を`W`、空き地を`.`と表した文字列|

## 運用
- レコードの追加のみ許可、変更・削除はしない
