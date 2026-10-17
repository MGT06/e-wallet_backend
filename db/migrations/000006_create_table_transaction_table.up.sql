CREATE TABLE transaction (
  id             INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  account_id     INT              NOT NULL REFERENCES accounts(id),
  type           transaction_type NOT NULL,
  amount         INT              NOT NULL,
  balance_before INT              NOT NULL ,
  balance_after  INT              NOT NULL CHECK (balance_after = balance_before + amount),
  topup_id       INT              REFERENCES topups(id),
  transfer_id    INT              REFERENCES transfers(id),
  delete_at      TIMESTAMPTZ,
  created_at     TIMESTAMPTZ      NOT NULL DEFAULT now(),
  CONSTRAINT chk_transaction_sign CHECK (
    (type = 'topup'    AND amount > 0) OR
    (type = 'transfer' AND amount < 0)
  )
);