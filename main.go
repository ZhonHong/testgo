package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

// 定義機器人用戶名
const botUsername = "TgGoAuth_bot"
const botToken = "7377067980:AAHPAw5RzXXe2N_kgWcYWBYsnEhtVOuxn-4"

type AuthData struct {
	AuthDate  string `json:"auth_date"`
	FirstName string `json:"first_name"`
	Id        string `json:"id"`
	LastName  string `json:"last_name"`
}

// checkTelegramAuthorization 驗證 Telegram 授權數據
func checkTelegramAuthorization(w http.ResponseWriter, r *http.Request) {
	authData := r.URL.Query() // 獲取查詢參數

	// 將查詢參數轉換為 map[string]string
	authDataMap := make(map[string]string)
	for key, values := range authData {
		if len(values) > 0 {
			authDataMap[key] = values[0]
		}
	}

	checkHash := authDataMap["hash"]
	delete(authDataMap, "hash")

	// 創建排序的數據檢查字符串
	dataCheckArr := make([]string, 0, len(authDataMap))
	for key, value := range authDataMap {
		dataCheckArr = append(dataCheckArr, key+"="+value)
	}
	sort.Strings(dataCheckArr)
	dataCheckString := ""
	for _, item := range dataCheckArr {
		if dataCheckString != "" {
			dataCheckString += "\n"
		}
		dataCheckString += item
	}

	// 計算哈希值並檢查
	secretKey := sha256.Sum256([]byte(botToken))
	hmacHash := hmac.New(sha256.New, secretKey[:])
	hmacHash.Write([]byte(dataCheckString))
	hash := hex.EncodeToString(hmacHash.Sum(nil))

	if hash != checkHash {
		http.Error(w, "Data is NOT from Telegram", http.StatusUnauthorized)
		return
	}

	// 檢查數據是否過期
	authDate, err := strconv.Atoi(authDataMap["auth_date"])
	if err != nil {
		http.Error(w, "Invalid auth_date", http.StatusBadRequest)
		return
	}
	if time.Now().Unix()-int64(authDate) > 86400 {
		http.Error(w, "Data is outdated", http.StatusUnauthorized)
		return
	}

	// 保存用戶數據
	authDataJSON, err := json.Marshal(authDataMap)
	if err != nil {
		http.Error(w, "Failed to save user data", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "tg_user",
		Value: url.QueryEscape(string(authDataJSON)),
		Path:  "/",
	})

	// 重定向到主頁面
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// getTelegramUserData 獲取 Telegram 用戶數據
func getTelegramUserData(r *http.Request) (AuthData, bool) {
	var authData AuthData
	cookie, err := r.Cookie("tg_user")
	if err != nil || cookie == nil {
		return authData, false
	}

	authDataJSON, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return authData, false
	}

	if err := json.Unmarshal([]byte(authDataJSON), &authData); err != nil {
		return authData, false
	}

	return authData, true
}

// handleRequest 處理請求
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 處理登出
	if r.URL.Query().Get("logout") == "1" {
		http.SetCookie(w, &http.Cookie{
			Name:   "tg_user",
			Value:  "",
			MaxAge: -1, // 立即過期
		})
		http.Redirect(w, r, "https://512b-36-235-198-9.ngrok-free.app/", http.StatusSeeOther)
		return
	}

	// 獲取 Telegram 用戶數據
	tgUser, ok := getTelegramUserData(r)
	fmt.Printf("%v\n", tgUser)
	var htmlContent string
	if ok {
		firstName := html.EscapeString(tgUser.FirstName)
		lastName := html.EscapeString(tgUser.LastName)
		fmt.Println(firstName, lastName)
		// if username, exists := tgUser["username"]; exists {
		// 	username = html.EscapeString(username)
		// 	htmlContent = `<h1>Hello, <a href="https://t.me/` + username + `">` + firstName + ` ` + lastName + `</a>!</h1>`
		// } else {
		// 	htmlContent = `<h1>Hello, ` + firstName + ` ` + lastName + `!</h1>`
		// }
		// if photoURL, exists := tgUser["photo_url"]; exists {
		// 	photoURL = html.EscapeString(photoURL)
		// 	htmlContent += `<img src="` + photoURL + `">`
		// }
		htmlContent += `<p><a href="?logout=1">Log out</a></p>`
	} else {
		htmlContent = `<h1>Hello, anonymous!</h1>`
		htmlContent += `<script async src="https://telegram.org/js/telegram-widget.js?2" data-telegram-login="` + botUsername + `" data-size="large" data-auth-url="https://512b-36-235-198-9.ngrok-free.app/check_authorization"></script>`
	}

	// 設置模板並生成 HTML
	tmpl := `<html>
        <head>
            <meta charset="utf-8">
            <title>Login Widget Example</title>
        </head>
        <body><center>{{ . }}</center></body>
    </html>`

	t := template.Must(template.New("webpage").Parse(tmpl))
	t.Execute(w, template.HTML(htmlContent))
}

func main() {
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/check_authorization", checkTelegramAuthorization)
	http.ListenAndServe(":8080", nil) // 監聽 8080 端口
}
