drop table if exists users;
create table users(
  id integer primary key default(uuid()),
  email text not null unique,
  name text not null,
  password_hash string not null,
  created_at datetime not null default(datetime('now'))
);

drop table if exists matches;
create table matches(
  id uuid primary key default(uuid()),
  name text not null,
  created_at datetime not null default(datetime('now'))
);

drop table if exists players;
create table players(
  player_id uuid,
  match_id uuid,
  primary key(player_id, match_id)
)
