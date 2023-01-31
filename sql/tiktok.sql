/*
 Navicat Premium Data Transfer

 Source Server         : 101.43.157.116
 Source Server Type    : MySQL
 Source Server Version : 80032
 Source Host           : 101.43.157.116:3306
 Source Schema         : tiktok

 Target Server Type    : MySQL
 Target Server Version : 80032
 File Encoding         : 65001

 Date: 30/01/2023 17:35:00
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comment_video
-- ----------------------------
DROP TABLE IF EXISTS `comment_video`;
CREATE TABLE `comment_video`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `identity` bigint(0) NULL DEFAULT NULL COMMENT '评论唯一标识',
  `video_identity` bigint(0) NULL DEFAULT NULL COMMENT '视频唯一标识',
  `user_identity` bigint(0) NULL DEFAULT NULL COMMENT '用户唯一标识',
  `text` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL COMMENT '评论内容',
  `comment_time` datetime(0) NULL DEFAULT NULL COMMENT '评论时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_comment_video_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for favourite_video
-- ----------------------------
DROP TABLE IF EXISTS `favourite_video`;
CREATE TABLE `favourite_video`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `video_identity` bigint(0) NULL DEFAULT NULL COMMENT '视频唯一标识',
  `user_identity` bigint(0) NULL DEFAULT NULL COMMENT '用户唯一标识',
  `status` tinyint(1) NULL DEFAULT NULL COMMENT '点赞状态（0：未点赞，1：已点赞）',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_favourite_video_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_basic
-- ----------------------------
DROP TABLE IF EXISTS `user_basic`;
CREATE TABLE `user_basic`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `identity` bigint(0) NULL DEFAULT NULL COMMENT '用户唯一标识',
  `username` varchar(36) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '用户名',
  `password` varchar(36) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '密码',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_basic_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for video_basic
-- ----------------------------
DROP TABLE IF EXISTS `video_basic`;
CREATE TABLE `video_basic`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `identity` bigint(0) NULL DEFAULT NULL COMMENT '视频唯一标识',
  `user_identity` bigint(0) NULL DEFAULT NULL COMMENT '用户唯一标识',
  `play_url` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '视频url',
  `cover_url` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '封面url',
  `title` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL COMMENT '标题',
  `publish_time` datetime(0) NULL DEFAULT NULL COMMENT '发布时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_video_basic_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
