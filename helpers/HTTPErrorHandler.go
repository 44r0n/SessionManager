package helpers

// HTTPErrorHandler type to handle http errors
type HTTPErrorHandler struct {
  Status int `json:"status"`
  Error string `json:"error"`
  Description string `json:"description"`
  Fields map[string]string `json:"fields"`
}
