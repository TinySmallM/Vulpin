package infra

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ldConfig struct {
	PlayerName string `json:"statusSettings.playerName"`
}

type LDClient struct {
	basePath    string
	mu          sync.RWMutex
	names       map[int]string
	cacheExpiry time.Time
}

func NewLDClient(ldPlayerBasePath string) *LDClient {
	return &LDClient{
		basePath: filepath.Join(ldPlayerBasePath, "vms", "config"),
		names:    make(map[int]string),
	}
}

// GetNames возвращает мапу: ADB порт -> имя эмулятора
func (cli *LDClient) GetNames() map[int]string {
	if cli.isCacheValid() {
		return cli.getCached()
	}
	return cli.refresh()
}

func (cli *LDClient) isCacheValid() bool {
	cli.mu.RLock()
	defer cli.mu.RUnlock()
	return time.Now().Before(cli.cacheExpiry)
}

func (cli *LDClient) getCached() map[int]string {
	cli.mu.RLock()
	defer cli.mu.RUnlock()
	return cli.names
}

func (cli *LDClient) refresh() map[int]string {
	cli.mu.Lock()
	defer cli.mu.Unlock()

	if time.Now().Before(cli.cacheExpiry) {
		return cli.names
	}

	files, _ := filepath.Glob(filepath.Join(cli.basePath, "leidian*.config"))
	names := make(map[int]string)

	for _, file := range files {
		index, err := extractVMIndex(file)

		if err != nil {
			continue
		}

		cfg, err := readLDConfig(file)

		if err != nil {
			continue
		}

		if cfg.PlayerName != "" {
			port := 5554 + (index * 2)
			names[port] = cfg.PlayerName
		}
	}

	cli.names = names
	cli.cacheExpiry = time.Now().Add(30 * time.Second)
	return names
}

// GetNameBySerial получает имя по serial (например, "emulator-5638")
func (cli *LDClient) GetNameBySerial(serial string) string {
	port := extractPort(serial)
	if port == 0 {
		return ""
	}

	names := cli.GetNames()
	if name, exists := names[port]; exists {
		return name
	}
	return ""
}

func extractVMIndex(filePath string) (int, error) {
	fileName := filepath.Base(filePath)
	name := strings.TrimSuffix(strings.TrimPrefix(fileName, "leidian"), ".config")
	return strconv.Atoi(name)
}

func readLDConfig(path string) (*ldConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ldConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func extractPort(serial string) int {
	parts := strings.Split(serial, "-")
	if len(parts) < 2 {
		return 0
	}
	port, _ := strconv.Atoi(parts[1])
	return port
}
