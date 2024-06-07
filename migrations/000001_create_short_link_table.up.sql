create table if not exists short_link (
    path varchar(12) primary key,
    link text not null
);
