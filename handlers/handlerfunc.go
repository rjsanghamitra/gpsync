package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rjsanghamitra/gpsync/database"
	"github.com/rjsanghamitra/gpsync/oauth"
	"golang.org/x/oauth2"
)

type Mdata struct {
	Ctime  string `json:"creationTime"`
	Height string `json:"height"`
	Width  string `json:"width"`
}

type Lib struct {
	Id     string `json:"id"`
	Purl   string `json:"productUrl"`
	Burl   string `json:"baseUrl"`
	Mtype  string `json:"mimeType"`
	Mmdata Mdata  `json:"mediaMetadata"`
	Fname  string `json:"filename"`
}

var Name = "photos"
var db = createDb(oauth.DPath + "/" + Name)

func Login(w http.ResponseWriter, r *http.Request) {
	url := oauth.GoogleOauthConfig.AuthCodeURL(oauth.RandomState) // the state that is passed in is used to identify the user request
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func OauthAndDownload(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauth.RandomState {
		fmt.Println("State is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	token, err := oauth.GoogleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	RedirectDueToError(w, r, err)
	DownloadLibraryPhotos(w, r, token)
	go DownloadAlbumPhotos(w, r, token)
	if err == nil {
		fmt.Fprintln(w, "Download completed. Please close this page.")
	}
}

func DownloadAlbumPhotos(w http.ResponseWriter, r *http.Request, token *oauth2.Token) {
	// albums links
	req, err1 := http.Get("https://photoslibrary.googleapis.com/v1/albums?access_token=" + token.AccessToken)
	RedirectDueToError(w, r, err1)

	var temp map[string]([]map[string]string)
	content, _ := io.ReadAll(req.Body)
	json.Unmarshal(content, &temp)

	// retreiving all the photos links from albums
	database.CreateAlbumsTable(*db, "Albums")
	createFolder("Album Photos")
	for _, i := range temp["albums"] {
		var res map[string]([]map[string]string)
		jsonText := url.Values{
			"pageSize": {"100"},
			"albumId":  {i["id"]},
		}
		resp, err := http.PostForm("https://photoslibrary.googleapis.com/v1/mediaItems:search?access_token="+token.AccessToken, jsonText)
		RedirectDueToError(w, r, err)
		resp.Header.Set("Content-type", "application/json")
		resp.Header.Set("Authorization", "Bearer "+token.AccessToken)
		json.NewDecoder(resp.Body).Decode(&res)
		for _, y := range res["mediaItems"] {
			database.InsertIntoAlbums(db, "Albums", y["filename"], y["id"])
			createFolder(oauth.DPath + "/Album Photos/" + i["title"])
			if !FileDownloaded(db, y["filename"]) {
				download(y["baseUrl"], oauth.DPath+"/Album Photos/"+i["title"]+"/"+y["filename"])
			} else {
				basePath := oauth.DPath + "/Library Photos"
				absPath, _ := FindAbsPath(basePath, y["filename"])
				err := os.Symlink(absPath, oauth.DPath+"/Album Photos/"+i["title"]+"/"+y["filename"])
				CheckError(err)
			}
		}
	}
	os.Exit(0)
}

func DownloadLibraryPhotos(w http.ResponseWriter, r *http.Request, token *oauth2.Token) {
	// first set of media items
	req, err := http.Get("https://photoslibrary.googleapis.com/v1/mediaItems?access_token=" + token.AccessToken)
	RedirectDueToError(w, r, err)
	resp, err := io.ReadAll(req.Body)
	CheckError(err)
	var temp2 map[string]([]map[string]string)
	json.Unmarshal(resp, &temp2)

	// rest of the media items
	database.CreateLibraryTable(*db, "LibraryPhotos")
	for {
		var temp1 map[string]string
		json.Unmarshal(resp, &temp1)
		tok := temp1["nextPageToken"]
		if tok == "" {
			break
		}
		req, err = http.Get("https://photoslibrary.googleapis.com/v1/mediaItems?access_token=" + token.AccessToken + "&pageToken=" + tok)
		var temp2 map[string]([]Lib)
		resp, _ = io.ReadAll(req.Body)
		json.Unmarshal(resp, &temp2)
		for _, b := range temp2["mediaItems"] {
			database.InsertIntoLibrary(db, "LibraryPhotos", b.Id, b.Mtype, b.Fname, b.Mmdata.Ctime)
			createFolder(oauth.DPath + "/Library Photos/" + b.Mmdata.Ctime[:4] + "/" + month(b.Mmdata.Ctime[5:7]) + "/" + b.Mmdata.Ctime[8:10])
			if b.Mtype == "image/jpeg" {
				download(b.Burl, oauth.DPath+"/Library Photos/"+b.Mmdata.Ctime[:4]+"/"+month(b.Mmdata.Ctime[5:7])+"/"+b.Mmdata.Ctime[8:10]+"/"+b.Fname)
			} else if b.Mtype == "video/mp4" {
				download(b.Burl+"=dv", oauth.DPath+"/Library Photos/"+b.Mmdata.Ctime[:4]+"/"+month(b.Mmdata.Ctime[5:7])+"/"+b.Mmdata.Ctime[8:10]+"/"+b.Fname)
			}
		}
		RedirectDueToError(w, r, err)
	}
}

func month(a string) string {
	num, _ := strconv.Atoi(a)
	return time.Month(num).String()
}

func download(url string, fname string) {
	filename := string(fname)
	resp, _ := http.Get(url)
	if resp.StatusCode != 200 {
		fmt.Println("Received non-200 response code")
	}
	file, err := os.Create(filename)
	CheckError(err)
	_, err = io.Copy(file, resp.Body)
	CheckError(err)
}

func createDb(name string) *sql.DB {
	_, err := os.Create(name + ".db")
	CheckError(err)
	db, err := sql.Open("sqlite3", name+".db")
	CheckError(err)
	return db
}

func createFolder(name string) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		err = os.MkdirAll(name, os.ModePerm)
		CheckError(err)
	}
}

func FileDownloaded(db *sql.DB, fname string) bool {
	stmt := "SELECT * FROM LibraryPhotos WHERE filename=?"
	rows, _ := db.Query(stmt, fname)
	if rows.Next() {
		return true
	} else {
		return false
	}
}

func FindAbsPath(base string, fname string) (string, error) {
	var absPath string
	ferr := filepath.Walk(base, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fname == info.Name() {
			absPath = path
		}
		return nil
	})
	if ferr != nil {
		return "", nil
	}
	return absPath, nil
}

func CheckError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

func RedirectDueToError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
