// -----------------------------------------------------------------------
// Service Configuration Model
// Minimal configuration structure for middleware use only
// -----------------------------------------------------------------------

package omnis

// ServiceConfig defines service metadata for middleware (minimal version)
type ServiceConfig struct {
	Version string // Service version (e.g., "1.0.0")
	Build   string // Build timestamp (e.g., "2025-08-27-15-30")
	Name    string // Service name (e.g., "my-api")
	Scope   string // Environment scope ("DEV", "PRD", etc.)
}
