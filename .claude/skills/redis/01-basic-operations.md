# Redis 基础操作

## 初始化 Redis 客户端

### 方法1：使用自定义库（推荐）
```go
import (
    "std-library-slim/redis"
    "std-library-slim/json"
    "github.com/beego/beego/v2/server/web"
)

func InitRedis() {
    opt := redis.Opt{}
    err := json.ParseE(web.AppConfig.DefaultString("REDIS_CONFIG", ""), &opt)
    if err != nil {
        panic(err)
    }
    redis.Init(&opt)
}

// 使用
redis.RDB().Set("key", "value", 0)
```

---

## 基础操作

### 1. String 操作

#### Set/Get
```go
import (
    "context"
    "time"
)

ctx := context.Background()

// 设置值（永不过期）
err := redis.RDB().Set(ctx, "user:1:name", "张三", 0).Err()

// 设置值（带过期时间）
err := redis.RDB().Set(ctx, "captcha:123456", "code_value", 5*time.Minute).Err()

// 获取值
val, err := redis.RDB().Get(ctx, "user:1:name").Result()
if err == redis.Nil {
    // Key 不存在
    logs.Warn("Key不存在")
} else if err != nil {
    // 其他错误
    logs.Error("Redis错误: %v", err)
} else {
    // 成功
    logs.Info("值: %s", val)
}
```

#### SetNX（不存在时设置）
```go
// 分布式锁场景
locked, err := redis.RDB().SetNX(ctx, "lock:order:123", "1", 10*time.Second).Result()
if !locked {
    logs.Warn("锁已被占用")
    return
}
defer redis.RDB().Del(ctx, "lock:order:123")  // 释放锁
```

#### SetEX（设置并指定过期时间）
```go
// 等同于 Set(key, value, duration)
err := redis.RDB().SetEX(ctx, "session:abc", "user_data", 30*time.Minute).Err()
```

#### Incr/Decr（自增/自减）
```go
// 点赞数 +1
newVal, err := redis.RDB().Incr(ctx, "post:123:likes").Result()

// 库存 -1
newVal, err := redis.RDB().Decr(ctx, "product:456:stock").Result()

// 自增指定值
delta := redis.RDB().IncrBy(ctx, "user:1:score", 10).Val()
```

---

### 2. Hash 操作

#### HSet/HGet
```go
// 设置单个字段
redis.RDB().HSet(ctx, "user:1", "name", "张三")

// 设置多个字段
redis.RDB().HSet(ctx, "user:1", map[string]interface{}{
    "name":  "张三",
    "age":   25,
    "email": "zhangsan@example.com",
})

// 获取单个字段
name := redis.RDB().HGet(ctx, "user:1", "name").Val()

// 获取所有字段
fields := redis.RDB().HGetAll(ctx, "user:1").Val()
// fields: map[string]string{"name":"张三", "age":"25", ...}
```

#### HExists
```go
exists := redis.RDB().HExists(ctx, "user:1", "email").Val()
if exists {
    logs.Info("字段存在")
}
```

#### HDel
```go
redis.RDB().HDel(ctx, "user:1", "age", "email")  // 删除多个字段
```

---

### 3. List 操作

#### LPush/RPush（左/右插入）
```go
// 左插入（头部）
redis.RDB().LPush(ctx, "notifications", "消息1", "消息2")

// 右插入（尾部）
redis.RDB().RPush(ctx, "tasks", "任务1", "任务2")
```

#### LPop/RPop（左/右弹出）
```go
// 左弹出（头部）
msg := redis.RDB().LPop(ctx, "notifications").Val()

// 右弹出（尾部）
task := redis.RDB().RPop(ctx, "tasks").Val()
```

#### LRange（范围查询）
```go
// 获取前10条消息
messages := redis.RDB().LRange(ctx, "notifications", 0, 9).Val()

// 获取所有
all := redis.RDB().LRange(ctx, "notifications", 0, -1).Val()
```

---

### 4. Set 操作

#### SAdd/SRem
```go
// 添加成员
redis.RDB().SAdd(ctx, "tags:123", "Go", "Redis", "MySQL")

// 删除成员
redis.RDB().SRem(ctx, "tags:123", "MySQL")
```

#### SMembers（获取所有成员）
```go
members := redis.RDB().SMembers(ctx, "tags:123").Val()
// []string{"Go", "Redis"}
```

#### SIsMember（判断是否存在）
```go
exists := redis.RDB().SIsMember(ctx, "tags:123", "Go").Val()
```

---

### 5. Sorted Set 操作

#### ZAdd（添加成员）
```go
// 排行榜：用户ID -> 分数
redis.RDB().ZAdd(ctx, "leaderboard", redis.Z{Score: 100, Member: "user:1"})
redis.RDB().ZAdd(ctx, "leaderboard", redis.Z{Score: 200, Member: "user:2"})
```

#### ZRange（范围查询）
```go
// 获取前10名（升序）
top10 := redis.RDB().ZRange(ctx, "leaderboard", 0, 9).Val()

// 获取前10名（降序）
top10 := redis.RDB().ZRevRange(ctx, "leaderboard", 0, 9).Val()
```

#### ZScore（获取分数）
```go
score := redis.RDB().ZScore(ctx, "leaderboard", "user:1").Val()
```

#### ZIncrBy（增加分数）
```go
newScore := redis.RDB().ZIncrBy(ctx, "leaderboard", 10, "user:1").Val()
```

---

## 通用操作

### Exists（判断 Key 是否存在）
```go
count := redis.RDB().Exists(ctx, "user:1", "user:2").Val()  // 返回存在的个数
```

### Del（删除 Key）
```go
redis.RDB().Del(ctx, "user:1", "user:2", "user:3")
```

### Expire（设置过期时间）
```go
redis.RDB().Expire(ctx, "session:abc", 30*time.Minute)
```

### TTL（查询剩余过期时间）
```go
ttl := redis.RDB().TTL(ctx, "session:abc").Val()
// -1: 永不过期
// -2: Key 不存在
// >0: 剩余秒数
```

### Keys（查找 Key）⚠️ 生产环境慎用
```go
keys := redis.RDB().Keys(ctx, "user:*").Val()  // 匹配所有 user:* 的 Key
```

### Scan（迭代 Key）推荐
```go
var cursor uint64
for {
    keys, cursor, err := redis.RDB().Scan(ctx, cursor, "user:*", 100).Result()
    // 处理 keys...
    
    if cursor == 0 {
        break
    }
}
```

---

## Pipeline（批量操作）

### 提高性能
```go
pipe := redis.RDB().Pipeline()

pipe.Set(ctx, "user:1:name", "张三", 0)
pipe.Set(ctx, "user:2:name", "李四", 0)
pipe.Incr(ctx, "counter")

// 执行批量命令
_, err := pipe.Exec(ctx)
```

---

## 事务（Watch）

### 乐观锁
```go
err := redis.RDB().Watch(ctx, func(tx *redis.Tx) error {
    // 读取当前值
    val, err := tx.Get(ctx, "balance").Int()
    if err != nil && err != redis.Nil {
        return err
    }
    
    // 修改值
    newVal := val - 100
    if newVal < 0 {
        return errors.New("余额不足")
    }
    
    // 提交事务
    _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
        pipe.Set(ctx, "balance", newVal, 0)
        return nil
    })
    
    return err
}, "balance")
```

---

## 错误处理

### 判断 Key 不存在
```go
val, err := redis.RDB().Get(ctx, "key").Result()
if err == redis.Nil {
    logs.Warn("Key 不存在")
} else if err != nil {
    logs.Error("Redis 错误: %v", err)
} else {
    logs.Info("值: %s", val)
}
```

---

## Context 管理

### 带超时的 Context
```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

val, err := redis.RDB().Get(ctx, "key").Result()
if err == context.DeadlineExceeded {
    logs.Error("Redis 操作超时")
}
```

---

## 最佳实践

### ✅ 推荐做法
1. **使用 Context 控制超时**
2. **批量操作使用 Pipeline**
3. **生产环境使用 Scan 而不是 Keys**
4. **设置合理的过期时间**，避免内存泄漏
5. **Key 命名遵循规范**（见 04-key-naming.md）

### ❌ 避免做法
1. ❌ 不要使用 `Keys *`（阻塞 Redis）
2. ❌ 不要存储过大的值（>10MB）
3. ❌ 不要忘记处理 `redis.Nil`
4. ❌ 不要在循环中逐条执行命令

---

## 常用场景

### 1. 验证码存储（5分钟过期）
```go
code := "123456"
redis.RDB().Set(ctx, fmt.Sprintf("LOGIN_CODE:%s", email), code, 5*time.Minute)

// 验证
storedCode := redis.RDB().Get(ctx, fmt.Sprintf("LOGIN_CODE:%s", email)).Val()
if storedCode == inputCode {
    redis.RDB().Del(ctx, fmt.Sprintf("LOGIN_CODE:%s", email))  // 验证后删除
}
```

### 2. Session 存储（30分钟）
```go
sessionData := map[string]interface{}{
    "user_id": 123,
    "role":    "admin",
}
redis.RDB().HSet(ctx, fmt.Sprintf("SESSION:%s", sessionID), sessionData)
redis.RDB().Expire(ctx, fmt.Sprintf("SESSION:%s", sessionID), 30*time.Minute)
```

### 3. 限流（令牌桶）
```go
key := fmt.Sprintf("RATE_LIMIT:%s", userIP)
count := redis.RDB().Incr(ctx, key).Val()

if count == 1 {
    redis.RDB().Expire(ctx, key, 1*time.Second)  // 1秒窗口
}

if count > 100 {
    logs.Warn("请求过于频繁")
    return
}
```

---

## 参考资料
- go-redis 文档: https://redis.uptrace.dev/
- Redis 命令参考: https://redis.io/commands

