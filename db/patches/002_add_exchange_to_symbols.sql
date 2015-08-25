begin;

alter table symbols add exchange character varying (100) not null;
alter table symbols add company_name character varying (255) not null;
alter table symbols add constraint symbols_symbol_exchange_unique unique (symbol, exchange);
create unique index on symbols (symbol);

commit;
