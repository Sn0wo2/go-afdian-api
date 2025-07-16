package payload

import (
	"net/http"
)

type APIBase struct {
	// --- INJECTED RAW ---
	RawResponse *http.Response `json:"-"`

	// --- PAYLOAD ---
	Base
}
