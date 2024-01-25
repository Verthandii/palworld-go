# palworld-go

_✨ 适用于palworld的进程守护+强力内存释放+内存不足自动重启服务端 ✨_

## 特别鸣谢+推荐

https://gist.github.com/Bluefissure/b0fcb05c024ee60cad4e23eb55463062

https://github.com/Hoshinonyaruko/palworld-go

## 使用方法

启动后配置（会继续完善）

打开 `\steamcmd\steamapps\common\PalServer\DefaultPalWorldSettings.ini` 配置文件

修改 `RCONEnabled=False`, 把 `False` 改为 `True` 启用 `Rcon`, 修改 `AdminPassword=""` 在 `""` 中设置你的管理员密码

修改完成后保存配置文件, 复制文档全部内容到

`\steamcmd\steamapps\common\PalServer\Pal\Saved\Config\WindowsServer\PalWorldSettings.ini`

保存配置文件

第一次启动`palworld-go-windows-amd64.exe`后会生成`config.json`配置文件

配置文件描述如下

| 配置项                       | 推荐值                                          | 备注                            |
|---------------------------|----------------------------------------------|-------------------------------|
| gamePath                  | "C:\\steamcmd\\steamapps\\common\\PalServer" | 游戏可执行文件路径 PalServer.exe 所处的位置 |
| address                   | "127.0.0.1:25575"                            | 服务器 IP 地址                     |
| rconPort                  | "25575"                                      | RCON 端口号                      |
| adminPassword             | "pwd"                                        | RCON 管理员密码                    |
| processName               | "PalServer.exe"                              | 进程名称 PalServer.exe            |
| checkInterval             | 30                                           | 进程存活检查时间（秒）                   |
| memoryCheckInterval       | 30                                           | 内存占用检测时间（秒）                   |
| memoryUsageThreshold      | 80                                           | 重启阈值（百分比）                     |
| memoryCleanupInterval     | 0                                            | 内存清理时间间隔（秒）                   |
| maintenanceWarningMessage | 服务器即将进行维护,你的存档已保存,请放心,请坐稳扶好,1分钟后重新登录。        | 维护警告消息                        |
| usePerfThreads            | true                                         | 多线程优化                         |

## 兼容性

`windows` 通过了测试，`linux` 有待测试

## 计划

- [x] 服务器进程保活
- [x] 服务器内存清理
- [ ] 通过页面修改游戏配置（如经验值倍率）

## 场景支持

内存不足的时候，通过 `rcon` 通知服务器成员，然后重启服务器

通过调用微软的 `rammap` 释放无用内存，并将有用内存转移至虚拟内存，实现一次释放 50%+ 内存
