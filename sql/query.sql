-- name: CreateOrders :exec
INSERT INTO Orders(
        ID,
        UserId,
        Amount,
        AccountNumber,
        Currency,
        Description,
        Status
    )
VALUES (?, ?, ?, ?, ?, ?, ?);
-- name: GetOrderByUserId :one
SELECT *
FROM Orders
WHERE UserId = ?
    AND ID = ?
LIMIT 1;
-- name: GetOrderByOrderId :one
SELECT *
FROM Orders
WHERE ID = ?
LIMIT 1;
-- name: UpdateOrderStatus :exec
UPDATE Orders
SET Status = ?
where ID = ?;