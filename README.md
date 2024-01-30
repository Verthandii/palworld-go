# palworld-go

**一款解决了《幻兽帕鲁》内存泄漏问题的游戏服务器启动器。**

## 特性

### 自动重启

每 **80** 秒检查一次内存，当内存达到阈值时，通过 `RCON` 通知所有在线玩家 **60** 秒后关闭服务器。

服务器关闭之后， **10** 秒后重启服务器以达到清理内存的目的。

![运行效果图1](/pic/palworld_go_1.png)

![运行效果图2](/pic/palworld_go_2.png)

![游戏内效果图](/pic/palworld_reboot.png)

### 自动清理无用内存

在 Windows，Linux 皆可定时清理无用内存，完美解决《幻兽帕鲁》服务器的内存泄漏问题。

**清理内存时需要使用 fork/exec, 此操作需要您使用管理员权限打开终端**

![运行效果图3](/pic/palworld_go_3.png)

#### Windows

通过调用微软的 [RAMMap](https://learn.microsoft.com/en-us/sysinternals/downloads/rammap) 释放无用内存。

#### Linux

通过 Linux 原生命令释放不再需要的缓存和清空交换空间，以回收内存资源。

### 自动备份

每经过配置好的时间间隔，对服务器所有数据进行备份，避免因为死档而导致游戏提前完结撒花。

![运行效果图4](/pic/palworld_go_4.png)

## 使用方法

**前提条件: 按[官方文档](https://tech.palworldgame.com/dedicated-server-guide)安装好所需文件**

### Windows Steam 客户端

1. 搜索 `pal`, 右键 `Palworld Dedicated Server` 如图所示 ![打开目录](/pic/windows_steam_start.png)
2. 将目录粘贴到 `config.json` 的 `gamePath` 中 ![打开目录](/pic/dir.png)
3. 以**管理员权限**打开终端，运行[下载](https://github.com/Verthandii/palworld-go/releases)好的可执行文件

### Windows SteamCMD

1. 将服务器目录粘贴到 `config.json` 的 `gamePath` 中
2. 以**管理员权限**打开终端，运行[下载](https://github.com/Verthandii/palworld-go/releases)好的可执行文件

### Linux SteamCMD

```shell
wget https://mirror.ghproxy.com/https://github.com/Verthandii/palworld-go/releases/download/v0.0.1/palworld-go-linux-amd64
wget https://mirror.ghproxy.com/https://github.com/Verthandii/palworld-go/releases/download/v0.0.1/config.json
vim ./config.json
chmod u+x palworld-go-linux-amd64 
./palworld-go-linux-amd64
```

## 配置文件描述

| 配置项                       | Windows 推荐值                                               | Linux 推荐值                                      | 备注                            |
|---------------------------|-----------------------------------------------------------|------------------------------------------------|-------------------------------|
| gamePath                  | "D:\Program Files (x86)\Steam\steamapps\common\PalServer" | "/home/steam/Steam/steamapps/common/PalServer" | 游戏可执行文件路径 PalServer.exe 所处的位置 |
| address                   | "127.0.0.1:25575"                                         | "127.0.0.1:25575"                              | 服务器地址 + RCON 端口               |
| adminPassword             | "WqB6oY7IzMffxF17Q8La"                                    | "WqB6oY7IzMffxF17Q8La"                         | RCON 管理员密码                    |
| processName               | "PalServer.exe"                                           | "PalServer.sh"                                 | 进程名称                          |
| processCheckInterval      | 5                                                         | 5                                              | 进程存活检查间隔（秒）                   |
| memoryUsageThreshold      | 75                                                        | 75                                             | 重启阈值（百分比）                     |
| memoryCleanupInterval     | 600                                                       | 600                                            | 内存清理间隔（秒）0 表示不清理内存            |
| backupPath                | "D:\Program Files (x86)\backup\PalServer"                 | "/home/steam/backup"                           | 备份路径                          |
| backupInterval            | 1800                                                      | 1800                                           | 备份间隔（秒） 0 表示不备份               |
| maintenanceWarningMessage | Memory_Not_Enough_The_Server_Will_Reboot                  | Memory_Not_Enough_The_Server_Will_Reboot       | 维护警告消息（不支持中文且不支持空格）           |
| usePerfThreads            | true                                                      | true                                           | 多线程优化                         |

## 计划

- [x] 服务器进程保活
- [x] 服务器内存清理
- [x] 自动备份
- [ ] 通过页面修改游戏配置（如经验值倍率）

## 特别鸣谢+推荐

https://gist.github.com/Bluefissure/b0fcb05c024ee60cad4e23eb55463062

https://github.com/Hoshinonyaruko/palworld-go
