package infra

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// RawDevice - сырые данные, которые мы получили от ADB CLI API.
type Device struct {
	Serial string
	Status string
}

// ADBClient отвечает только за взаимодействие с системной утилитой adb.
type ADBClient struct {
	adbPath       string // Путь к исполняемому файлу
	ldConsolePath string // Путь к утилите эмулятора
}

// NewADBClient создает новый клиент.
func NewADBClient() *ADBClient {
	return &ADBClient{
		adbPath:       `G:\LDPlayer\LDPlayer9\adb.exe`,
		ldConsolePath: `G:\LDPlayer\LDPlayer9\ldconsole.exe`, //Это путь к утилите LDPlayer9 для определение локальных имен.
	}
}

// GetLDPlayerNames загружает список всех эмуляторов
func (cli *ADBClient) GetLDPlayerNames() map[int]string {
	names := make(map[int]string)

	cmd := exec.Command(cli.ldConsolePath, "list")
	output, err := cmd.Output()

	if err != nil {
		return names
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		vmNumber, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		name := strings.Join(parts[2:], " ")
		names[vmNumber] = name
	}

	return names
}

// GetModel - получить модель устройства
func (cli *ADBClient) GetModel(serial string) string {
	model, err := cli.shellCommand(serial, "getprop ro.product.model")

	if err != nil {
		return "Unknown"
	}

	return model
}

// GetNameByMap получает имя эмулятора по serial, используя заранее загруженную мапу.
func (cli *ADBClient) GetNameByMap(serial string, namesMap map[int]string) string {
	port := extractPort(serial)
	if port == 0 {
		return ""
	}

	// Формула LDPlayer: vmNumber = (port - 5554) / 2
	vmNumber := (port - 5554) / 2

	if name, exists := namesMap[vmNumber]; exists {
		return name
	}

	return ""
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

// extractPort извлекает порт из serial
func extractPort(serial string) int {
	parts := strings.Split(serial, "-")

	if len(parts) < 2 {
		return 0
	}

	port, _ := strconv.Atoi(parts[1])
	return port
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
