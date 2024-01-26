//go:build !windows

package supervisor

import (
	"log"
	"os/exec"
	"strings"
)

func (s *supervisor) isAlive() bool {
	out, err := exec.Command("pgrep", "-f", s.config.ProcessName).Output()
	if err != nil {
		log.Printf("【Supervisor】健康检查失败【%v】\n", err)
		return false
	}

	// 检查输出结果，如果结果不为空，则至少存在一个进程
	output := strings.TrimSpace(string(out))
	return output != ""
}
