create table if not exists t_roles (
    code varchar(255) primary key not null,
    name varchar(255) not null
);

insert into t_roles values ('USER', 'Пользователь'), ('BLOGGER', 'Блоггер'), ('MODERATOR', 'Модератор'), ('ADMIN', 'Администратор');

create table if not exists t_users (
    login varchar(255) primary key not null,
    password varchar(255) not null,
    email varchar(255) not null unique,
    first_name varchar(255) not null,
    middle_name varchar(255),
    last_name varchar(255) not null,
    birthdate date not null,
    role varchar(255) not null references t_roles(code) default 'USER'
);

create table if not exists t_blogs (
    id bigserial primary key not null,
    title varchar(255) not null,
    content varchar(500),
    author varchar(255) not null references t_users(login)
)