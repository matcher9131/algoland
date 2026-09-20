# p2_inputs - 問題「ポンプ」のユーザ入力

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|created_at|TIMESTAMPZ|T|now()|-|-|
|game_id|uuid|T|-|FK (games.id)|-|
|updated_at|TIMESTAMPZ|T|now()|-|-|
|value|text|T|-|-|-|

## 運用
- レコードの追加許可、削除禁止
- 更新はvalueカラムのみ可能、更新時にupdated_atを現在時刻で更新する

## 備考
他の問題と異なり時系列でのデータが重要ではないため上書きを基本方針とする