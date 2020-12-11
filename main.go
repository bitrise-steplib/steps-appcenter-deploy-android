package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bitrise-io/appcenter"
	"github.com/bitrise-io/appcenter/client"
	"github.com/bitrise-io/appcenter/model"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-steputils/tools"
	"github.com/bitrise-io/go-utils/log"
)

const statusEnvKey = "APPCENTER_DEPLOY_STATUS"

type config struct {
	Debug              bool            `env:"debug,required"`
	AppPath            string          `env:"app_path,file"`
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

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}
	stepconf.Print(cfg)
	fmt.Println()

	app := model.App{
		Owner:   cfg.OwnerName,
		AppName: cfg.AppName,
	}

	releaseOptions := model.ReleaseOptions{
		GroupNames:    strings.Split(cfg.DistributionGroup, "\n"),
		Mandatory:     cfg.Mandatory,
		NotifyTesters: cfg.NotifyTesters,
		FilePath:      cfg.AppPath,
		App:           app,
	}

	api := client.CreateAPIWithClientParams(string(cfg.APIToken), cfg.Debug)
	appAPI := appcenter.CreateApplicationAPI(api, releaseOptions)

	log.Infof("Uploading binary")

	release, err := appAPI.NewRelease()
	if err != nil {
		failf("Failed to create new release, error: %s", err)
	}

	log.Donef("- Done")
	fmt.Println()

	releaseAPI := appcenter.CreateReleaseAPI(api, release, releaseOptions)

	log.Infof("Setting distribution group(s)")

	err = releaseAPI.AddGroupsToRelease(releaseOptions.GroupNames)
	if err != nil {
		failf("Failed to set groups on the release %s, groups: %s, error: %s", release.ID, releaseOptions.GroupNames, err)
	}

	if len(cfg.MappingPath) > 0 {
		log.Infof("Uploading mapping file")
		if err := releaseAPI.UploadSymbol(cfg.MappingPath); err != nil {
			failf("Failed to upload symbol file(%s), error: %s", cfg.MappingPath, err)
		}
		log.Donef("- Done")
		fmt.Println()
	}

	if len(cfg.ReleaseNotes) > 0 {
		log.Infof("Setting release notes")
		if err := releaseAPI.SetReleaseNote(cfg.ReleaseNotes); err != nil {
			failf("Failed to set release note, error: %s", err)
		}
		log.Donef("- Done")
		fmt.Println()
	}

	log.Infof("Gatehering public group(s)")

	var publicGroup []string
	for _, groupName := range releaseOptions.GroupNames {
		groupName = strings.TrimSpace(groupName)
		if len(groupName) == 0 {
			continue
		}

		log.Printf("- %s", groupName)

		group, err := appAPI.Groups(groupName)
		if err != nil {
			failf("Failed to fetch group with name: (%s), error: %s", groupName, err)
		}

		if group.IsPublic {
			publicGroup = append(publicGroup, groupName)
		}
	}

	log.Donef("- Done")
	fmt.Println()

	log.Infof("Setting distribution store(s)")

	for _, storeName := range strings.Split(cfg.DistributionStore, "\n") {
		storeName = strings.TrimSpace(storeName)

		if len(storeName) == 0 {
			continue
		}

		log.Printf("- %s", storeName)

		store, err := appAPI.Stores(storeName)
		if err != nil {
			failf("Failed to fetch store with name: (%s), error: %s", storeName, err)
		}

		if err := releaseAPI.AddStore(store); err != nil {
			failf("Failed to add store(%s) to the release, error: %s", storeName, err)
		}
	}

	log.Donef("- Done")
	fmt.Println()

	log.Infof("Setting distribution tester(s)")

	for _, email := range strings.Split(cfg.DistributionTester, "\n") {
		email = strings.TrimSpace(email)

		if len(email) == 0 {
			continue
		}

		log.Printf("- %s", email)

		if err := releaseAPI.AddTester(email); err != nil {
			failf("Failed to add tester(%s) to the release, error: %s", email, err)
		}
	}

	log.Donef("- Done")
	fmt.Println()

	log.Infof("Exporting outputs")

	var outputs = map[string]string{
		statusEnvKey:                    "success",
		"APPCENTER_DEPLOY_INSTALL_URL":  release.InstallURL,
		"APPCENTER_DEPLOY_DOWNLOAD_URL": release.DownloadURL,
		"APPCENTER_RELEASE_PAGE_URL":    fmt.Sprintf("https://appcenter.ms/orgs/%s/apps/%s/distribute/releases/%d", cfg.OwnerName, cfg.AppName, release.ID),
		"APPCENTER_DEPLOY_RELEASE_ID":   strconv.Itoa(release.ID),
	}

	if len(publicGroup) == 1 {
		outputs["APPCENTER_PUBLIC_INSTALL_PAGE_URL"] = fmt.Sprintf("https://install.appcenter.ms/users/%s/apps/%s/distribution_groups/%s", cfg.OwnerName, cfg.AppName, publicGroup[0])
	}

	if len(publicGroup) > 1 {
		var groupUrls []string
		for _, groupName := range publicGroup {
			groupUrls = append(groupUrls, fmt.Sprintf("https://install.appcenter.ms/users/%s/apps/%s/distribution_groups/%s", cfg.OwnerName, cfg.AppName, groupName))
		}

		outputs["APPCENTER_PUBLIC_INSTALL_PAGE_URLS"] = strings.Join(groupUrls, ", ")
	}

	for key, value := range outputs {
		log.Printf("- %s: %s", key, value)
		if err := tools.ExportEnvironmentWithEnvman(key, value); err != nil {
			failf("Failed to export environment variable: %s with value: %s. Error: %s", key, value, err)
		}
	}

	log.Donef("- Done")
}

func failf(f string, args ...interface{}) {
	log.Errorf(f, args...)

	if err := tools.ExportEnvironmentWithEnvman(statusEnvKey, "failed"); err != nil {
		log.Errorf("Failed to export environment variable: %s with value: %s. Error: %s", statusEnvKey, "failed", err)
	}

	os.Exit(1)
}
