package cmd

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	semver "github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultVersion string = "no-version"
	ghUser         string = "willfantom"
	ghRepo         string = "overseerr-exporter"
)

var version string = defaultVersion

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of OverCLIrr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Overseerr Exporter Version: %s\n", version)
	},
}

var overseerrVersion = &cobra.Command{
	Use:   "overseerr-version",
	Short: "Print the version of the connected Overseerr instance",
	PreRun: func(cmd *cobra.Command, args []string) {
		setOverseer()
	},
	Run: func(cmd *cobra.Command, args []string) {
		status, err := overseerr.Status()
		if err != nil {
			logrus.WithField("message", err.Error()).Fatalln("Could not get version info from Overseerr")
		}
		fmt.Printf("Overseerr Version: %s\n", status.Version)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check if Overseerr Exporter is up to date",
	Run: func(cmd *cobra.Command, args []string) {
		available, err := checkForUpdate()
		if err != nil {
			logrus.WithField("message", err.Error()).Fatalln("Could not get check for overseerr exporter updates")
		}
		if available {
			fmt.Printf("Overseerr Exporter update available\n")
		}
		fmt.Printf("Overseerr Exporter is the latest version\n")
	},
}

func checkForUpdate() (bool, error) {
	semverVersion, err := semver.NewSemver(version)
	if err != nil {
		logrus.WithField("extended", err.Error()).
			Errorln("could not determine the version of overseerr exporter")
		return false, err
	}
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), ghUser, ghRepo)
	if err != nil {
		logrus.WithField("extended", err.Error()).
			Errorln("could get repository release info from github")
		return false, err
	}
	ghVer, err := semver.NewSemver(*release.TagName)
	if err != nil {
		logrus.WithField("extended", err.Error()).
			Errorln("latest github release tag not semver compliant")
		return false, err
	}
	if ghVer.GreaterThan(semverVersion) {
		logrus.Infoln("found a more recent release on github")
		return true, nil
	}
	logrus.Infoln("overseerr-exporter found to be latest version")
	return false, nil
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(updateCmd)
	RootCmd.AddCommand(overseerrVersion)
}
