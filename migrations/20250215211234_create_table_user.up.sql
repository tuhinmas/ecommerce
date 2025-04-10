CREATE TABLE `user` (
  `id` varchar(50) NOT NULL,
  `warehouse_id` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `phone` varchar(13) NOT NULL,
  `password` text NOT NULL,
  `gender` enum('male','female') NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);