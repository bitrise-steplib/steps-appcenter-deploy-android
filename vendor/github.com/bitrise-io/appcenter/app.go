package appcenter

import (
	"fmt"

	"github.com/bitrise-io/appcenter/client"
	"github.com/bitrise-io/appcenter/model"
)

// AppAPI ...
type AppAPI struct {
	API            client.API
	ReleaseOptions model.ReleaseOptions
}

// CreateApplicationAPI ...
func CreateApplicationAPI(api client.API, releaseOptions model.ReleaseOptions) AppAPI {
	return AppAPI{
		API:            api,
		ReleaseOptions: releaseOptions,
	}
}

// NewRelease ...
func (a AppAPI) NewRelease() (model.Release, error) {
	releaseID, err := a.API.CreateRelease(a.ReleaseOptions)
	if err != nil {
		return model.Release{},
			fmt.Errorf("failed to create new release on app: %s, owner: %s, %v",
				a.ReleaseOptions.App.AppName,
				a.ReleaseOptions.App.Owner,
				err)
	}

	return a.API.GetAppReleaseDetails(a.ReleaseOptions.App, releaseID)
}

// Groups ...
func (a AppAPI) Groups(name string) (model.Group, error) {
	return a.API.GetGroupByName(name, a.ReleaseOptions.App)
}

// Stores ...
func (a AppAPI) Stores(name string) (model.Store, error) {
	return a.API.GetStore(name, a.ReleaseOptions.App)
}
