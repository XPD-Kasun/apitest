CREATE TABLE IF NOT EXISTS assignment (
    id int primary key auto_increment,
    user_id int not null,
    task_id int not null,
    foreign key (user_id) references appuser(id) on delete cascade,
    foreign key (task_id) references task(id) on delete cascade
);