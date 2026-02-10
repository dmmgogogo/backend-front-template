package admin

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// RolePermission 角色-权限关联表
type RolePermission struct {
	ID           int64 `json:"id" orm:"pk;column(id);auto"`
	RoleID       int64 `json:"role_id" orm:"column(role_id);index"`
	PermissionID int64 `json:"permission_id" orm:"column(permission_id);index"`
	CreatedTime  int64 `json:"created_time" orm:"column(created_time)"`
}

func init() {
	orm.RegisterModel(new(RolePermission))
}

func (rp *RolePermission) TableName() string {
	return "app_role_permissions"
}

// AssignPermissions 为角色分配权限 (批量)
func (rp *RolePermission) AssignPermissions(roleID int64, permissionIDs []int64) error {
	db := orm.NewOrm()

	// 先删除该角色的所有权限
	_, err := db.QueryTable(rp.TableName()).
		Filter("role_id", roleID).
		Delete()

	if err != nil {
		return err
	}

	// 批量插入新权限
	var rolePermissions []RolePermission
	now := time.Now().Unix()
	for _, permID := range permissionIDs {
		rolePermissions = append(rolePermissions, RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
			CreatedTime:  now,
		})
	}

	if len(rolePermissions) > 0 {
		_, err = db.InsertMulti(len(rolePermissions), rolePermissions)
	}

	return err
}

// RemovePermission 移除角色的单个权限
func (rp *RolePermission) RemovePermission(roleID int64, permissionID int64) error {
	db := orm.NewOrm()
	_, err := db.QueryTable(rp.TableName()).
		Filter("role_id", roleID).
		Filter("permission_id", permissionID).
		Delete()
	return err
}

// GetRolePermissions 查询角色的所有权限ID
func (rp *RolePermission) GetRolePermissions(roleID int64) ([]int64, error) {
	db := orm.NewOrm()
	var rolePermissions []RolePermission
	_, err := db.QueryTable(rp.TableName()).
		Filter("role_id", roleID).
		All(&rolePermissions)

	if err != nil {
		return nil, err
	}

	var permissionIDs []int64
	for _, rp := range rolePermissions {
		permissionIDs = append(permissionIDs, rp.PermissionID)
	}

	return permissionIDs, nil
}

// GetUserPermissions 查询用户的所有权限ID (通过用户的角色)
func (rp *RolePermission) GetUserPermissions(userID int64) ([]int64, error) {
	db := orm.NewOrm()

	// 先查询用户的所有角色ID
	var userRoles []UserRole

	userRoleModel := &UserRole{}
	userRoles, err := userRoleModel.GetUserRoles(userID)

	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []int64{}, nil
	}

	// 提取角色ID列表
	var roleIDs []int64
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	// 查询这些角色的所有权限ID
	var rolePermissions []RolePermission
	_, err = db.QueryTable(rp.TableName()).
		Filter("role_id__in", roleIDs).
		All(&rolePermissions)

	if err != nil {
		return nil, err
	}

	// 去重权限ID
	permissionIDMap := make(map[int64]bool)
	for _, rp := range rolePermissions {
		permissionIDMap[rp.PermissionID] = true
	}

	var permissionIDs []int64
	for permID := range permissionIDMap {
		permissionIDs = append(permissionIDs, permID)
	}

	return permissionIDs, nil
}
