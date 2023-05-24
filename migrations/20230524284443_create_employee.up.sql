CREATE TABLE employees
(
    id        bigserial not null primary key,
    name      varchar not null,
    surname   varchar not null,
    phone     varchar not null unique,
    companyId INT,
    passportId INT,
    departmentId INT,

    CONSTRAINT companyId FOREIGN KEY (companyId) REFERENCES public.companies (id),
    CONSTRAINT passportId FOREIGN KEY (passportId) REFERENCES public.passports (id),
    CONSTRAINT departmentId FOREIGN KEY (departmentId) REFERENCES public.departments (id)
);