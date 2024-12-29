insert into users (id, email, name, password_hash, token)
values (
    '67b31dee-ac79-453a-a68b-75ef55da5a52',
    'bob@aol.com',
    'Bob Smith',
    readfile('password_hash.dat'),
    'kihLmzc0hrW2IzkYko3sJlIXXDv5h/EwkC6DgHzMfeg='
);

insert into users (id, email, name, password_hash, token)
values (
    '67b31dee-ac79-453a-a68b-75ef55da5a53',
    'tom@aol.com',
    'Tom Carter',
    readfile('password_hash.dat'),
    'kihLmzc0hrW2IzkYko3sJlIXXDv5h/EwkC6DgHzMfeh='
);

insert into lobbies (id, name, owner_id)
values (
    'ab7b8a66-c6cd-48a1-9d02-ecafe255ed3e',
    "Bob's Lobby",
    '67b31dee-ac79-453a-a68b-75ef55da5a52'
);

insert into lobby_users (lobby_id, user_id) 
values ('ab7b8a66-c6cd-48a1-9d02-ecafe255ed3e','67b31dee-ac79-453a-a68b-75ef55da5a52');

insert into lobby_users (lobby_id, user_id) 
values ('ab7b8a66-c6cd-48a1-9d02-ecafe255ed3e','67b31dee-ac79-453a-a68b-75ef55da5a53');
