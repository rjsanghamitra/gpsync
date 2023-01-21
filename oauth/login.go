package oauth

import (
	"encoding/json"
	"flag"
	"io"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Auth struct {
	Client_id     string `json:"client_id"`
	Pid           string `json:"project_id"`
	Au            string `json:"auth_uri"`
	Tu            string `json:"token_uri"`
	Ap            string `json:"auth_provider_x509_cert_url"`
	Client_Secret string `json:"client_secret"`
	Ru            string `json:"redirect_uris"`
}

var (
	GoogleOauthConfig *oauth2.Config
	RandomState       string
	DPath             = GetDownloadPath()
)

func Client() {
	temp, _ := os.UserHomeDir()
	text, _ := os.Open(temp + "/client_secret.json")
	bv, _ := io.ReadAll(text)
	var res map[string]Auth
	json.Unmarshal([]byte(bv), &res)
	var Cid string = res["web"].Client_id
	var Csecret string = res["web"].Client_Secret
	defer text.Close()
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     Cid,
		ClientSecret: Csecret,
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
		Endpoint:     google.Endpoint,
	}
}

func GetDownloadPath() string {
	def, _ := os.UserHomeDir()
	PathFlag := flag.String("path", def, "Enter the path to the directory where you want the photos to be backed up.")
	flag.Parse()
	return *PathFlag
}
