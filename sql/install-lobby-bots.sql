-- Additive installation for existing databases. Stop the server first.
pragma foreign_keys = on;
begin immediate;
create table if not exists lobby_bots (
    id text primary key,
    lobby_id text not null,
    name text not null,
    foreign key(lobby_id) references lobbies (id) on delete cascade
);
create index if not exists lobby_bots_lobby on lobby_bots(lobby_id);
commit;
