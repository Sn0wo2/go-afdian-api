package payload

import (
	"net/http"
)

// Ping payload
type Ping struct {
	// --- INJECTED RAW ---
	RawResponse *http.Response `json:"-"` // Raw HTTP response

	// --- PAYLOAD ---
	APIBase
	Data *struct {
		APIDataBase
		UID  string `json:"uid,omitempty"`  // User ID
		Page any    `json:"page,omitempty"` // Page number
	} `json:"data,omitempty"`
}
