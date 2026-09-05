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
func (ser *ADBService) ScanDevices() error {
	ser.monitor.SetScanning(true)
	defer ser.monitor.SetScanning(false)

	rawDevices, err := ser.client.GetDevices()
	if err != nil {
		return fmt.Errorf("infra scan failed: %w", err)
	}

	// 2. Маппим и кладем в стейт
	for _, raw := range rawDevices {
		// Берем серийник (raw.Serial) и спрашиваем у клиента модель
		model := ser.client.GetModel(raw.Serial)

		device := &state.ADBDevice{
			Serial: raw.Serial,
			Status: mapStatus(raw.Status),
			Model:  model,
		}

		ser.monitor.UpdateDevice(device)
	}

	return nil
}

// StartPeriodicScan запускает фоновый опрос в отдельной горутине.
func (ser *ADBService) StartPeriodicScan(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := ser.ScanDevices(); err != nil {
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
