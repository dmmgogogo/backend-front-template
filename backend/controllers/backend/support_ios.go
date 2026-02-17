package backend

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"e-woms/conf"
	dto "e-woms/dto/backend"
	apiModel "e-woms/models/backend"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// VerifyIOSSupportPurchase iOS 内购验单并累加赞助金额
// 客户端（Flutter in_app_purchase）购买成功后上传 product_id、transaction_id、receipt_data，服务端向 Apple 验单后写入 app_support_orders 并更新用户 support_total_amount / support_level
// @Title iOS内购验单
// @Description 验证 App Store receipt 并累加用户赞助金额
// @Tags 用户-赞助
// @Accept json
// @Produce json
// @Param body body dto.VerifyIOSSupportReq true "请求参数"
// @Success 200 {object} dto.VerifyIOSSupportResp "验单成功"
// @router /api/backend/support/ios/verify [post]
func (c *UserController) VerifyIOSSupportPurchase() {
	userID := c.GetCurrentUserID()
	if userID == 0 {
		c.Error(conf.UNAUTHORIZED, "请先登录")
		return
	}

	var req dto.VerifyIOSSupportReq
	if err := c.ParseJson(&req); err != nil {
		c.Error(conf.PARAMS_ERROR, "参数解析失败")
		return
	}

	if req.ProductID == "" || req.TransactionID == "" || req.ReceiptData == "" {
		c.Error(conf.PARAMS_ERROR, "缺少必要参数")
		return
	}
	logs.Info("[VerifyIOSSupportPurchase] user=%d product=%s tx=%s receipt_len=%d", userID, req.ProductID, req.TransactionID, len(req.ReceiptData))

	if strings.Count(req.ReceiptData, ".") == 2 {
		logs.Warn("[VerifyIOSSupportPurchase] jws receipt detected, client should send base64 receipt")
		c.Error(conf.PARAMS_ERROR, "receipt_data 为 JWS 格式，请使用 base64 收据")
		return
	}
	trimmedReceipt := strings.TrimSpace(req.ReceiptData)
	if strings.HasPrefix(trimmedReceipt, "{") || strings.HasPrefix(trimmedReceipt, "[") {
		logs.Warn("[VerifyIOSSupportPurchase] json receipt detected")
		c.Error(conf.PARAMS_ERROR, "receipt_data 应为 App Store base64 收据")
		return
	}

	if _, err := base64.StdEncoding.DecodeString(req.ReceiptData); err != nil {
		c.Error(conf.PARAMS_ERROR, "receipt_data 格式错误")
		return
	}

	amount, ok := iosProductAmount(req.ProductID)
	if !ok {
		c.Error(conf.PARAMS_ERROR, "未知商品")
		return
	}

	verified, err := verifyAppleReceipt(req.ReceiptData)
	if err != nil {
		logs.Error("[VerifyIOSSupportPurchase] 验单失败: %v", err)
		c.Error(conf.SERVER_ERROR, "验单失败")
		return
	}
	if !verified {
		logs.Warn("[VerifyIOSSupportPurchase] receipt not verified user=%d product=%s tx=%s", userID, req.ProductID, req.TransactionID)
		c.Error(conf.PARAMS_ERROR, "验单未通过")
		return
	}

	user, created, err := apiModel.AddSupportByTransaction(
		userID,
		"ios",
		req.ProductID,
		req.TransactionID,
		amount,
		req.ReceiptData,
	)
	if err != nil {
		logs.Error("[VerifyIOSSupportPurchase] 写入赞助失败: %v", err)
		c.Error(conf.SERVER_ERROR, "写入赞助失败")
		return
	}

	if !created {
		logs.Info("[VerifyIOSSupportPurchase] transaction reused: %s", req.TransactionID)
	}
	logs.Info("[VerifyIOSSupportPurchase] success user=%d product=%s tx=%s amount=%.2f total=%.2f level=%d", userID, req.ProductID, req.TransactionID, amount, user.SupportTotalAmount, user.SupportLevel)

	c.Success(dto.VerifyIOSSupportResp{
		SupportTotalAmount: user.SupportTotalAmount,
		SupportLevel:       user.SupportLevel,
		Amount:             amount,
		TransactionID:      req.TransactionID,
	})
}

// iosProductAmount 商品 ID 对应金额（与 App Store Connect 内购商品一致，可按项目修改）
func iosProductAmount(productID string) (float64, bool) {
	productMap := map[string]float64{
		"com.yourapp.support.1":    1,
		"com.yourapp.support.5":    5,
		"com.yourapp.support.10":   10,
		"com.yourapp.support.50":   50,
		"com.yourapp.support.100":  100,
		"com.yourapp.support.300":  300,
		"com.yourapp.support.500":  500,
		"com.yourapp.support.1000": 1000,
	}
	amount, ok := productMap[productID]
	return amount, ok
}

// verifyAppleReceipt 向 Apple 验单（生产 + 沙盒回退）
func verifyAppleReceipt(receiptData string) (bool, error) {
	start := time.Now()
	sharedSecret := web.AppConfig.DefaultString("ios_iap_shared_secret", "")
	if strings.TrimSpace(sharedSecret) == "" {
		return false, fmt.Errorf("ios_iap_shared_secret is empty")
	}
	payload := map[string]interface{}{
		"receipt-data":             receiptData,
		"password":                 sharedSecret,
		"exclude-old-transactions": true,
	}

	verify := func(url string) (int, error) {
		stepStart := time.Now()
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return -1, err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 12 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return -1, err
		}
		defer resp.Body.Close()

		var result struct {
			Status int `json:"status"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return -1, err
		}
		logs.Info("[verifyAppleReceipt] url=%s status=%d elapsed_ms=%d", url, result.Status, time.Since(stepStart).Milliseconds())
		return result.Status, nil
	}

	status, err := verify("https://buy.itunes.apple.com/verifyReceipt")
	if err != nil {
		return false, err
	}

	if status == 21007 {
		logs.Info("[verifyAppleReceipt] got 21007, fallback sandbox")
		status, err = verify("https://sandbox.itunes.apple.com/verifyReceipt")
		if err != nil {
			return false, err
		}
	}
	logs.Info("[verifyAppleReceipt] final_status=%d total_elapsed_ms=%d", status, time.Since(start).Milliseconds())
	return status == 0, nil
}
