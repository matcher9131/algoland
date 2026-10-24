# p3_patterns - 問題「ダンジョン」のパターン

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|game_id|uuid|T|-|PK, FK (games.id)|-|
|pattern|int8|T|-|-|-|

## 運用
- レコードの追加のみ許可、変更・削除はしない
