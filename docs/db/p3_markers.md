# p3_markers - 問題「ダンジョン」におけるマーキング

## カラム

|物理名|データ型|Not null|Default|Key|備考|
|---|---|:-:|:-:|---|---|
|game_id|uuid|T|-|PK<br>FK (games.id)|複合主キー|
|cell|int8|T|-|PK|複合主キー|
|type|enum (p3_markers_type)|T|-|-|-|

## 列挙体

|p3_markers_type|
|---|
|down_arrow|
|right_arrow|
|up_arrow|
|left_arrow|
|red_flag|
|green_flag|
|blue_flag|
|yellow_flag|
|purple_flag|

## 運用
- 追加・更新・削除すべて可