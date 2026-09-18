import Link from "next/link";

const stack = [
  { layer: "Frontend", detail: "Next.js 16.3.5 / React 19.3.0 / Tailwind CSS 4.3.3" },
  { layer: "Backend", detail: "Go 1.27.1 container image, reached over a Vercel service binding" },
  { layer: "Database", detail: "Supabase Postgres via pgx 5.11.0" },
  { layer: "Auth", detail: "Supabase Auth, JWT verified in Go against the project JWKS" },
];

export default function Home() {
  return (
    <main className="mx-auto flex w-full max-w-2xl flex-1 flex-col justify-center gap-8 px-6 py-16">
      <header className="flex flex-col gap-2">
        <h1 className="font-mono text-3xl font-semibold tracking-tight">algoland</h1>
        <p className="text-sm text-foreground/60">
          Next.js フロントエンドと Go バックエンドを 1 つの Vercel プロジェクトにまとめた構成です。
        </p>
      </header>

      <dl className="flex flex-col divide-y divide-foreground/10 border-y border-foreground/10">
        {stack.map(({ layer, detail }) => (
          <div key={layer} className="flex flex-col gap-1 py-3 sm:flex-row sm:gap-6">
            <dt className="w-24 shrink-0 text-sm font-medium">{layer}</dt>
            <dd className="text-sm text-foreground/70">{detail}</dd>
          </div>
        ))}
      </dl>

      <Link
        href="/status"
        className="self-start rounded-md border border-foreground/20 px-4 py-2 text-sm transition-colors hover:bg-foreground/5"
      >
        バックエンドの疎通を確認する →
      </Link>
    </main>
  );
}
