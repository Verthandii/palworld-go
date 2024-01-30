//go:build windows

package supervisor

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func (s *supervisor) isAlive() bool {
	out, err := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", s.config.ProcessName)).Output()
	if err != nil {
		log.Printf("【Supervisor】健康检查失败【%v】\n", err)
		return false
	}
	return strings.Contains(string(out), s.config.ProcessName)
}

func (s *supervisor) Close() {

}
