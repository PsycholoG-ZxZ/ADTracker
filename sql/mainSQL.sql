show databases;

create database adtracker;

use adtracker;

DROP table Users;

CREATE table Users (
    user_id int auto_increment primary key, 
    email   Nvarchar(129) NOT NULL, 
    verification bool NOT NULL Default 0,
    verification_code int(10),
    UNIQUE(email)
);

show tables ;



ALTER TABLE Users
MODIFY COLUMN verification_code Nvarchar(10);

ALTER table items
add url Nvarchar(1500) NOT NULL;
SHOW CREATE TABLE items ;
ALTER TABLE items
MODIFY COLUMN price Nvarchar(12);

update adtracker.items SET price = '4 900' where price = '4 800';


SELECT * FROM  users;
SELECT * FROM  items;
SELECT * FROM  users_items;
DELETE FROM users
WHERE email='dolokhov99_v2@11mail.ru';
DELETE FROM items
WHERE price='4 800';

Select email from users join users_items ON (users_items.user_id = users.user_id) AND (users_items.item_id = 6);

Select URL from items;

CREATE TABLE items (
     item_id    int auto_increment PRIMARY KEY,
     price      int(12)
);

CREATE TABLE users_items (
user_id int NOT NULL,
item_id int NOT NULL,
PRIMARY KEY (user_id, item_id),
CONSTRAINT FK_Usr_Itm_Usr FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE,
CONSTRAINT FK_Usr_Itm_Itm FOREIGN KEY (item_id) REFERENCES items (item_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

insert into adtracker.users (email, verification_code) values ('dimka_volnyi@mail.ru', '1111111111');

