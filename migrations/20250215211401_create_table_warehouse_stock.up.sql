CREATE TABLE `warehouse_stock` (
  `id` varchar(50) NOT NULL,
  `warehouse_id` varchar(50) NOT NULL,
  `sku_id` varchar(50) NOT NULL,
  `stock` int DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);