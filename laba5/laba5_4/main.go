package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"sync"
	"time"
)

// Пример ядра свёртки (3x3) для размытия
var kernel = [][]float64{
	{0.0625, 0.125, 0.0625},
	{0.125, 0.25, 0.125},
	{0.0625, 0.125, 0.0625},
}

// convolvePixel применяет ядро свёртки к пикселю (x, y) исходного изображения src
// и возвращает новый цвет, который записывается в результирующее изображение dst.
func convolvePixel(src draw.RGBA64Image, x, y int) color.RGBA64 {
	var sumR, sumG, sumB float64
	var sumA uint16 // альфа обычно не «размывают», но можно и обрабатывать аналогично

	bounds := src.Bounds()

	kernelSize := len(kernel) // у нас 3

	// Определяем «середину» ядра; для 3x3 это 1
	offset := kernelSize / 2

	for ky := 0; ky < kernelSize; ky++ {
		for kx := 0; kx < kernelSize; kx++ {
			// Координаты пикселя исходного изображения, которые участвуют в свёртке
			ix := x + (kx - offset)
			iy := y + (ky - offset)

			// Проверяем границы (если вышли за границы, можно либо пропускать, либо clamp-ить)
			if ix < bounds.Min.X || ix >= bounds.Max.X ||
				iy < bounds.Min.Y || iy >= bounds.Max.Y {
				continue
			}

			// Считываем цвет из исходного изображения
			c := src.RGBA64At(ix, iy)

			weight := kernel[ky][kx]

			sumR += float64(c.R) * weight
			sumG += float64(c.G) * weight
			sumB += float64(c.B) * weight
			// Альфа — по желанию: можно усреднять, а можно просто брать исходное значение
			// Здесь для примера возьмём ту же схему
			sumA = c.A
		}
	}

	// Округляем и «обрезаем», чтобы было в пределах 0..65535
	newR := uint16(math.Min(math.Max(sumR, 0.0), 65535.0))
	newG := uint16(math.Min(math.Max(sumG, 0.0), 65535.0))
	newB := uint16(math.Min(math.Max(sumB, 0.0), 65535.0))
	// Альфа без изменений или усреднить (здесь оставляем, что было в последней итерации)
	newA := uint16(sumA)

	return color.RGBA64{R: newR, G: newG, B: newB, A: newA}
}

// convolveRow обрабатывает одну строку (y) и результат записывает в dstImg
func convolveRow(srcImg draw.RGBA64Image, dstImg draw.RGBA64Image, y int) {
	bounds := srcImg.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		newColor := convolvePixel(srcImg, x, y)
		dstImg.SetRGBA64(x, y, newColor)
	}
}

func main() {
	// Открываем исходное изображение
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

	// Пробуем привести к draw.RGBA64Image
	srcRGBA64, ok := srcImg.(draw.RGBA64Image)
	if !ok {
		fmt.Println("Не удалось преобразовать к draw.RGBA64Image")
		return
	}

	// Создаём новое изображение тех же размеров для хранения результата
	bounds := srcRGBA64.Bounds()
	dstRGBA64 := image.NewRGBA64(bounds) // *RGBA64 реализует интерфейс draw.RGBA64Image

	start := time.Now()

	// Параллельно обрабатываем построчно
	var wg sync.WaitGroup
	height := bounds.Max.Y - bounds.Min.Y
	wg.Add(height)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		go func(row int) {
			defer wg.Done()
			convolveRow(srcRGBA64, dstRGBA64, row)
		}(y)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Время параллельной свёртки: %v\n", elapsed)

	// Сохраняем результат
	outputFile, err := os.Create("output_convolved.png")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, dstRGBA64); err != nil {
		fmt.Println("Ошибка сохранения PNG:", err)
		return
	}

	fmt.Println("Изображение после свёртки сохранено в output_convolved.png")
}
