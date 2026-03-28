-- name: CreateAccount :one
INSERT INTO accounts (document_number)
VALUES ($1)
RETURNING id, document_number, created_at;

-- name: GetAccount :one
SELECT id, document_number, created_at
FROM accounts
WHERE id = $1;