create table data (
    id serial primary key not null,
    item_id integer references item(id) on delete cascade not null,
    format text not null,
    data bytea not null
);