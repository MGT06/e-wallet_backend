CREATE TABLE payment_methods (
  id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name        VARCHAR(100)         UNIQUE NOT NULL,
  type        payment_method_type  NOT NULL,
  xendit_code VARCHAR(50),
  icon_url    TEXT,
  tax_percent NUMERIC(5,2)         NOT NULL DEFAULT 0 CHECK (tax_percent >= 0),
  created_at  TIMESTAMPTZ          NOT NULL DEFAULT now()
);