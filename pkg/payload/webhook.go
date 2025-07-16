package payload

import (
	"net/http"
)

// WebHook payload
//
// WARNING: Do NOT convert amounts directly to float when performing operations on money (due to floating-point precision errors)!
// Use github.com/shopspring/decimal for handling monetary values.
type WebHook struct {
	// --- RAW ---
	RawRequest *http.Request `json:"-"`

	// --- PAYLOAD ---
	Base
	Data *struct {
		Type  string `json:"type,omitempty"`
		Order *struct {
			OutTradeNo    string `json:"out_trade_no,omitempty"`
			PlanTitle     string `json:"plan_title,omitempty"`
			UserPrivateId string `json:"user_private_id,omitempty"`
			UserId        string `json:"user_id,omitempty"`
			PlanId        string `json:"plan_id,omitempty"`
			Title         string `json:"title,omitempty"`
			Month         int    `json:"month,omitempty"`
			TotalAmount   string `json:"total_amount,omitempty"`
			ShowAmount    string `json:"show_amount,omitempty"`
			Status        int    `json:"status,omitempty"`
			Remark        string `json:"remark,omitempty"`
			RedeemId      string `json:"redeem_id,omitempty"`
			ProductType   int    `json:"product_type,omitempty"`
			Discount      string `json:"discount,omitempty"`
			SkuDetail     []struct {
				SkuId   string `json:"sku_id,omitempty"`
				Count   int    `json:"count,omitempty"`
				Name    string `json:"name,omitempty"`
				AlbumId string `json:"album_id,omitempty"`
				Pic     string `json:"pic,omitempty"`
				Stock   string `json:"stock,omitempty"`
				PostId  string `json:"post_id,omitempty"`
			} `json:"sku_detail,omitempty"`
			AddressPerson  string `json:"address_person,omitempty"`
			AddressPhone   string `json:"address_phone,omitempty"`
			AddressAddress string `json:"address_address,omitempty"`
		} `json:"order,omitempty"`
		Sign string `json:"sign,omitempty"`
	} `json:"data,omitempty"`
}
