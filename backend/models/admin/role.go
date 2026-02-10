package admin

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Role 企业角色表
type Role struct {
	ID          int64  `json:"id" orm:"pk;column(id);auto"`
	MerchantID  int64  `json:"merchant_id" orm:"column(merchant_id);index"`
	RoleName    string `json:"role_name" orm:"column(role_name)"`
	RoleCode    string `json:"role_code" orm:"column(role_code)"`
	IsSystem    int    `json:"is_system" orm:"column(is_system)"` // 1-系统预置, 0-自定义
	Description string `json:"description" orm:"column(description)"`
	Status      int    `json:"status" orm:"column(status)"` // 0-禁用, 1-启用
	CreatedTime int64  `json:"created_time" orm:"column(created_time)"`
	UpdatedTime int64  `json:"updated_time" orm:"column(updated_time)"`
}

func init() {
	orm.RegisterModel(new(Role))
}

func (r *Role) TableName() string {
	return "app_roles"
}

// Create 创建角色
func (r *Role) Create() error {
	r.CreatedTime = time.Now().Unix()
	r.UpdatedTime = time.Now().Unix()

	db := orm.NewOrm()
	_, err := db.Insert(r)
	return err
}

// GetByID 根据ID查询角色
func (r *Role) GetByID(roleId int64, merchantID int64) error {
	db := orm.NewOrm()
	err := db.QueryTable(r.TableName()).
		Filter("id", roleId).
		Filter("merchant_id", merchantID).
		One(r)
	return err
}

// Update 更新角色
func (r *Role) Update() error {
	r.UpdatedTime = time.Now().Unix()

	db := orm.NewOrm()
	_, err := db.Update(r, "RoleName", "Description", "Status", "UpdatedTime")
	return err
}

// Delete 删除角色
func (r *Role) Delete() error {
	// 系统角色不可删除
	if r.IsSystem == 1 {
		return orm.ErrNoRows // 用错误表示不可删除
	}

	db := orm.NewOrm()
	_, err := db.Delete(r)
	return err
}

// List 查询角色列表
func (r *Role) List(keyword string, status int, page int, pageSize int) ([]Role, int64, error) {
	db := orm.NewOrm()
	qs := db.QueryTable(r.TableName())

	// 关键词搜索
	if keyword != "" {
		qs = qs.Filter("role_name__icontains", keyword)
	}

	// 状态筛选
	if status >= 0 {
		qs = qs.Filter("status", status)
	}

	// 总数
	total, _ := qs.Count()

	// 分页
	var roles []Role
	offset := (page - 1) * pageSize
	_, err := qs.OrderBy("-created_time").Limit(pageSize, offset).All(&roles)

	return roles, total, err
}

// CheckRoleCodeExists 检查角色代码是否已存在
func (r *Role) CheckRoleCodeExists(roleCode string, excludeID int64) (bool, error) {
	db := orm.NewOrm()
	qs := db.QueryTable(r.TableName()).
		Filter("role_code", roleCode)

	if excludeID > 0 {
		qs = qs.Exclude("id", excludeID)
	}

	exists := qs.Exist()
	return exists, nil
}
