package infra

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// RawDevice - сырые данные, которые мы получили от ADB CLI.
type RawDevice struct {
	Serial string
	Status string
}

// ADBClient отвечает только за взаимодействие с системной утилитой adb.
type ADBClient struct {
	adbPath string // Путь к исполняемому файлу (по умолчанию "adb")
}

// NewADBClient создает новый клиент.
// В будущем сюда можно будет добавить настройку пути к adb или таймаутов.
func NewADBClient() *ADBClient {
	return &ADBClient{
		adbPath: "adb", // По умолчанию ищем adb в PATH системы
	}
}

// GetDevices выполняет команду "adb devices" и возвращает список сырых устройств.
func (сlient *ADBClient) GetDevices() ([]RawDevice, error) {
	// Выполняем команду
	cmd := exec.Command(сlient.adbPath, "devices")

	output, err := cmd.Output()
	if err != nil {
		// Если команда вернула ошибку (например, adb не установлен или демон умер)
		return nil, fmt.Errorf("failed to execute adb command: %w", err)
	}

	var devices []RawDevice
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.Contains(line, "List of devices") {
			continue
		}

		// strings.Fields разбивает строку по любым пробельным символам (пробелы, табы)
		// Пример: "emulator-5554    device" -> ["emulator-5554", "device"]
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		devices = append(devices, RawDevice{
			Serial: parts[0],
			Status: parts[1],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading adb output: %w", err)
	}

	return devices, nil
}
