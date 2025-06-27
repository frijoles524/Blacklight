package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procFreeConsole      = kernel32.NewProc("FreeConsole")
	procAttachConsole    = kernel32.NewProc("AttachConsole")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")
	user32               = syscall.NewLazyDLL("user32.dll")
	procShowWindow       = user32.NewProc("ShowWindow")
)

const (
	ATTACH_PARENT_PROCESS = ^uintptr(0)
	SW_HIDE               = 0
)

func DetachConsole() error {
	_, _, err := procFreeConsole.Call()
	if err != nil && err.Error() != "The operation completed successfully." {
		return fmt.Errorf("FreeConsole failed: %v", err)
	}
	hwnd, _, _ := procGetConsoleWindow.Call()
	if hwnd != 0 {
		procShowWindow.Call(hwnd, SW_HIDE)
	}
	return nil
}

func InstallRequirements(requirementsFile string) error {
	defer ShutdownPython()
	file, err := os.Open(requirementsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fmt.Printf("Installing %s ...\n", line)
		err := RunString(fmt.Sprintf("import pip._internal; pip._internal.main(['install', '%s', '--target=python312runtime/Lib/site-packages'])", line))
		if err != nil {
			return fmt.Errorf("failed to install %s: %w", line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func ContainsGuiFlag(slice []string) bool {
	for i, v := range slice {
		if v == "-gui" {
			slice[i] = ""
			return true
		}
	}
	return false
}
