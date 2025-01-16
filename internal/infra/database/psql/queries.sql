-- name: ListAll :many
SELECT * from orders;

-- name: Save :exec
INSERT INTO orders (id, price, tax, final_price) VALUES ($1, $2, $3, $4);