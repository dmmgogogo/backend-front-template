# Beego ORM v2 基础操作

## 初始化 ORM

### 使用自定义库初始化（推荐）
```go
import (
    "std-library-slim/dbase"
    "std-library-slim/json"
    "github.com/beego/beego/v2/server/web"
)

func InitMysql() {
    opt := dbase.Opt{}
    err := json.ParseE(web.AppConfig.DefaultString("MYSQL_CONFIG", ""), &opt)
    if err != nil {
        panic(err)
    }
    dbase.Init(&opt)
}
```

---

## Model 定义

### 标准 Model 结构
```go
package models

import "github.com/beego/beego/v2/client/orm"

type User struct {
    ID          int64  `orm:"pk;auto" json:"id"`
    Email       string `orm:"column(email);size(100);unique" json:"email"`
    Password    string `orm:"column(password);size(255)" json:"-"`
    Nickname    string `orm:"column(nickname);size(50)" json:"nickname"`
    Avatar      string `orm:"column(avatar);size(255)" json:"avatar"`
    Status      int    `orm:"column(status);default(1)" json:"status"`
    CreatedTime int64  `orm:"column(created_time)" json:"created_time"`
    UpdatedTime int64  `orm:"column(updated_time)" json:"updated_time"`
}

func init() {
    // 注册模型
    orm.RegisterModel(new(User))
}

func (u *User) TableName() string {
    return "app_users"
}
```

### 字段标签说明

| 标签 | 说明 | 示例 |
|------|------|------|
| `pk` | 主键 | `orm:"pk;auto"` |
| `auto` | 自增 | `orm:"pk;auto"` |
| `column(name)` | 字段名 | `orm:"column(user_id)"` |
| `size(n)` | 字段长度 | `orm:"size(100)"` |
| `unique` | 唯一索引 | `orm:"unique"` |
| `index` | 普通索引 | `orm:"index"` |
| `default(v)` | 默认值 | `orm:"default(1)"` |
| `null` | 允许 NULL | `orm:"null"` |
| `-` | 忽略字段 | `orm:"-"` |

---

## CRUD 操作

### 1. 插入（Create）

#### 单条插入
```go
db := orm.NewOrm()

user := &User{
    Email:       "test@example.com",
    Password:    "hashed_password",
    Nickname:    "Test User",
    Status:      1,
    CreatedTime: time.Now().Unix(),
}

id, err := db.Insert(user)
if err != nil {
    logs.Error("插入失败: %v", err)
    return err
}

logs.Info("插入成功，ID: %d", id)
```

#### 批量插入
```go
users := []User{
    {Email: "user1@example.com", Nickname: "User1"},
    {Email: "user2@example.com", Nickname: "User2"},
}

successNums, err := db.InsertMulti(100, users)  // 100为批量大小
```

---

### 2. 查询（Read）

#### 查询单条记录
```go
db := orm.NewOrm()

// 按主键查询
user := &User{ID: 1}
err := db.Read(user)
if err == orm.ErrNoRows {
    // 记录不存在
} else if err != nil {
    // 其他错误
}

// 按其他字段查询
user := &User{Email: "test@example.com"}
err := db.Read(user, "Email")  // 指定查询字段
```

#### 查询多条记录
```go
db := orm.NewOrm()

var users []*User
num, err := db.QueryTable("app_users").
    Filter("status", 1).
    All(&users)

logs.Info("查询到 %d 条记录", num)
```

#### 带条件查询
```go
var users []*User

_, err := db.QueryTable("app_users").
    Filter("status", 1).                    // 等于
    Filter("level__gte", 5).                // 大于等于
    Filter("nickname__icontains", "test").  // 包含（不区分大小写）
    Exclude("email__endswith", "@temp.com"). // 排除
    OrderBy("-id").                         // 降序
    Limit(20).                              // 限制条数
    Offset(0).                              // 偏移量
    All(&users)
```

#### 分页查询
```go
page := 1
pageSize := 20

var users []*User
query := db.QueryTable("app_users").Filter("status", 1)

// 查询总数
total, _ := query.Count()

// 查询列表
_, err := query.
    OrderBy("-id").
    Limit(pageSize).
    Offset((page - 1) * pageSize).
    All(&users)
```

---

### 3. 更新（Update）

#### 更新指定字段（推荐）
```go
db := orm.NewOrm()

user := &User{ID: 1}
user.Nickname = "新昵称"
user.UpdatedTime = time.Now().Unix()

// 只更新 Nickname 和 UpdatedTime 字段
_, err := db.Update(user, "Nickname", "UpdatedTime")
```

#### 更新所有字段
```go
user := &User{ID: 1}
db.Read(user)  // 先读取

user.Nickname = "新昵称"
user.Avatar = "new_avatar.jpg"

// 更新所有字段（除主键）
_, err := db.Update(user)
```

#### 批量更新
```go
num, err := db.QueryTable("app_users").
    Filter("status", 0).
    Update(orm.Params{
        "status": 1,
        "updated_time": time.Now().Unix(),
    })

logs.Info("更新了 %d 条记录", num)
```

---

### 4. 删除（Delete）

#### 物理删除
```go
db := orm.NewOrm()

user := &User{ID: 1}
_, err := db.Delete(user)
```

#### 逻辑删除（推荐）
```go
user := &User{ID: 1}
user.Status = 0  // 标记为删除
user.UpdatedTime = time.Now().Unix()

_, err := db.Update(user, "Status", "UpdatedTime")
```

#### 批量删除
```go
num, err := db.QueryTable("app_users").
    Filter("status", 0).
    Filter("created_time__lt", oldTimestamp).
    Delete()
```

---

## 查询条件操作符

| 操作符 | 说明 | 示例 |
|--------|------|------|
| (无) | 等于 | `Filter("id", 1)` |
| `__exact` | 等于（显式） | `Filter("id__exact", 1)` |
| `__gt` | 大于 | `Filter("age__gt", 18)` |
| `__gte` | 大于等于 | `Filter("age__gte", 18)` |
| `__lt` | 小于 | `Filter("age__lt", 60)` |
| `__lte` | 小于等于 | `Filter("age__lte", 60)` |
| `__in` | IN 查询 | `Filter("id__in", []int{1,2,3})` |
| `__contains` | 包含（区分大小写） | `Filter("name__contains", "test")` |
| `__icontains` | 包含（不区分大小写） | `Filter("name__icontains", "test")` |
| `__startswith` | 以...开头 | `Filter("email__startswith", "admin")` |
| `__endswith` | 以...结尾 | `Filter("email__endswith", "@qq.com")` |
| `__isnull` | 是否为 NULL | `Filter("deleted_at__isnull", true)` |

---

## 高级查询

### 关联查询
```go
type Order struct {
    ID     int64  `orm:"pk;auto"`
    UserID int64  `orm:"column(user_id)"`
    User   *User  `orm:"rel(fk)"` // 外键关联
}

// 查询订单并加载用户
order := &Order{ID: 1}
db.Read(order)
db.LoadRelated(order, "User")  // 加载关联的用户
```

### 原生 SQL 查询
```go
var maps []orm.Params
num, err := db.Raw("SELECT * FROM app_users WHERE status = ?", 1).Values(&maps)

for _, m := range maps {
    logs.Info("ID: %s, Email: %s", m["id"], m["email"])
}
```

### QueryRow（单行）
```go
var name string
var age int
err := db.Raw("SELECT name, age FROM users WHERE id = ?", 1).QueryRow(&name, &age)
```

---

## 最佳实践

### ✅ 推荐做法
1. **使用 ORM 而不是原生 SQL**（除非必要）
2. **更新时指定字段名**，避免误更新
3. **使用逻辑删除**，保留数据历史
4. **分页查询先 Count 再 All**
5. **为常用查询字段添加索引**

### ❌ 避免做法
1. ❌ 不要在循环中执行 SQL
2. ❌ 不要使用 `Select("*")`（性能问题）
3. ❌ 不要直接拼接 SQL（SQL注入风险）
4. ❌ 不要忘记处理 `orm.ErrNoRows`

---

## 错误处理

### 判断记录不存在
```go
user := &User{ID: 999}
err := db.Read(user)

if err == orm.ErrNoRows {
    logs.Warn("用户不存在")
} else if err != nil {
    logs.Error("查询失败: %v", err)
}
```

### 判断唯一键冲突
```go
_, err := db.Insert(user)
if err != nil {
    if strings.Contains(err.Error(), "Duplicate entry") {
        logs.Warn("邮箱已存在")
    }
}
```

---

## 参考资料
- Beego ORM 文档: https://beego.wiki/docs/mvc/model/overview/
- MySQL 官方文档: https://dev.mysql.com/doc/

