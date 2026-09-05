// internal/state/adb_state.go
package state

import (
	"sync"
	"time"
)

// DeviceStatus - статус ADB устройства
type DeviceStatus string

const (
	StatusOnline       DeviceStatus = "Online"
	StatusOffline      DeviceStatus = "Offline"
	StatusUnauthorized DeviceStatus = "Unauthorized"
	StatusConnecting   DeviceStatus = "Connecting..."
)

// ADBDevice - информация об Android устройстве
type ADBDevice struct {
	Serial   string       `json:"serial"`    // Серийный номер (уникальный ID)
	Model    string       `json:"model"`     // Модель устройства
	Name     string       `json:"name"`      // Имя эмулятора
	Status   DeviceStatus `json:"status"`    // Статус подключения
	LastSeen time.Time    `json:"last_seen"` // Последнее подключение
}

func (d *ADBDevice) String() string {
	panic("unimplemented")
}

// ADBMonitor - глобальное состояние ADB монитора
type ADBMonitor struct {
	// devicesMu защищает поля devices, isScanning и lastScan
	// от одновременного чтения и записи разными горутинами.
	devicesMu sync.RWMutex

	devices    map[string]*ADBDevice
	isScanning bool
	lastScan   time.Time
}

func NewADBMonitor() *ADBMonitor {
	return &ADBMonitor{
		devices: make(map[string]*ADBDevice),
	}
}

// --- Методы для работы с устройствами ---

// UpdateDevice - добавить или обновить устройство
func (mon *ADBMonitor) UpdateDevice(device *ADBDevice) {
	mon.devicesMu.Lock()
	defer mon.devicesMu.Unlock()

	device.LastSeen = time.Now()
	mon.devices[device.Serial] = device
}

// RemoveDevice - удалить устройство
func (mon *ADBMonitor) RemoveDevice(serial string) {
	mon.devicesMu.Lock()
	defer mon.devicesMu.Unlock()

	delete(mon.devices, serial)
}

// GetDevices - получить все устройства (копия)
func (mon *ADBMonitor) GetDevices() []*ADBDevice {
	mon.devicesMu.RLock()
	defer mon.devicesMu.RUnlock()

	devices := make([]*ADBDevice, 0, len(mon.devices))
	for _, device := range mon.devices {
		devices = append(devices, device)
	}
	return devices
}

// GetDevice - получить конкретное устройство
func (mon *ADBMonitor) GetDevice(serial string) *ADBDevice {
	mon.devicesMu.RLock()
	defer mon.devicesMu.RUnlock()

	return mon.devices[serial]
}

// SetScanning - установить статус сканирования
func (mon *ADBMonitor) SetScanning(scanning bool) {
	mon.devicesMu.Lock()
	defer mon.devicesMu.Unlock()

	mon.isScanning = scanning
	if scanning {
		mon.lastScan = time.Now()
	}
}

// IsScanning - проверить, идет ли сканирование
func (mon *ADBMonitor) IsScanning() bool {
	mon.devicesMu.RLock()
	defer mon.devicesMu.RUnlock()

	return mon.isScanning
}

// GetLastScan - получить время последнего сканирования
func (mon *ADBMonitor) GetLastScan() time.Time {
	mon.devicesMu.RLock()
	defer mon.devicesMu.RUnlock()

	return mon.lastScan
}

// ClearAll - очистить все устройства (для теста)
func (mon *ADBMonitor) ClearAll() {
	mon.devicesMu.Lock()
	defer mon.devicesMu.Unlock()

	mon.devices = make(map[string]*ADBDevice)
}
