create table project
(
    id   uuid not null,
    name varchar(255) unique,
    primary key (id)
);

create table version
(
    id         uuid not null,
    name       varchar(255) unique,
    project_id uuid,
    primary key (id),
    constraint UniqueNameAndProject unique (project_id, name)
);


create table build
(
    id         uuid not null,
    name       varchar(255) unique,
    result     varchar(255) check (result in ('SUCCESS', 'FAILURE')),
    timestamp  bigint,
    duration   bigint,
    hash       varchar(255),
    version_id uuid,
    primary key (id),
    constraint UniqueNameAndVersion unique (version_id, name)
);


create table commit
(
    id          uuid not null,
    author      varchar(255),
    email       varchar(255),
    description varchar(255),
    hash        varchar(255),
    timestamp   bigint,
    build_id    uuid,
    primary key (id)
);


create table file
(
    id             uuid not null,
    content_type   varchar(255),
    file_extension varchar(255),
    build_id       uuid unique,
    primary key (id)
);


create table creation_state
(
    id             uuid not null,
    file_extension varchar(255),
    build_id       uuid unique,
    primary key (id)
);

alter table version
    add constraint fkey_version_project
        foreign key (project_id)
            references project;

alter table build
    add constraint fkey_build_version
        foreign key (version_id)
            references version;

alter table commit
    add constraint fkey_commit_build
        foreign key (build_id)
            references build;

alter table file
    add constraint fkey_file_build
        foreign key (build_id)
            references build;

alter table creation_state
    add constraint fkey_creation_state_build
        foreign key (build_id)
            references build;
