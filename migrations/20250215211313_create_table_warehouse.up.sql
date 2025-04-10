CREATE TABLE `warehouse` (
  `id` varchar(50) NOT NULL,
  `shop_id` varchar(50) NOT NULL,
  `location` varchar(100) DEFAULT NULL,
  `address` text NOT NULL,
  `is_active` tinyint DEFAULT '1',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);