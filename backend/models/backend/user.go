package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

// User 用户表（普通投资用户）
type User struct {
	ID            int64  `json:"id" orm:"pk;column(id);auto"`
	Uid           int64  `json:"uid" orm:"column(uid);unique"`                  // 用户ID
	Username      string `json:"username" orm:"column(username);unique"`        // 用户名
	Email         string `json:"email" orm:"column(email);unique"`              // 邮箱
	Password      string `json:"-" orm:"column(password)"`                      // 不返回给前端
	Nickname      string `json:"nickname" orm:"column(nickname)"`               // 昵称
	Avatar        string `json:"avatar" orm:"column(avatar)"`                   // 头像
	Status        int    `json:"status" orm:"column(status);default(1)"`        // 状态 1:正常 0:禁用
	LastLoginTime int64  `json:"last_login_time" orm:"column(last_login_time)"` // 最后登录时间
	CreatedTime   int64  `json:"created_time" orm:"column(created_time);index"` // 创建时间
	UpdatedTime   int64  `json:"updated_time" orm:"column(updated_time)"`       // 更新时间
}

func init() {
	orm.RegisterModel(new(User))
}

func (u *User) TableName() string {
	return "app_users"
}

// EncryptPassword 加密密码
func EncryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password + "e-woms_salt_2026"))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateInviteCode 生成唯一邀请码（NX + 6位随机字符）
func GenerateInviteCode() (string, error) {
	db := orm.NewOrm()

	// 最多尝试10次生成唯一邀请码
	for i := 0; i < 10; i++ {
		// 生成随机8位字符
		bytes := make([]byte, 4)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}
		code := hex.EncodeToString(bytes) // 生成8位十六进制字符
		code = strings.ToUpper(code)      // 转大写

		// 检查是否已存在
		exists := db.QueryTable("app_users").
			Filter("invite_code", code).
			Exist()

		if !exists {
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique invite code")
}

// GenerateUIDRandom 生成10位随机数
func GenerateUIDRandom() int64 {
	// 注意：Go 1.20+ 已自动初始化随机数生成器，无需手动调用 rand.Seed()
	// 生成10位随机数
	return rand.Int63n(9000000000) + 1000000000
}

// GenerateUID 生成唯一用户ID（U + 8位数字）
func GenerateUID() (int64, error) {
	db := orm.NewOrm()

	// 生成 UID: U + 8位数字（从 maxID+1 开始）
	uid := GenerateUIDRandom()

	// 检查是否已存在（防止并发问题）
	exists := db.QueryTable("app_users").Filter("uid", uid).Exist()
	if exists {
		// 如果存在，再尝试几次
		for i := 0; i < 5; i++ {
			uid = GenerateUIDRandom()
			exists = db.QueryTable("app_users").Filter("uid", uid).Exist()
			if !exists {
				return uid, nil
			}
		}
		return 0, fmt.Errorf("failed to generate unique uid")
	}

	return uid, nil
}

// Create 创建用户
func (u *User) Create() error {
	u.CreatedTime = time.Now().Unix()
	u.UpdatedTime = time.Now().Unix()
	u.Password = EncryptPassword(u.Password)

	// 生成唯一UID
	if u.Uid == 0 {
		uid, err := GenerateUID()
		if err != nil {
			return err
		}
		u.Uid = uid
	}

	db := orm.NewOrm()
	_, err := db.Insert(u)
	return err
}

// GetByID 根据ID查询
func (u *User) GetByID(id int64) error {
	db := orm.NewOrm()
	u.ID = id
	err := db.QueryTable(u.TableName()).
		Filter("id", id).
		One(u)
	return err
}

// GetByEmail 根据邮箱查询
func (u *User) GetByEmail(email string) error {
	db := orm.NewOrm()
	err := db.QueryTable(u.TableName()).
		Filter("email", email).
		One(u)
	return err
}

// GetByUsername 根据用户名查询
func (u *User) GetByUsername(username string) error {
	db := orm.NewOrm()
	err := db.QueryTable(u.TableName()).
		Filter("username", username).
		One(u)
	return err
}

// CheckEmailExists 检查邮箱是否已存在
func CheckEmailExists(email string) (bool, error) {
	db := orm.NewOrm()
	exists := db.QueryTable("app_users").
		Filter("email", email).
		Exist()
	return exists, nil
}

// CheckUsernameExists 检查用户名是否已存在
func CheckUsernameExists(username string) (bool, error) {
	db := orm.NewOrm()
	exists := db.QueryTable("app_users").
		Filter("username", username).
		Exist()
	return exists, nil
}

// Login 用户登录（验证用户名密码）
func (u *User) Login(username, password string) error {
	db := orm.NewOrm()

	// 根据用户名查询用户
	err := db.QueryTable(u.TableName()).
		Filter("username", username).
		One(u)

	if err != nil {
		return err
	}

	// 验证密码
	if u.Password != EncryptPassword(password) {
		return orm.ErrNoRows // 密码错误返回用户不存在错误
	}

	// 检查账号状态
	if u.Status != 1 {
		return fmt.Errorf("account disabled")
	}

	// 更新最后登录时间
	u.LastLoginTime = time.Now().Unix()
	_, _ = db.Update(u, "LastLoginTime")

	return nil
}

// CheckMinerIDExists 检查矿工ID是否存在
func (u *User) CheckMinerIDExists(minerID string) bool {
	db := orm.NewOrm()
	exists := db.QueryTable("app_users").
		Filter("miner_id", minerID).
		Exist()
	return exists
}

// UpdatePassword 更新登录密码
func (u *User) UpdatePassword(newPassword string) error {
	db := orm.NewOrm()
	u.Password = EncryptPassword(newPassword)
	u.UpdatedTime = time.Now().Unix()
	_, err := db.Update(u, "Password", "UpdatedTime")
	return err
}

// CreateUserByAdmin 管理员创建用户
func CreateUserByAdmin(email, password, username, nickname string, status int) (*User, error) {
	db := orm.NewOrm()
	now := time.Now().Unix()

	// 生成 UID（8-10位随机数字）
	uid, err := GenerateUID()
	if err != nil {
		logs.Error("[CreateUserByAdmin] Generate UID error: %v", err)
		return nil, err
	}

	// 加密登录密码
	encryptedPassword := EncryptPassword(password)

	// 创建用户
	user := &User{
		Uid:         uid,
		Email:       email,
		Password:    encryptedPassword,
		Username:    username,
		Nickname:    nickname,
		Status:      status,
		CreatedTime: now,
		UpdatedTime: now,
	}

	// 开启事务 (Beego ORM v2 Begin() 返回 TxOrmer)
	txOrm, err := db.Begin()
	if err != nil {
		logs.Error("[CreateUserByAdmin] Begin transaction error: %v", err)
		return nil, err
	}

	// 使用 defer 确保事务最终被提交或回滚
	success := false
	defer func() {
		if success {
			err := txOrm.Commit()
			if err != nil {
				logs.Error("[CreateUserByAdmin] Commit transaction error: %v", err)
			} else {
				logs.Info("[CreateUserByAdmin] Transaction committed successfully")
			}
		} else {
			txOrm.Rollback()
			logs.Warn("[CreateUserByAdmin] Transaction rolled back")
		}
	}()

	// 插入用户
	userID, err := txOrm.Insert(user)
	if err != nil {
		logs.Error("[CreateUserByAdmin] Insert user error: %v", err)
		return nil, err
	}
	user.ID = userID

	// 标记成功，defer 中会提交事务
	success = true

	logs.Info("[CreateUserByAdmin] User created successfully: id=%d, email=%s", userID, email)
	return user, nil
}
