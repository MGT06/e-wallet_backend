CREATE TYPE status_transfer AS ENUM ('pending', 'success', 'failed');
CREATE TYPE status_topups AS ENUM ('pending', 'success', 'failed', 'expired');
CREATE TYPE transaction_type AS ENUM ('topup', 'transfer');
CREATE TYPE payment_method_type AS ENUM ('bank', 'ewallet');