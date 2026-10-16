CREATE TABLE transfers (
  id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  reference_no VARCHAR(32)        NOT NULL UNIQUE,
  sender_id    BIGINT             NOT NULL REFERENCES accounts(id),
  receiver_id  BIGINT             NOT NULL REFERENCES accounts(id),
  amount       BIGINT             NOT NULL CHECK (amount > 0),
  notes        VARCHAR(255),
  status       status_transfer NOT NULL DEFAULT 'pending',
  created_at   TIMESTAMPTZ        NOT NULL DEFAULT now(),
  CONSTRAINT chk_transfers_parties CHECK (sender_id <> receiver_id)
);