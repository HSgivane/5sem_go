package main

import (
	"database/sql"
	"fmt"
	"log"
	"proekt/database"
	"time"

	"gopkg.in/telebot.v3"
)

const AdminPassword = "купить админку"

func main() {
	db := database.InitDB()
	defer db.Close()

	botToken := "7418686248:AAGhuHY_c2oA0i5JNuZ01PdwRXYuV5HSleU"
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}

	replyMenu := &telebot.ReplyMarkup{}
	btnIdea := replyMenu.Text("Предложить идею")
	btnViewIdeas := replyMenu.Text("Посмотреть идеи")

	replyMenu.Reply(
		replyMenu.Row(btnIdea),
		replyMenu.Row(btnViewIdeas),
	)

	// стврт
	bot.Handle("/start", func(c telebot.Context) error {
		userID := c.Sender().ID
		database.AddUser(db, userID) // Регистрация пользователя
		return c.Send("Добро пожаловать!", replyMenu)
	})

	// предложить идею
	bot.Handle(&btnIdea, func(c telebot.Context) error {
		userID := c.Sender().ID
		database.SetProposingIdea(db, userID, 1) // Устанавливаем флаг
		return c.Send("Введите вашу идею:")
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID
		message := c.Text()

		// Проверка на ввод пароля
		if message == AdminPassword {
			database.SetAdmin(db, userID)
			return c.Send("Теперь вы админ.")
		}

		// проверка флага
		isProposing := database.IsProposingIdea(db, userID)
		if isProposing {
			database.SaveIdea(db, userID, message)
			database.SetProposingIdea(db, userID, 0) // Сбрасываем флаг
			return c.Send("Спасибо за вашу идею!")
		}

		return c.Send("Команда не распознана.")
	})

	// сохранение картинки
	bot.Handle(telebot.OnPhoto, func(c telebot.Context) error {
		userID := c.Sender().ID
		isProposing := database.IsProposingIdea(db, userID)

		if isProposing {
			// Сохраняем изображение
			photo := c.Message().Photo
			fileName := fmt.Sprintf("uploads/%d_%s.jpg", userID, time.Now().Format("20060102150405"))
			err := bot.Download(photo.MediaFile(), fileName)
			if err != nil {
				return c.Send("Ошибка загрузки изображения.")
			}

			// Получаем текст сообщения, если он есть
			text := c.Message().Caption
			if text == "" {
				text = "None"
			}

			database.SaveIdeaWithImage(db, userID, text, fileName)
			database.SetProposingIdea(db, userID, 0)
			return c.Send("Ваша идея с изображением сохранена.")
		}

		return c.Send("Отправка изображений разрешена только в режиме добавления идеи.")
	})

	// посмотреть идеи
	bot.Handle(&btnViewIdeas, func(c telebot.Context) error {
		userID := c.Sender().ID
		if database.CheckPermission(db, userID) == "admin" {
			idea, imagePath, err := database.PopIdeaWithImage(db)
			if err != nil {
				if err == sql.ErrNoRows {
					return c.Send("Идей больше нет.")
				}
				log.Printf("Ошибка извлечения идеи: %v", err)
				return c.Send("Произошла ошибка при получении идеи.")
			}

			if imagePath != "" {
				photo := &telebot.Photo{File: telebot.FromDisk(imagePath)}
				// Сначала отправляем изображение
				err := c.Send(photo)
				if err != nil {
					log.Printf("Ошибка отправки изображения: %v", err)
					return c.Send("Ошибка отправки изображения.")
				}
				// Затем отправляем текст
				return c.Send(fmt.Sprintf("Идея: %s", idea))
			}

			// Если изображения нет, отправляем только текст
			return c.Send(fmt.Sprintf("Идея: %s", idea))
		}
		return c.Send("У вас нет прав для просмотра идей.")
	})

	bot.Handle("/офнись", func(c telebot.Context) error {
		bot.Stop()
		return c.Send("Бот выключается...")
	})

	bot.Start()
}
