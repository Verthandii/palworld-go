package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/ini.v1"
)

type Config struct {
	GamePath                  string  `json:"gamePath"`                  // 游戏可执行文件路径 PalServer.exe 所处的位置
	Address                   string  `json:"address"`                   // 服务器地址 + RCON 端口
	AdminPassword             string  `json:"adminPassword"`             // RCON 管理员密码
	ProcessName               string  `json:"processName"`               // 进程名称 PalServer.exe
	ProcessCheckInterval      int     `json:"processCheckInterval"`      // 进程存活检查间隔（秒）
	MemoryUsageThreshold      float64 `json:"memoryUsageThreshold"`      // 重启阈值（百分比）
	MemoryCleanupInterval     int     `json:"memoryCleanupInterval"`     // 内存清理间隔（秒）
	BackupPath                string  `json:"backupPath"`                // 备份路径
	BackupInterval            int     `json:"backupInterval"`            // 备份间隔（秒）
	MaintenanceWarningMessage string  `json:"maintenanceWarningMessage"` // 维护警告消息（不支持中文且不支持空格）
	UsePerfThreads            bool    `json:"usePerfThreads"`            // 多线程优化
}

func (c *Config) PrintLog() {
	log.Printf("【Config】游戏服务器目录【%s】\n", c.GamePath)
	log.Printf("【Config】服务器地址 + RCON 端口【%s】\n", c.Address)
	log.Printf("【Config】RCON 管理员密码【%s】\n", c.AdminPassword)
	log.Printf("【Config】进程名称【%s】\n", c.ProcessName)
	log.Printf("【Config】进程存活检查间隔【%d】秒\n", c.ProcessCheckInterval)
	log.Printf("【Config】重启阈值【%.2f】%%\n", c.MemoryUsageThreshold)
	log.Printf("【Config】内存清理间隔【%d】秒\n", c.MemoryCleanupInterval)
	log.Printf("【Config】备份路径【%s】\n", c.BackupPath)
	log.Printf("【Config】备份间隔【%d】秒\n", c.BackupInterval)
	log.Printf("【Config】维护警告消息【%s】\n", c.MaintenanceWarningMessage)
	log.Printf("【Config】多线程优化【%v】\n", c.UsePerfThreads)
}

// 默认配置
var defaultConfig = &Config{
	GamePath:                  "",
	Address:                   "127.0.0.1:25575",
	AdminPassword:             "WqB6oY7IzMffxF17Q8La",
	ProcessName:               processName,
	ProcessCheckInterval:      5,
	MemoryUsageThreshold:      80,
	MemoryCleanupInterval:     3600,
	BackupPath:                "",
	BackupInterval:            3600,
	MaintenanceWarningMessage: "Memory_Not_Enough_The_Server_Will_Reboot",
	UsePerfThreads:            true,
}

const (
	gameDefaultConfigFile = "DefaultPalWorldSettings.ini"
	configFile            = "config.json"
	processName           = "PalServer.exe"
)

func Init() *Config {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("【Config】无法读取配置文件, 正在创建默认配置...\n")
		createDefaultConfig()
		return defaultConfig
	}

	var config *Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Printf("【Config】配置解析失败, 正在使用默认配置...\n")
		return defaultConfig
	}

	fix(config)
	log.Printf("【Config】配置文件已生效\n")

	return config
}

func createDefaultConfig() {
	data, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		log.Printf("【Config】无法创建默认配置文件【%v】\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(configFile, data, 0666)
	if err != nil {
		log.Printf("【Config】无法写入默认配置文件【%v】\n", err)
		os.Exit(1)
	}

	log.Printf("【Config】默认配置文件创建成功\n")
}

func fix(config *Config) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("【Config】工作目录获取失败【%v】\n", err)
		os.Exit(1)
	}

	gamePath := config.GamePath
	if gamePath == "" {
		gamePath = currentDir
	}
	gamePath = filepath.Join(currentDir, config.ProcessName)

	if _, err = os.Stat(gamePath); os.IsNotExist(err) {
		log.Printf("【Config】当前目录未找到 %s 文件, 请将程序放置在 %s 同目录下\n", config.ProcessName, config.ProcessName)
		os.Exit(1)
	}

	if config.GamePath == "" {
		config.GamePath = currentDir
	}
	if config.Address == "" {
		config.Address = defaultConfig.Address
	}
	if config.AdminPassword == "" {
		log.Printf("【Config】配置文件错误: RCON 密码未填写\n")
		os.Exit(1)
	}
	if config.ProcessName == "" {
		config.ProcessName = defaultConfig.ProcessName
	}
	if config.ProcessCheckInterval <= 0 {
		config.ProcessCheckInterval = defaultConfig.ProcessCheckInterval
	}
	if config.MemoryUsageThreshold <= 0 {
		config.MemoryUsageThreshold = defaultConfig.MemoryUsageThreshold
	}
	if config.MemoryCleanupInterval < 0 {
		config.MemoryCleanupInterval = defaultConfig.MemoryCleanupInterval
	}
	if config.BackupPath == "" {
		config.BackupPath = filepath.Join(currentDir, "backup")
	}
	if config.BackupInterval < 0 {
		config.BackupInterval = defaultConfig.BackupInterval
	}
	if config.MaintenanceWarningMessage == "" {
		config.MaintenanceWarningMessage = defaultConfig.MaintenanceWarningMessage
	}

	_, rconPort, err := net.SplitHostPort(config.Address)
	if err != nil {
		log.Printf("【Config】配置文件错误: address 填写错误\n")
		os.Exit(1)
	}

	copyGameConfig(config, false)
	configMap := parseGameConfig(config)
	configMap["RCONEnabled"] = "True"
	configMap["RCONPort"] = rconPort
	configMap["AdminPassword"] = fmt.Sprintf(`"%s"`, config.AdminPassword)
	err = os.WriteFile(filepath.Join(config.GamePath, gameConfigFile), marshalGameConfig(configMap), 0666)
	if err != nil {
		log.Printf("【Config】更新游戏配置文件失败【%v】\n", err)
		os.Exit(1)
	}
}

// copyGameConfig 如果没有游戏配置文件，则将默认的配置文件复制过去
func copyGameConfig(c *Config, force bool) {
	filePath := filepath.Join(c.GamePath, gameConfigFile)
	dir, _ := path.Split(filePath)

	stat, err := os.Stat(dir)
	if err == nil {
		if !stat.IsDir() {
			log.Printf("【Config】游戏目录损坏, 请重新下载游戏\n")
			os.Exit(1)
		}
	}

	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0777); err != nil {
			log.Printf("【Config】创建游戏目录失败【%v】\n", err)
			os.Exit(1)
		}
	}

	// 目录存在 或者 被新建出来了
	stat, err = os.Stat(filePath)
	if err == nil && !force {
		// 存在游戏配置，不用 copy
		return
	}

	defaultSetting, err := os.ReadFile(filepath.Join(c.GamePath, gameDefaultConfigFile))
	if err != nil {
		log.Printf("【Config】读取游戏默认配置失败【%v】\n", err)
		os.Exit(1)
	}

	if err = os.WriteFile(filePath, defaultSetting, 0666); err != nil {
		log.Printf("【Config】生成游戏配置文件失败【%v】\n", err)
		os.Exit(1)
	}

	log.Printf("【Config】生成游戏配置文件成功\n")
}

func parseGameConfig(c *Config) map[string]string {
	f, err := ini.Load(filepath.Join(c.GamePath, gameConfigFile))
	if err != nil {
		log.Printf("【Config】加载游戏配置文件失败【%v】\n", err)
		os.Exit(1)
	}
	kvs := f.Section("/Script/Pal.PalGameWorldSettings").Key("OptionSettings").Strings(",")
	if len(kvs) < 2 {
		log.Printf("【Config】游戏配置文件损坏, 重新生成\n")
		copyGameConfig(c, true)
		return parseGameConfig(c)
	}

	res := make(map[string]string)
	for _, kv := range kvs {
		pair := strings.SplitN(kv, "=", 2)
		if len(pair) != 2 {
			log.Printf("【Config】游戏配置文件损坏, 重新生成\n")
			copyGameConfig(c, true)
			return parseGameConfig(c)
		}
		k := pair[0]
		v := pair[1]
		if k[0] == '(' {
			k = k[1:]
		}
		if v[len(v)-1] == ')' {
			v = v[:len(v)-1]
		}
		res[k] = v
	}
	return res
}

func marshalGameConfig(configMap map[string]string) []byte {
	arr := make([]string, 0)
	for k, v := range configMap {
		arr = append(arr, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(arr)
	buf := bytes.NewBufferString("[/Script/Pal.PalGameWorldSettings]\nOptionSettings=(")
	buf.WriteString(strings.Join(arr, ","))
	buf.WriteString(")")
	return buf.Bytes()
}
