-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = $1;

-- name: GetAccountForUpdate :one
SELECT *
FROM accounts
WHERE id = $1 FOR NO KEY
UPDATE;

-- name: ListAccountsByOwner :many
SELECT *
FROM accounts
WHERE owner = $1
ORDER BY id ASC LIMIT $2
OFFSET $3;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = $2
WHERE id = $1 RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = $1 RETURNING *;

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE id = $1;