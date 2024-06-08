alter table build add column ready int not null default (1) check (ready IN (0, 1));
