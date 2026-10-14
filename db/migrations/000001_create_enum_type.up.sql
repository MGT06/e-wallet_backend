CREATE TYPE status_transaction AS ENUM ('pending', 'success', 'failed');
CREATE TYPE transaction_type AS ENUM ('topup', 'transfer');
CREATE TYPE payment_method_type AS ENUM ('bank', 'ewallet');