package model

// Release ...
type Release struct {
	ID                            int      `json:"id"`
	AppName                       string   `json:"app_name"`
	AppDisplayName                string   `json:"app_display_name"`
	AppOs                         string   `json:"app_os"`
	Version                       string   `json:"version"`
	Origin                        string   `json:"origin"`
	ShortVersion                  string   `json:"short_version"`
	ReleaseNotes                  string   `json:"release_notes"`
	ProvisioningProfileName       string   `json:"provisioning_profile_name"`
	ProvisioningProfileType       string   `json:"provisioning_profile_type"`
	ProvisioningProfileExpiryDate string   `json:"provisioning_profile_expiry_date"`
	IsProvisioningProfileSyncing  bool     `json:"is_provisioning_profile_syncing"`
	Size                          int      `json:"size"`
	MinOs                         string   `json:"min_os"`
	DeviceFamily                  string   `json:"device_family"`
	AndroidMinAPILevel            string   `json:"android_min_api_level"`
	BundleIdentifier              string   `json:"bundle_identifier"`
	PackageHashes                 []string `json:"package_hashes"`
	Fingerprint                   string   `json:"fingerprint"`
	UploadedAt                    string   `json:"uploaded_at"`
	DownloadURL                   string   `json:"download_url"`
	AppIconURL                    string   `json:"app_icon_url"`
	InstallURL                    string   `json:"install_url"`
	DestinationType               string   `json:"destination_type"`
	DistributionGroups            []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"distribution_groups"`
	DistributionStores []struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		Type             string `json:"type"`
		PublishingStatus string `json:"publishing_status"`
	} `json:"distribution_stores"`
	Destinations []struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		IsLatest         bool   `json:"is_latest"`
		Type             string `json:"type"`
		PublishingStatus string `json:"publishing_status"`
		DestinationType  string `json:"destination_type"`
		DisplayName      string `json:"display_name"`
	} `json:"destinations"`
	IsUdidProvisioned bool `json:"is_udid_provisioned"`
	CanResign         bool `json:"can_resign"`
	Build             struct {
		BranchName    string `json:"branch_name"`
		CommitHash    string `json:"commit_hash"`
		CommitMessage string `json:"commit_message"`
	} `json:"build"`
	Enabled         bool   `json:"enabled"`
	Status          string `json:"status"`
	IsExternalBuild bool   `json:"is_external_build"`
	Error           Error  `json:"error"`
}
