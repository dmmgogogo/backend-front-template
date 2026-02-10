package admin

import (
	"encoding/json"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

// OperationLog 管理员操作日志
type OperationLog struct {
	ID            int64  `json:"id" orm:"pk;column(id);auto"`
	AdminUserID   int64  `json:"admin_user_id" orm:"column(admin_user_id);index"`        // 管理员用户ID
	AdminUsername string `json:"admin_username" orm:"column(admin_username)"`            // 管理员用户名
	OperationType string `json:"operation_type" orm:"column(operation_type);index"`      // 操作类型：create/update/delete/export
	Module        string `json:"module" orm:"column(module);index"`                      // 操作模块
	Action        string `json:"action" orm:"column(action)"`                            // 具体操作
	TargetType    string `json:"target_type" orm:"column(target_type)"`                  // 目标类型
	TargetID      int64  `json:"target_id" orm:"column(target_id)"`                      // 目标ID
	RequestPath   string `json:"request_path" orm:"column(request_path)"`                // 请求路径
	RequestMethod string `json:"request_method" orm:"column(request_method)"`            // 请求方法
	RequestParams string `json:"request_params" orm:"column(request_params);type(text)"` // 请求参数
	IPAddress     string `json:"ip_address" orm:"column(ip_address)"`                    // IP地址
	UserAgent     string `json:"user_agent" orm:"column(user_agent)"`                    // User-Agent
	Status        int    `json:"status" orm:"column(status)"`                            // 1=成功 0=失败
	ErrorMsg      string `json:"error_msg" orm:"column(error_msg);type(text)"`           // 错误信息
	CreatedTime   int64  `json:"created_time" orm:"column(created_time);index"`          // 创建时间
}

func init() {
	orm.RegisterModel(new(OperationLog))
}

func (l *OperationLog) TableName() string {
	return "app_admin_operation_logs"
}

// LogOperationParams 记录操作日志的参数
type LogOperationParams struct {
	AdminUserID   int64                  // 管理员用户ID
	AdminUsername string                 // 管理员用户名
	OperationType string                 // create/update/delete/export
	Module        string                 // 模块名称
	Action        string                 // 具体操作
	TargetType    string                 // 目标类型
	TargetID      int64                  // 目标ID
	RequestPath   string                 // 请求路径
	RequestMethod string                 // 请求方法
	RequestParams map[string]interface{} // 请求参数
	IPAddress     string                 // IP地址
	UserAgent     string                 // User-Agent
	Status        int                    // 1=成功 0=失败
	ErrorMsg      string                 // 错误信息
}

// LogOperation 记录管理员操作日志
func LogOperation(params LogOperationParams) error {
	db := orm.NewOrm()

	// 序列化请求参数为 JSON
	requestParamsJSON := ""
	if params.RequestParams != nil {
		jsonBytes, err := json.Marshal(params.RequestParams)
		if err != nil {
			logs.Warn("[LogOperation] Failed to marshal request params: %v", err)
		} else {
			requestParamsJSON = string(jsonBytes)
		}
	}

	log := &OperationLog{
		AdminUserID:   params.AdminUserID,
		AdminUsername: params.AdminUsername,
		OperationType: params.OperationType,
		Module:        params.Module,
		Action:        params.Action,
		TargetType:    params.TargetType,
		TargetID:      params.TargetID,
		RequestPath:   params.RequestPath,
		RequestMethod: params.RequestMethod,
		RequestParams: requestParamsJSON,
		IPAddress:     params.IPAddress,
		UserAgent:     params.UserAgent,
		Status:        params.Status,
		ErrorMsg:      params.ErrorMsg,
		CreatedTime:   time.Now().Unix(),
	}

	_, err := db.Insert(log)
	if err != nil {
		logs.Error("[LogOperation] Insert log failed: %v", err)
		return err
	}

	logs.Info("[LogOperation] Admin operation logged: admin=%s(%d), action=%s, target=%s:%d, status=%d",
		params.AdminUsername, params.AdminUserID, params.Action, params.TargetType, params.TargetID, params.Status)
	return nil
}

// GetOperationLogs 查询操作日志列表
func GetOperationLogs(page, pageSize int, adminUserID int64, adminUsername string, operationType, module string) ([]OperationLog, int64, error) {
	db := orm.NewOrm()
	var logList []OperationLog

	qb := db.QueryTable(new(OperationLog))

	// 筛选条件
	if adminUserID > 0 {
		qb = qb.Filter("admin_user_id", adminUserID)
	}
	if adminUsername != "" {
		qb = qb.Filter("admin_username", adminUsername)
	}
	if operationType != "" {
		qb = qb.Filter("operation_type", operationType)
	}
	if module != "" {
		qb = qb.Filter("module", module)
	}

	// 统计总数
	total, err := qb.Count()
	if err != nil {
		logs.Error("[GetOperationLogs] Count error: %v", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	_, err = qb.OrderBy("-created_time").Limit(pageSize, offset).All(&logList)
	if err != nil {
		logs.Error("[GetOperationLogs] Query error: %v", err)
		return nil, 0, err
	}

	return logList, total, nil
}
