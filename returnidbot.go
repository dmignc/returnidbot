package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api" //get -u github.com/go-telegram-bot-api/telegram-bot-api
)

var (
	// переменная с токеном
	telegramBotToken string
)

func init() {
	// принимаем на входе флаг -telegrambottoken
	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.Parse()

	// без него не запускаемся
	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}

	fmt.Println("platform:", runtime.GOOS)
}

func main() {

	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)
	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		readingmessageFunc(update, bot)
	}
}

//функция чтения сообщения в чате
func readingmessageFunc(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

	var reply string             //переменная с текстом ответа
	var telegram_id string       //переменная с id отправителя
	var telegram_username string //переменная с username отправителя

	// логируем от кого какое сообщение пришло
	log.Printf("[%s:%s] %s", update.Message.From.UserName, strconv.Itoa(update.Message.From.ID), update.Message.Text)
	//WriteToLogfile("returnidbot.log", update.Message.From.UserName+" | "+strconv.Itoa(update.Message.From.ID)+" | "+update.Message.Text)
	// свитч на обработку комманд
	// комманда - сообщение, начинающееся с "/"
	switch update.Message.Command() {
	case "start":

		//telegram_id = string(update.Message.Chat.ID)
		telegram_username = string(update.Message.Chat.UserName)
		telegram_id = strconv.FormatInt(update.Message.Chat.ID, 10)
		reply = "Your ID: " + telegram_id + "\nYour UserName: @" + telegram_username

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
		log.Printf("[%s:%s] %s", update.Message.From.UserName, strconv.Itoa(update.Message.From.ID), update.Message.Text)
		//WriteToLogfile("returnidbot.log", time.Now().String()+" | "+telegram_id+" | "+telegram_username)
	}
}

func WriteToLogfile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}
