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

INSERT IGNORE INTO ecommerce.warehouse
(id, shop_id, location, address, is_active, created_at, updated_at, deleted_at)
VALUES('01950a6d-4531-76f3-ad87-c6daac630274', '01950a6a-4816-7e58-a15a-fbd229ed3e72', 'Tangerang', 'Jl. Tangerang', 1, '2025-02-15 16:24:48', NULL, NULL);

INSERT IGNORE INTO ecommerce.warehouse
(id, shop_id, location, address, is_active, created_at, updated_at, deleted_at)
VALUES('3c4d1ceb-93d4-4233-aea7-700bf51e5051', '35ffcc6c-be96-4a35-a1a6-d4108c313b45', 'bandung', 'Jl Cibaduyut', 1, '2025-02-13 03:37:43', '2025-02-13 08:07:22', NULL);

INSERT IGNORE INTO ecommerce.warehouse
(id, shop_id, location, address, is_active, created_at, updated_at, deleted_at)
VALUES('5daf7427-fb6a-4ccb-8d54-2c496503f04c', '35ffcc6c-be96-4a35-a1a6-d4108c313b45', 'Sumedang', 'Jl Sumedang', 1, '2025-02-14 09:17:14', NULL, NULL);
