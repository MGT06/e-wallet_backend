CREATE TABLE transfers (
  id           INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  reference_no VARCHAR(32)     NOT NULL UNIQUE,
  sender_id    INT             NOT NULL REFERENCES accounts(id),
  receiver_id  INT             NOT NULL REFERENCES accounts(id),
  amount       INT             NOT NULL CHECK (amount > 0),
  notes        VARCHAR(255),
  status       status_transfer NOT NULL DEFAULT 'pending',
  created_at   TIMESTAMPTZ        NOT NULL DEFAULT now(),
  CONSTRAINT chk_transfers_parties CHECK (sender_id <> receiver_id)
);