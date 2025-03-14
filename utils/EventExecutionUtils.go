package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// ExecuteCommands 依次执行命令
func ExecuteCommands(commands []string) error {
	for _, cmdStr := range commands {
		err := executeCommand(cmdStr)
		if err != nil {
			return fmt.Errorf("执行命令失败: %s, 错误: %v", cmdStr, err)
		}
	}
	return nil
}

// executeCommand 执行单个命令
func executeCommand(cmdStr string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr) // Windows 使用 cmd.exe
	} else {
		cmd = exec.Command("sh", "-c", cmdStr) // Linux/macOS 使用 sh
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("执行命令失败: %s, 错误: %v", cmdStr, err)
	}
	fmt.Printf("成功执行命令: %s\n", cmdStr)
	return nil
}
