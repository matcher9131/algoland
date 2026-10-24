# API（ダンジョン）

## 型

### DungeonResponse

ダンジョンのゲーム状態

```ts
type DungeonResponse = {
    // ゲーム結果
    readonly result: "ongoing" | "correct" | "incorrect";

    // 残りHP
    readonly hp: number; // 正確にはint

    // 座標
    readonly x: number; // 正確にはint
    readonly y: number; // 正確にはint

    // どの方向に壁があるか（Lv1, 2のみで使用）
    readonly wall: number; // 正確には4bitの整数（下位bitから順に左、下、上、右）

    // TODO: Lv3の迷路本体、倒したモンスター
}
```
