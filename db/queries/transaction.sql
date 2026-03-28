-- name: CreateTransaction :one
INSERT INTO transactions (account_id, operation_type, amount)
VALUES ($1, $2, $3)
RETURNING id, account_id, operation_type, amount, event_date;