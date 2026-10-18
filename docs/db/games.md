# games - ゲームセッション

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|created_at|TIMESTAMPZ|T|now()|-|-|
|id|uuid|T|uuidv7()|PK|自動採番|
|pattern|int8|F|NULL|-|問題「ダンジョン」におけるパターン番号（BEで使用）|
|problem_id|int8|T|-|FK (problems.id)|-|

## 運用
- レコードの追加のみ許可、変更・削除はしない
