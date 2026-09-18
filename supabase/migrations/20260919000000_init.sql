-- Initial schema.
--
-- The Go backend connects to Postgres directly with pgx rather than through
-- PostgREST, so every table must rely on row-level security for
-- authorization. db.AsUser sets request.jwt.claims and the authenticated role
-- for the duration of a transaction, which is what makes auth.uid() resolve
-- below.

create table public.profiles (
  id uuid primary key references auth.users (id) on delete cascade,
  display_name text not null check (char_length(display_name) between 1 and 60),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

comment on table public.profiles is 'Public profile data, one row per auth user.';

alter table public.profiles enable row level security;

create policy "profiles are readable by any authenticated user"
  on public.profiles
  for select
  to authenticated
  using (true);

create policy "users insert their own profile"
  on public.profiles
  for insert
  to authenticated
  with check (auth.uid() = id);

create policy "users update their own profile"
  on public.profiles
  for update
  to authenticated
  using (auth.uid() = id)
  with check (auth.uid() = id);

create policy "users delete their own profile"
  on public.profiles
  for delete
  to authenticated
  using (auth.uid() = id);

-- Keep updated_at honest without trusting the client to send it.
create function public.set_updated_at()
returns trigger
language plpgsql
security definer
set search_path = ''
as $$
begin
  new.updated_at = now();
  return new;
end;
$$;

create trigger profiles_set_updated_at
  before update on public.profiles
  for each row
  execute function public.set_updated_at();
