CREATE TABLE IF NOT EXISTS stores (
  id uuid PRIMARY KEY,
  created_at timestamp with time zone NULL,
  updated_at timestamp with time zone NULL,
  name character varying NOT NULL,
  status character varying(20) NOT NULL,
  description character varying(255) NULL,
  account_id uuid NOT NULL,
  category_id uuid NOT NULL,
  user_id uuid NOT NULL,
  tags text [] NULL,
  lat numeric(10, 8) NULL,
  lng numeric(11, 8) NULL
);

CREATE TABLE IF NOT EXISTS categories (
  id uuid PRIMARY KEY,
  created_at timestamp with time zone NULL,
  updated_at timestamp with time zone NULL,
  name character varying NOT NULL,
  status character varying(20) NULL
);

CREATE TABLE IF NOT EXISTS accounts (
  id uuid PRIMARY KEY,
  created_at timestamp with time zone NULL,
  updated_at timestamp with time zone NULL,
  balance numeric(20, 8) NULL
);

ALTER TABLE stores ADD FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE stores ADD FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE RESTRICT ON UPDATE CASCADE;

CREATE INDEX ON stores (account_id);
CREATE INDEX ON stores (category_id);
CREATE INDEX ON stores (user_id);

COMMENT ON COLUMN stores.status IS 'must be pending, active, disable or block';
COMMENT ON COLUMN stores.lat IS 'must be latitude';
COMMENT ON COLUMN stores.lng IS 'must be longitude';
COMMENT ON COLUMN categories.status IS 'must be pending, active or disable';
COMMENT ON COLUMN accounts.balance IS 'must be a positive';
