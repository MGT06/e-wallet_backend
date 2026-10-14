CREATE TABLE transaction (
  id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  account_id     BIGINT           NOT NULL REFERENCES accounts(id),
  type           transaction_type NOT NULL,
  amount         BIGINT           NOT NULL,
  balance_before BIGINT           NOT NULL ,
  balance_after  BIGINT           NOT NULL CHECK (balance_after = balance_before + amount),
  topup_id       BIGINT           REFERENCES topups(id),
  transfer_id    BIGINT           REFERENCES transfers(id),
  delete_at      TIMESTAMPTZ,
  created_at     TIMESTAMPTZ      NOT NULL DEFAULT now(),
  CONSTRAINT chk_transaction_sign CHECK (
    (type = 'topup'    AND amount > 0) OR
    (type = 'transfer' AND amount < 0)
  )
);