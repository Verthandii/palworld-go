package rcon

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorcon/rcon"

	"github.com/Verthandii/palworld-go/config"
)

// CmdName https://tech.palworldgame.com/server-commands
type CmdName string

const (
	Shutdown         CmdName = "Shutdown"         // /Shutdown {Seconds} {MessageText}
	DoExit           CmdName = "DoExit"           // /DoExit
	Broadcast        CmdName = "Broadcast"        // /Broadcast {MessageText}
	KickPlayer       CmdName = "KickPlayer"       // /KickPlayer {SteamID}
	BanPlayer        CmdName = "BanPlayer"        // /BanPlayer {SteamID}
	TeleportToPlayer CmdName = "TeleportToPlayer" // /TeleportToPlayer {SteamID}
	TeleportToMe     CmdName = "TeleportToMe"     // /TeleportToMe {SteamID}
	ShowPlayers      CmdName = "ShowPlayers"      // /ShowPlayers
	Info             CmdName = "Info"             // /Info
	Save             CmdName = "Save"             // /Save
)

// Client .
type Client struct {
	c    *config.Config
	conn *rcon.Conn
}

// New .
func New(c *config.Config) (*Client, error) {
	conn, err := rcon.Dial(c.Address, c.AdminPassword)
	if err != nil {
		return nil, err
	}

	return &Client{
		c:    c,
		conn: conn,
	}, nil
}

// Close .
func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Printf("【RCON】关闭连接时发生错误: %v\n", err)
	}
}

// HandleMemoryUsage 发广播 重启维护
func (c *Client) HandleMemoryUsage(threshold float64) {
	c.Broadcast(fmt.Sprintf("broadcast Memory_Is_Above_%v%%", threshold))
	c.Broadcast(c.c.MaintenanceWarningMessage)
	c.Save()
	c.Shutdown("60", "Reboot_In_60_Seconds")
}

// Shutdown 通知玩家且 seconds 秒后关闭服务器
func (c *Client) Shutdown(seconds, message string) {
	c.exec(Shutdown, seconds, message)
}

// DoExit 强制关闭服务器
func (c *Client) DoExit() {
	c.exec(DoExit)
}

// Broadcast 广播
func (c *Client) Broadcast(message string) {
	c.exec(Broadcast, message)
}

// KickPlayer 下线该玩家
func (c *Client) KickPlayer(steamId string) {
	c.exec(KickPlayer, steamId)
}

// BanPlayer BAN该玩家
func (c *Client) BanPlayer(steamId string) {
	c.exec(BanPlayer, steamId)
}

// TeleportToPlayer 传送到该玩家所在地点
func (c *Client) TeleportToPlayer(steamId string) {
	c.exec(TeleportToPlayer, steamId)
}

// TeleportToMe 将该玩家传送到自身所在地点
func (c *Client) TeleportToMe(steamId string) {
	c.exec(TeleportToMe, steamId)
}

// ShowPlayers 展示所有在线玩家
func (c *Client) ShowPlayers() {
	c.exec(ShowPlayers)
}

// Info 展示服务器信息
func (c *Client) Info() {
	c.exec(Info)
}

// Save 存档
func (c *Client) Save() {
	c.exec(Save)
}

func (c *Client) exec(cmd CmdName, args ...string) {
	argStr := strings.Join(args, " ")
	cmdStr := string(cmd)
	if argStr != "" {
		cmdStr = fmt.Sprintf("%s %s", cmd, argStr)
	}

	if _, err := c.conn.Execute(cmdStr); err != nil {
		log.Printf("【RCON】执行命令【%s】时发生错误【%v】\n", cmdStr, err)
	} else {
		log.Printf("【RCON】执行命令【%s】\n", cmdStr)
	}
}
