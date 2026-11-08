package payload

// QueryPlan payload
//
// WARNING: Do NOT convert amounts directly to float when performing operations on money (due to floating-point precision errors)!
// Use github.com/shopspring/decimal for handling monetary values.
type QueryPlan struct {
	// --- PAYLOAD ---
	APIBase
	Data *struct {
		APIDataBase
		// --- NORMAL ---
		Plan *struct {
			PlanID              string `json:"plan_id,omitempty"`
			Price               string `json:"price,omitempty"`
			Name                string `json:"name,omitempty"`
			ProductType         int    `json:"product_type,omitempty"` // 0 for subscription, 1 for product, 2 for bundle, 3 for custom, 4 for ticket
			Desc                string `json:"desc,omitempty"`
			ReplyContent        string `json:"reply_content,omitempty"`
			ReplayRandomContent string `json:"replay_random_content,omitempty"` // Field name follows the official doc ("replay")
			Independent         int    `json:"independent,omitempty"`           // 0 non-independent, 1 independent
			Permanent           int    `json:"permanent,omitempty"`             // 0 non-permanent, 1 permanent
			PayMonth            int    `json:"pay_month,omitempty"`             // 1 monthly, 3 quarterly, 12 yearly
			Skus                []struct {
				SkuID              string `json:"sku_id,omitempty"`
				PlanID             string `json:"plan_id,omitempty"`
				Name               string `json:"name,omitempty"`
				Desc               string `json:"desc,omitempty"`
				Stock              string `json:"stock,omitempty"`
				Price              string `json:"price,omitempty"`
				ReplyContent       string `json:"reply_content,omitempty"`
				ReplyRandomContent string `json:"reply_random_content,omitempty"`
			} `json:"skus,omitempty"` // Absent for subscription plans
		} `json:"plan,omitempty"`
	} `json:"data,omitempty"`
}
