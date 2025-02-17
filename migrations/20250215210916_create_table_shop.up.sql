CREATE TABLE `shop` (
  `id` varchar(50) NOT NULL,
  `name` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);

INSERT IGNORE INTO ecommerce.shop
(id, name, created_at, updated_at, deleted_at)
VALUES('01950a6a-4816-7e58-a15a-fbd229ed3e72', 'shop 2', '2025-02-15 16:21:32', NULL, NULL);

INSERT IGNORE INTO ecommerce.shop
(id, name, created_at, updated_at, deleted_at)
VALUES('35ffcc6c-be96-4a35-a1a6-d4108c313b45', 'shop 1', '2025-02-13 03:36:57', '2025-02-13 03:41:14', NULL);
