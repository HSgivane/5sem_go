package main

import (
	"fmt"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sync"
	"time"
)

// filterParallel обрабатывает только одну строку пикселей
func filterParallel(img draw.RGBA64Image, y int) {
	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		origColor := img.RGBA64At(x, y)
		gray := uint16((uint32(origColor.R) + uint32(origColor.G) + uint32(origColor.B)) / 3)
		newColor := color.RGBA64{
			R: gray,
			G: gray,
			B: gray,
			A: origColor.A,
		}
		img.SetRGBA64(x, y, newColor)
	}
}

func main() {
	inputFile, err := os.Open("wakeup.png")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer inputFile.Close()

	srcImg, err := png.Decode(inputFile)
	if err != nil {
		fmt.Println("Ошибка декодирования PNG:", err)
		return
	}

	drawImg, ok := srcImg.(draw.RGBA64Image)
	if !ok {
		fmt.Println("Преобразование к draw.RGBA64Image не удалось")
		return
	}

	// Время начала
	start := time.Now()

	bounds := drawImg.Bounds()
	height := bounds.Max.Y - bounds.Min.Y

	// Создаём WaitGroup на количество строк
	var wg sync.WaitGroup
	wg.Add(height)

	// Запускаем горутину на каждую строку
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		go func(row int) {
			// Не забываем вызвать Done() после окончания работы горутины
			defer wg.Done()
			filterParallel(drawImg, row)
		}(y)
	}

	// Ждём, пока все горутины завершатся
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Время параллельной обработки: %v\n", elapsed)

	outputFile, err := os.Create("output_parallel.png")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, drawImg); err != nil {
		fmt.Println("Ошибка сохранения PNG:", err)
		return
	}

	fmt.Println("Параллельная обработка завершена. Результат сохранён в output_parallel.png")
}
