-- CREATE TABLE `test_table` (
--     `id` INT NOT NULL AUTO_INCREMENT,
--     `name` VARCHAR(255) NOT NULL,
--     `age` INT NOT NULL,
--     `email` VARCHAR(255) UNIQUE NOT NULL,
--     `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
--     PRIMARY KEY (`id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- Create room_type table first
CREATE TABLE `room_type` (
   `room_type_id` int NOT NULL AUTO_INCREMENT,
   `type_name` varchar(255) NOT NULL,
   PRIMARY KEY (`room_type_id`),
   UNIQUE KEY `type_name` (`type_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create rooms table
CREATE TABLE `rooms` (
   `room_id` int NOT NULL AUTO_INCREMENT,
   `room_type_id` int NOT NULL,
   `room_description` text,
   `price` decimal(10,0) DEFAULT NULL,
   `availability` tinyint DEFAULT '0',
   PRIMARY KEY (`room_id`),
   KEY `fk_room_type` (`room_type_id`),
   CONSTRAINT `fk_room_type` FOREIGN KEY (`room_type_id`) REFERENCES `room_type` (`room_type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create amenity table
CREATE TABLE `amenity` (
   `amenity_id` int NOT NULL AUTO_INCREMENT,
   `amenity_name` varchar(255) NOT NULL,
   PRIMARY KEY (`amenity_id`),
   UNIQUE KEY `amenity_name` (`amenity_name`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create room_amenity table
CREATE TABLE `room_amenity` (
   `room_id` int NOT NULL,
   `amenity_id` int NOT NULL,
   KEY `room_id` (`room_id`),
   KEY `amenity_id` (`amenity_id`),
   CONSTRAINT `room_amenity_ibfk_1` FOREIGN KEY (`room_id`) REFERENCES `rooms` (`room_id`),
   CONSTRAINT `room_amenity_ibfk_2` FOREIGN KEY (`amenity_id`) REFERENCES `amenity` (`amenity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create guest table
CREATE TABLE `guest` (
   `guest_id` varchar(255) NOT NULL,
   `uuid` varchar(36) NOT NULL,
   `guest_name` varchar(255) DEFAULT NULL,
   `email` varchar(255) DEFAULT NULL,
   `phone_number` varchar(20) DEFAULT NULL,
   `password` varchar(255) DEFAULT NULL,
   `updated_at` timestamp NULL DEFAULT NULL,
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `token` varchar(255) DEFAULT NULL,
   PRIMARY KEY (`uuid`),
   UNIQUE KEY `guest_id` (`guest_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create bookings table
CREATE TABLE `bookings` (
   `booking_id` int NOT NULL AUTO_INCREMENT,
   `room_id` int NOT NULL,
   `guest_id` varchar(255) NOT NULL,
   `num_of_adults` int DEFAULT NULL,
   `num_of_children` int DEFAULT NULL,
   `checkin_date` date DEFAULT NULL,
   `checkout_date` date DEFAULT NULL,
   `checkin_time` time DEFAULT NULL,
   `checkout_time` time DEFAULT NULL,
   `price` decimal(10,0) DEFAULT NULL,
   `status` varchar(50) DEFAULT NULL,
   `booking_date` date DEFAULT NULL,
   PRIMARY KEY (`booking_id`),
   KEY `room_id` (`room_id`),
   KEY `guest_id` (`guest_id`),
   CONSTRAINT `bookings_ibfk_1` FOREIGN KEY (`room_id`) REFERENCES `rooms` (`room_id`),
   CONSTRAINT `bookings_ibfk_2` FOREIGN KEY (`guest_id`) REFERENCES `guest` (`guest_id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Create users table
CREATE TABLE `users` (
   `user_id` varchar(50) NOT NULL,
   `user_name` varchar(50) NOT NULL,
   `email` varchar(50) NOT NULL,
   `phone_number` varchar(50) DEFAULT NULL,
   `updated_at` timestamp NULL DEFAULT NULL,
   `created_at` timestamp NULL DEFAULT NULL,
   `password` varchar(255) NOT NULL,
   `api_key` varchar(255) DEFAULT NULL,
   PRIMARY KEY (`user_id`),
   UNIQUE KEY `user_id` (`user_id`),
   UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;