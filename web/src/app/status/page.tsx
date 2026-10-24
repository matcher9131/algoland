import { connection } from "next/server";

import { apiFetch } from "@/lib/api";

type Health = { status: string; env: string };

async function readHealth(): Promise<
  { ok: true; body: Health } | { ok: false; message: string }
> {
  try {
    const response = await apiFetch("/healthz");
    if (!response.ok) {
      return {
        ok: false,
        message: `バックエンドが ${response.status} を返しました`,
      };
    }
    return { ok: true, body: (await response.json()) as Health };
  } catch (error) {
    return {
      ok: false,
      message: error instanceof Error ? error.message : String(error),
    };
  }
}

export default async function StatusPage() {
  // The service binding resolves only at request time, so this page must
  // never be prerendered.
  await connection();

  const result = await readHealth();

  return (
    <main className="mx-auto flex w-full max-w-2xl flex-1 flex-col justify-center gap-6 px-6 py-16">
      <h1 className="font-mono text-2xl font-semibold tracking-tight">
        Backend status
      </h1>

      {result.ok ? (
        <dl className="flex flex-col divide-y divide-foreground/10 border-y border-foreground/10">
          <div className="flex gap-6 py-3">
            <dt className="w-24 shrink-0 text-sm font-medium">status</dt>
            <dd className="font-mono text-sm text-green-600 dark:text-green-400">
              {result.body.status}
            </dd>
          </div>
          <div className="flex gap-6 py-3">
            <dt className="w-24 shrink-0 text-sm font-medium">env</dt>
            <dd className="font-mono text-sm text-foreground/70">
              {result.body.env}
            </dd>
          </div>
        </dl>
      ) : (
        <p className="rounded-md border border-red-500/30 bg-red-500/5 px-4 py-3 font-mono text-sm text-red-600 dark:text-red-400">
          {result.message}
        </p>
      )}
    </main>
  );
}
