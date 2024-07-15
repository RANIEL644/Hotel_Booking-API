create database HB;
use HB;
show tables;

describe rooms;
describe amenity;
describe room_amenity;
describe room_type;
describe users;
describe bookings;

select * from rooms;
select * from amenity;
select * from room_amenity;
select * from room_type;
select * from users;
select * from  bookings;

update rooms
set to_date = "2024-07-19" where room_id = 1;
ALTER TABLE rooms
MODIFY COLUMN availability BOOL;



-- SELECT 
--     r.room_id, 
--     rt.type_name, 
--     r.room_description, 
--     GROUP_CONCAT(a.amenity_name SEPARATOR ', ') AS amenities 
-- FROM 
--     rooms r
-- JOIN 
--     room_type rt ON r.room_type_id = rt.room_type_id
-- LEFT JOIN 
--     room_amenity ra ON r.room_id = ra.room_id
-- LEFT JOIN 
--     amenity a ON ra.amenity_id = a.amenity_id
-- GROUP BY 
--     r.room_id;

-- SELECT user_id, user_name, email, phone_number, updated_at, created_at, password
--               FROM users WHERE email = "raniel@gmail.com";
              
              
SET SQL_SAFE_UPDATES=0;


CREATE TABLE Rooms (
    room_id INT AUTO_INCREMENT PRIMARY KEY,
    room_number int NOT NULL ,
    room_type_id int not null,
    price DECIMAL(10, 2) NOT NULL,
    room_description text,
    foreign key (room_type_id) references room_type(room_type_id)
     -- wifi BOOLEAN DEFAULT FALSE,
    -- ac BOOLEAN DEFAULT FALSE,
--     non_ac BOOLEAN DEFAULT FALSE,
--     facilities TEXT
);

alter table rooms add column to_date date;

SELECT amenities.* FROM amenity
JOIN room_amenity ON amenities.id = room_amenity.amenity_id
WHERE room_amenity.room_id IN (1);


INSERT INTO Rooms (room_number, room_type_id, room_description, price )
VALUES (101, 1,'Standard Room with city view',150);
select * from rooms;
SELECT room_description FROM Rooms WHERE room_id = 2;
delete from rooms where room_id = 2;
ALTER TABLE amenity AUTO_INCREMENT = 1;


Create table Amenity (
amenity_id int not null auto_increment,
amenity_name varchar(255) not null unique,
primary key (amenity_id)
);


insert into amenity (amenity_name) values ("AC");
delete from amenity where amenity_id = 1;


create table Room_Amenity (
room_id int not null,
amenity_id int not null, 
foreign key (room_id) references rooms(room_id),
foreign key (amenity_id) references amenity (amenity_id)
);


create table Room_type(
room_type_id int not null primary key,
type_name varchar(255) not null unique
);
insert into room_type values ("1", "Single");


create table Room_Availability(
availability_id int not null primary key unique,
room_id int not null unique,
start_date date,
end_date date,
foreign key (room_id) references rooms(room_id)
);


create table Bookings(
booking_id int not null primary key unique,
room_id int not null unique,
guest_id int not null,
num_of_adults int ,
num_of_children int,
checkin_date date,
checkout_date date,
foreign key (room_id) references rooms(room_id),
foreign key (guest_id) references guest(guest_id)
);


create table Guest(
guest_id int not null primary key unique,
guest_name varchar(50) not null,
email varchar(50) not null,
phone_number int not null,
password varchar(255) not null
);

create table user(
user_id int not null primary key unique,
user_name varchar(50) not null,
email varchar(50) not null,
phone_number int not null,
password varchar(255) not null unique
);



create table Booking_details(
id int not null,
booking_id int not null primary key unique,
room_id int not null unique,
guest_id int not null,
guest_name int not null,
num_of_adults int ,
num_of_children int,
amount decimal,
foreign key (booking_id) references bookings(booking_id),
foreign key (guest_id) references guest(guest_id)
);






