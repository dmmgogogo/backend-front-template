package admin

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

// User
type User struct {
	ID            int64  `json:"id" orm:"pk;column(id);auto"`
	Username      string `json:"username" orm:"column(username)"`               // 邮箱登录
	Password      string `json:"-" orm:"column(password)"`                      // 不返回给前端
	RealName      string `json:"real_name" orm:"column(real_name)"`             // 真实姓名
	Email         string `json:"email" orm:"column(email);unique"`              // 全局唯一邮箱
	Phone         string `json:"phone" orm:"column(phone)"`                     // 手机号
	Status        int    `json:"status" orm:"column(status);index"`             // 0-禁用, 1-启用
	FirstLogin    int    `json:"first_login" orm:"column(first_login)"`         // 首次登录, 需要修改密码, 0=未修改 1=已修改
	VerifyCode    string `json:"verify_code" orm:"column(verify_code)"`         // Google验证码
	LastLoginTime int64  `json:"last_login_time" orm:"column(last_login_time)"` // 最后登录时间
	CreatedTime   int64  `json:"created_time" orm:"column(created_time);index"` // 创建时间
	UpdatedTime   int64  `json:"updated_time" orm:"column(updated_time)"`       // 更新时间
}

// UserInfoRes 用户信息响应（不包含密码）
type UserInfoRes struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	RealName      string `json:"real_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Status        int    `json:"status"`
	LastLoginTime int64  `json:"last_login_time"`
	CreatedTime   int64  `json:"created_time"`
	UpdatedTime   int64  `json:"updated_time"`
	FirstLogin    int    `json:"first_login"`
}

func init() {
	orm.RegisterModel(new(User))
}

func (u *User) TableName() string {
	return "app_admin_users"
}

// EncryptPassword 加密密码
func EncryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password + "super_admin_salt_2026"))
	return hex.EncodeToString(h.Sum(nil))
}

// Create 创建用户
func (u *User) Create() error {
	u.CreatedTime = time.Now().Unix()
	u.UpdatedTime = time.Now().Unix()
	u.Password = EncryptPassword(u.Password)

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

// GetByUsername 根据用户名查询
func (u *User) GetByUsername(username string) error {
	db := orm.NewOrm()
	err := db.QueryTable(u.TableName()).
		Filter("username", username).
		One(u)
	return err
}

// CheckUsernameExists 检查用户名是否已存在
func (u *User) CheckUsernameExists(username string, excludeID int64) (bool, error) {
	db := orm.NewOrm()
	qs := db.QueryTable(u.TableName()).
		Filter("username", username)

	if excludeID > 0 {
		qs = qs.Exclude("id", excludeID)
	}

	exists := qs.Exist()
	return exists, nil
}

// ToUserInfoRes 转换为响应格式（去除密码等敏感信息）
func (u *User) ToUserInfoRes() *UserInfoRes {
	return &UserInfoRes{
		ID:            u.ID,
		Username:      u.Username,
		RealName:      u.RealName,
		Email:         u.Email,
		Phone:         u.Phone,
		Status:        u.Status,
		LastLoginTime: u.LastLoginTime,
		CreatedTime:   u.CreatedTime,
		UpdatedTime:   u.UpdatedTime,
		FirstLogin:    u.FirstLogin,
	}
}

// GetByEmail 根据邮箱查询用户（全局唯一）
func (u *User) GetByEmail(email string) error {
	db := orm.NewOrm()
	err := db.QueryTable(u.TableName()).
		Filter("email", email).
		One(u)
	return err
}

// CheckEmailExists 检查邮箱是否已存在（全局唯一）
func (u *User) CheckEmailExists(email string, excludeID int64) (bool, error) {
	db := orm.NewOrm()
	qs := db.QueryTable(u.TableName()).
		Filter("email", email)

	if excludeID > 0 {
		qs = qs.Exclude("id", excludeID)
	}

	exists := qs.Exist()
	return exists, nil
}

func (u *User) LoginByUsername(username, password string) error {
	db := orm.NewOrm()
	err := db.QueryTable(u.TableName()).
		Filter("username", username).
		Filter("status", 1). // 只查询启用的用户
		One(u)
	if err != nil {
		return err
	}

	// 验证密码
	encryptedPassword := EncryptPassword(password)
	logs.Debug("[User][LoginByUsername] encryptedPassword: %s, u.Password: %s", encryptedPassword, u.Password)
	if u.Password != encryptedPassword {
		return orm.ErrNoRows // 密码错误返回用户不存在错误
	}

	// 更新最后登录时间
	u.LastLoginTime = time.Now().Unix()
	_, _ = db.Update(u, "LastLoginTime")

	return nil
}

// ChangePassword 修改密码（验证旧密码）
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	// 验证旧密码
	encryptedOldPassword := EncryptPassword(oldPassword)
	if u.Password != encryptedOldPassword {
		return orm.ErrNoRows // 旧密码错误
	}

	// 更新密码和first_login状态
	db := orm.NewOrm()
	u.Password = EncryptPassword(newPassword)
	u.FirstLogin = 1 // 标记已完成首次登录密码修改
	u.UpdatedTime = time.Now().Unix()
	_, err := db.Update(u, "Password", "FirstLogin", "UpdatedTime")
	return err
}
