---
contentVersion: 1
title: "MySQL B+ 树索引原理"
slug: "mysql-btree-index"
description: "深入理解 MySQL InnoDB 中 B+ 树索引的数据结构、查询过程和优化技巧。"
date: "2026-06-20"
updated: "2026-06-22"
category: "数据库"
tags:
  - MySQL
  - 索引
  - B+树
  - 数据结构
cover: ""
draft: false
pinned: false
comment: true
toc: true
---

# MySQL B+ 树索引原理

索引是数据库中最重要的性能优化手段之一。MySQL InnoDB 存储引擎使用 **B+ 树**作为索引的底层数据结构，理解其原理对于写出高效的 SQL 至关重要。

## 一、为什么是 B+ 树而不是 B 树

B+ 树是 B 树的变体，和 B 树的区别在于：

1. **非叶子节点只存键值，不存数据**
   - 同样大小的磁盘页能存更多键值
   - 树更矮，磁盘 I/O 更少

2. **叶子节点包含全部数据，且通过指针串联**
   - 全表扫描只需遍历叶子链表
   - 范围查询非常高效

3. **查询稳定性**
   - 所有查找最终都必须到叶子节点
   - 查询时间复杂度严格 O(log n)

```
B 树:
       [key|data]
      /          \
[key|data]    [key|data]

B+ 树:
        [key]
       /     \
    [key]    [key]
    /    \   /   \
[data] [data] [data] [data]  <- 叶子链表
```

## 二、InnoDB 中的索引类型

### 聚簇索引（Clustered Index）

- 以主键构建 B+ 树
- 叶子节点存储**完整行数据**
- 一个表只能有一个聚簇索引

### 辅助索引（Secondary Index）

- 以非主键列构建 B+ 树
- 叶子节点存储**主键值**
- 查询需要**回表**：先在辅助索引找到主键，再用主键去聚簇索引找完整数据

```sql
-- 假设 id 是主键，name 上有索引
SELECT * FROM users WHERE name = '张三';

-- 执行过程:
-- 1. 在 name 的辅助索引中找到 '张三' -> id=100
-- 2. 用 id=100 去聚簇索引找完整行数据（回表）
```

## 三、最左前缀原则

联合索引遵循**最左前缀原则**：

```sql
-- 创建联合索引
CREATE INDEX idx_a_b_c ON table(a, b, c);

-- 能用到索引
WHERE a = 1           -- 匹配最左列，✅
WHERE a = 1 AND b = 2 -- 匹配前两列，✅
WHERE a = 1 AND c = 3 -- 只用 a，c 断掉了，部分命中 ⚠️

-- 无法用到索引
WHERE b = 2           -- 没有 a，❌
WHERE b = 2 AND c = 3 -- 没有 a，❌
```

理解最左前缀的关键是：B+ 树按 (a, b, c) 的优先级排序。如果没有 a，就无法确定数据在树中的位置。

## 四、覆盖索引

如果查询的列都包含在索引中，就不需要回表，这就是**覆盖索引**。

```sql
-- 假设有索引 idx_name_age(name, age)
SELECT name, age FROM users WHERE name = '张三';

-- 查询的 name 和 age 都在索引中，不需要回表
-- Using index（性能最好）
```

常用优化技巧：
- SELECT 只查需要的列，尽量使用覆盖索引
- 使用 `EXPLAIN` 检查执行计划的 Extra 列
- 看到 `Using index` 说明是覆盖索引

## 五、索引下推（ICP）

MySQL 5.6+ 支持索引下推（Index Condition Pushdown）：

```sql
-- 假设联合索引 idx_name_age(name, age)
SELECT * FROM users WHERE name LIKE '张%' AND age = 25;
```

没有 ICP 时：从索引中找到所有 name LIKE '张%' 的主键，全部回表，再过滤 age=25。
有 ICP 时：在索引层就过滤掉 age≠25 的主键，减少回表次数。

## 六、索引优化的几个建议

1. **选择性高的列建索引**
   - 选择性 = DISTINCT 值数量 / 总行数
   - 选择性接近 1 的列最适合建索引

2. **避免在索引列上做计算**
   ```sql
   -- ❌ 不好
   WHERE YEAR(created_at) = 2026
   -- ✅ 好
   WHERE created_at >= '2026-01-01' AND created_at < '2027-01-01'
   ```

3. **注意索引数量**：索引不是越多越好，写操作需要维护索引

4. **使用 EXPLAIN 分析**：关注 type 列（ref 好于 range 好于 ALL）

## 小结

B+ 树索引是 MySQL 性能优化的核心知识。理解其结构和工作方式，才能在实际项目中写出高效的 SQL。重点关注：聚簇索引与辅助索引的区别、最左前缀原则、覆盖索引优化。
