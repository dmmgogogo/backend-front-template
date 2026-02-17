package api

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// SupportOrder iOS 内购订单流水（幂等去重，按 transaction_id 唯一）
// 用于赞助/打赏等场景，后端验 Apple 收据后写入
func init() {
	orm.RegisterModel(new(SupportOrder))
}

// SupportOrder 赞助订单
type SupportOrder struct {
	ID            int64     `orm:"pk;auto;column(id)" json:"id"`
	UserID        int64     `orm:"column(user_id);index" json:"user_id"`
	Platform      string    `orm:"column(platform);size(20);default(ios)" json:"platform"`
	ProductID     string    `orm:"column(product_id);size(128)" json:"product_id"`
	TransactionID string    `orm:"column(transaction_id);size(128);unique" json:"transaction_id"`
	Amount        float64   `orm:"column(amount);digits(12);decimals(2);default(0)" json:"amount"`
	ReceiptData   string    `orm:"column(receipt_data);type(text);null" json:"receipt_data"`
	CreatedAt     time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
}

func (s *SupportOrder) TableName() string {
	return "app_support_orders"
}
