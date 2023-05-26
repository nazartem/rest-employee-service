CREATE companies
(
    id   bigserial not null primary key,
    name varchar not null unique
);

CREATE departments
(
    id     bigserial not null primary key,
    company_id INT,
    name   varchar not null,
    phone varchar not null unique,

    CONSTRAINT company_id FOREIGN KEY (company_id) REFERENCES public.companies (id)
);

CREATE employees
(
    id        bigserial not null primary key,
    name      varchar not null,
    surname   varchar not null,
    phone     varchar not null unique,
    company_id INT,
    passport_type varchar not null,
    passport_number varchar not null unique,
    department_id INT,

    CONSTRAINT company_id FOREIGN KEY (company_id) REFERENCES public.companies (id),
    CONSTRAINT department_id FOREIGN KEY (department_id) REFERENCES public.departments (id)
);

INSERT INTO companies (name) VALUES ('Astra'), ('Google'), ('Apple'), ('Sber'), ('Yandex');