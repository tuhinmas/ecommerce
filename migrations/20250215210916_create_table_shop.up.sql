CREATE TABLE `shop` (
    `id` varchar(50) NOT NULL,
    `name` varchar(255) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
)