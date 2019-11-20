package main

import (
	"os"
	"strings"

	"github.com/bitrise-io/appcenter"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
)

type config struct {
	Debug              bool            `env:"debug,required"`
	ApkPath            string          `env:"apk_path,file"`
	AppName            string          `env:"app_name,required"`
	APIToken           stepconf.Secret `env:"api_token,required"`
	OwnerName          string          `env:"owner_name,required"`
	Mandatory          bool            `env:"mandatory,required"`
	MappingPath        string          `env:"mapping_path"`
	ReleaseNotes       string          `env:"release_notes"`
	NotifyTesters      bool            `env:"notify_testers,required"`
	DistributionGroup  string          `env:"distribution_group"`
	DistributionStore  string          `env:"distribution_store"`
	DistributionTester string          `env:"distribution_tester"`
}

func failf(f string, args ...interface{}) {
	log.Errorf(f, args...)
	os.Exit(1)
}

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	app := appcenter.NewClient(string(cfg.APIToken), cfg.Debug).Apps(cfg.OwnerName, cfg.AppName)

	release, err := app.NewRelease(cfg.ApkPath)
	if err != nil {
		failf("Failed to create new release, error: %s", err)
	}

	if len(cfg.MappingPath) > 0 {
		if err := release.UploadSymbol(cfg.MappingPath); err != nil {
			failf("Failed to upload symbol file(%s), error: %s", cfg.MappingPath, err)
		}
	}

	if len(cfg.ReleaseNotes) > 0 {
		if err := release.SetReleaseNote(cfg.ReleaseNotes); err != nil {
			failf("Failed to set release note, error: %s", err)
		}
	}

	for _, groupName := range strings.Split(cfg.DistributionGroup, "\n") {
		groupName = strings.TrimSpace(groupName)

		if len(groupName) == 0 {
			continue
		}

		group, err := app.Groups(groupName)
		if err != nil {
			failf("Failed to fetch group with name: (%s), error: %s", groupName, err)
		}

		if err := release.AddGroup(group, cfg.Mandatory, cfg.NotifyTesters); err != nil {
			failf("Failed to add group(%s) to the release, error: %s", groupName, err)
		}
	}

	for _, storeName := range strings.Split(cfg.DistributionStore, "\n") {
		storeName = strings.TrimSpace(storeName)

		if len(storeName) == 0 {
			continue
		}

		store, err := app.Stores(storeName)
		if err != nil {
			failf("Failed to fetch store with name: (%s), error: %s", storeName, err)
		}

		if err := release.AddStore(store); err != nil {
			failf("Failed to add store(%s) to the release, error: %s", storeName, err)
		}
	}

	for _, email := range strings.Split(cfg.DistributionTester, "\n") {
		email = strings.TrimSpace(email)

		if len(email) == 0 {
			continue
		}

		if err := release.AddTester(email, cfg.Mandatory, cfg.NotifyTesters); err != nil {
			failf("Failed to add tester(%s) to the release, error: %s", email, err)
		}
	}
}
