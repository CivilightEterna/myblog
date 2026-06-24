---
contentVersion: 1
title: "Java 线程池拒绝策略详解"
slug: "java-threadpool-rejection"
description: "深入分析 Java 线程池的四种拒绝策略，以及在真实业务场景中的最佳实践。"
date: "2026-06-15"
updated: "2026-06-15"
category: "后端"
tags:
  - Java
  - 线程池
  - 并发编程
  - Java面试
cover: ""
draft: false
pinned: false
comment: true
toc: true
---

# Java 线程池拒绝策略详解

线程池是 Java 并发编程中最常用的工具之一。当线程池和队列都满时，新提交的任务就会触发**拒绝策略**。理解四种内置的拒绝策略以及如何自定义拒绝策略，是线程池使用的关键。

## 一、线程池执行流程回顾

```
新任务提交
    ↓
核心线程数是否已满？
    ↓ 未满
  创建核心线程执行任务
    ↓ 已满
工作队列是否已满？
    ↓ 未满
  任务加入工作队列等待
    ↓ 已满
最大线程数是否已满？
    ↓ 未满
  创建非核心线程执行任务
    ↓ 已满
触发拒绝策略 ← 我们在这里
```

## 二、四种内置拒绝策略

### 1. AbortPolicy（默认）

直接抛出 `RejectedExecutionException`。

```java
ThreadPoolExecutor pool = new ThreadPoolExecutor(
    2, 4, 60, TimeUnit.SECONDS,
    new LinkedBlockingQueue<>(2),
    new ThreadPoolExecutor.AbortPolicy()
);

try {
    pool.execute(task);
} catch (RejectedExecutionException e) {
    log.error("任务被拒绝: {}", e.getMessage());
    // 通常记录日志、告警
}
```

**适用场景**：需要立即感知到线程池过载的情况，配合监控告警系统。

### 2. CallerRunsPolicy

由提交任务的线程（调用者线程）自己执行这个任务。

```java
new ThreadPoolExecutor.CallerRunsPolicy()
```

**优点**：相当于一种天然的限流——调用者执行任务时不会继续提交新任务。
**缺点**：可能拖慢主线程，影响整体吞吐量。
**适用场景**：主线程对延迟不太敏感，但不希望任务丢失。

### 3. DiscardPolicy

直接丢弃被拒绝的任务，不抛异常。

```java
new ThreadPoolExecutor.DiscardPolicy()
```

**危险**：任务静默丢失，非常不推荐在核心业务中使用。
**适用场景**：日志上报、非关键指标采集等即使丢失也影响不大的场景。

### 4. DiscardOldestPolicy

丢弃队列中最旧的任务（即将要执行的那个），然后重试提交当前任务。

```java
new ThreadPoolExecutor.DiscardOldestPolicy()
```

**注意**：丢弃的是队首的任务，也就是等待最久的那个。
**适用场景**：优先处理最新数据的场景，如实时数据流处理。

## 三、自定义拒绝策略

四种内置策略满足不了所有场景。自定义拒绝策略需要实现 `RejectedExecutionHandler` 接口：

```java
public class AlertAndRetryPolicy implements RejectedExecutionHandler {

    @Override
    public void rejectedExecution(Runnable r, ThreadPoolExecutor executor) {
        // 1. 发送告警
        alertService.send("线程池过载，触发拒绝策略");

        // 2. 尝试重新入队（阻塞方式）
        try {
            if (!executor.isShutdown()) {
                executor.getQueue().put(r); // 阻塞等待
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        // 3. 持久化到数据库（最终兜底）
        saveToDb(r);
    }
}
```

## 四、常见业务场景处理

### 场景一：高并发秒杀

```java
ThreadPoolExecutor pool = new ThreadPoolExecutor(
    10, 50, 60, TimeUnit.SECONDS,
    new ArrayBlockingQueue<>(200),
    new CallerRunsPolicy() // 让调用者执行，天然限流
);
```

选择 CallerRunsPolicy 配合有界队列，当流量超出处理能力时，让 Tomcat/Netty 的线程参与计算，形成天然的背压。

### 场景二：异步下单

```java
ThreadPoolExecutor pool = new ThreadPoolExecutor(
    5, 10, 30, TimeUnit.SECONDS,
    new LinkedBlockingQueue<>(100),
    (r, executor) -> {
        // 持久化到 DB，后续补偿
        savePendingTask(r);
        alert("下单任务被拒绝，已落库待补单");
    }
);
```

核心业务不能丢任务。拒绝时持久化到数据库，定时任务扫表补偿。

### 场景三：日志处理

```java
ThreadPoolExecutor pool = new ThreadPoolExecutor(
    2, 4, 60, TimeUnit.SECONDS,
    new LinkedBlockingQueue<>(1000),
    new ThreadPoolExecutor.DiscardPolicy()
);
```

日志属于可丢弃的辅助性任务，直接丢弃。

## 五、线程池参数设置建议

1. **核心线程数**：CPU 密集型设为 `N+1`，IO 密集型设为 `2N`（N 为 CPU 核数）
2. **队列**：推荐使用有界队列（ArrayBlockingQueue / 有界 LinkedBlockingQueue），避免 OOM
3. **拒绝策略**：核心业务用 CallerRunsPolicy 或自定义持久化；非核心用 DiscardPolicy
4. **监控**：生产环境建议接入线程池监控，关注 activeCount、queueSize、completedTaskCount

```java
// 简单的监控日志
ScheduledExecutorService monitor = Executors.newSingleThreadScheduledExecutor();
monitor.scheduleAtFixedRate(() -> {
    log.info("线程池状态 - 活跃: {}, 队列: {}, 完成: {}",
        pool.getActiveCount(),
        pool.getQueue().size(),
        pool.getCompletedTaskCount()
    );
}, 10, 30, TimeUnit.SECONDS);
```

## 小结

拒绝策略是线程池的"兜底"机制。选择的优先级应该是：**核心业务保证不丢 > 系统稳定性 > 吞吐量**。AbortPolicy 适合需要感知异常的场景，CallerRunsPolicy 适合需要天然限流的场景，自定义策略适合有完善补偿机制的场景。
