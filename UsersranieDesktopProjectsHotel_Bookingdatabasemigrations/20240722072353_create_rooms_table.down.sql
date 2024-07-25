-- CREATE TABLE `rooms` (
--    `room_id` int NOT NULL AUTO_INCREMENT,
--    `room_type_id` int NOT NULL,
--    `room_description` text,
--    `price` decimal(10,0) DEFAULT NULL,
--    `availability` tinyint DEFAULT '0',
--    PRIMARY KEY (`room_id`),
--    KEY `fk_room_type` (`room_type_id`),
--    CONSTRAINT `fk_room_type` FOREIGN KEY (`room_type_id`) REFERENCES `room_type` (`room_type_id`),
--    CONSTRAINT `rooms_ibfk_1` FOREIGN KEY (`room_type_id`) REFERENCES `room_type` (`room_type_id`)
--  ) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


-- drop table test_table;

-- Drop the table that has foreign key constraints first
DROP TABLE IF EXISTS `room_amenity`;
DROP TABLE IF EXISTS `bookings`;

-- Drop the tables that are referenced by the above tables
DROP TABLE IF EXISTS `rooms`;
DROP TABLE IF EXISTS `guest`;

-- Drop the tables that are referenced by `rooms` and `room_amenity`
DROP TABLE IF EXISTS `amenity`;

-- Finally, drop the remaining tables
DROP TABLE IF EXISTS `room_type`;
DROP TABLE IF EXISTS `users`;
