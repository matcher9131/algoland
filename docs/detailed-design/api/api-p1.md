# API（呑兵衛）

## 型

### DrinkerResponse

呑兵衛のゲーム状態

```ts
type DrinkerResponse = {
    // ゲーム結果
    readonly result: "ongoing" | "correct" | "incorrect";

    // 残り回数
    readonly rest: number; // 正確にはint

    // 入力と出力の履歴
    readonly trials: readonly {
        readonly input: number;
        readonly output: boolean;
    }[];
}
```
