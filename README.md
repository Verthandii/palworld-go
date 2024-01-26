# palworld-go

_✨ 适用于palworld的进程守护+强力内存释放+内存不足自动重启服务端 ✨_

## 场景支持

### 自动重启

当内存打到阈值时，通过 `rcon` 通知所有在线玩家，然后重启服务器以清理内存

### 自动清理无用内存（仅支持 windows）

通过调用微软的 `rammap` 释放无用内存，完美解决《幻兽帕鲁》服务器的内存泄漏问题

## 特别鸣谢+推荐

https://gist.github.com/Bluefissure/b0fcb05c024ee60cad4e23eb55463062

https://github.com/Hoshinonyaruko/palworld-go

## 使用方法

### Windows Steam 客户端

1. 搜索 `pal`, 右键 `Palworld Dedicated Server` 如图所示。![打开目录](/pic/windows_steam_start.png)
2. 复制目录到 `config.json` 的 `gamePath` 中，再按需修改 `config.json` 中其他配置项
3. 将[下载](https://github.com/Verthandii/palworld-go/release)好的可执行文件移动到此目录下，在命令行中运行即可

### Windows SteamCMD

### Linux SteamCMD

## 配置文件描述

| 配置项                       | 推荐值                                                       | 备注                            |
|---------------------------|-----------------------------------------------------------|-------------------------------|
| gamePath                  | "D:\Program Files (x86)\Steam\steamapps\common\PalServer" | 游戏可执行文件路径 PalServer.exe 所处的位置 |
| address                   | "127.0.0.1:25575"                                         | 服务器 IP 地址                     |
| adminPassword             | "WqB6oY7IzMffxF17Q8La"                                    | RCON 管理员密码                    |
| processName               | "PalServer.exe"                                           | 进程名称 PalServer.exe            |
| checkInterval             | 5                                                         | 进程存活检查时间（秒）                   |
| rconPort                  | "25575"                                                   | RCON 端口号                      |
| memoryCheckInterval       | 70                                                        | 内存占用检测时间（秒）                   |
| memoryUsageThreshold      | 80                                                        | 重启阈值（百分比）                     |
| memoryCleanupInterval     | 3600                                                      | 内存清理时间间隔（秒）                   |
| maintenanceWarningMessage | Memory_Not_Enough_The_Server_Will_Reboot                  | 维护警告消息（不支持中文且不支持空格）           |
| usePerfThreads            | true                                                      | 多线程优化                         |

## 兼容性

`windows` 通过了测试，`linux` 有待测试

## 计划

- [x] 服务器进程保活
- [x] 服务器内存清理
- [ ] 通过页面修改游戏配置（如经验值倍率）

