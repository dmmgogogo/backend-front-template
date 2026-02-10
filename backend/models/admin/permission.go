package admin

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Permission 企业权限表
type Permission struct {
	ID               int64  `json:"id" orm:"pk;column(id);auto"`
	PermissionName   string `json:"permission_name" orm:"column(permission_name)"`
	PermissionNameEn string `json:"permission_name_en" orm:"column(permission_name_en)"`
	PermissionCode   string `json:"permission_code" orm:"column(permission_code)"`
	APIRoute         string `json:"api_route" orm:"column(api_route)"`
	HTTPMethod       string `json:"http_method" orm:"column(http_method)"`
	Module           string `json:"module" orm:"column(module)"`
	Description      string `json:"description" orm:"column(description)"`
	DescriptionEn    string `json:"description_en" orm:"column(description_en)"`
	CreatedTime      int64  `json:"created_time" orm:"column(created_time)"`
	UpdatedTime      int64  `json:"updated_time" orm:"column(updated_time)"`
}

func init() {
	orm.RegisterModel(new(Permission))
}

func (p *Permission) TableName() string {
	return "app_permissions"
}

// Create 创建权限
func (p *Permission) Create() error {
	p.CreatedTime = time.Now().Unix()
	p.UpdatedTime = time.Now().Unix()

	db := orm.NewOrm()
	_, err := db.Insert(p)
	return err
}

// GetByID 根据ID查询权限
func (p *Permission) GetByID(id int64) error {
	db := orm.NewOrm()
	p.ID = id
	err := db.Read(p)
	return err
}

// List 查询权限列表
func (p *Permission) List(module string, keyword string, page int, pageSize int) ([]Permission, int64, error) {
	db := orm.NewOrm()
	qs := db.QueryTable(p.TableName())

	// 构建基础条件
	baseCond := orm.NewCondition()

	// 模块筛选
	if module != "" {
		baseCond = baseCond.And("module", module)
	}

	// 关键词搜索（支持中英文名称）
	if keyword != "" {
		keywordCond := orm.NewCondition().Or("permission_name__icontains", keyword).Or("permission_name_en__icontains", keyword)
		baseCond = baseCond.AndCond(keywordCond)
	}

	if module != "" || keyword != "" {
		qs = qs.SetCond(baseCond)
	}

	// 总数
	total, _ := qs.Count()

	// 分页
	var permissions []Permission
	offset := (page - 1) * pageSize
	_, err := qs.OrderBy("module", "id").Limit(pageSize, offset).All(&permissions)

	return permissions, total, err
}

// GetByIDs 根据ID列表批量查询
func (p *Permission) GetByIDs(ids []int64) ([]Permission, error) {
	db := orm.NewOrm()
	var permissions []Permission
	_, err := db.QueryTable(p.TableName()).Filter("id__in", ids).All(&permissions)
	return permissions, err
}

// Update 更新权限
func (p *Permission) Update() error {
	p.UpdatedTime = time.Now().Unix()

	db := orm.NewOrm()
	_, err := db.Update(p, "PermissionName", "PermissionNameEn", "PermissionCode", "APIRoute", "HTTPMethod", "Module", "Description", "DescriptionEn", "UpdatedTime")
	return err
}

// Delete 删除权限
func (p *Permission) Delete() error {
	db := orm.NewOrm()
	_, err := db.Delete(p)
	return err
}

// CheckCodeExists 检查权限代码是否已存在
func (p *Permission) CheckCodeExists(code string, excludeID int64) (bool, error) {
	db := orm.NewOrm()
	qs := db.QueryTable(p.TableName()).Filter("permission_code", code)

	if excludeID > 0 {
		qs = qs.Exclude("id", excludeID)
	}

	exists := qs.Exist()
	return exists, nil
}
