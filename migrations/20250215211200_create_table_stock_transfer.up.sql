CREATE TABLE `stock_transfer` (
  `id` varchar(50) NOT NULL,
  `from_warehouse_id` varchar(50) NOT NULL,
  `to_warehouse_id` varchar(50) NOT NULL,
  `sku_id` varchar(50) NOT NULL,
  `quantity` int DEFAULT NULL,
  `created_by` varchar(50) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(50) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)