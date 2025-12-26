CREATE TABLE IF NOT EXISTS appuser (
    id int primary key auto_increment,
    username varchar(10),
    password varchar(100),
    firstname varchar(10),
    lastname varchar(10),
    email varchar(100)
);