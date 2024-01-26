package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/ini.v1"
)

type Config struct {
	GamePath                  string  `json:"gamePath"`                  // 游戏可执行文件路径PalServer.exe所处的位置
	Address                   string  `json:"address"`                   // 服务器 IP 地址
	AdminPassword             string  `json:"adminPassword"`             // RCON 管理员密码
	ProcessName               string  `json:"processName"`               // 进程名称 PalServer.exe
	CheckInterval             int     `json:"checkInterval"`             // 进程存活检查时间（秒）
	RCONPort                  string  `json:"rconPort"`                  // RCON 端口号
	MemoryCheckInterval       int     `json:"memoryCheckInterval"`       // 内存占用检测时间（秒）
	MemoryUsageThreshold      float64 `json:"memoryUsageThreshold"`      // 重启阈值（百分比）
	MemoryCleanupInterval     int     `json:"memoryCleanupInterval"`     // 内存清理时间间隔（秒）
	MaintenanceWarningMessage string  `json:"maintenanceWarningMessage"` // 维护警告消息
	UsePerfThreads            bool    `json:"usePerfThreads"`            // 多线程优化
}

// 默认配置
var defaultConfig = &Config{
	GamePath:                  "",
	Address:                   "127.0.0.1:25575",
	AdminPassword:             "default_password",
	ProcessName:               processName,
	CheckInterval:             30, // 30 秒
	RCONPort:                  "25575",
	MemoryCheckInterval:       30,                                      // 30 秒
	MemoryUsageThreshold:      80,                                      // 80%
	MemoryCleanupInterval:     0,                                       // 内存清理时间间隔，设为半小时（1800秒）0代表不清理
	MaintenanceWarningMessage: "服务器即将进行维护,你的存档已保存,请放心,请坐稳扶好,1分钟后重新登录。", // 默认的维护警告消息
	UsePerfThreads:            true,                                    // 默认启用多线程优化
}

const (
	gameDefaultConfigFile = "DefaultPalWorldSettings.ini"
	configFile            = "config.json"
	processName           = "PalServer.exe"
)

func Init() *Config {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Println("无法读取配置文件, 正在创建默认配置...")
		createDefaultConfig()
		return defaultConfig
	}

	var config *Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("配置解析失败, 正在使用默认配置...")
		return defaultConfig
	}

	fix(config)
	return config
}

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
		gamePath = filepath.Join(currentDir, config.ProcessName)
	}

	if _, err = os.Stat(gamePath); os.IsNotExist(err) {
		log.Printf("当前目录未找到 %s 文件, 请将程序放置在 %s 同目录下\n", config.ProcessName, config.ProcessName)
		os.Exit(1)
	}

	if config.GamePath == "" {
		config.GamePath = currentDir
	}
	if config.Address == "" {
		config.Address = defaultConfig.Address
	}
	if config.AdminPassword == "" {
		log.Printf("配置文件错误: RCON 密码未填写\n")
		os.Exit(1)
	}
	if config.ProcessName == "" {
		config.ProcessName = processName
	}
	if config.CheckInterval < 0 {
		config.CheckInterval = defaultConfig.CheckInterval
	}
	if config.RCONPort == "" {
		config.RCONPort = defaultConfig.RCONPort
	}
	if config.MemoryCheckInterval < 0 {
		config.MemoryCheckInterval = defaultConfig.MemoryCheckInterval
	}
	if config.MemoryUsageThreshold <= 0 {
		config.MemoryUsageThreshold = defaultConfig.MemoryUsageThreshold
	}
	if config.MemoryCleanupInterval < 0 {
		config.MemoryCleanupInterval = 0
	}
	if config.MaintenanceWarningMessage == "" {
		config.MaintenanceWarningMessage = "服务器即将进行维护,你的存档已保存,请放心,请坐稳扶好,1分钟后重新登录。"
	}

	copyGameConfig(config, false)
	configMap := parseGameConfig(config)
	configMap["RCONEnabled"] = "True"
	configMap["AdminPassword"] = config.AdminPassword
	err = os.WriteFile(path.Join(config.GamePath, gameConfigFile), marshalGameConfig(configMap), 0666)
	if err != nil {
		log.Printf("更新游戏配置文件失败【%v】", err)
		os.Exit(1)
	}
}

// copyGameConfig 如果没有游戏配置文件，则将默认的配置文件复制过去
func copyGameConfig(c *Config, force bool) {
	filePath := path.Join(c.GamePath, gameConfigFile)
	dir, _ := path.Split(filePath)

	stat, err := os.Stat(dir)
	if err == nil {
		if !stat.IsDir() {
			log.Printf("游戏目录损坏, 请重新下载游戏")
			os.Exit(1)
		}
	}

	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0777); err != nil {
			log.Printf("创建游戏目录失败【%v】", err)
			os.Exit(1)
		}
	}

	// 目录存在 或者 被新建出来了
	stat, err = os.Stat(filePath)
	if err == nil && !force {
		// 存在游戏配置，不用 copy
		return
	}

	defaultSetting, err := os.ReadFile(path.Join(c.GamePath, "DefaultPalWorldSettings.ini"))
	if err != nil {
		log.Printf("读取游戏默认配置失败【%v】", err)
		os.Exit(1)
	}

	if err = os.WriteFile(filePath, defaultSetting, 0666); err != nil {
		log.Printf("生成游戏配置文件失败【%v】", err)
		os.Exit(1)
	}

	log.Printf("生成游戏配置文件成功\n")
}

func parseGameConfig(c *Config) map[string]string {
	f, err := ini.Load(path.Join(c.GamePath, gameConfigFile))
	if err != nil {
		log.Printf("加载游戏配置文件失败【%v】", err)
		os.Exit(1)
	}
	kvs := f.Section("/Script/Pal.PalGameWorldSettings").Key("OptionSettings").Strings(",")
	if len(kvs) < 2 {
		log.Printf("游戏配置文件损坏, 重新生成")
		copyGameConfig(c, true)
		return parseGameConfig(c)
	}

	res := make(map[string]string)
	for _, kv := range kvs {
		pair := strings.SplitN(kv, "=", 2)
		if len(pair) != 2 {
			log.Printf("游戏配置文件损坏, 重新生成")
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
