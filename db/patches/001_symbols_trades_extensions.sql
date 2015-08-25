begin;

create extension if not exists hstore;
create extension if not exists pgcrypto;

create table symbols (
  id uuid primary key default gen_random_uuid(),
  symbol character varying(10) not null,
  properties hstore,
  created_at timestamp without time zone not null default (now() at time zone 'UTC'),
  updated_at timestamp without time zone not null default (now() at time zone 'UTC')
);

create table trades (
  symbol_id uuid,
  time timestamp without time zone not null,
  price numeric(15, 10) not null,
  volume integer not null,
  properties hstore,
  created_at timestamp without time zone not null default (now() at time zone 'UTC'),
  updated_at timestamp without time zone not null default (now() at time zone 'UTC')
);

alter table trades add primary key (symbol_id, time);
alter table trades add foreign key (symbol_id) references symbols (id) on delete cascade;

commit;
