package main

import (
	"errors"
	"fmt"
	"math"
)

// Задание 1
// formatIP принимает массив из 4 байтов и возвращает строку в формате IP.
func formatIP(ip [4]byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// listEven принимает два числа и возвращает срез чётных чисел и ошибку.
func listEven(start, end int) ([]int, error) {
	if start > end {
		return nil, errors.New("левая граница больше правой")
	}

	var evens []int
	for i := start; i <= end; i++ {
		if i%2 == 0 {
			evens = append(evens, i)
		}
	}
	return evens, nil
}

// Задание 2. Карты.
// countChars принимает строку и возвращает карту с количеством вхождений каждого символа.
func countChars(s string) map[rune]int {
	charCount := make(map[rune]int)
	for _, char := range s {
		charCount[char]++
	}
	return charCount
}

// Задание 3. Структуры, методы и интерфейсы.

type Point struct {
	X, Y float64
}

type Segment struct {
	Start, End Point
}

// Метод для вычисления длины отрезка.
func (s Segment) Length() float64 {
	dx := s.End.X - s.Start.X
	dy := s.End.Y - s.Start.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type Triangle struct {
	A, B, C Point
}

// Метод для вычисления площади треугольника.
func (t Triangle) Area() float64 {
	return math.Abs((t.A.X*(t.B.Y-t.C.Y) + t.B.X*(t.C.Y-t.A.Y) + t.C.X*(t.A.Y-t.B.Y)) / 2)
}

type Circle struct {
	Center Point
	Radius float64
}

// Метод для вычисления площади круга.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Интерфейс Shape с методом Area.
type Shape interface {
	Area() float64
}

// Функция для печати площади фигуры.
func printArea(s Shape) {
	result := s.Area()
	fmt.Printf("Площадь фигуры: %.2f\n", result)
}

// Задание 4. Функциональное программирование.

// Map применяет функцию ко всем элементам среза.
func Map(slice []float64, fn func(float64) float64) []float64 {
	result := make([]float64, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func main() {
	// Пример для задания 1.
	ip := [4]byte{127, 0, 0, 1}
	fmt.Println("IP-адрес:", formatIP(ip))

	// Пример вызова listEven.
	evens, err := listEven(1, 10)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Чётные числа:", evens)
	}

	// Пример для задания 2.
	text := "hello world"
	fmt.Println("Подсчёт символов:", countChars(text))

	// Пример для задания 3.
	triangle := Triangle{Point{0, 0}, Point{4, 0}, Point{0, 3}}
	circle := Circle{Point{0, 0}, 5}

	printArea(triangle)
	printArea(circle)

	// Пример для задания 4.
	slice := []float64{1.0, 2.0, 3.0, 4.0}
	square := func(x float64) float64 { return x * x }
	fmt.Println("Исходный срез:", slice)
	fmt.Println("Срез после применения Map:", Map(slice, square))
}
