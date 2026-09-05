package service

import (
	"fmt"
	"time"

	"github.com/TinySmallM/vulpin/internal/infra"
	"github.com/TinySmallM/vulpin/internal/state"
)

type ADBService struct {
	adbClient *infra.ADBClient
	ldClient  *infra.LDClient
	monitor   *state.ADBMonitor
}

func NewADBService(
	adbClient *infra.ADBClient,
	ldClient *infra.LDClient,
	monitor *state.ADBMonitor,
) *ADBService {
	return &ADBService{
		adbClient: adbClient,
		ldClient:  ldClient,
		monitor:   monitor,
	}
}

// ScanDevices делает один проход: опрашивает ADB и обновляет стейт.
func (ser *ADBService) ScanDevices() error {
	ser.monitor.SetScanning(true)
	defer ser.monitor.SetScanning(false)

	rawDevices, err := ser.adbClient.GetDevices()

	if err != nil {
		return fmt.Errorf("infra scan failed: %w", err)
	}

	for _, raw := range rawDevices {
		model := ser.adbClient.GetModel(raw.Serial)
		name := ser.ldClient.GetNameBySerial(raw.Serial)

		device := &state.ADBDevice{
			Serial: raw.Serial,
			Status: mapStatus(raw.Status),
			Name:   name,
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
