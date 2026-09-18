-- Applied after the migrations, mirroring the table privileges Supabase hands
-- to the API roles. Without these, a policy could allow a row that the role
-- still has no privilege to touch, and the test would fail for the wrong
-- reason.

grant all on all tables in schema public to anon, authenticated, service_role;
grant all on all sequences in schema public to anon, authenticated, service_role;
grant all on all functions in schema public to anon, authenticated, service_role;
