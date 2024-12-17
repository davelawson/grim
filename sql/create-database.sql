drop table if exists users;
create table users (
    id text primary key,
    email text not null unique,
    name text not null,
    password_hash binary(44) not null,
    created_at datetime not null default (datetime('now')),
    token text default (null)
);

drop table if exists matches;
create table matches (
    id text primary key,
    name text not null,
    created_at datetime not null default (datetime('now'))
);

drop table if exists players;
create table players (
    player_id text,
    match_id text,
    primary key (player_id, match_id)
);

drop table if exists lobbies;
create table lobbies (
    id text primary key,
    name text not null,
    owner_id text
    created_at datetime not null default (datetime('now'))
);
