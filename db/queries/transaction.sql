-- name: CreateTransaction :one
INSERT INTO transactions (account_id, operation_type, amount, balance)
VALUES ($1, $2, $3, $4)
RETURNING id, account_id, operation_type, amount, balance, event_date;

-- name: FetchAllDebitTransactionsByAccountID :many
SELECT * from transactions
WHERE account_id = $1 AND balance < 0
ORDER BY event_date, id
FOR UPDATE;

-- name: UpdateTransactionByID :exec
UPDATE transactions 
SET balance = $2
WHERE id = $1;