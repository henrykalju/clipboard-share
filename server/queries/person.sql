-- name: GetPerson :one
select
    *
from person
where username = $1;

-- name: InsertPerson :exec
insert into person (
    username,
    password
) values (
    $1,
    $2
) returning *;

-- name: UpdatePassword :exec
update person set
password = $2
where username = $1;