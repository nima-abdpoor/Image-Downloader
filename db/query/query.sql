-- name: CreateImageResult :one
INSERT INTO image_result (query_id,
                          image_url,
                          image_data)
values ($1, $2, $3) RETURNING *;

-- name: CreateQuery :one
INSERT INTO query (query,
                   status,
                   per_page,
                   page)
values ($1, $2, $3, $4) RETURNING *;

-- name: GetQueryByStatus :many
SELECT *
FROM query
WHERE status = $1;

-- name: UpdateQuery :one
UPDATE query
set status = $2, updated_at = $3
WHERE id = $1 RETURNING *;