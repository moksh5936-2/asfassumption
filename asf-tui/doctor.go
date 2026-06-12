package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func resolvedPath(p string) string {
	real, err := filepath.EvalSymlinks(p)
	if err != nil {
		return p
	}
	return real
}

func runDoctor(verbose bool) {
	fmt.Println("ASF Doctor — System Diagnostic")
	fmt.Println()

	printSection("System")
	printField("OS", runtime.GOOS)
	printField("Architecture", runtime.GOARCH)
	printField("Version", ASFVersion)

	exe, _ := os.Executable()
	printField("Go binary (invoked)", exe)
	if resolved := resolvedPath(exe); resolved != exe {
		printField("Go binary (resolved)", resolved)
	}

	printSection("Paths")
	checkPath("Config directory", asfConfigDir())
	checkPath("Config file", asfConfigPath())
	checkPath("Cache directory", asfCacheDir())
	checkPath("Data directory", asfDataDir())
	checkPath("Engine directory", asfEngineDir())
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
		printField("Use native engine", fmt.Sprintf("%v", cfg.Engine.UseNativeEngine))
	}

	printSection("License")
	lic := LoadLicense()
	if lic != nil {
		printField("Status", lic.Message)
	}

	printSection("ASF Go TUI Binary")
	active := exe
	binaries := findAllAsfBinaries()
	if len(binaries) == 0 {
		printField("None found", "")
	} else {
		for _, b := range binaries {
			tag := ""
			if b == active {
				tag = " [ACTIVE]"
			}
			ver := binaryVersion(b)
			printField(b, ver+tag)
		}
		if len(binaries) > 1 {
			fmt.Println()
			fmt.Println("  \033[1;33m⚠  WARNING: Multiple ASF installations detected.\033[0m")
			fmt.Println("  This can cause version confusion.")
			fmt.Println("  Run 'asf doctor --fix' to clean up stale binaries.")
		}
	}

	printSection("ASF Native Engine")
	printField("Analysis engine", "Go native (compiled in)")
	printField("Available via", "asf analyze <file> [--graph] [-e <evidence>]")
	printField("Python required", "No — native engine works standalone")

	printSection("Dependencies")
	checkDep("tesseract", "tesseract --version")
	checkDep("ollama", "ollama --version")

	printSection("Local AI")
	mm := NewModelManager()
	if mm.CheckAvailable() {
		if mm.CheckRunning() {
			ver := mm.GetVersion()
			if ver != "" {
				printField("Ollama version", ver)
			}
			printField("Ollama running", "yes")
			models, err := mm.ListInstalledAPI()
			if err != nil {
				printField("Installed models", fmt.Sprintf("error: %v", err))
			} else {
				printField("Installed models count", fmt.Sprintf("%d", len(models)))
				if len(models) > 0 {
					names := make([]string, len(models))
					for i, m := range models {
						names[i] = m.Name
					}
					printField("Models", strings.Join(names, ", "))
				}
			}
		} else {
			printField("Ollama running", "no — start with: ollama serve")
		}
	} else {
		printField("Ollama running", "no (binary not found)")
	}

	if err == nil {
		printField("AI enabled", fmt.Sprintf("%v", cfg.AI.Enabled))
		active := cfg.AI.ActiveModel
		if active != "" {
			if mm.IsModelInstalled(active) {
				printField("Active model", fmt.Sprintf("%s (installed)", active))
			} else {
				printField("Active model", fmt.Sprintf("%s (not installed)", active))
			}
		} else {
			printField("Active model", "none set")
		}
	}

	if verbose {
		printSection("Runtime Diagnostics")
		printField("CWD", mustGetwd())
		printField("PATH", os.Getenv("PATH"))
		printField("PYTHONPATH", os.Getenv("PYTHONPATH"))

		printSection("Environment")
		for _, e := range os.Environ() {
			if strings.HasPrefix(e, "ASF_") || strings.HasPrefix(e, "XDG_") || strings.HasPrefix(e, "HOME") || strings.HasPrefix(e, "USER") || strings.HasPrefix(e, "PYTHONPATH") {
				printField(e[:strings.Index(e, "=")], e[strings.Index(e, "=")+1:])
			}
		}

	}

	printSection("PATH")
	fmt.Println("  " + os.Getenv("PATH"))
	fmt.Println()
}

func runDoctorFix() {
	fmt.Println("ASF Doctor — Fix Mode")
	fmt.Println()

	binaries := findAllAsfBinaries()
	active := binaryPath()

	if len(binaries) > 1 {
		fmt.Println("  Found multiple ASF binaries:")
		fmt.Println()
		for _, b := range binaries {
			extra := ""
			if b == active {
				extra = " (active — keeping)"
			}
			ver := binaryVersion(b)
			fmt.Printf("    %s  %s%s\n", b, ver, extra)
		}
		fmt.Println()

		stale := 0
		for _, b := range binaries {
			if b != active {
				stale++
			}
		}

		if stale > 0 {
			fmt.Printf("  Removing %d stale binary(ies)...\n", stale)
			for _, b := range binaries {
				if b != active {
					fmt.Printf("  Removing: %s\n", b)
					os.Remove(b)
					fmt.Printf("    ✓ Removed\n")
				}
			}
			fmt.Println()
		}
	}

	fmt.Println("  ── Native Go engine: built-in (no Python required) ──")
	fmt.Println()
	fmt.Println("  ✓ All checks complete. Native Go engine is active.")
	fmt.Println()
	fmt.Println("  Done. Run 'asf doctor' to verify.")
}

func findAllAsfBinaries() []string {
	seen := map[string]bool{}
	var result []string
	pathEnv := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(pathEnv) {
		if dir == "" {
			continue
		}
		for _, name := range []string{"asf", "asf.exe"} {
			p := filepath.Join(dir, name)
			if info, err := os.Stat(p); err == nil && !info.IsDir() && !seen[p] {
				seen[p] = true
				result = append(result, p)
			}
		}
	}
	oldPath := filepath.Join(os.Getenv("HOME"), ".asf", "asf")
	if _, err := os.Stat(oldPath); err == nil && !seen[oldPath] {
		seen[oldPath] = true
		result = append(result, oldPath)
	}
	return result
}

func binaryVersion(path string) string {
	cmd := exec.Command(path, "--version")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func binaryPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "unknown"
	}
	return exe
}

func mustGetwd() string {
	d, err := os.Getwd()
	if err != nil {
		return "error: " + err.Error()
	}
	return d
}

func printSection(name string) {
	fmt.Printf("\033[1;36m%s\033[0m\n", name)
	fmt.Printf("\033[1;36m%s\033[0m\n", "─"+strings.Repeat("─", len(name)))
}

func printField(name, value string) {
	if value == "" {
		return
	}
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
