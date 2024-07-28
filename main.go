package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

var (
	conf = &oauth2.Config{
		ClientID:     "20635273",
		ClientSecret: "-----BEGIN RSA PUBLIC KEY-----MIIBCgKCAQEAyMEdY1aR+sCR3ZSJrtztKTKqigvO/vBfqACJLZtS7QMgCGXJ6XIRyy7mx66W0/sOFa7/1mAZtEoIokDP3ShoqF4fVNb6XeqgQfaUHd8wJpDWHcR2OFwvplUUI1PLTktZ9uW2WE23b+ixNwJjJGwBDJPQEQFBE+vfmH0JP503wr5INS1poWg/j25sIWeYPHYeOrFp/eXaqhISP6G+q2IeTaWTXpwZj4LzXq5YOpk4bYEQ6mvRq7D1aHWfYmlEGepfaYR8Q0YqvvhYtMte3ITnuSJs171+GDqpdKcSwHnd6FudwGO4pcCOj4WcDuXc2CTHgH8gFTNhp/Y8/SpDOhvn9QIDAQAB-----END RSA PUBLIC KEY-----",
		Scopes:       []string{"profile"},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://oauth.telegram.org/token",
			AuthURL:  "https://oauth.telegram.org/auth",
		},
		RedirectURL: "https://88c8-36-235-202-251.ngrok-free.app//callback",
	}
	oauthStateString = "random" // 用於防止CSRF攻擊
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleMain)
	r.HandleFunc("/login", handleTelegramLogin)
	r.HandleFunc("/callback", handleTelegramCallback)

	http.Handle("/", r)
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	html := `<html><body><a href="/login">Login with Telegram</a></body></html>`
	fmt.Fprintf(w, html)
}

func handleTelegramLogin(w http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleTelegramCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("code exchange failed: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// 使用 token 訪問用戶數據
	fmt.Fprintf(w, "Access Token: %s\n", token.AccessToken)
}
