-- name: GetItemsByPerson :many
select
    *
from item
where person_id = $1;

-- name: GetItemByIdAndPerson :one
select
    *
from item
where person_id = $1
    and id = $2;

-- name: GetDataByItem :many
select
    *
from data
where item_id = $1;

-- name: InsertItem :one
insert into item (person_id, content, type)
values ($1, $2, $3) returning *;

-- name: InsertData :one
insert into data (item_id, format, data)
values ($1, $2, $3) returning *;

-- name: DeleteItem :exec
delete from item
where id = $1;