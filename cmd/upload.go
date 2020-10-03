package cmd

import (
	"context"
	"fmt"
	"os"

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
var serviceAccount string

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
	cmdUpload.Flags().StringVarP(&serviceAccount, "credential", "c", "", "Service account file")

	err := cmdUpload.MarkFlagRequired("package")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.AddCommand(cmdUpload)
}

func upload(filename string) error {
	ctx := context.Background()

	if serviceAccount == "" {
		serviceAccount = os.Getenv("ANDROID-APP-UPLOADER-CREDENTIAL")
	}

	// client := new(http.Client)
	// client.Timeout = time.Minute

	// Create the AndroidPublisherService
	options := []option.ClientOption{
		option.WithCredentialsFile(serviceAccount),
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

	// Upload new apk to developer console
	bundle, err := publisherService.Edits.Bundles.Upload(packageName, edit.Id).Media(file, googleapi.ContentType("application/octet-stream")).Do()
	if err != nil {
		return err
	}

	// Assign apk to production track
	track := "production" // "alpha", "beta", "internal" or "production"
	_, err = publisherService.Edits.Tracks.Update(packageName, edit.Id, track, &androidpublisher.Track{
		Releases: []*androidpublisher.TrackRelease{
			{VersionCodes: []int64{bundle.VersionCode}},
		},
	}).Do()
	if err != nil {
		return err
	}

	// Commit changes for edit (publish)
	// publisherService.Edits.Commit(packageName, edit.Id)

	return nil
}