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
	RawResponse *http.Response `json:"-"`

	// --- PAYLOAD ---
	APIBase
	Data *struct {
		TotalCount int `json:"total_count,omitempty"`
		TotalPage  int `json:"total_page,omitempty"`
		List       []struct {
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
			SponsorPlans []struct {
				PlanId         string `json:"plan_id,omitempty"`
				Rank           int    `json:"rank,omitempty"`
				UserId         string `json:"user_id,omitempty"`
				Status         int    `json:"status,omitempty"`
				Name           string `json:"name,omitempty"`
				Pic            string `json:"pic,omitempty"`
				Desc           string `json:"desc,omitempty"`
				Price          string `json:"price,omitempty"`
				UpdateTime     int    `json:"update_time,omitempty"`
				PayMonth       int    `json:"pay_month,omitempty"`
				ShowPrice      string `json:"show_price,omitempty"`
				Independent    int    `json:"independent,omitempty"`
				Permanent      int    `json:"permanent,omitempty"`
				CanBuyHide     int    `json:"can_buy_hide,omitempty"`
				NeedAddress    int    `json:"need_address,omitempty"`
				ProductType    int    `json:"product_type,omitempty"`
				SaleLimitCount int    `json:"sale_limit_count,omitempty"`
				NeedInviteCode bool   `json:"need_invite_code,omitempty"`
				ExpireTime     int    `json:"expire_time,omitempty"`
				SkuProcessed   []any  `json:"sku_processed,omitempty"`
				RankType       int    `json:"rankType,omitempty"`
			} `json:"sponsor_plans,omitempty"`
			CurrentPlan *struct {
				Name           string        `json:"name,omitempty"`
				PlanId         string        `json:"plan_id,omitempty"`
				Rank           int           `json:"rank,omitempty"`
				UserId         string        `json:"user_id,omitempty"`
				Status         int           `json:"status,omitempty"`
				Pic            string        `json:"pic,omitempty"`
				Desc           string        `json:"desc,omitempty"`
				Price          string        `json:"price,omitempty"`
				UpdateTime     int           `json:"update_time,omitempty"`
				PayMonth       int           `json:"pay_month,omitempty"`
				ShowPrice      string        `json:"show_price,omitempty"`
				Independent    int           `json:"independent,omitempty"`
				Permanent      int           `json:"permanent,omitempty"`
				CanBuyHide     int           `json:"can_buy_hide,omitempty"`
				NeedAddress    int           `json:"need_address,omitempty"`
				ProductType    int           `json:"product_type,omitempty"`
				SaleLimitCount int           `json:"sale_limit_count,omitempty"`
				NeedInviteCode bool          `json:"need_invite_code,omitempty"`
				ExpireTime     int           `json:"expire_time,omitempty"`
				SkuProcessed   []interface{} `json:"sku_processed,omitempty"`
				RankType       int           `json:"rankType,omitempty"`
			} `json:"current_plan,omitempty"`
			AllSumAmount string `json:"all_sum_amount,omitempty"`
			CreateTime   int    `json:"create_time,omitempty"`
			LastPayTime  int    `json:"last_pay_time,omitempty"`
			User         *struct {
				UserId string `json:"user_id,omitempty"`
				Name   string `json:"name,omitempty"`
				Avatar string `json:"avatar,omitempty"`
			} `json:"user,omitempty"`
			FirstPayTime int `json:"first_pay_time,omitempty"`
		} `json:"list,omitempty"`
	} `json:"data,omitempty"`
}
