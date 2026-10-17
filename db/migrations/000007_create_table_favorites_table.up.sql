CREATE TABLE favorites (
  account_id          INT  NOT NULL REFERENCES accounts(id),
  favorite_account_id INT  NOT NULL REFERENCES accounts(id)  CHECK (account_id <> favorite_account_id),
  created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (account_id, favorite_account_id)
);