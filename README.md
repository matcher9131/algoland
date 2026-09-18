# algoland

Next.js フロントエンドと Go バックエンドを **1 つの Vercel プロジェクト**にまとめ、Supabase をデータベースとして使う構成です。

## 構成

```
algoland/
├── vercel.json              # Vercel Services 定義（web / api とルーティング）
├── Makefile                 # 開発タスク
├── package.json             # Supabase CLI / Vercel CLI をピン留め
├── web/                     # Next.js（ネイティブランタイム）
│   ├── e2e/                 # Playwright
│   ├── vitest.config.mts    # browser / server の 2 プロジェクト構成
│   └── src/
│       ├── app/
│       ├── lib/
│       │   ├── api.ts       # service binding 経由で Go を呼ぶヘルパー
│       │   └── supabase/    # browser / server / proxy クライアント
│       ├── test/            # MSW サーバー
│       └── proxy.ts         # セッション更新（旧 middleware.ts）
├── api/                     # Go（コンテナイメージ）
│   ├── Dockerfile.vercel
│   ├── cmd/server/
│   └── internal/
│       ├── auth/            # Supabase JWT を JWKS で検証
│       ├── config/
│       ├── db/              # pgx プール + RLS スコープ実行
│       ├── handler/         # ルーティング（依存はインターフェース）
│       ├── httpx/
│       ├── store/           # handler.Store の Postgres 実装
│       └── testdb/          # testcontainers による使い捨て DB
└── supabase/
    └── migrations/
```

### デプロイの形

`web` は Next.js のネイティブランタイムで動くため、プリレンダリング済みページの CDN 配信・ISR・Image Optimization がそのまま効きます。`api` はコンテナイメージとしてビルドされ、**公開ルートを持ちません**。`vercel.json` の service binding によって `web` のサーバーサイドにだけ `API_URL` が注入され、そこからのみ到達できます。

```
インターネット ──> [ web : Next.js ] ──binding(API_URL)──> [ api : Go コンテナ ]
                                                                    │
                                                                    └──> Supabase Postgres (pgx)
```

## バージョン

すべて 2026-09-19 時点の最新安定版で固定しています。

| 対象 | バージョン | 備考 |
|---|---|---|
| Next.js | 16.3.5 | |
| React / React DOM | 19.3.0 | |
| TypeScript | 6.0.3 | 下記「TypeScript」参照 |
| Tailwind CSS | 4.3.3 | |
| ESLint | 9.39.5 | 下記「ESLint」参照 |
| eslint-config-next | 16.3.5 | |
| @supabase/supabase-js | 2.116.0 | |
| @supabase/ssr | 0.12.7 | |
| Go | 1.27.1 | |
| jackc/pgx | v5.11.0 | |
| golang-jwt/jwt | v5.3.1 | |
| MicahParks/keyfunc | v3.8.2 | |
| Supabase CLI | 2.117.0 | |
| Vercel CLI | 59.23.2 | |
| Vitest | 5.0.1 | |
| React Testing Library | 16.3.3 | |
| MSW | 2.15.0 | |
| Playwright | 1.63.0 | |
| stretchr/testify | v1.12.1 | |
| testcontainers-go | v0.44.0 | |

npm は両方の `package.json` で完全一致のバージョンを指定し、`.npmrc` に `save-exact=true` を置いています。Go は `go.mod` と `Dockerfile.vercel` の両方で 1.27.1 を指定しています。

### TypeScript

npm の `latest` は **7.0.2**（Go 移植版コンパイラ）ですが、`eslint-config-next@16.3.5` が依存する `typescript-eslint@8.70.0` の peer 要件が `typescript >=4.8.4 <6.1.0` のため、7 系を入れると lint が動作しません。安定版リリースは 5.9.3 → 6.0.2 / 6.0.3 → 7.0.x という並びなので、**ツールチェーン全体が通る最新安定版として 6.0.3** を採用しています。`typescript-eslint` が 7 系に対応したら引き上げてください。

### ESLint

`latest` は **10.11.0** ですが、`eslint-config-next@16.3.5` が同梱する `eslint-plugin-react` / `eslint-plugin-import` / `eslint-plugin-jsx-a11y` が ESLint 10 未対応で、実行すると `contextOrFilename.getFilename is not a function` で落ちます。そのため **9.39.5**（9 系の最新）に固定しています。

## セットアップ

### 1. 依存関係

```bash
make setup
```

Go は `~/.local/go` にインストール済みです。PATH を通してください。

```bash
echo 'export PATH="$HOME/.local/go/bin:$PATH"' >> ~/.bashrc && source ~/.bashrc
```

### 2. Supabase

ローカルスタックを起動し、マイグレーションを適用します（Docker が必要です）。

```bash
make db-start    # 出力される API URL と anon key を控える
make db-reset    # supabase/migrations を適用
```

### 3. 環境変数

```bash
cp web/.env.example web/.env.local
cp api/.env.example api/.env
```

`make db-start` が出力した値を埋めてください。本番では Vercel のプロジェクト設定に同じ変数を登録します（`API_URL` は binding が自動注入するので**登録不要**）。

### 4. 開発サーバー

```bash
make dev         # vercel dev で web と api を同時起動（binding も再現される）
```

個別に動かす場合：

```bash
make dev-api     # :8080
make dev-web     # :3000
```

## テスト

| レイヤ | ツール | 何を検証するか |
|---|---|---|
| Go ユニット | 標準 testing + testify 1.12.1 | 設定読み込み、JWT ミドルウェア、HTTP ハンドラ（フェイク使用・DB 不要） |
| Go 統合 | testcontainers-go 0.44.0 | **RLS ポリシーの強制**、`AsUser` のスコープと漏洩防止、トリガ、pgx プール |
| FE ユニット | Vitest 5.0.1 + React Testing Library 16.3.3 | Client / 同期 Server Component の描画 |
| FE サーバー側 | Vitest (node) + MSW 2.15.0 | `apiFetch` の URL 組み立て・トークン付与・キャッシュ無効化 |
| E2E | Playwright 1.63.0 | async Server Component、proxy、Go バックエンドとの結合 |

```bash
make test             # Go + FE のユニット/統合（E2E 以外）
make test-api         # Go のみ（Postgres コンテナを自動起動）
make test-api-short   # Docker 不要の Go テストのみ
make test-web         # Vitest のみ
make coverage         # 両方のカバレッジ
```

### Go の統合テスト

`api/internal/testdb` が `postgres:17-alpine` のコンテナを起動し、Supabase の shim（`auth` スキーマ・`auth.uid()`・ロール）を当ててから `supabase/migrations` を順に適用します。**ローカルの Supabase スタックは不要**で、`go test ./...` だけで完結します。テストごとにテーブルは truncate され、コンテナはテストバイナリ終了時に testcontainers の reaper が破棄します。

RLS は Go ではなく Postgres が実行するものなので、認可の検証は必ずこの実 DB 経由で行っています。

### E2E

Vitest は **async Server Component をテストできない**（Next.js 公式の制約）ため、`/status` のようなページは Playwright が担当します。

```bash
export $(npx supabase status -o env | tr -d '"' | xargs)
make test-e2e
```

`playwright.config.ts` が Next.js と Go サーバーの両方を起動します。`DATABASE_URL` と `SUPABASE_URL` が無い場合はバックエンド依存の spec が自動で skip され、フロントだけの spec は動きます。

#### 初回のみ: Chromium のシステムライブラリ

Chromium の実行には `libnss3` / `libnspr4` / `libasound2` が必要です。公式手順 (`sudo npx playwright install-deps chromium`) は root を要求し、かつ headed モード用の X 関連パッケージまで 20 以上入れますが、**headless で実際に不足するのは 4 ファイルだけ**です。root 無しで済ませられます。

```bash
make e2e-deps
```

これは 3 つの .deb（計 2MB）を `apt-get download` し、`~/.local/playwright-deps` に展開します。`playwright.config.ts` はこのディレクトリがあれば自動で `LD_LIBRARY_PATH` に加えるため、以後は追加の環境変数なしで `make test-e2e` が動きます。システムパッケージが入っている環境や CI では何もしません。

## デプロイ前に必要な手動設定

以下は自動化できないため、Vercel のダッシュボードで設定してください。

1. **Beta 機能の有効化** — この構成は [Services](https://vercel.com/docs/services) と [Container Images](https://vercel.com/docs/functions/container-images) を使います。どちらも Beta で、アカウント側での有効化が必要です。
2. **`PORT=8080` の設定** — Vercel はコンテナへのトラフィックを既定でポート 80 に流します。distroless の非 root ユーザーは 80 を bind できないため、プロジェクトの環境変数に `PORT=8080` を設定してください（`Dockerfile.vercel` の `ENV PORT` は参照されません）。
3. **環境変数** — `NEXT_PUBLIC_SUPABASE_URL` / `NEXT_PUBLIC_SUPABASE_ANON_KEY` / `DATABASE_URL` / `SUPABASE_URL`。`DATABASE_URL` は Supabase の Connection Pooler（ポート 6543）の URI を使ってください。

## 設計上の注意

- **binding はランタイム限定** — `API_URL` はビルド時にも proxy（middleware）でも解決されません。静的生成中に Go を呼ぶページは作れないため、`/status` では `connection()` でプリレンダリングを打ち切っています。
- **RLS** — Go は PostgREST を経由せず pgx で直接 Postgres に繋ぐため、認可はすべて RLS が担います。`db.AsUser` がトランザクション内で `request.jwt.claims` と `authenticated` ロールを `SET LOCAL` 相当で設定するので、`auth.uid()` を参照するポリシーがそのまま機能します。接続プールの使い回しでも設定は漏れません。
- **セッション Cookie とキャッシュ** — `@supabase/ssr` 0.12.7 の `setAll` は第 2 引数でキャッシュ抑止ヘッダーを渡します。これを応答に反映しないと CDN が他ユーザーのセッションを配信しうるため、`web/src/lib/supabase/proxy.ts` で必ず設定しています。
- **Secure Compute / Static IP** — コンテナイメージでは未対応です。Supabase 側で IP 許可リストを使う予定がある場合はこの構成では成立しません。
- **ハンドラの依存はインターフェース** — `handler.Deps` は `Store` と `auth.Verifier` を受け取ります。具象の `*pgxpool.Pool` や JWKS 実装に依存していないため、ハンドラのテストは DB もネットワークも使いません。
- **`npm audit`** — ルートの警告は `vercel` CLI が内部で持つ `undici` に由来する開発専用の依存です。`npm audit fix --force` は Vercel CLI を 54 系へダウングレードするため、実行しないでください。

## コマンド一覧

```bash
make help
```
