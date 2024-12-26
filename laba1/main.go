package main

import (
	"errors"
	"fmt"
	"math"
)

func main() {
	fmt.Println("Hello, World!")

	fmt.Println("\n1 задача")
	hello()

	fmt.Println("\n2 задача")
	fmt.Println("\nTest 1")
	err := printEven(5, 12)
	if err != nil {
		fmt.Println("Ошибка: 0", err)
	}

	fmt.Println("\nTest 2")
	err2 := printEven(10, 1)
	if err2 != nil {
		fmt.Println("Ошибка: 0", err2)
	}

	fmt.Println("\n3 задача")
	result(3, 5, "+")

	result(3, 5, "-")

	result(7, 10, "*")

	result(25, 5, "/")

	result(3, 5, "#")
}

// Первая задача
func hello() {
	var name string
	fmt.Print("Введите ваше имя: ")
	fmt.Scanln(&name)
	fmt.Printf("Привет, %s! \n", name)
}

// Вторая задача
func printEven(a, b int64) error {
	if a > b {
		return errors.New("левая граница диапазона больше правой")
	}

	for i := a; i <= b; i++ {
		if i%2 == 0 {
			fmt.Println(i)
		}
	}
	return nil
}

// Третья задача
func apply(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("деление на ноль")
		}
		return a / b, nil
	default:
		return 0, errors.New("действие не поддерживается")
	}
}

func result(a, b float64, operator string) {
	res, err := apply(a, b, operator)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		// Проверяем, является ли результат целым числом
		if res == math.Floor(res) {
			fmt.Printf("Результат: %d\n", int(res))
		} else {
			fmt.Printf("Результат: %.2f\n", res)
		}
	}
}
