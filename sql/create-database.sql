-- Explicit database reset only. The server never runs this script.
pragma foreign_keys = off;
drop table if exists launch_receipts;
drop table if exists players;
drop table if exists lobby_users;
drop table if exists matches;
drop table if exists lobbies;
drop table if exists users;
create table users (
    id text primary key,
    email text not null unique,
    name text not null,
    password_hash binary(44) not null,
    created_at datetime not null default (datetime('now')),
    token text default null,
    admin integer not null default 0 check (admin in (0, 1))
);

create table matches (
    id text primary key,
    lobby_id text not null unique,
    name text not null,
    created_at datetime not null default (datetime('now')),
    deleted_at datetime default null,
    status text not null check (status in ('setup', 'finished')),
    revision integer not null check (revision >= 0),
    outcome text,
    snapshot text not null,
    foreign key(lobby_id) references lobbies (id),
    check ((status = 'setup' and outcome is null) or
           (status = 'finished' and outcome is not null))
);

create table players (
    player_id text not null,
    match_id text not null,
    primary key (player_id, match_id),
    foreign key(player_id) references users (id),
    foreign key(match_id) references matches (id)
);

create table lobbies (
    id text primary key,
    name text not null,
    owner_id text not null,
    created_at datetime not null default (datetime('now')),
    deleted_at datetime default null,
    status text not null default 'open' check (status in ('open', 'closed')),
    match_id text unique,
    foreign key(owner_id) references users (id),
    foreign key(match_id) references matches (id),
    check ((status = 'open' and match_id is null) or
           (status = 'closed' and match_id is not null))
);

create table lobby_users (
    lobby_id text not null,
    user_id text not null,
    ready integer not null default 0 check (ready in (0, 1)),
    primary key(lobby_id, user_id),
    foreign key(lobby_id) references lobbies (id) on delete cascade,
    foreign key(user_id) references users (id)
);

create table launch_receipts (
    lobby_id text not null,
    actor_id text not null,
    request_id text not null,
    match_id text not null,
    command text not null,
    response text not null,
    primary key(lobby_id, actor_id, request_id),
    foreign key(lobby_id) references lobbies (id),
    foreign key(actor_id) references users (id),
    foreign key(match_id) references matches (id)
);

pragma foreign_keys = on;

