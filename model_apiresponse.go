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
	Result        interface{}            `json:"result,omitempty"`
	Error         string                 `json:"error,omitempty"`
	Stack         []string               `json:"stack,omitempty"`
	Request       map[string]interface{} `json:"request,omitempty"`
	Log           map[string]string      `json:"log,omitempty"`
}
