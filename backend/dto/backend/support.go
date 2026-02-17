package backend

// VerifyIOSSupportReq iOS 内购验单请求（App Store 收据校验）
// 对应 Flutter in_app_purchase 购买成功后上传服务端验单
type VerifyIOSSupportReq struct {
	ProductID     string `json:"product_id" valid:"Required"`
	TransactionID string `json:"transaction_id" valid:"Required"`
	ReceiptData   string `json:"receipt_data" valid:"Required"`
}

// VerifyIOSSupportResp iOS 内购验单响应
type VerifyIOSSupportResp struct {
	SupportTotalAmount float64 `json:"support_total_amount"`
	SupportLevel       int     `json:"support_level"`
	Amount             float64 `json:"amount"`
	TransactionID      string  `json:"transaction_id"`
}
