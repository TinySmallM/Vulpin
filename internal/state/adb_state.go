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
	Serial       string       `json:"serial"`        // Серийный номер (уникальный ID)
	Model        string       `json:"model"`         // Модель устройства
	AndroidVer   string       `json:"android_ver"`   // Версия Android
	Status       DeviceStatus `json:"status"`        // Статус подключения
	LastSeen     time.Time    `json:"last_seen"`     // Последнее подключение
	BatteryLevel int          `json:"battery_level"` // Уровень батареи (0-100)
}

func (a *ADBDevice) String() string {
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
func (monitor *ADBMonitor) UpdateDevice(device *ADBDevice) {
	monitor.devicesMu.Lock()
	defer monitor.devicesMu.Unlock()

	device.LastSeen = time.Now()
	monitor.devices[device.Serial] = device
}

// RemoveDevice - удалить устройство
func (monitor *ADBMonitor) RemoveDevice(serial string) {
	monitor.devicesMu.Lock()
	defer monitor.devicesMu.Unlock()

	delete(monitor.devices, serial)
}

// GetDevices - получить все устройства (копия)
func (monitor *ADBMonitor) GetDevices() []*ADBDevice {
	monitor.devicesMu.RLock()
	defer monitor.devicesMu.RUnlock()

	devices := make([]*ADBDevice, 0, len(monitor.devices))
	for _, device := range monitor.devices {
		devices = append(devices, device)
	}
	return devices
}

// GetDevice - получить конкретное устройство
func (monitor *ADBMonitor) GetDevice(serial string) *ADBDevice {
	monitor.devicesMu.RLock()
	defer monitor.devicesMu.RUnlock()

	return monitor.devices[serial]
}

// SetScanning - установить статус сканирования
func (monitor *ADBMonitor) SetScanning(scanning bool) {
	monitor.devicesMu.Lock()
	defer monitor.devicesMu.Unlock()

	monitor.isScanning = scanning
	if scanning {
		monitor.lastScan = time.Now()
	}
}

// IsScanning - проверить, идет ли сканирование
func (monitor *ADBMonitor) IsScanning() bool {
	monitor.devicesMu.RLock()
	defer monitor.devicesMu.RUnlock()

	return monitor.isScanning
}

// GetLastScan - получить время последнего сканирования
func (monitor *ADBMonitor) GetLastScan() time.Time {
	monitor.devicesMu.RLock()
	defer monitor.devicesMu.RUnlock()

	return monitor.lastScan
}

// ClearAll - очистить все устройства (для теста)
func (monitor *ADBMonitor) ClearAll() {
	monitor.devicesMu.Lock()
	defer monitor.devicesMu.Unlock()

	monitor.devices = make(map[string]*ADBDevice)
}
