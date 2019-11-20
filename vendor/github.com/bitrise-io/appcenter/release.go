package appcenter

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

// ReleaseOptions ...
type ReleaseOptions struct {
	BuildVersion string `json:"build_version,omitempty"`
	BuildNumber  string `json:"build_number,omitempty"`
	ReleaseID    int    `json:"release_id,omitempty"`
}

// Release ...
type Release struct {
	app                           App
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
}

// AddGroup ...
func (r Release) AddGroup(g Group, mandatoryUpdate, notifyTesters bool) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/groups", baseURL, r.app.owner, r.app.name, r.ID)
		postRequest = struct {
			ID              string `json:"id"`
			MandatoryUpdate bool   `json:"mandatory_update"`
			NotifyTesters   bool   `json:"notify_testers"`
		}{
			ID:              g.ID,
			MandatoryUpdate: mandatoryUpdate,
			NotifyTesters:   notifyTesters,
		}
	)

	statusCode, err := r.app.client.jsonRequest(http.MethodPost, postURL, postRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// AddStore ...
func (r Release) AddStore(s Store) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/stores", baseURL, r.app.owner, r.app.name, r.ID)
		postRequest = struct {
			ID string `json:"id"`
		}{
			ID: s.ID,
		}
	)

	statusCode, err := r.app.client.jsonRequest(http.MethodPost, postURL, postRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// AddTester ...
func (r Release) AddTester(email string, mandatoryUpdate, notifyTesters bool) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/testers", baseURL, r.app.owner, r.app.name, r.ID)
		postRequest = struct {
			Email           string `json:"email"`
			MandatoryUpdate bool   `json:"mandatory_update"`
			NotifyTesters   bool   `json:"notify_testers"`
		}{
			Email:           email,
			MandatoryUpdate: mandatoryUpdate,
			NotifyTesters:   notifyTesters,
		}
	)

	statusCode, err := r.app.client.jsonRequest(http.MethodPost, postURL, postRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// SetReleaseNote ...
func (r Release) SetReleaseNote(releaseNote string) error {
	var (
		putURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d", baseURL, r.app.owner, r.app.name, r.ID)
		putRequest = struct {
			ReleaseNotes string `json:"release_notes,omitempty"`
		}{
			ReleaseNotes: releaseNote,
		}
	)

	statusCode, err := r.app.client.jsonRequest(http.MethodPut, putURL, putRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, putURL)
	}

	return nil
}

// UploadSymbol - build and version is required for Android and optional for iOS
func (r Release) UploadSymbol(filePath string) error {
	var symbolType = SymbolTypeDSYM
	if r.AppOs == "Android" {
		symbolType = SymbolTypeMapping
	}

	// send file upload request
	var (
		postURL  = fmt.Sprintf("%s/v0.1/apps/%s/%s/symbol_uploads", baseURL, r.app.owner, r.app.name)
		postBody = struct {
			SymbolType SymbolType `json:"symbol_type"`
			FileName   string     `json:"file_name,omitempty"`
			Build      string     `json:"build,omitempty"`
			Version    string     `json:"version,omitempty"`
		}{
			FileName:   filepath.Base(filePath),
			Build:      r.ShortVersion,
			Version:    r.Version,
			SymbolType: symbolType,
		}
		postResponse struct {
			SymbolUploadID string    `json:"symbol_upload_id"`
			UploadURL      string    `json:"upload_url"`
			ExpirationDate time.Time `json:"expiration_date"`
		}
	)

	statusCode, err := r.app.client.jsonRequest(http.MethodPost, postURL, postBody, &postResponse)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, postURL, postBody)
	}

	// upload file to {upload_url}
	statusCode, err = r.app.client.uploadFile(postResponse.UploadURL, filePath)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postResponse.UploadURL)
	}

	var (
		patchURL  = fmt.Sprintf("%s/v0.1/apps/%s/%s/symbol_uploads/%s", baseURL, r.app.owner, r.app.name, postResponse.SymbolUploadID)
		patchBody = map[string]string{
			"status": "committed",
		}
	)

	statusCode, err = r.app.client.jsonRequest(http.MethodPatch, patchURL, patchBody, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, patchURL)
	}

	return nil
}
