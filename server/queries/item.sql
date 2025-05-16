-- name: GetItemsByPerson :many
select
    *
from item
where person_id = $1
order by created_at;

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

-- name: CheckSizes :exec
delete from item i
where i.id in (
    select
        item_id
    from (
        select
            item_id
            ,size
            ,sum(size) over (order by item_id desc) as total
        from (
            select
                i.id as item_id
                ,sum(pg_column_size(d.data)) as size
            from item i
            join data d on i.person_id = $1
                and d.item_id = i.id
            group by i.id
            order by created_at desc
        )
    )
    where total > sqlc.arg(threshold)::int
    offset 1
);
