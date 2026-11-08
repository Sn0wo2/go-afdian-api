package payload

// UpdatePlanReply payload for the /api/open/update-plan-reply endpoint.
type UpdatePlanReply struct {
	// --- PAYLOAD ---
	APIBase
	Data *struct {
		APIDataBase
	} `json:"data,omitempty"`
}
