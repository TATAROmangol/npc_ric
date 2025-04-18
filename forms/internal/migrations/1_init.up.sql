create table institutions (
    id serial primary key,
    name text not null,
    inn int not null,
    columns text not null
);

create table users (
    id serial primary key,
    info text not null,
    institution_id int not null references institutions(id),
    mentor_id text not null references mentors(id)
);

create table documents (
    institution_id int primary key references institutions(id),
    users_id int[] not null 
);