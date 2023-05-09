/*
 Navicat Premium Data Transfer

 Source Server         : Klaatoo
 Source Server Type    : MariaDB
 Source Server Version : 101102 (10.11.2-MariaDB-1:10.11.2+maria~ubu2204)
 Source Host           : 192.168.219.107:13306
 Source Schema         : Inae

 Target Server Type    : MariaDB
 Target Server Version : 101102 (10.11.2-MariaDB-1:10.11.2+maria~ubu2204)
 File Encoding         : 65001

 Date: 28/04/2023 10:40:17
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for eth_address
-- ----------------------------
DROP TABLE IF EXISTS `eth_address`;
CREATE TABLE `eth_address` (
  `address` varchar(50) NOT NULL DEFAULT '',
  `prikey` varchar(70) NOT NULL DEFAULT '',
  `pubkey` varchar(70) NOT NULL DEFAULT '',
  `create_dt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  `active` enum('Y','N') NOT NULL DEFAULT 'Y',
  PRIMARY KEY (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of eth_address
-- ----------------------------
BEGIN;
INSERT INTO `eth_address` (`address`, `prikey`, `pubkey`, `create_dt`, `active`) VALUES ('0x02484faf04d8c50b18f280b4c2323accae253c70', '1', '1', '2023-04-27 11:20:31', 'Y');
INSERT INTO `eth_address` (`address`, `prikey`, `pubkey`, `create_dt`, `active`) VALUES ('0x89df83a8dddca8e30d88a863af9a6a71a400df64', '1', '1', '2023-04-27 11:20:48', 'Y');
COMMIT;

-- ----------------------------
-- Table structure for eth_history
-- ----------------------------
DROP TABLE IF EXISTS `eth_history`;
CREATE TABLE `eth_history` (
  `id` int(1) unsigned NOT NULL AUTO_INCREMENT,
  `block_number` int(1) unsigned NOT NULL DEFAULT 0,
  `block_hash` varchar(70) NOT NULL DEFAULT '',
  `tx_hash` varchar(70) NOT NULL DEFAULT '',
  `from_address` varchar(50) NOT NULL DEFAULT '',
  `to_address` varchar(50) NOT NULL DEFAULT '',
  `value` varchar(50) NOT NULL,
  `create_dt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of eth_history
-- ----------------------------
BEGIN;
INSERT INTO `eth_history` (`id`, `block_number`, `block_hash`, `tx_hash`, `from_address`, `to_address`, `value`, `create_dt`) VALUES (1, 956, '0x73a4bc5f3af5a80c359155f47ee8b20512631425297add7e203f1c6e2780b8df', '0x4f6ee55f7493838aedfabc8942382daaa15a7d58a2a5b1cc5b43a9bb04e2db3b', '0x02484faf04d8c50b18f280b4c2323accae253c70', '0x89df83a8dddca8e30d88a863af9a6a71a400df64', '10000000000000000000', '2023-04-27 08:51:33');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
