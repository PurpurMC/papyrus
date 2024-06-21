create table metadata
(
    id       uuid         not null,
    name     varchar(255) not null,
    value    text         not null,
    build_id uuid         not null,
    primary key (id),
    constraint UniqueNameAndBuild unique (build_id, name)
);

alter table metadata
    add constraint fkey_metadata_build
        foreign key (build_id)
            references build;
