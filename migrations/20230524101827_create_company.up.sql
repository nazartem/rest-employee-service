CREATE TABLE companies
(
    id   bigserial not null primary key,
    name varchar not null unique
);