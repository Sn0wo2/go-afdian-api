package payload

import (
	"net/http"
)

// QuerySponsor payload
//
// WARNING: Do NOT convert amounts directly to float when performing operations on money (due to floating-point precision errors)!
// Use github.com/shopspring/decimal for handling monetary values.
type QuerySponsor struct {
	// --- INJECTED RAW ---
	RawResponse *http.Response `json:"-"` // Raw HTTP response

	// --- PAYLOAD ---
	APIBase
	Data struct {
		APIDataBase
		TotalCount int `json:"total_count,omitempty"` // Total number of sponsors
		TotalPage  int `json:"total_page,omitempty"`  // Total number of pages, default 20 per page
		List       []struct {
			SponsorPlans []struct {
				PlanID         string `json:"plan_id,omitempty"`
				Rank           int    `json:"rank,omitempty"`
				UserID         string `json:"user_id,omitempty"`
				Status         int    `json:"status,omitempty"`
				Name           string `json:"name,omitempty"` // Plan name
				Pic            string `json:"pic,omitempty"`  // Picture URL
				Desc           string `json:"desc,omitempty"` // Description
				Price          string `json:"price,omitempty"`
				UpdateTime     int    `json:"update_time,omitempty"`  // Last update time
				PayMonth       int    `json:"pay_month,omitempty"`    // Sponsored months
				ShowPrice      string `json:"show_price,omitempty"`   // Display price
				Independent    int    `json:"independent,omitempty"`  // Whether it's an independent plan
				Permanent      int    `json:"permanent,omitempty"`    // Whether it's a permanent plan
				CanBuyHide     int    `json:"can_buy_hide,omitempty"` // Whether a hidden plan can be purchased
				NeedAddress    int    `json:"need_address,omitempty"` // Whether a shipping address is required
				ProductType    int    `json:"product_type,omitempty"`
				SaleLimitCount int    `json:"sale_limit_count,omitempty"`
				NeedInviteCode bool   `json:"need_invite_code,omitempty"` // Whether an invitation code is required
				ExpireTime     int    `json:"expire_time,omitempty"`
				// Unknown
				SkuProcessed []any `json:"sku_processed,omitempty"`
				RankType     int   `json:"rankType,omitempty"`
			} `json:"sponsor_plans,omitempty"`
			CurrentPlan *struct {
				Name           string `json:"name,omitempty"` // Plan name
				PlanID         string `json:"plan_id,omitempty"`
				Rank           int    `json:"rank,omitempty"` // Rank
				UserID         string `json:"user_id,omitempty"`
				Status         int    `json:"status,omitempty"`
				Pic            string `json:"pic,omitempty"`  // Picture URL
				Desc           string `json:"desc,omitempty"` // Description
				Price          string `json:"price,omitempty"`
				UpdateTime     int    `json:"update_time,omitempty"`
				PayMonth       int    `json:"pay_month,omitempty"`    // Sponsored months
				ShowPrice      string `json:"show_price,omitempty"`   // Display price
				Independent    int    `json:"independent,omitempty"`  // Whether it's an independent plan
				Permanent      int    `json:"permanent,omitempty"`    // Whether it's a permanent plan
				CanBuyHide     int    `json:"can_buy_hide,omitempty"` // Whether a hidden plan can be purchased
				NeedAddress    int    `json:"need_address,omitempty"` // Whether a shipping address is required
				ProductType    int    `json:"product_type,omitempty"`
				SaleLimitCount int    `json:"sale_limit_count,omitempty"`
				NeedInviteCode bool   `json:"need_invite_code,omitempty"` // Whether an invitation code is required
				ExpireTime     int    `json:"expire_time,omitempty"`
				// Unknown
				SkuProcessed []any `json:"sku_processed,omitempty"`
				RankType     int   `json:"rankType,omitempty"`
			} `json:"current_plan,omitempty"` // Current sponsored plan. If only name is "", it means no plan.
			AllSumAmount string `json:"all_sum_amount,omitempty"` // Total sponsored amount (pre-discount).
			FirstPayTime int    `json:"first_pay_time,omitempty"` // Timestamp of the first payment
			CreateTime   int    `json:"create_time,omitempty"`    // Timestamp of when the user became a sponsor (first sponsorship)
			LastPayTime  int    `json:"last_pay_time,omitempty"`  // Timestamp of the last sponsorship
			User         *struct {
				UserID        string `json:"user_id,omitempty"`         // Unique user ID
				Name          string `json:"name,omitempty"`            // Nickname (not unique)
				Avatar        string `json:"avatar,omitempty"`          // Avatar URL
				UserPrivateID string `json:"user_private_id,omitempty"` // Unique ID for each user
			} `json:"user,omitempty"`
		} `json:"list,omitempty"`
	} `json:"data,omitempty"`
}

func (s *QuerySponsor) SetRawResponse(resp *http.Response) {
	s.RawResponse = resp
}
