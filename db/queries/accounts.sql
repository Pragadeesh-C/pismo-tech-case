-- name: CreateAccount :one
INSERT INTO accounts (document_number)
VALUES ($1)
RETURNING id, document_number, created_at;