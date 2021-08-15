-- name: CreateSet :one
INSERT INTO sets (set_uid)
VALUES ($1)
RETURNING *;

-- name: CreateSetData :exec
INSERT INTO sets_datas(data, set_uid)
VALUES ($1, $2);

-- name: GetSetDatas :many
SELECT * 
FROM sets_datas
WHERE sets_datas.set_uid = $1;

-- name: ListSetsDatas :many
SELECT *
FROM sets_datas
ORDER BY set_uid;

-- name: SetIntersectingSet :exec
INSERT INTO intersecting_sets(set_uid, intersectingset_uid)
VALUES ($1, $2);

-- name: GetIntersectingSet :one
SELECT * FROM intersecting_sets
WHERE set_uid = $1
ORDER BY intersectingset_uid;

-- name: ListIntersectingSets :many
SELECT * FROM intersecting_sets
ORDER BY set_uid;