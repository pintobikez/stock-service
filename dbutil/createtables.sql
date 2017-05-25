CREATE DATABASE stockservice;

USE stockservice;

CREATE TABLE IF NOT EXISTS `stock` (
  `sku` varchar(16) NOT NULL,
  `warehouse` varchar(45) NOT NULL,
  `quantity` int(6) NOT NULL DEFAULT '0',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  KEY `skuWarehouse` (`sku`,`warehouse`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `reservation` (
  `sku` varchar(16) NOT NULL,
  `warehouse` varchar(45) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  KEY `skuWarehouse2` (`warehouse`,`sku`) USING BTREE,
  KEY `res_create_at` (`created_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;