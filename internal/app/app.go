package app

import (
	"fmt"
	"time"

	"github.com/TinySmallM/vulpin/internal/helper"
	"github.com/TinySmallM/vulpin/internal/infra"
	"github.com/TinySmallM/vulpin/internal/service"
	"github.com/TinySmallM/vulpin/internal/state"
)

type App struct {
	monitor    *state.ADBMonitor
	client     *infra.ADBClient
	adbService *service.ADBService
}

func NewApp() *App {
	// Создаем зависимости один раз
	monitor := state.NewADBMonitor()
	client := infra.NewADBClient()

	return &App{
		monitor:    monitor,
		client:     client,
		adbService: service.NewADBService(client, monitor),
	}
}

func (app *App) GetMonitor() *state.ADBMonitor {
	return app.monitor
}

func (app *App) StartMonitoring() error {
	fmt.Println("🚀 Запуск мониторинга...")

	if err := app.adbService.ScanDevices(); err != nil {
		return fmt.Errorf("ошибка сканирования: %w", err)
	}

	app.adbService.StartPeriodicScan(5 * time.Second)
	fmt.Println("✅ Мониторинг запущен!")
	return nil
}

// PrintDevices - выводит текущие устройства в консоль
func (app *App) PrintDevices() {
	devices := app.monitor.GetDevices()
	fmt.Printf("\n📱 Найдено устройств: %d\n", len(devices))
	for _, device := range devices {
		helper.Inspect(device)
	}
}
