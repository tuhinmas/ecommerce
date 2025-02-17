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

INSERT IGNORE INTO ecommerce.`user`
(id, warehouse_id, name, phone, password, gender, created_at, updated_at)
VALUES('0194fd66-dca4-7432-9c4b-9eb37661d2e2', '3c4d1ceb-93d4-4233-aea7-700bf51e5051', 'dedi', '08982510066', '3c9909afec25354d551dae21590bb26e38d53f2173b8d3dc3eee4c047e7ab1c1eb8b85103e3be7ba613b31bb5c9c36214dc9f14a42fd7a2fdb84856bca5c44c2', 'male', '2025-02-13 03:42:44', NULL);

INSERT IGNORE INTO ecommerce.`user`
(id, warehouse_id, name, phone, password, gender, created_at, updated_at)
VALUES('01950b6a-882d-7c6d-a48c-a4b09fe55f87', '5daf7427-fb6a-4ccb-8d54-2c496503f04c', 'andi', '08982510077', '3c9909afec25354d551dae21590bb26e38d53f2173b8d3dc3eee4c047e7ab1c1eb8b85103e3be7ba613b31bb5c9c36214dc9f14a42fd7a2fdb84856bca5c44c2', 'male', '2025-02-15 21:01:25', NULL);
