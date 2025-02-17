CREATE TABLE `admin` (
  `id` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `username` varchar(100) NOT NULL,
  `password` text NOT NULL,
  `gender` enum('male','female') NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);

INSERT INTO admin (id, name, username, password, gender) VALUES ("a5d44fd7-2f76-4abf-8880-818c63c4e94f","dedi","dedih",SHA2("123", 512),"male");