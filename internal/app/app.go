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
	monitor *state.ADBMonitor
	service *service.ADBService
}

func NewApp() *App {
	// Создаем зависимости один раз
	adbClient := infra.NewADBClient()
	ldClient := infra.NewLDClient(`G:\LDPlayer\LDPlayer9`)
	monitor := state.NewADBMonitor()

	svc := service.NewADBService(adbClient, ldClient, monitor)

	return &App{
		service: svc,
		monitor: monitor,
	}
}

func (app *App) GetMonitor() *state.ADBMonitor {
	return app.monitor
}

func (app *App) StartMonitoring() error {
	fmt.Println("🚀 Запуск мониторинга...")

	if err := app.service.ScanDevices(); err != nil {
		return fmt.Errorf("ошибка сканирования: %w", err)
	}

	app.service.StartPeriodicScan(5 * time.Second)
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
