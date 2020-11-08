create table user (id  bigserial not null, created timestamp, email varchar(40) not null, age int not null, name varchar(50) not null, sername varchar(50) not null, password varchar(50) not null, primary key (id));

create table camera (id bigserial not null, created timestamp, active bit not null, loc varchar(50)not null, token_session_camera varchar(1000), token_session_consumer varchar(1000), id_camera int not null, user int not null, primary key(id));

alter table if exists camera add constraint FKae1ky8v52w8jkmpojhg9daq03 foreign key (user) references user;

insert into user (id, name, sername, age, email, password)
values(1, 'Mauro', 'Garcia',29, 'mauropuravida1@gmail.com', 'test1');

insert into user (id, name, sername, age, email, password)
values(1, 'Mauro', 'Garcia',29, 'mauropuravida2@gmail.com', 'test2');

insert into user (id, name, sername, age, email, password)
values(1, 'Mauro', 'Garcia',29, 'mauropuravida3@gmail.com', 'test3');

insert into user (id, name, sername, age, email, password)
values(1, 'Mauro', 'Garcia',29, 'mauropuravida4@gmail.com', 'test4');

insert into camera(id, user, active, loc, token_session_camera, token_session_consumer, id_camera)
values(1, 1, 'cochabamba 1745', '', '', 1);
