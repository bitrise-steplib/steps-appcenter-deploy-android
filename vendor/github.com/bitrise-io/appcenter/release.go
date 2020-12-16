package appcenter

import (
	"strings"

	"github.com/bitrise-io/appcenter/client"
	"github.com/bitrise-io/appcenter/model"
)

// ReleaseAPI ...
type ReleaseAPI struct {
	API            client.API
	Release        model.Release
	ReleaseOptions model.ReleaseOptions
}

// CreateReleaseAPI ...
func CreateReleaseAPI(api client.API, release model.Release, releaseOptions model.ReleaseOptions) ReleaseAPI {
	return ReleaseAPI{
		API:            api,
		Release:        release,
		ReleaseOptions: releaseOptions,
	}
}

// AddGroup ...
func (r ReleaseAPI) AddGroup(g model.Group) error {
	return r.API.AddReleaseToGroup(g, r.Release.ID, r.ReleaseOptions)
}

// AddGroupsToRelease ...
func (r ReleaseAPI) AddGroupsToRelease(groupNames []string) error {
	if len(groupNames) > 0 {
		for _, groupName := range groupNames {
			if len(strings.TrimSpace(groupName)) == 0 {
				continue
			}
			group, err := r.API.GetGroupByName(groupName, r.ReleaseOptions.App)
			if err != nil {
				return err
			}

			err = r.AddGroup(group)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AddStore ...
func (r ReleaseAPI) AddStore(s model.Store) error {
	return r.API.AddReleaseToStore(s, r.Release.ID, r.ReleaseOptions)
}

// AddTester ...
func (r ReleaseAPI) AddTester(email string) error {
	return r.API.AddTesterToRelease(email, r.Release.ID, r.ReleaseOptions)
}

// SetReleaseNote ...
func (r ReleaseAPI) SetReleaseNote(releaseNote string) error {
	return r.API.SetReleaseNoteOnRelease(releaseNote, r.Release.ID, r.ReleaseOptions)
}

// UploadSymbol - build and version is required for Android and optional for iOS
func (r ReleaseAPI) UploadSymbol(filePath string) error {
	return r.API.UploadSymbolToRelease(filePath, r.Release, r.ReleaseOptions)
}
