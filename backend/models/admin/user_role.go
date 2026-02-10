package admin

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// UserRole 用户-角色关联表
type UserRole struct {
	ID          int64 `json:"id" orm:"pk;column(id);auto"`
	UserID      int64 `json:"user_id" orm:"column(user_id)"`
	RoleID      int64 `json:"role_id" orm:"column(role_id)"`
	CreatedTime int64 `json:"created_time" orm:"column(created_time)"`
}

func init() {
	orm.RegisterModel(new(UserRole))
}

func (ur *UserRole) TableName() string {
	return "app_user_roles"
}

// AssignRole 为用户分配角色
func (ur *UserRole) AssignRole(userID int64, roleID int64) error {
	db := orm.NewOrm()

	// 检查是否已分配
	existing := &UserRole{}
	err := db.QueryTable(ur.TableName()).
		Filter("user_id", userID).
		Filter("role_id", roleID).
		One(existing)

	// 如果已存在,直接返回成功
	if err == nil {
		return nil
	}

	// 插入新的用户-角色关联
	ur.UserID = userID
	ur.RoleID = roleID
	ur.CreatedTime = time.Now().Unix()

	_, err = db.Insert(ur)
	return err
}

// RemoveRole 移除用户的角色
func (ur *UserRole) RemoveRole(userID int64, roleID int64) error {
	db := orm.NewOrm()
	_, err := db.QueryTable(ur.TableName()).
		Filter("user_id", userID).
		Filter("role_id", roleID).
		Delete()
	return err
}

// GetUserRoles 查询用户的所有角色ID
func (ur *UserRole) GetUserRoles(userID int64) ([]UserRole, error) {
	db := orm.NewOrm()
	var userRoles []UserRole
	_, err := db.QueryTable(ur.TableName()).
		Filter("user_id", userID).
		All(&userRoles)

	if err != nil {
		return nil, err
	}

	return userRoles, nil
}

// GetUserRolesWithDetail 查询用户的所有角色(包含角色详情)
func (ur *UserRole) GetUserRolesWithDetail(userID int64) ([]Role, error) {
	// 先获取角色ID列表
	userRoles, err := ur.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []Role{}, nil
	}

	var roleIDs []int64
	for _, userRole := range userRoles {
		roleIDs = append(roleIDs, userRole.RoleID)
	}

	// 批量查询角色详情
	db := orm.NewOrm()
	var roles []Role
	roleModel := &Role{}
	_, err = db.QueryTable(roleModel.TableName()).
		Filter("id__in", roleIDs).
		All(&roles)

	if err != nil {
		return nil, err
	}

	return roles, nil
}

// 获取指定用户在指定企业的角色列表
func (ur *UserRole) GetUserRolesByUserIDAndMerchantID(userID int64) (roles []string, err error) {
	var permissionIDs []int64
	var permissions []Permission

	// 如果是管理员则默认查企业全部
	if userID > 0 {
		// 查询用户的权限ID列表
		rolePermissionModel := &RolePermission{}
		permissionIDs, err = rolePermissionModel.GetUserPermissions(userID)
		if err != nil {
			return
		}

		// 如果没有权限,返回空列表
		if len(permissionIDs) == 0 {
			return
		}

		// 批量查询权限详情
		permissionModel := &Permission{}
		permissions, err = permissionModel.GetByIDs(permissionIDs)
		if err != nil {
			return
		}
	} else {
		// 查询当前企业下面的全部权限
		permissionModel := &Permission{}
		permissions, _, err = permissionModel.List("", "", 1, 10000)
		if err != nil {
			return
		}
	}

	for _, permission := range permissions {
		roles = append(roles, permission.PermissionCode)
	}
	return
}

func (ur *UserRole) RemoveAllRoles(userID int64) error {
	db := orm.NewOrm()
	_, err := db.QueryTable(ur.TableName()).
		Filter("user_id", userID).
		Delete()
	return err
}
