This is a tool written in Golang that creates a backup of Google Photos in your local system. 
- It uses Goroutines to download Album Photos.
- It does not create duplicates of photos.
- It creates symlinks of photos in Album Photos if they already exist in Library Photos.

# How to use?
Create a Google OAuth Web Client ID. Download the json file and store it in your home directory as "client_secret.json".

You can follow the steps in this webpage to create a Google OAuth Client ID 
https://developers.google.com/identity/gsi/web/guides/get-google-api-clientid

Clone this repo:
`git clone https://github.com/rjsanghamitra/gpsync.git`

Run the command to create the binary:
`go install gpsync`

Run the following command to download the photos:
`gpsync --path <path>`
You can specify the path where you want to store the photos as a command line argument.

**Disclaimer:** Since this tool uses the Google Photos API, which is known to be slightly buggy, it might not download all the photos in your account. 
