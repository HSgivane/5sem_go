package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// админка
func SetAdmin(db *sql.DB, userID int64) {
	_, err := db.Exec(`
		UPDATE users SET permission = 'admin' WHERE id = ?
	`, userID)
	if err != nil {
		log.Printf("Ошибка назначения прав admin: %v", err)
	}
}

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "bot_database.db")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			permission TEXT DEFAULT 'default',
			proposing_idea INTEGER DEFAULT 0
		)
	`)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы пользователей: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ideas (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			idea TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы идей: %v", err)
	}

	return db
}

// флаг
func IsProposingIdea(db *sql.DB, userID int64) bool {
	var proposing int
	err := db.QueryRow(`
		SELECT proposing_idea FROM users WHERE id = ?
	`, userID).Scan(&proposing)
	if err != nil {
		log.Printf("Ошибка проверки флага proposing_idea: %v", err)
		return false
	}
	return proposing == 1
}

// рег
func AddUser(db *sql.DB, userID int64) {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO users (id) VALUES (?)
	`, userID)
	if err != nil {
		log.Printf("Ошибка добавления пользователя: %v", err)
	}
}

// флаг
func SetProposingIdea(db *sql.DB, userID int64, proposing int) {
	_, err := db.Exec(`
		UPDATE users SET proposing_idea = ? WHERE id = ?
	`, proposing, userID)
	if err != nil {
		log.Printf("Ошибка обновления флага: %v", err)
	}
}

// сохранение идеи
func SaveIdea(db *sql.DB, userID int64, idea string) {
	_, err := db.Exec(`
		INSERT INTO ideas (user_id, idea) VALUES (?, ?)
	`, userID, idea)
	if err != nil {
		log.Printf("Ошибка сохранения идеи: %v", err)
	}
}

// получение идей
func GetAllIdeas(db *sql.DB) []string {
	rows, err := db.Query(`
		SELECT idea FROM ideas
	`)
	if err != nil {
		log.Printf("Ошибка получения идей: %v", err)
		return nil
	}
	defer rows.Close()

	var ideas []string
	for rows.Next() {
		var idea string
		if err := rows.Scan(&idea); err != nil {
			log.Printf("Ошибка чтения строки: %v", err)
		}
		ideas = append(ideas, idea)
	}

	return ideas
}

func SaveIdeaWithImage(db *sql.DB, userID int64, idea string, imagePath string) {
	_, err := db.Exec(`
        INSERT INTO ideas (user_id, idea, image) VALUES (?, ?, ?)
    `, userID, idea, imagePath)
	if err != nil {
		log.Printf("Ошибка сохранения идеи с картинкой: %v", err)
	}
}

// получение картинки
func PopIdeaWithImage(db *sql.DB) (string, string, error) {
	var idea string
	var image sql.NullString
	var id int

	err := db.QueryRow(`
        SELECT id, idea, image FROM ideas LIMIT 1
    `).Scan(&id, &idea, &image)
	if err != nil {
		return "", "", err
	}

	// Обработка NULL для image
	imagePath := ""
	if image.Valid {
		imagePath = image.String
	}

	_, err = db.Exec(`
        DELETE FROM ideas WHERE id = ?
    `, id)
	if err != nil {
		return "", "", err
	}
	return idea, imagePath, nil
}

// получение одно идеи
func PopIdea(db *sql.DB) (string, error) {
	var idea string
	var id int
	err := db.QueryRow(`
        SELECT id, idea FROM ideas LIMIT 1
    `).Scan(&id, &idea)
	if err != nil {
		return "", err
	}
	_, err = db.Exec(`
        DELETE FROM ideas WHERE id = ?
    `, id)
	if err != nil {
		return "", err
	}
	return idea, nil
}

// Проверка прав
func CheckPermission(db *sql.DB, userID int64) string {
	var permission string
	err := db.QueryRow(`
		SELECT permission FROM users WHERE id = ?
	`, userID).Scan(&permission)
	if err != nil {
		log.Printf("Ошибка проверки прав: %v", err)
		return "default"
	}
	return permission
}
