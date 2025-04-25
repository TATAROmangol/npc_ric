create table institutions (
    id serial primary key,
    name text not null,
    inn int not null,
    columns text not null
);

create table mentors (
    id serial primary key,
    name text not null
);

create table forms (
    id serial primary key,
    info text not null,
    institution_id int not null references institutions(id) ON DELETE CASCADE,
    mentor_id int not null references mentors(id) ON DELETE CASCADE
);