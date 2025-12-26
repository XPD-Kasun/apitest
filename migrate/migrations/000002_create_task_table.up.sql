CREATE TABLE IF NOT EXISTS task (
    id int primary key auto_increment,
    taskname varchar(20),
    description varchar(20),
    created_at datetime not null default current_timestamp,
    due datetime not null
);