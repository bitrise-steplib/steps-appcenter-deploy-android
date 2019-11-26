package appcenter

// Group ...
type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Origin      string `json:"origin"`
	IsPublic    bool   `json:"is_public"`
	Error       Error  `json:"error"`
}
