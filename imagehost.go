package main

import (
	"fmt"
	"html/template"

	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/stacktic/dropbox"
)

func uploadFile(file multipart.File) string {
	clientid := os.Getenv("DROPBOX_KEY")
	clientsecret := os.Getenv("DROPBOX_TOKEN")
	token := os.Getenv("DROPBOX_ACCESS_TOKEN")

	DB := dropbox.NewDropbox()
	DB.SetAppInfo(clientid, clientsecret)
	DB.SetAccessToken(token)

	filename := fmt.Sprintf("%s.png", strconv.FormatInt(rand.Int63(), 32))
	filepath := fmt.Sprintf("/Public/Screenshots/%s", filename)

	var tmpstr []byte
	length, _ := file.Read(tmpstr)
	DB.FilesPut(file, int64(length), filepath, true, "")

	account, _ := DB.GetAccountInfo()
	uid := strconv.FormatInt(int64(account.UID), 10)

	return fmt.Sprintf("https://dl.dropbox.com/u/%s/Screenshots/%s", uid, filename)
}

func main() {
	// http.HandleFunc("/", timelineHandler)
	staticFs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", staticFs))

	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/upload", uploadHandler)

	http.ListenAndServe(":3001", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("templates", "layout.html")
	var p string

	if r.URL.Path == "/" {
		p = "index.html"
	} else {
		p = r.URL.Path
	}

	fp := path.Join("templates", p)

	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, _ := r.FormFile("file")
	publicURL := uploadFile(file)
	io.WriteString(w, publicURL)
}
