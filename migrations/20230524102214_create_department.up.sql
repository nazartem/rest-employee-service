CREATE TABLE departments
(
    id     bigserial not null primary key,
    company_id INT,
    name   varchar not null,
    phone varchar not null unique,

    CONSTRAINT company_id FOREIGN KEY (company_id) REFERENCES public.companies (id)
);