package appcenter

const (
	baseURL = `https://api.appcenter.ms`
)

// Error ...
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) String() string {
	return e.Code + ": " + e.Message
}
