CREATE TABLE `order` (
  `id` varchar(50) NOT NULL,
  `user_id` varchar(50) NOT NULL,
  `payment_method` varchar(20) NOT NULL,
  `amount` int NOT NULL,
  `status` enum('Pending','Paid','Cancelled') DEFAULT NULL,
  `address` text NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)