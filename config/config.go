package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	GamePath                  string  `json:"gamePath"`                  // 游戏可执行文件路径PalServer.exe所处的位置
	Address                   string  `json:"address"`                   // 服务器 IP 地址
	RCONPort                  string  `json:"rconPort"`                  // RCON 端口号
	AdminPassword             string  `json:"adminPassword"`             // RCON 管理员密码
	ProcessName               string  `json:"processName"`               // 进程名称 PalServer
	CheckInterval             int     `json:"checkInterval"`             // 进程存活检查时间（秒）
	MemoryCheckInterval       int     `json:"memoryCheckInterval"`       // 内存占用检测时间（秒）
	MemoryUsageThreshold      float64 `json:"memoryUsageThreshold"`      // 重启阈值（百分比）
	MemoryCleanupInterval     int     `json:"memoryCleanupInterval"`     // 内存清理时间间隔（秒）
	MaintenanceWarningMessage string  `json:"maintenanceWarningMessage"` // 维护警告消息
}

// 默认配置
var defaultConfig = &Config{
	GamePath:                  "",
	Address:                   "127.0.0.1:25575",
	AdminPassword:             "default_password",
	ProcessName:               "PalServer",
	CheckInterval:             30, // 30 秒
	RCONPort:                  "25575",
	MemoryCheckInterval:       30,                                      // 30 秒
	MemoryUsageThreshold:      80,                                      // 80%
	MemoryCleanupInterval:     0,                                       // 内存清理时间间隔，设为半小时（1800秒）0代表不清理
	MaintenanceWarningMessage: "服务器即将进行维护,你的存档已保存,请放心,请坐稳扶好,1分钟后重新登录。", // 默认的维护警告消息
}

// 配置文件路径
const configFile = "config.json"

func init() {
	var config *Config
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Println("无法读取配置文件, 正在创建默认配置...")
		createDefaultConfig()
		return
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("配置解析失败, 正在使用默认配置...")
		return
	}

	fix(config)
}

func CFG() *Config { return defaultConfig }

func createDefaultConfig() {
	data, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		log.Println("无法创建默认配置文件:", err)
		os.Exit(1)
	}

	err = os.WriteFile(configFile, data, 0666)
	if err != nil {
		log.Println("无法写入默认配置文件:", err)
		os.Exit(1)
	}

	log.Println("默认配置文件已创建:", configFile)
}

func fix(config *Config) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("工作目录获取失败【%v】\n", err)
		os.Exit(1)
	}

	gamePath := config.GamePath
	if gamePath == "" {
		gamePath = filepath.Join(currentDir, "PalServer.exe")
	}

	if _, err = os.Stat(gamePath); os.IsNotExist(err) {
		log.Println("当前目录未找到 PalServer.exe 文件, 请将程序放置在 PalServer.exe 同目录下")
		os.Exit(1)
	}

	if config.GamePath == "" {
		config.GamePath = currentDir
	}
	if config.Address == "" {
		config.Address = defaultConfig.Address
	}
	if config.AdminPassword == "" {
		config.AdminPassword = defaultConfig.AdminPassword
	}
	if config.MemoryCleanupInterval == 0 {
		config.MemoryCleanupInterval = 0
	}
	if config.MaintenanceWarningMessage == "" {
		config.MaintenanceWarningMessage = "服务器即将进行维护,你的存档已保存,请放心,请坐稳扶好,1分钟后重新登录。"
	}
}
