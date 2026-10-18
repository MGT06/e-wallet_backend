CREATE TABLE topups (
  id                INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  reference_no      VARCHAR(32)    NOT NULL UNIQUE,
  account_id        INT            NOT NULL REFERENCES accounts(id),
  payment_method_id INT            NOT NULL REFERENCES payment_methods(id),
  amount            INT            NOT NULL CHECK (amount > 0),
  tax               INT            NOT NULL DEFAULT 0 CHECK (tax >= 0),
  total_amount      INT            NOT NULL CHECK (total_amount >= amount + tax),
  xendit_invoice_id VARCHAR(100)   UNIQUE,
  payment_url       TEXT,
  expires_at        TIMESTAMPTZ,
  status            status_topups  NOT NULL DEFAULT 'pending',
  paid_at           TIMESTAMPTZ,
  created_at        TIMESTAMPTZ    NOT NULL DEFAULT now()
);