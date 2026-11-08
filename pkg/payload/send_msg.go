package payload

// SendMsg payload for the /api/open/send-msg endpoint (rate limited to 10/s and 1000/h by afdian).
type SendMsg struct {
	// --- PAYLOAD ---
	APIBase
	Data *struct {
		APIDataBase
	} `json:"data,omitempty"`
}
