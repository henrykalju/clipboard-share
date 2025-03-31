create type clipboard_type as enum ('X11', 'WINDOWS');

create table item (
    id serial primary key not null,
    person_id integer references person(id) not null,
    type clipboard_type not null,
    content text not null,
    created_at timestamp default current_timestamp not null
);