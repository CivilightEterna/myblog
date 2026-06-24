---
contentVersion: 1
title: "Redis 缓存击穿怎么解决"
slug: "redis-cache-breakdown"
description: "整理 Redis 缓存击穿、缓存穿透、缓存雪崩的区别和解决方案。"
date: "2026-06-24"
updated: "2026-06-24"
category: "后端"
tags:
  - Redis
  - 缓存
  - Java面试
cover: ""
draft: false
pinned: true
comment: true
toc: true
---

# Redis 缓存击穿怎么解决

在高并发场景下，Redis 缓存的三个经典问题常常被面试官问到：**缓存穿透**、**缓存击穿**和**缓存雪崩**。它们的名字很相近但含义不同，本文梳理清楚它们的区别和解决方案。

## 一、缓存穿透

**缓存穿透**是指查询一个数据库中也不存在的数据。

由于缓存不命中，请求会穿透缓存直接打到数据库。如果有人恶意构造大量不存在的 key 发起请求，就会给数据库带来很大的压力。

### 解决方案

1. **布隆过滤器（Bloom Filter）**
   - 将所有可能存在的数据哈希到一个足够大的 bitmap 中
   - 一个一定不存在的数据会被这个 bitmap 拦截
   - 存在一定的误判率

```java
// 使用 Guava 的布隆过滤器
BloomFilter<String> filter = BloomFilter.create(
    Funnels.stringFunnel(Charset.defaultCharset()),
    1_000_000,  // 预计插入数量
    0.01        // 误判率
);
```

2. **缓存空值**
   - 当查询返回 null 时，缓存一个空值，设置较短的过期时间
   - 简单有效，但会占用一些内存空间

```java
public String getData(String key) {
    String value = redis.get(key);
    if (value != null) {
        return "NULL".equals(value) ? null : value;
    }
    String dbValue = db.query(key);
    if (dbValue == null) {
        redis.setex(key, 60, "NULL"); // 空值缓存 60 秒
    } else {
        redis.set(key, dbValue);
    }
    return dbValue;
}
```

## 二、缓存击穿

**缓存击穿**是指一个热点 key 在过期的一瞬间，大量请求同时访问这个 key，在缓存过期的瞬间，所有请求全部打到数据库。

### 解决方案

1. **互斥锁（Mutex）**
   - 在缓存失效时，不是所有线程都去请求数据库
   - 使用分布式锁保证只有一个线程去加载数据

```java
public String getDataWithLock(String key) {
    String value = redis.get(key);
    if (value != null) {
        return value;
    }
    // 使用 setnx 实现分布式锁
    String lockKey = "lock:" + key;
    try {
        if (redis.setnx(lockKey, "1") == 1) {
            redis.expire(lockKey, 10);
            String dbValue = db.query(key);
            redis.set(key, dbValue);
            return dbValue;
        } else {
            Thread.sleep(100);
            return getDataWithLock(key); // 递归重试
        }
    } finally {
        redis.del(lockKey);
    }
}
```

2. **永不过期**
   - 对热点 key 不设置过期时间
   - 通过异步更新的方式刷新缓存

3. **逻辑过期**
   - 在 value 中存储过期时间戳
   - 读取时判断是否逻辑过期，若过期则异步刷新

## 三、缓存雪崩

**缓存雪崩**是指大量缓存 key 在同一时间过期，或者 Redis 服务宕机，导致所有请求打到数据库。

### 解决方案

1. **过期时间加随机值**：避免大量 key 同时过期
2. **Redis 集群**：使用主从 + 哨兵或集群保证高可用
3. **限流降级**：使用 Hystrix 或 Sentinel 进行服务降级
4. **多级缓存**：本地缓存 + Redis 多级缓存

```java
// 过期时间加随机值
int expireTime = 3600 + new Random().nextInt(600);
redis.setex(key, expireTime, value);
```

## 四、三者对比总结

| 问题类型 | 原因 | 核心解决思路 |
|---------|------|------------|
| 缓存穿透 | 查不存在的数据 | 布隆过滤器 / 缓存空值 |
| 缓存击穿 | 热点 key 过期 | 互斥锁 / 永不过期 |
| 缓存雪崩 | 大量 key 同时过期 | 随机过期时间 / 高可用 |

## 小结

这三个问题是 Redis 缓存架构中的经典难题。在实际项目中，往往需要组合使用多种方案。对于大部分中小型项目，**缓存空值 + 互斥锁 + 随机过期时间**的组合已经足够应对了。
