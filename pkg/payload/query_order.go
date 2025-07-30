package payload

import (
	"net/http"
)

// QueryOrder payload
//
// WARNING: Do NOT convert amounts directly to float when performing operations on money (due to floating-point precision errors)!
// Use github.com/shopspring/decimal for handling monetary values.
type QueryOrder struct {
	// --- PAYLOAD ---
	APIBase
	Data *struct {
		APIDataBase
		List []struct {
			OutTradeNo    string `json:"out_trade_no,omitempty"` // Order number
			CustomOrderID string `json:"custom_order_id"`
			UserID        string `json:"user_id,omitempty"`      // User ID of the buyer
			PlanID        string `json:"plan_id,omitempty"`      // Plan ID, empty if it's a custom amount
			Title         string `json:"title,omitempty"`        // Order title
			Month         int    `json:"month,omitempty"`        // Number of months sponsored
			TotalAmount   string `json:"total_amount,omitempty"` // Actual payment amount. 0.00 if a redeem code is used.
			ShowAmount    string `json:"show_amount,omitempty"`  // Display amount, pre-discount
			Status        int    `json:"status,omitempty"`       // 2 for successful transaction. Only this type is pushed currently.
			Remark        string `json:"remark,omitempty"`       // Order remark
			RedeemID      string `json:"redeem_id,omitempty"`
			ProductType   int    `json:"product_type,omitempty"` // 0 for regular plan, 1 for product for sale
			Discount      string `json:"discount,omitempty"`
			SkuDetail     []struct {
				SkuID   string `json:"sku_id,omitempty"`
				Count   int    `json:"count,omitempty"` // Quantity
				Name    string `json:"name,omitempty"`
				AlbumID string `json:"album_id,omitempty"`
				Pic     string `json:"pic,omitempty"` // Picture URL
			} `json:"sku_detail,omitempty"` // Details of SKUs if it's a product for sale
			CreateTime     int    `json:"create_time,omitempty"`      // Creation time of the order
			UserName       string `json:"user_name,omitempty"`        // Username of the buyer
			PlanTitle      string `json:"plan_title,omitempty"`       // Title of the plan
			UserPrivateID  string `json:"user_private_idm,omitempty"` // Unique ID for each user
			AddressPerson  string `json:"address_person,omitempty"`   // Recipient's name
			AddressPhone   string `json:"address_phone,omitempty"`    // Recipient's phone number
			AddressAddress string `json:"address_address,omitempty"`  // Recipient's address
		} `json:"list,omitempty"`
		TotalCount int `json:"total_count,omitempty"` // Total number of orders
		TotalPage  int `json:"total_page,omitempty"`  // Total number of pages, default 50 per page.
	} `json:"data,omitempty"`
}

func (o *QueryOrder) SetRawResponse(resp *http.Response) {
	o.RawResponse = resp
}
