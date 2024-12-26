package main

import (
	"fmt"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"time"
)

func filter(img draw.RGBA64Image) {
	// Получаем границы изображения
	bounds := img.Bounds()
	// Перебираем все пиксели
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Считываем исходный цвет
			origColor := img.RGBA64At(x, y)
			// Простейший перевод в градации серого:
			// (R + G + B) / 3
			gray := uint16((uint32(origColor.R) + uint32(origColor.G) + uint32(origColor.B)) / 3)

			// Создаём новый цвет (с учётом альфа-канала)
			newColor := color.RGBA64{
				R: gray,
				G: gray,
				B: gray,
				A: origColor.A, // альфа оставим без изменений
			}
			// Устанавливаем цвет в точку (x, y)
			img.SetRGBA64(x, y, newColor)
		}
	}
}

func main() {
	// 1. Открываем файл с картинкой
	inputFile, err := os.Open("wakeup.png")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer inputFile.Close()

	// 2. Декодируем PNG в переменную типа image.Image
	srcImg, err := png.Decode(inputFile)
	if err != nil {
		fmt.Println("Ошибка декодирования PNG:", err)
		return
	}

	// 3. Преобразуем (type assertion) к draw.RGBA64Image
	drawImg, ok := srcImg.(draw.RGBA64Image)
	if !ok {
		fmt.Println("Преобразование к draw.RGBA64Image не удалось")
		return
	}

	// 4. Замеряем время начала
	start := time.Now()

	// 5. Вызываем filter
	filter(drawImg)

	// 6. Замеряем время окончания и считаем дельту
	elapsed := time.Since(start)
	fmt.Printf("Время обработки изображения: %v\n", elapsed)

	// 7. Создаём файл для сохранения результата
	outputFile, err := os.Create("output.png")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer outputFile.Close()

	// 8. Сохраняем изображение
	if err := png.Encode(outputFile, drawImg); err != nil {
		fmt.Println("Ошибка сохранения PNG:", err)
		return
	}

	fmt.Println("Изображение успешно обработано и сохранено в output.png")
}
