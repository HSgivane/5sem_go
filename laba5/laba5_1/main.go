package main

import (
	"fmt"
	"time"
)

// count читает числа из канала и, например, выводит их квадраты.
// Чтение будет продолжаться, пока канал не закрыт.
func count(ch <-chan int) {
	for num := range ch {
		fmt.Printf("Получено число: %d, его квадрат: %d\n", num, num*num)
	}
	fmt.Println("Канал закрыт, горутина count завершила работу.")
}

func main() {
	// Создаём канал для передачи int
	ch := make(chan int)

	// Запускаем функцию count в отдельной горутине
	go count(ch)

	// Отправляем несколько чисел в канал
	for i := 1; i <= 5; i++ {
		ch <- i
	}

	// Закрываем канал
	close(ch)

	// Чтобы горутина count успела вывести все значения,
	// добавляем небольшую паузу или используем sync.WaitGroup (на выбор).
	time.Sleep(time.Second)

	fmt.Println("Функция main завершила работу.")
}
