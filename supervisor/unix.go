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
		if err1, ok := err.(*exec.ExitError); ok && err1.ExitCode() == 1 {
			// 命令执行成功，但没有找到任何匹配的进程
			return false
		}

		log.Printf("【Supervisor】健康检查失败【%v】\n", err)
		return false
	}

	// 检查输出结果，如果结果不为空，则至少存在一个进程
	output := strings.TrimSpace(string(out))
	return output != ""
}

func (s *supervisor) Close() {
}
