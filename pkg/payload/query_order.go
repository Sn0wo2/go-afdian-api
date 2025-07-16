package payload

import (
	"net/http"
)

// QueryOrder payload
//
// WARNING: Do NOT convert amounts directly to float when performing operations on money (due to floating-point precision errors)!
// Use github.com/shopspring/decimal for handling monetary values.
type QueryOrder struct {
	// --- INJECTED RAW ---
	RawResponse *http.Response `json:"-"`

	// --- PAYLOAD ---
	APIBase
	Data *struct {
		// --- ERROR ---
		Explain string `json:"explain,omitempty"`
		Debug   *struct {
			KvString string `json:"kv_string,omitempty"`
		} `json:"debug,omitempty"`
		Request *struct {
			UserId string `json:"user_id,omitempty"`
			Params string `json:"params,omitempty"`
			Ts     int    `json:"ts,omitempty"`
			Sign   string `json:"sign,omitempty"`
		} `json:"request,omitempty"`

		// --- NORMAL ---
		List []struct {
			OutTradeNo    string `json:"out_trade_no,omitempty"`
			CustomOrderId string `json:"custom_order_id,omitempty"`
			UserId        string `json:"user_id,omitempty"`
			UserPrivateId string `json:"user_private_id,omitempty"`
			PlanId        string `json:"plan_id,omitempty"`
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
			} `json:"sku_detail,omitempty"`
			AddressPerson  string `json:"address_person,omitempty"`
			AddressPhone   string `json:"address_phone,omitempty"`
			AddressAddress string `json:"address_address,omitempty"`
		} `json:"list,omitempty"`
		TotalCount int `json:"total_count,omitempty"`
		TotalPage  int `json:"total_page,omitempty"`
	} `json:"data,omitempty"`
}
