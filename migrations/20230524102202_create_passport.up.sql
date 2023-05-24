CREATE TABLE passports
(
    id     bigserial not null primary key,
    type   varchar not null,
    number varchar not null unique
);