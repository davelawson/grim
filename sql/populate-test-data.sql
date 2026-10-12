-- Test users share password123; each scrypt hash uses that user's email as salt.
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
    x'eea36129f29812d5481acee5d34548ed56ddafd67891b4d9c8fe01db26607276',
    'kihLmzc0hrW2IzkYko3sJlIXXDv5h/EwkC6DgHzMfeh='
);

insert into users (id, email, name, password_hash, token)
values (
    '67b31dee-ac79-453a-a68b-75ef55da5a54',
    'jim@aol.com',
    'Jimmy John',
    x'2d7480f129188cbd393f5bf6596f495280bd0e428eb27e60ed5d4347a0b965c0',
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
