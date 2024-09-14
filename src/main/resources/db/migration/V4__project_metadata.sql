create table project_metadata
(
    id       uuid         not null,
    name     varchar(255) not null,
    p_value  text         not null, // value is a reserved value
    project_id uuid         not null,
    primary key (id),
    constraint UniqueMetadataNameAndProject unique (project_id, name)
);

alter table project_metadata
    add constraint fkey_project_metadata_build
        foreign key (project_id)
            references project;
