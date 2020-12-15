package client

import (
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/bitrise-io/appcenter/util"

	"github.com/bitrise-io/appcenter/model"
)

type fileAssetResponse struct {
	ReleaseID       string `json:"id"`
	PackageAssetID  string `json:"package_asset_id"`
	Token           string `json:"token"`
	UploadDomain    string `json:"upload_domain"`
	URLEncodedToken string `json:"url_encoded_token"`
}

// API ...
type API struct {
	Client Client
}

// CreateAPIWithClientParams ...
func CreateAPIWithClientParams(token string, debug bool) API {
	return API{
		Client: NewClient(token, debug),
	}
}

// GetAppReleaseDetails ...
func (api API) GetAppReleaseDetails(app model.App, releaseID int) (model.Release, error) {
	//fetch releases and find the latest
	var (
		releaseShowURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%s", baseURL, app.Owner, app.AppName, strconv.Itoa(releaseID))
		release        model.Release
	)

	statusCode, err := api.Client.jsonRequest(http.MethodGet, releaseShowURL, nil, &release)
	if err != nil {
		return model.Release{}, err
	}

	if statusCode != http.StatusOK {
		return model.Release{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, releaseShowURL, release)
	}

	return release, err
}

// GetGroupByName ...
func (api API) GetGroupByName(groupName string, app model.App) (model.Group, error) {
	var (
		getURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/distribution_groups/%s", baseURL, app.Owner, app.AppName, groupName)
		getResponse model.Group
	)

	statusCode, err := api.Client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
	if err != nil {
		return model.Group{}, err
	}

	if statusCode != http.StatusOK {
		return model.Group{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
	}

	return getResponse, err
}

// GetStore ...
func (api API) GetStore(storeName string, app model.App) (model.Store, error) {
	var (
		getURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/distribution_stores/%s", baseURL, app.Owner, app.AppName, storeName)
		getResponse model.Store
	)

	statusCode, err := api.Client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
	if err != nil {
		return model.Store{}, err
	}

	if statusCode != http.StatusOK {
		return model.Store{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
	}

	return getResponse, nil
}

// AddReleaseToGroup ...
func (api API) AddReleaseToGroup(g model.Group, releaseID int, opts model.ReleaseOptions) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/groups", baseURL, opts.App.Owner, opts.App.AppName, releaseID)
		postRequest = struct {
			ID              string `json:"id"`
			MandatoryUpdate bool   `json:"mandatory_update"`
			NotifyTesters   bool   `json:"notify_testers"`
		}{
			ID:              g.ID,
			MandatoryUpdate: opts.Mandatory,
			NotifyTesters:   opts.NotifyTesters,
		}
	)

	body, err := api.Client.MarshallContent(postRequest)
	if err != nil {
		return err
	}

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, body, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// AddReleaseToStore ...
func (api API) AddReleaseToStore(s model.Store, releaseID int, opts model.ReleaseOptions) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/stores", baseURL, opts.App.Owner, opts.App.AppName, releaseID)
		postRequest = struct {
			ID string `json:"id"`
		}{
			ID: s.ID,
		}
	)

	body, err := api.Client.MarshallContent(postRequest)
	if err != nil {
		return err
	}

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, body, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// AddTesterToRelease ...
func (api API) AddTesterToRelease(email string, releaseID int, opts model.ReleaseOptions) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/testers", baseURL, opts.App.Owner, opts.App.AppName, releaseID)
		postRequest = struct {
			Email           string `json:"email"`
			MandatoryUpdate bool   `json:"mandatory_update"`
			NotifyTesters   bool   `json:"notify_testers"`
		}{
			Email:           email,
			MandatoryUpdate: opts.Mandatory,
			NotifyTesters:   opts.NotifyTesters,
		}
	)

	body, err := api.Client.MarshallContent(postRequest)
	if err != nil {
		return err
	}

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, body, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// SetReleaseNoteOnRelease ...
func (api API) SetReleaseNoteOnRelease(releaseNote string, releaseID int, opts model.ReleaseOptions) error {
	var (
		putURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d", baseURL, opts.App.Owner, opts.App.AppName, releaseID)
		putRequest = struct {
			ReleaseNotes string `json:"release_notes,omitempty"`
		}{
			ReleaseNotes: releaseNote,
		}
	)

	body, err := api.Client.MarshallContent(putRequest)
	if err != nil {
		return err
	}

	statusCode, err := api.Client.jsonRequest(http.MethodPut, putURL, body, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, putURL)
	}

	return nil
}

// UploadSymbolToRelease - build and version is required for Android and optional for iOS
func (api API) UploadSymbolToRelease(filePath string, release model.Release, opts model.ReleaseOptions) error {
	var symbolType = model.SymbolTypeDSYM
	if release.AppOs == "Android" {
		symbolType = model.SymbolTypeMapping
	}

	// send file upload request
	var (
		postURL  = fmt.Sprintf("%s/v0.1/apps/%s/%s/symbol_uploads", baseURL, opts.App.Owner, opts.App.AppName)
		postBody = struct {
			SymbolType model.SymbolType `json:"symbol_type"`
			FileName   string           `json:"file_name,omitempty"`
			Build      string           `json:"build,omitempty"`
			Version    string           `json:"version,omitempty"`
		}{
			FileName:   filepath.Base(filePath),
			Build:      release.Version,
			Version:    release.ShortVersion,
			SymbolType: symbolType,
		}
		postResponse struct {
			SymbolUploadID string    `json:"symbol_upload_id"`
			UploadURL      string    `json:"upload_url"`
			ExpirationDate time.Time `json:"expiration_date"`
		}
	)

	body, err := api.Client.MarshallContent(postBody)
	if err != nil {
		return err
	}

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, body, &postResponse)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, postURL, postBody)
	}

	// upload file to {upload_url}
	statusCode, err = api.Client.uploadFile(postResponse.UploadURL, filePath)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postResponse.UploadURL)
	}

	var (
		patchURL  = fmt.Sprintf("%s/v0.1/apps/%s/%s/symbol_uploads/%s", baseURL, opts.App.Owner, opts.App.AppName, postResponse.SymbolUploadID)
		patchBody = map[string]string{
			"status": "committed",
		}
	)

	body, err = api.Client.MarshallContent(patchBody)
	if err != nil {
		return err
	}

	statusCode, err = api.Client.jsonRequest(http.MethodPatch, patchURL, body, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, patchURL)
	}

	return nil
}

// CreateRelease ...
func (api API) CreateRelease(opts model.ReleaseOptions) (int, error) {
	var (
		assetsURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/uploads/releases",
			baseURL,
			opts.App.Owner,
			opts.App.AppName)
		assetResponse fileAssetResponse
	)

	statusCode, err := api.Client.jsonRequest(http.MethodPost, assetsURL, nil, &assetResponse)
	if err != nil {
		return -1, err
	}

	if statusCode != http.StatusCreated {
		return -1, fmt.Errorf("invalid status code: %d, url: %s", statusCode, assetsURL)
	}

	fmt.Println("")
	fmt.Println(fmt.Sprintf("Asset release ID: %s", assetResponse.ReleaseID))
	fmt.Println(fmt.Sprintf("File asset ID: %s", assetResponse.PackageAssetID))

	file := util.LocalFile{FilePath: opts.FilePath}
	err = file.OpenFile()
	if err != nil {
		return -1, err
	}

	fileName := file.FileName()
	fileSize := file.FileSize()

	fmt.Println("")
	fmt.Println("Uploading file with metadata:")
	fmt.Println(fmt.Sprintf("- File name: %s", fileName))
	fmt.Println(fmt.Sprintf("- File size: %s", strconv.Itoa(fileSize)))

	var (
		metadataURL = fmt.Sprintf("%s/upload/set_metadata/%s?file_name=%s&file_size=%s&token=%s",
			assetResponse.UploadDomain,
			assetResponse.PackageAssetID,
			fileName,
			strconv.Itoa(fileSize),
			assetResponse.URLEncodedToken)
		metadataResponse struct {
			ID             string `json:"id"`
			ChunkSize      int    `json:"chunk_size"`
			ChunkList      []int  `json:"chunk_list"`
			BlobPartitions int    `json:"blob_partitions"`
		}
	)

	statusCode, err = api.Client.jsonRequest(http.MethodPost, metadataURL, nil, &metadataResponse)
	if err != nil {
		return -1, err
	}

	if statusCode != http.StatusOK {
		return -1, fmt.Errorf("invalid status code: %d, url: %s", statusCode, metadataURL)
	}

	fmt.Println("")
	fmt.Println("Upload information:")
	fmt.Println(fmt.Sprintf("Chunk size: %d bytes", metadataResponse.ChunkSize))
	fmt.Println(fmt.Sprintf("Chunk number: %d", len(metadataResponse.ChunkList)))

	fmt.Println("")
	fmt.Println("Uploading chunks ...")

	fileChunks := file.MakeChunks(metadataResponse.ChunkSize)

	err = api.uploadChunksParallelly(fileChunks, metadataResponse.ChunkList, assetResponse)
	if err != nil {
		return -1, err
	}

	fmt.Println("")
	fmt.Println("Chunk upload finished...")

	var (
		uploadFinishedURL = fmt.Sprintf("%s/upload/finished/%s?token=%s",
			assetResponse.UploadDomain,
			assetResponse.PackageAssetID,
			assetResponse.URLEncodedToken)
		finishedResponse interface{}
	)

	statusCode, err = api.Client.jsonRequest(http.MethodPost, uploadFinishedURL, nil, &finishedResponse)
	if err != nil {
		return -1, err
	}

	if statusCode != http.StatusOK {
		return -1, fmt.Errorf("invalid status code: %d, url: %s", statusCode, uploadFinishedURL)
	}

	fmt.Println("")
	fmt.Println("Upload finished...")

	//patch release

	var (
		releasePatchURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/uploads/releases/%s",
			baseURL,
			opts.App.Owner,
			opts.App.AppName,
			assetResponse.ReleaseID)
		releaseBody = struct {
			UploadStatus string `json:"upload_status"`
		}{
			UploadStatus: "uploadFinished",
		}
		releasePatchResponse interface{}
	)

	body, err := api.Client.MarshallContent(releaseBody)
	if err != nil {
		return -1, err
	}

	statusCode, err = api.Client.jsonRequest(http.MethodPatch, releasePatchURL, body, &releasePatchResponse)
	if err != nil {
		return -1, err
	}

	if statusCode != http.StatusOK {
		return -1, fmt.Errorf("invalid status code: %d, url: %s", statusCode, releasePatchURL)
	}

	fmt.Println("")
	fmt.Println("Release patched...")

	fmt.Println("")
	fmt.Println("Waiting for the AppCenter release to getting ready...")

	uploadStatus := "commited"
	releaseDistinctID := -1
	attempts := 1
	for ok := true; ok; ok = !uploadIsReadyForDeploy(uploadStatus) {
		fmt.Println(fmt.Sprintf("Attempt(s): %d", attempts))

		var (
			getURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/uploads/releases/%s",
				baseURL,
				opts.App.Owner,
				opts.App.AppName,
				assetResponse.ReleaseID)
			getResponse struct {
				ID                string `json:"id"`
				ReleaseDistinctID int    `json:"release_distinct_id,omitempty"`
				UploadStatus      string `json:"upload_status"`
			}
		)

		statusCode, err = api.Client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
		if err != nil {
			return -1, err
		}

		if statusCode != http.StatusOK {
			return -1, fmt.Errorf("invalid status code: %d, url: %s", statusCode, getURL)
		}

		uploadStatus = getResponse.UploadStatus

		if uploadIsReadyForDeploy(uploadStatus) {
			releaseDistinctID = getResponse.ReleaseDistinctID
		} else {
			attempts++

			sleepDuration := generateRandomIntBetweenRange(5, 10)
			fmt.Println(fmt.Sprintf("Waiting for %d second(s), current status: %s", sleepDuration, uploadStatus))

			time.Sleep(time.Duration(sleepDuration) * time.Second)
		}
	}

	fmt.Println("")
	fmt.Println(fmt.Sprintf("Release created with ID: %d", releaseDistinctID))

	return releaseDistinctID, nil
}

func (api API) uploadChunksParallelly(fileChunks [][]byte, chunkIDs []int, assetResponse fileAssetResponse) (retErr error) {
	var wg sync.WaitGroup
	wg.Add(len(fileChunks))

	for idx, chunkID := range chunkIDs {
		chunk := fileChunks[idx]

		go func(chunk []byte, ID int) {
			defer wg.Done()

			fmt.Println(fmt.Sprintf("Uploading chunk with ID: %d, size: %d", ID, len(chunk)))

			var (
				chunkUploadURL = fmt.Sprintf("%s/upload/upload_chunk/%s?block_number=%s&token=%s",
					assetResponse.UploadDomain,
					assetResponse.PackageAssetID,
					strconv.Itoa(ID),
					assetResponse.URLEncodedToken)
				chunkUploadResponse struct {
					Error     bool   `json:"error"`
					ErrorCode string `json:"error_code"`
				}
			)

			statusCode, err := api.Client.jsonRequest(http.MethodPost, chunkUploadURL, chunk, &chunkUploadResponse)
			if err != nil {
				retErr = err

				return
			}
			if chunkUploadResponse.Error {
				retErr = fmt.Errorf("failed to upload chunk, chunk id: %d, error code: %s",
					ID,
					chunkUploadResponse.ErrorCode)
				return
			}

			if statusCode != http.StatusOK {
				retErr = fmt.Errorf("invalid status code: %d, url: %s", statusCode, chunkUploadURL)
				return
			}

			fmt.Println(fmt.Sprintf("Uploading finished, ID: %d", ID))
		}(chunk, chunkID)
	}

	wg.Wait()

	return
}

func generateRandomIntBetweenRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func uploadIsReadyForDeploy(status string) bool {
	return status == "readyToBePublished"
}
