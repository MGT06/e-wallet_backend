CREATE TABLE topups (
  id                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  reference_no      VARCHAR(32)        NOT NULL UNIQUE,
  account_id        BIGINT             NOT NULL REFERENCES accounts(id),
  payment_method_id INT                NOT NULL REFERENCES payment_methods(id),
  amount            BIGINT             NOT NULL CHECK (amount > 0),
  tax               BIGINT             NOT NULL DEFAULT 0 CHECK (tax >= 0),
  total_amount      BIGINT             NOT NULL CHECK (total_amount >= amount),
  status            status_transaction NOT NULL DEFAULT 'pending',
  paid_at           TIMESTAMPTZ,
  created_at        TIMESTAMPTZ        NOT NULL DEFAULT now()
);