CREATE TABLE accounts (
  id                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  fullname          VARCHAR(100)  NOT NULL,
  image_url         TEXT,
  phone             VARCHAR(15)   UNIQUE,
  email             VARCHAR(255)  NOT NULL UNIQUE,
  password          VARCHAR(255)  NOT NULL,
  pin               VARCHAR(255),
  balance           BIGINT        NOT NULL DEFAULT 0 CHECK (balance >= 0),
  created_at        TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ   NOT NULL DEFAULT now()
);