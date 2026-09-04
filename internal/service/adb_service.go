package service

import (
	"fmt"
	"time"

	"github.com/TinySmallM/vulpin/internal/infra"
	"github.com/TinySmallM/vulpin/internal/state"
)

type ADBService struct {
	client  *infra.ADBClient
	monitor *state.ADBMonitor
}

func NewADBService(client *infra.ADBClient, monitor *state.ADBMonitor) *ADBService {
	return &ADBService{
		client:  client,
		monitor: monitor,
	}
}

// ScanDevices делает один проход: опрашивает ADB и обновляет стейт.
func (k *ADBService) ScanDevices() error {
	k.monitor.SetScanning(true)
	defer k.monitor.SetScanning(false)

	// 1. Получаем сырые данные из инфры
	rawDevices, err := k.client.GetDevices()
	if err != nil {
		return fmt.Errorf("infra scan failed: %w", err)
	}

	// 2. Маппим и кладем в стейт
	for _, raw := range rawDevices {
		device := &state.ADBDevice{
			Serial: raw.Serial,
			Status: mapStatus(raw.Status),
			Model:  "Unknown", // Пока заглушка, потом научимся тянуть через adb shell
		}
		k.monitor.UpdateDevice(device)
	}

	return nil
}

// StartPeriodicScan запускает фоновый опрос в отдельной горутине.
func (k *ADBService) StartPeriodicScan(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := k.ScanDevices(); err != nil {
				fmt.Printf("⚠️ Ошибка сканирования: %v\n", err)
			}
		}
	}()
}

// mapStatus переводит сырой статус ADB в наш типизированный статус.
func mapStatus(raw string) state.DeviceStatus {
	switch raw {
	case "device":
		return state.StatusOnline
	case "unauthorized":
		return state.StatusUnauthorized
	case "offline":
		return state.StatusOffline
	default:
		return state.StatusConnecting
	}
}
