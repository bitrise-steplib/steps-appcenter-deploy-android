package appcenter

// Store ...
type Store struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Track         string `json:"track"`
	IntuneDetails struct {
		TargetAudience struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"target_audience"`
		AppCategory struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"app_category"`
	} `json:"intune_details"`
	ServiceConnectionID string `json:"service_connection_id"`
	CreatedBy           string `json:"created_by"`
	Error               Error  `json:"error"`
}
