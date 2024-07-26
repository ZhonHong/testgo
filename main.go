package main

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// 设置路由处理器
	// go startTelegramBot() // 使用 goroutine 启动 Telegram 机器人监听
	http.HandleFunc("/login", serveLoginHTML)

	// 启动 HTTP 服务器，监听 8080 端口
	port := ":8080"
	println("Server started on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

func startTelegramBot() {
	// 从环境变量中获取Telegram Bot的Token
	bot, err := tgbotapi.NewBotAPI("7229697482:AAF695rOVrb1ew0ncva6I8nLAe0x1LD03Wg")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true // 开启调试模式

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// 设置一个更新配置
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5

	// 获取Bot的更新信息
	updates, err := bot.GetUpdatesChan(u)

	// 处理接收到的消息
	for update := range updates {
		if update.Message == nil { // 忽略非消息事件
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// 回复消息示例
		reply := "I got your message: " + update.Message.Text
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}

func serveLoginHTML(w http.ResponseWriter, r *http.Request) {
	// 设置Content-Type为text/html，告诉浏览器返回的是HTML页面
	w.Header().Set("Content-Type", "text/html")

	// HTML内容
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Login Successful</title>
		</head>
		<body>
			<h1>登入成功</h1>
			<p>您已成功登入。</p>
		</body>
		</html>
	`

	// 将HTML内容写入到ResponseWriter中，以返回给客户端
	w.Write([]byte(htmlContent))
}
