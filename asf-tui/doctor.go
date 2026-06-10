package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func runDoctor() {
	fmt.Println("ASF Doctor — System Diagnostic")
	fmt.Println()

	printSection("System")
	printField("OS", runtime.GOOS)
	printField("Architecture", runtime.GOARCH)
	printField("Binary", binaryPath())
	printField("Version", ASFVersion)

	printSection("Paths")
	checkPath("Config directory", asfConfigDir())
	checkPath("Config file", asfConfigPath())
	checkPath("Cache directory", asfCacheDir())
	checkPath("Data directory", asfDataDir())
	checkPath("License file", asfLicensePath())

	printSection("Permissions")
	checkPerm("Config dir write", asfConfigDir())
	checkPerm("Cache dir write", asfCacheDir())

	printSection("Configuration")
	cfg, err := LoadConfig(asfConfigPath())
	if err != nil {
		printField("Config load", fmt.Sprintf("error: %v", err))
	} else {
		printField("Theme", cfg.Appearance.Theme)
		printField("Analysis depth", cfg.Analysis.Depth)
		printField("Export format", cfg.Output.Default)
		printField("Export directory", cfg.Output.Directory)
		printField("AI enabled", fmt.Sprintf("%v", cfg.AI.Enabled))
		printField("Active model", cfg.AI.ActiveModel)
	}

	printSection("License")
	lic := LoadLicense()
	if lic != nil {
		printField("Status", lic.Message)
	}

	printSection("Python Engine")
	pyPath := findPython()
	if pyPath == "" {
		printField("Python", "not found")
	} else {
		printField("Python binary", pyPath)
	}
	enginePath := findAsfEngine()
	if enginePath == "" {
		printField("ASF engine", "not found")
	} else {
		printField("ASF engine", enginePath)
	}

	printSection("Dependencies")
	checkDep("tesseract", "tesseract --version")
	checkDep("ollama", "ollama --version")
	checkDep("python3", "python3 --version")

	printSection("PATH")
	fmt.Println("  " + os.Getenv("PATH"))
	fmt.Println()
}

func binaryPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "unknown"
	}
	return exe
}

func findPython() string {
	for _, name := range []string{"python3", "python"} {
		path, err := exec.LookPath(name)
		if err == nil {
			return path
		}
	}
	return ""
}

func findAsfEngine() string {
	exe, err := os.Executable()
	if err == nil {
		binDir := filepath.Dir(exe)
		enginePath := filepath.Join(binDir, "engine")
		if _, err := os.Stat(enginePath); err == nil {
			return enginePath
		}
		enginePath = filepath.Join(binDir, "asf")
		if _, err := os.Stat(enginePath + ".py"); err == nil {
			return enginePath + ".py"
		}
	}
	if _, err := exec.LookPath("asf"); err == nil {
		return "asf (in PATH)"
	}
	py := findPython()
	if py != "" {
		cmd := exec.Command(py, "-m", "asf.cli.main", "--help")
		if cmd.Run() == nil {
			return fmt.Sprintf("python3 -m asf.cli.main (%s)", py)
		}
	}
	return ""
}

func printSection(name string) {
	fmt.Printf("\033[1;36m%s\033[0m\n", name)
	fmt.Printf("\033[1;36m%s\033[0m\n", "─"+string(rune(0x2500)))
}

func printField(name, value string) {
	fmt.Printf("  \033[1m%-18s\033[0m %s\n", name+":", value)
}

func checkPath(name, path string) {
	info, err := os.Stat(path)
	if err == nil {
		mode := "exists"
		if info.IsDir() {
			mode = "directory"
		}
		printField(name, fmt.Sprintf("%s (%s)", path, mode))
	} else if os.IsNotExist(err) {
		printField(name, fmt.Sprintf("%s (not found)", path))
	} else {
		printField(name, fmt.Sprintf("%s (%v)", path, err))
	}
}

func checkPerm(name, dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		printField(name, fmt.Sprintf("cannot create: %v", err))
		return
	}
	tmpFile := filepath.Join(dir, ".asf_write_test")
	err = os.WriteFile(tmpFile, []byte("test"), 0644)
	if err != nil {
		printField(name, fmt.Sprintf("not writable: %v", err))
		return
	}
	os.Remove(tmpFile)
	printField(name, "writable")
}

func checkDep(name, cmd string) {
	parts := splitCmd(cmd)
	if len(parts) == 0 {
		printField(name, "unknown")
		return
	}
	_, err := exec.LookPath(parts[0])
	if err != nil {
		printField(name, "not found")
		return
	}
	c := exec.Command(parts[0], parts[1:]...)
	if c.Run() == nil {
		printField(name, "available")
	} else {
		printField(name, "found but error")
	}
}

func splitCmd(cmd string) []string {
	var parts []string
	var buf []rune
	inQuote := false
	for _, ch := range cmd {
		if ch == ' ' && !inQuote {
			if len(buf) > 0 {
				parts = append(parts, string(buf))
				buf = nil
			}
		} else {
			buf = append(buf, ch)
		}
	}
	if len(buf) > 0 {
		parts = append(parts, string(buf))
	}
	return parts
}
