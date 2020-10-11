CREATE TABLE users (
    user_id             int PRIMARY KEY,
    email               varchar(129) NOT NULL, 
    verification        BOOLEAN NOT NULL DEFAULT FALSE,
    verification_code   int(10),
    UNIQUE(email)
);
/*
    email 129 символов
*/

CREATE TABLE items (
     item_id    integer PRIMARY KEY,
     URL        varchar(1500) not null,
     price      numeric NOT NULL
);

CREATE TABLE users_items (
    user_id integer REFERENCES users,
    item_id INT REFERENCES orders,
    PRIMARY KEY (user_id, item_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items(id) ON UPDATE CASCADE
);