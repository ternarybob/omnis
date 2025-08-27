// -----------------------------------------------------------------------
// API Response Model
// Minimal response structure for middleware use only
// -----------------------------------------------------------------------

package omnis

// ApiResponse represents the structured API response format (minimal version)
type ApiResponse struct {
	Version       string                 `json:"version,omitempty"`
	Build         string                 `json:"build,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Support       string                 `json:"support,omitempty"`
	Status        int                    `json:"status"`
	Scope         string                 `json:"scope,omitempty"`
	CorrelationId string                 `json:"correlationid,omitempty"`
	Log           map[string]interface{} `json:"log,omitempty"`
	Result        interface{}            `json:"result"`
	Error         string                 `json:"error,omitempty"`
	Stack         []string               `json:"stack,omitempty"`
	Request       map[string]interface{} `json:"request,omitempty"`
}
