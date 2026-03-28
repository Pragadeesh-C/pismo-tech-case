CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    operation_type SMALLINT NOT NULL CHECK (operation_type BETWEEN 1 AND 4),
    amount NUMERIC(12,2) NOT NULL CHECK (amount <> 0),
    event_date TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_account_id
ON transactions(account_id);