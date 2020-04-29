/*
Navicat MySQL Data Transfer

Source Server         : 127.0.0.1
Source Server Version : 50553
Source Host           : 127.0.0.1:3306
Source Database       : seckill

Target Server Type    : MYSQL
Target Server Version : 50553
File Encoding         : 65001

Date: 2018-07-07 15:23:46
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for activity
-- ----------------------------
DROP TABLE IF EXISTS `activity`;
CREATE TABLE `activity` (
  `activity_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '活动Id',
  `activity_name` varchar(50) NOT NULL DEFAULT '' COMMENT '活动名称',
  `product_id` int(11) unsigned NOT NULL COMMENT '商品Id',
  `start_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '活动开始时间',
  `end_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '活动结束时间',
  `total` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '商品数量',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '活动状态',
  `sec_speed` int(5) unsigned NOT NULL DEFAULT '0' COMMENT '每秒限制多少个商品售出',
  `buy_limit` int(5) unsigned NOT NULL COMMENT '购买限制',
  `buy_rate` decimal(2,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '购买限制',
  PRIMARY KEY (`activity_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='@活动数据表';

-- ----------------------------
-- Records of activity
-- ----------------------------
INSERT INTO `activity` VALUES ('1', '香蕉大甩卖', '1', '530871061', '530872061', '20', '0', '1', '1', '0.20');
INSERT INTO `activity` VALUES ('2', '苹果大甩卖', '2', '530871061', '530872061', '20', '0', '1', '1', '0.20');
INSERT INTO `activity` VALUES ('3', '桃子大甩卖', '3', '1530928052', '1530989052', '20', '0', '1', '1', '0.20');
INSERT INTO `activity` VALUES ('4', '梨子大甩卖', '4', '1530928052', '1530989052', '20', '0', '1', '1', '0.20');

-- ----------------------------
-- Table structure for product
-- ----------------------------
DROP TABLE IF EXISTS `product`;
CREATE TABLE `product` (
  `product_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '商品Id',
  `product_name` varchar(50) NOT NULL DEFAULT '' COMMENT '商品名称',
  `total` int(5) unsigned NOT NULL DEFAULT '0' COMMENT '商品数量',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '商品状态',
  PRIMARY KEY (`product_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='@商品数据表';

-- ----------------------------
-- Records of product
-- ----------------------------
INSERT INTO `product` VALUES ('1', '香蕉', '100', '1');
INSERT INTO `product` VALUES ('2', '苹果', '100', '1');
INSERT INTO `product` VALUES ('3', '桃子', '100', '1');
INSERT INTO `product` VALUES ('4', '梨子', '100', '1');
