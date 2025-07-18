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
	RawRequest *http.Request `json:"-"` // Raw HTTP request

	// --- PAYLOAD ---
	Base
	Data *struct {
		// --- ERROR ---
		Explain string `json:"explain,omitempty"` // Error description
		Debug   *struct {
			KvString string `json:"kv_string,omitempty"` // For debugging signature errors
		} `json:"debug,omitempty"`
		Request *struct {
			UserID string `json:"user_id,omitempty"`
			Params string `json:"params,omitempty"` // Request parameters
			Ts     int    `json:"ts,omitempty"`     // Timestamp
			Sign   string `json:"sign,omitempty"`   // Signature
		} `json:"request,omitempty"`

		// --- NORMAL ---
		Type  string `json:"type,omitempty"` // Webhook type, e.g., "order"
		Order *struct {
			OutTradeNo    string `json:"out_trade_no,omitempty"` // Order number
			PlanTitle     string `json:"plan_title,omitempty"`
			UserPrivateID string `json:"user_private_id,omitempty"`
			UserID        string `json:"user_id,omitempty"`
			PlanID        string `json:"plan_id,omitempty"`
			Title         string `json:"title,omitempty"`        // Order title
			Month         int    `json:"month,omitempty"`        // Sponsored months
			TotalAmount   string `json:"total_amount,omitempty"` // Total amount paid
			ShowAmount    string `json:"show_amount,omitempty"`  // Display amount
			Status        int    `json:"status,omitempty"`       // Order status (2 for success)
			Remark        string `json:"remark,omitempty"`
			RedeemID      string `json:"redeem_id,omitempty"`    // Redeem code ID
			ProductType   int    `json:"product_type,omitempty"` // Product type (0 for plan, 1 for product)
			Discount      string `json:"discount,omitempty"`
			SkuDetail     []struct {
				SkuID   string `json:"sku_id,omitempty"`
				Count   int    `json:"count,omitempty"` // Quantity
				Name    string `json:"name,omitempty"`
				AlbumID string `json:"album_id,omitempty"`
				Pic     string `json:"pic,omitempty"` // Picture URL
				Stock   string `json:"stock,omitempty"`
				PostID  string `json:"post_id,omitempty"`
			} `json:"sku_detail,omitempty"`
			AddressPerson  string `json:"address_person,omitempty"`  // Recipient's name
			AddressPhone   string `json:"address_phone,omitempty"`   // Recipient's phone number
			AddressAddress string `json:"address_address,omitempty"` // Recipient's address
		} `json:"order,omitempty"`
		Sign string `json:"sign,omitempty"` // Signature
	} `json:"data,omitempty"`
}
