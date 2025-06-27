package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func InstallRequirements(requirementsFile string) error {
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
		err := RunString(fmt.Sprintf("import pip; pip.main(['install', '%s', '--target=python312runtime/Lib/site-packages'])", line))
		if err != nil {
			return fmt.Errorf("failed to install %s: %w", line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
