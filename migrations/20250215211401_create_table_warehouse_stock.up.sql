CREATE TABLE `warehouse_stock` (
  `id` varchar(50) NOT NULL,
  `warehouse_id` varchar(50) NOT NULL,
  `sku_id` varchar(50) NOT NULL,
  `stock` int DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);

INSERT IGNORE INTO ecommerce.warehouse_stock
(id, warehouse_id, sku_id, stock, created_at, updated_at)
VALUES('0195099c-6081-72d9-bd7e-d103ee0e99d8', '5daf7427-fb6a-4ccb-8d54-2c496503f04c', 'b32b86b0-aa0a-4bea-806c-61f0250a8c28', 20, '2025-02-15 12:36:38', '2025-02-15 12:56:54');

INSERT IGNORE INTO ecommerce.warehouse_stock
(id, warehouse_id, sku_id, stock, created_at, updated_at)
VALUES('879ff585-45bd-4cb8-987c-106dec98b2ab', '3c4d1ceb-93d4-4233-aea7-700bf51e5051', '502a7a6a-ef9a-46ce-b9ff-cb33b21fdbec', 40, '2025-02-13 03:40:40', '2025-02-14 09:46:57');

INSERT IGNORE INTO ecommerce.warehouse_stock
(id, warehouse_id, sku_id, stock, created_at, updated_at)
VALUES('b32b86b0-aa0a-4bea-806c-61f0250a8c28', '3c4d1ceb-93d4-4233-aea7-700bf51e5051', 'b32b86b0-aa0a-4bea-806c-61f0250a8c28', 56, '2025-02-13 04:03:14', '2025-02-15 12:52:05');
