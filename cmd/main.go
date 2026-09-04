package main

import (
	"fmt"

	"github.com/TinySmallM/vulpin/internal/app"
)

func main() {
	application := app.NewApp()

	if err := application.StartMonitoring(); err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}

	// Выводим данные сразу после запуска
	application.PrintDevices()

	// Держим программу запущенной
	select {}
}
