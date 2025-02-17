CREATE TABLE `order_item` (
  `id` varchar(50) NOT NULL,
  `order_id` varchar(50) NOT NULL,
  `sku_id` varchar(50) NOT NULL,
  `quantity` int NOT NULL,
  `price` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)