package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Задание 1
func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")

	if name == "" || age == "" {
		http.Error(w, "Отсутствует параметр имени или возраста", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("Меня зовут %s, мне %s лет", name, age)
	fmt.Fprintln(w, response)
	log.Println()
}

// Задание 2
func addHandler(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")
	a, errA := strconv.Atoi(aStr)
	b, errB := strconv.Atoi(bStr)

	if errA != nil || errB != nil {
		http.Error(w, "Недопустимые параметры. Пожалуйста, укажите два допустимых целых числа.", http.StatusBadRequest)
		return
	}

	result := a + b
	fmt.Fprintf(w, "Результат сложения: %d", result)
	log.Printf("Выполненное сложение: %d + %d = %d\n", a, b, result)
}

func subHandler(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")
	a, errA := strconv.Atoi(aStr)
	b, errB := strconv.Atoi(bStr)

	if errA != nil || errB != nil {
		http.Error(w, "Недопустимые параметры. Пожалуйста, укажите два допустимых целых числа.", http.StatusBadRequest)
		return
	}

	result := a - b
	fmt.Fprintf(w, "Результат вычитания: %d", result)
	log.Printf("Выполненное вычитание: %d - %d = %d\n", a, b, result)
}

func mulHandler(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")
	a, errA := strconv.Atoi(aStr)
	b, errB := strconv.Atoi(bStr)

	if errA != nil || errB != nil {
		http.Error(w, "Недопустимые параметры. Пожалуйста, укажите два допустимых целых числа.", http.StatusBadRequest)
		return
	}

	result := a * b
	fmt.Fprintf(w, "Результат умножения: %d", result)
	log.Printf("Выполненное умножение: %d * %d = %d\n", a, b, result)
}

func divHandler(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")
	a, errA := strconv.Atoi(aStr)
	b, errB := strconv.Atoi(bStr)

	if errA != nil || errB != nil {
		http.Error(w, "Недопустимые параметры. Пожалуйста, укажите два допустимых целых числа.", http.StatusBadRequest)
		return
	}

	if b == 0 {
		http.Error(w, "Деление на ноль не допускается.", http.StatusBadRequest)
		return
	}

	result := a / b
	fmt.Fprintf(w, "Результат деления: %d", result)
	log.Printf("Выполненное делениен: %d / %d = %d\n", a, b, result)
}

// Задание 3
type RequestBody struct {
	Text string `json:"text"`
}

func charCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestBody RequestBody
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	charCount := make(map[rune]int)
	for _, char := range requestBody.Text {
		charCount[char]++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(charCount)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

	log.Printf("Character count for input: %s\n", requestBody.Text)
}

func main() {
	http.HandleFunc("/", handler) // 1

	http.HandleFunc("/add", addHandler) // 2
	http.HandleFunc("/sub", subHandler)
	http.HandleFunc("/mul", mulHandler)
	http.HandleFunc("/div", divHandler)

	http.HandleFunc("/charcount", charCountHandler) // 3

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
