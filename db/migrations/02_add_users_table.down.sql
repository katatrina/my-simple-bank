ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT "unique_owner_currency";

ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT "accounts_owner_fkey";

DROP TABLE IF EXISTS "users";