package model

// ReleaseOptions ...
type ReleaseOptions struct {
	BuildVersion  string
	BuildNumber   string
	GroupNames    []string
	Mandatory     bool
	NotifyTesters bool
	FilePath      string
	App           App
}
