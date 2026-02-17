package api

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// AddSupportByTransaction 累加赞助金额（按 transaction_id 幂等）
// 用于 iOS 内购验单成功后更新用户赞助总额与等级
func AddSupportByTransaction(userID int64, platform, productID, transactionID string, amount float64, receiptData string) (*User, bool, error) {
	o := orm.NewOrm()

	exist := &SupportOrder{}
	err := o.QueryTable(new(SupportOrder).TableName()).
		Filter("transaction_id", transactionID).
		One(exist)
	if err == nil {
		user := &User{ID: userID}
		if e := o.Read(user); e != nil {
			return nil, false, e
		}
		return user, false, nil
	}
	if err != orm.ErrNoRows {
		return nil, false, err
	}

	tx, err := o.Begin()
	if err != nil {
		return nil, false, err
	}

	user := &User{ID: userID}
	if err := tx.Read(user); err != nil {
		_ = tx.Rollback()
		return nil, false, err
	}

	order := &SupportOrder{
		UserID:        userID,
		Platform:      platform,
		ProductID:     productID,
		TransactionID: transactionID,
		Amount:        amount,
		ReceiptData:   receiptData,
	}
	if _, err := tx.Insert(order); err != nil {
		_ = tx.Rollback()
		return nil, false, err
	}

	user.SupportTotalAmount += amount
	user.SupportLevel = resolveSupportLevel(user.SupportTotalAmount)
	if user.SupportTotalAmount > 0 && user.Vip == 0 {
		user.Vip = 1
	}
	user.UpdatedTime = time.Now().Unix()

	if _, err := tx.Update(user, "SupportTotalAmount", "SupportLevel", "Vip", "UpdatedTime"); err != nil {
		_ = tx.Rollback()
		return nil, false, err
	}

	if err := tx.Commit(); err != nil {
		return nil, false, err
	}

	return user, true, nil
}

func resolveSupportLevel(totalAmount float64) int {
	switch {
	case totalAmount >= 500:
		return 5
	case totalAmount >= 300:
		return 4
	case totalAmount >= 100:
		return 3
	case totalAmount >= 50:
		return 2
	case totalAmount >= 5:
		return 1
	default:
		return 0
	}
}
