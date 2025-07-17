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
		// --- ERROR ---
		Explain string `json:"explain,omitempty"`
		Debug   *struct {
			KvString string `json:"kv_string,omitempty"`
		} `json:"debug,omitempty"`
		Request *struct {
			UserID string `json:"user_id,omitempty"`
			Params string `json:"params,omitempty"`
			Ts     int    `json:"ts,omitempty"`
			Sign   string `json:"sign,omitempty"`
		} `json:"request,omitempty"`

		// --- NORMAL ---
		Type  string `json:"type,omitempty"`
		Order *struct {
			OutTradeNo    string `json:"out_trade_no,omitempty"`
			PlanTitle     string `json:"plan_title,omitempty"`
			UserPrivateID string `json:"user_private_id,omitempty"`
			UserID        string `json:"user_id,omitempty"`
			PlanID        string `json:"plan_id,omitempty"`
			Title         string `json:"title,omitempty"`
			Month         int    `json:"month,omitempty"`
			TotalAmount   string `json:"total_amount,omitempty"`
			ShowAmount    string `json:"show_amount,omitempty"`
			Status        int    `json:"status,omitempty"`
			Remark        string `json:"remark,omitempty"`
			RedeemID      string `json:"redeem_id,omitempty"`
			ProductType   int    `json:"product_type,omitempty"`
			Discount      string `json:"discount,omitempty"`
			SkuDetail     []struct {
				SkuID   string `json:"sku_id,omitempty"`
				Count   int    `json:"count,omitempty"`
				Name    string `json:"name,omitempty"`
				AlbumID string `json:"album_id,omitempty"`
				Pic     string `json:"pic,omitempty"`
				Stock   string `json:"stock,omitempty"`
				PostID  string `json:"post_id,omitempty"`
			} `json:"sku_detail,omitempty"`
			AddressPerson  string `json:"address_person,omitempty"`
			AddressPhone   string `json:"address_phone,omitempty"`
			AddressAddress string `json:"address_address,omitempty"`
		} `json:"order,omitempty"`
		Sign string `json:"sign,omitempty"`
	} `json:"data,omitempty"`
}
