create table person (
    id serial primary key not null,
    name varchar(20) not null
);

insert into person (name) values ('test');