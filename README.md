# android-app-publisher

![release](https://img.shields.io/github/v/release/safu9/android-app-publisher)

Upload your app to play store from terminal or ci environment

## Install

```
go get -u github.com/safu9/android-app-publisher
```

## Usage

1. [Create service account from Play Console and add permissions](https://play.google.com/console/u/0/developers/api-access)
2. [Create and download JSON key file for the service account](https://console.cloud.google.com/apis/credentials/serviceaccountkey)
3. Run `android-app-publisher upload <app-file> --package <package-name> --credentials <key-file>`  
Or, you can use `ANDROID-APP-UPLOADER-CREDENTIALS` environment variable to specify the key file.
4. Done!
