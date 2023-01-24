This is a tool written in Golang that creates a backup of Google Photos in your local system. It uses Goroutines to download album photos.

# How to use?

Create a Google OAuth Web Client ID. Download the json file and store it in your home directory.

You can follow the steps in this webpage to create a Google OAuth Client ID 
https://support.google.com/workspacemigrate/answer/9222992?hl=en

Download this
https://github.com/rjsanghamitra/gpsync/releases/download/v1.0.0/gpsync_1.0.0_Linux_x86_64.tar.gz
Extract the file contents and open terminal in that directory.
Run the following command:
`./gpsync --path <path>`
You can specify the path where you want to store the photos as a command line argument.
