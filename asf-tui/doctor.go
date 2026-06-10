package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func runDoctor(verbose bool) {
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
		printField("Python path (config)", cfg.Engine.PythonPath)
	}

	printSection("License")
	lic := LoadLicense()
	if lic != nil {
		printField("Status", lic.Message)
	}

	printSection("Python Engine Error")
	pyPath := findPython()
	if pyPath == "" {
		printField("Python binary", "not found")
	} else {
		pyVer := pythonVersion(pyPath)
		printField("Python binary", pyPath)
		printField("Python version", pyVer)
	}
	enginePath := findAsfEngine()
	if enginePath == "" {
		printField("ASF engine", "not found")
	} else {
		printField("ASF engine", enginePath)
	}

	editable := findEditableInstall()
	if editable != "" {
		printField("WARNING", "")
		printField("Editable install", editable)
		printField("Remediation", "pip uninstall asf-validator && pip install -e /correct/path")
	}

	printSection("Dependencies")
	checkDep("tesseract", "tesseract --version")
	checkDep("ollama", "ollama --version")
	checkDep("python3", "python3 --version")

	binaries := findAllAsfBinaries()
	printSection("ASF Binaries in PATH")
	if len(binaries) == 0 {
		printField("None found", "")
	} else {
		active := binaryPath()
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

	if verbose {
		printSection("Runtime Diagnostics")
		exe, _ := os.Executable()
		printField("Executable", exe)
		printField("CWD", mustGetwd())
		printField("PATH", os.Getenv("PATH"))

		printSection("Environment")
		for _, e := range os.Environ() {
			if strings.HasPrefix(e, "ASF_") || strings.HasPrefix(e, "XDG_") || strings.HasPrefix(e, "HOME") || strings.HasPrefix(e, "USER") {
				printField(e[:strings.Index(e, "=")], e[strings.Index(e, "=")+1:])
			}
		}

		printSection("Python Discovery Candidates")
		for _, c := range pythonCandidates(cfg) {
			ok := "✓"
			info, err := os.Stat(c)
			if err != nil || info.IsDir() {
				ok = "✗"
			}
			printField(ok, c)
		}

		printSection("Python Package Check")
		if pyPath != "" {
			err2 := exec.Command(pyPath, "-c", "import asf").Run()
			if err2 == nil {
				printField("asf package", "importable")
			} else {
				printField("asf package", "NOT importable")
			}
			out, _ := exec.Command(pyPath, "-m", "pip", "list", "--format=columns").Output()
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "asf") || strings.Contains(line, "asf-validator") {
					printField("    pip", strings.TrimSpace(line))
				}
			}
		}

		printSection("Editable Install Files")
		findEditableInstallsVerbose()
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

	if len(binaries) <= 1 {
		fmt.Println("  No duplicate binaries found. Nothing to clean up.")
		return
	}

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

	if stale == 0 {
		fmt.Println("  No stale binaries to remove.")
		return
	}

	fmt.Printf("  Would remove %d stale binary(ies).\n", stale)
	fmt.Println()

	for _, b := range binaries {
		if b != active {
			fmt.Printf("  Removing: %s\n", b)
			os.Remove(b)
			fmt.Printf("    ✓ Removed\n")
		}
	}

	fmt.Println()
	fmt.Println("  Done. Only the active binary remains.")
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

func findPython() string {
	for _, name := range []string{"python3", "python"} {
		path, err := exec.LookPath(name)
		if err == nil {
			return path
		}
	}
	return ""
}

func pythonVersion(path string) string {
	cmd := exec.Command(path, "--version")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
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

func findEditableInstall() string {
	py := findPython()
	if py == "" {
		return ""
	}
	out, err := exec.Command(py, "-m", "pip", "list", "--format=columns").Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "asf") {
			fields := strings.Fields(line)
			if len(fields) >= 3 && fields[1] == "0.1.0" {
				out2, _ := exec.Command(py, "-m", "pip", "show", "-f", fields[0]).Output()
				for _, l := range strings.Split(string(out2), "\n") {
					if strings.Contains(l, "Editable") {
						parts := strings.SplitN(l, ":", 2)
						if len(parts) == 2 {
							return strings.TrimSpace(parts[1])
						}
					}
					if strings.Contains(l, "Location") && strings.Contains(l, "/Users/") {
						parts := strings.SplitN(l, ":", 2)
						if len(parts) == 2 {
							loc := strings.TrimSpace(parts[1])
							return loc
						}
					}
				}
			}
		}
	}
	return ""
}

func findEditableInstallsVerbose() {
	py := findPython()
	if py == "" {
		return
	}
	searchDirs := []string{
		filepath.Join(os.Getenv("HOME"), ".asf"),
	}
	for _, p := range filepath.SplitList(os.Getenv("PATH")) {
		if p != "" {
			searchDirs = append(searchDirs, p)
		}
	}
	// Check for .egg-link files
	for _, d := range searchDirs {
		if d == "" {
			continue
		}
		found := false
		filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if strings.HasSuffix(info.Name(), ".egg-link") || strings.HasSuffix(info.Name(), ".pth") {
				if !found {
					printField("    Found", path)
					found = true
				}
			}
			return nil
		})
	}
	// Check for editable install marker
	out, _ := exec.Command(py, "-c", `
import site, sys
for d in site.getsitepackages():
    import glob
    for f in glob.glob(d + "/*.egg-link") + glob.glob(d + "/*.pth"):
        print(f)
`).Output()
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			printField("    Site-pkg", line)
		}
	}
}

func pythonCandidates(cfg *Config) []string {
	var candidates []string
	if cfg != nil && cfg.Engine.PythonPath != "" {
		candidates = append(candidates, cfg.Engine.PythonPath+" (config)")
	}
	exe, err := os.Executable()
	binDir := ""
	if err == nil {
		binDir = filepath.Dir(exe)
	}
	if binDir != "" {
		candidates = append(candidates,
			filepath.Join(binDir, "engine", "bin", "python3"),
			filepath.Join(binDir, "engine", "bin", "python"),
			filepath.Join(binDir, "asf"),
		)
	}
	candidates = append(candidates, filepath.Join(asfDataDir(), "venv", "bin", "python3"))
	candidates = append(candidates, filepath.Join(asfDataDir(), "venv", "bin", "python"))
	for _, name := range []string{"asf", "asf.py", "asf-cli"} {
		if p, err2 := exec.LookPath(name); err2 == nil {
			candidates = append(candidates, p+" (PATH)")
		}
	}
	for _, name := range []string{"python3", "python"} {
		if p, err2 := exec.LookPath(name); err2 == nil {
			candidates = append(candidates, p+" (PATH)")
		}
	}
	candidates = append(candidates, "python3 (fallback)")
	return candidates
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
