package infra

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// RawDevice - сырые данные, которые мы получили от ADB CLI API.
type Device struct {
	Serial string
	Status string
	Model  string
}

// ADBClient отвечает только за взаимодействие с системной утилитой adb.
type ADBClient struct {
	adbPath string // Путь к исполняемому файлу (по умолчанию "adb")
}

// NewADBClient создает новый клиент.
func NewADBClient() *ADBClient {
	return &ADBClient{
		adbPath: `G:\LDPlayer\LDPlayer9\adb.exe`,
	}
}

// shellCommand - приватный хелпер для выполнения команд
func (cli *ADBClient) shellCommand(serial, command string) (string, error) {
	cmd := exec.Command(cli.adbPath, "-s", serial, "shell", command)
	output, err := cmd.Output()

	if err != nil {
		return "", err // Возвращаем ошибку наверх
	}
	return strings.TrimSpace(string(output)), nil
}

// GetModel - получить модель устройства
func (cli *ADBClient) GetModel(serial string) string {
	model, err := cli.shellCommand(serial, "getprop ro.product.model")

	if err != nil {
		return "Unknown"
	}
	return model
}

// GetDevices выполняет команду "adb devices" и возвращает список сырых устройств.
func (cli *ADBClient) GetDevices() ([]Device, error) {
	cmd := exec.Command(cli.adbPath, "devices")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute adb command: %w", err)
	}

	var devices []Device
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.Contains(line, "List of devices") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		devices = append(devices, Device{
			Serial: parts[0],
			Status: parts[1],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading adb output: %w", err)
	}

	return devices, nil
}
