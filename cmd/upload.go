package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// References
// https://developers.google.com/android-publisher
// https://developers.google.com/android-publisher/api-ref/rest
// https://godoc.org/google.golang.org/api/androidpublisher/v3
// https://github.com/googleapis/google-api-go-client

var packageName string
var credentialsFile string
var track string

var cmdUpload = &cobra.Command{
	Use:   "upload [file]",
	Short: "Upload app file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		err := upload(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Uploaded app to Play Console!")
	},
}

func init() {
	cmdUpload.Flags().StringVarP(&packageName, "package", "p", "", "Package name of app")
	cmdUpload.Flags().StringVarP(&credentialsFile, "credentials", "c", "", "Service account credentials file")
	cmdUpload.Flags().StringVarP(&track, "track", "t", "production", "Track to upload app\n\"production\", \"alpha\", \"beta\" or \"internal\"")

	err := cmdUpload.MarkFlagRequired("package")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.AddCommand(cmdUpload)
}

func upload(filename string) error {
	ctx := context.Background()

	if credentialsFile == "" {
		credentialsFile = os.Getenv("ANDROID_APP_UPLOADER_CREDENTIALS")
	}

	// client := new(http.Client)
	// client.Timeout = time.Minute

	// Create the AndroidPublisherService
	options := []option.ClientOption{
		option.WithCredentialsFile(credentialsFile),
		option.WithScopes(androidpublisher.AndroidpublisherScope),
		// option.WithHTTPClient(client),
	}

	publisherService, err := androidpublisher.NewService(ctx, options...)
	if err != nil {
		return err
	}

	// Create a new edit to make changes to your listing
	edit, err := publisherService.Edits.Insert(packageName, nil).Do()
	if err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var versionCode int64

	// Upload new apk to developer console
	switch filepath.Ext(filename) {
	case ".apk":
		apk, err := publisherService.Edits.Apks.Upload(packageName, edit.Id).Media(file, googleapi.ContentType("application/octet-stream")).Do()
		if err != nil {
			return err
		}
		versionCode = apk.VersionCode

	case ".aab":
		bundle, err := publisherService.Edits.Bundles.Upload(packageName, edit.Id).Media(file, googleapi.ContentType("application/octet-stream")).Do()
		if err != nil {
			return err
		}
		versionCode = bundle.VersionCode

	default:
		return errors.New("This file type is not supported.")
	}

	// Assign apk to production track
	_, err = publisherService.Edits.Tracks.Update(packageName, edit.Id, track, &androidpublisher.Track{
		Releases: []*androidpublisher.TrackRelease{
			{
				VersionCodes: []int64{versionCode},
				Status:       "draft", // "draft", "inProgress", "halted" or "completed"
			},
		},
	}).Do()
	if err != nil {
		return err
	}

	// Commit changes for edit
	_, err = publisherService.Edits.Commit(packageName, edit.Id).Do()
	if err != nil {
		return err
	}

	return nil
}
