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
    created_at datetime not null default (datetime('now')),
    foreign key(owner_id) references users (id) on delete cascade
);

drop table if exists lobby_users;
create table lobby_users (
    lobby_id text,
    user_id text,
    primary key(lobby_id, user_id),
    foreign key(lobby_id) references lobbies (id) on delete cascade,
    foreign key(user_id) references users (id) on delete cascade
);

