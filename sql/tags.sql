-- 标签分类表
CREATE TABLE `tag_categories` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `description` varchar(200) DEFAULT NULL COMMENT '分类描述',
  `parent_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '父分类ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_parent_name` (`parent_id`, `name`),
  CONSTRAINT `fk_tag_category_parent` FOREIGN KEY (`parent_id`) REFERENCES `tag_categories` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签分类表';

-- 修改标签表，添加分类关联
ALTER TABLE `problem_tags` 
  ADD COLUMN `category_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '所属分类ID' AFTER `color`,
  ADD COLUMN `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间' AFTER `created_at`,
  ADD CONSTRAINT `fk_tag_category` FOREIGN KEY (`category_id`) REFERENCES `tag_categories` (`id`) ON DELETE RESTRICT;

-- 如果之前没有 problem_tag_relations 表，需要创建（用于题目和标签的多对多关系）
CREATE TABLE IF NOT EXISTS `problem_tag_relations` (
  `problem_id` bigint(20) UNSIGNED NOT NULL COMMENT '题目ID',
  `tag_id` bigint(20) UNSIGNED NOT NULL COMMENT '标签ID',
  PRIMARY KEY (`problem_id`, `tag_id`),
  KEY `idx_tag_id` (`tag_id`),
  CONSTRAINT `fk_ptr_problem` FOREIGN KEY (`problem_id`) REFERENCES `problems` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_ptr_tag` FOREIGN KEY (`tag_id`) REFERENCES `problem_tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='题目标签关联表';
