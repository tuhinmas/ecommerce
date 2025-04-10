CREATE TABLE `product` (
  `id` varchar(50) NOT NULL,
  `shop_id` varchar(50) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `created_by` varchar(50) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(50) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);

-- INSERT IGNORE INTO ecommerce.product
-- (`id`, `shop_id`, `name`, `created_by`, `created_at`, `updated_by`, `updated_at`, `deleted_at`)
-- VALUES('01950a76-db4d-7665-9a97-f84b62e96212', '01950a6a-4816-7e58-a15a-fbd229ed3e72', 'Sabun Colek', 'a5d44fd7-2f76-4abf-8880-818c63c4e94f', '2025-02-15 16:35:16', NULL, NULL, NULL);

-- INSERT IGNORE INTO ecommerce.product
-- (`id`, `shop_id`, `name`, `created_by`, `created_at`, `updated_by`, `updated_at`, `deleted_at`)
-- VALUES('7b8759d5-ebc0-440b-b7af-1c367bd12b60', '35ffcc6c-be96-4a35-a1a6-d4108c313b45', 'Kopi Abc', 'a5d44fd7-2f76-4abf-8880-818c63c4e94f', '2025-02-13 03:38:31', NULL, '2025-02-15 21:15:43', NULL);
