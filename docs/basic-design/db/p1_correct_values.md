# p1_correct_values - 問題「呑兵衛」の正解値

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|game_id|uuid|T|-|PK, FK (games.id)|-|
|correct_value|double|T|-|-|-|

## 運用
- レコードの追加のみ許可、変更・削除はしない

## 備考
- 呑兵衛Lv2でのみ使用する
  - Lv1はユーザ入力により $X$ の値を変動させるためDBで保持しない
