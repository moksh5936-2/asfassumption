package main

import (
	"fmt"
	"io"
	"net/http"
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

func asfModuleVersion(pyPath string) string {
	cmd := exec.Command(pyPath, "-c", "import asf; print(asf.__version__)")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func engineBundleURL(version string) string {
	return fmt.Sprintf("https://github.com/moksh5936-2/asfassumption/releases/download/v%s/asf-python-engine-v%s.tar.gz", version, version)
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

	pyPath := findPython()

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
		printField("Python path (config)", cfg.Engine.PythonPath)
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

	printSection("ASF Python Engine")
	if pyPath == "" {
		printField("Python binary", "not found")
	} else {
		pyVer := pythonVersion(pyPath)
		printField("Python binary", pyPath)
		printField("Python version", pyVer)
	}
	engineDir := asfEngineDir()
	if bundledEngineExists() {
		printField("Engine directory", engineDir+" ✓")
		asfVer := asfModuleVersion(pyPath)
		if asfVer != "" {
			printField("ASF module version", asfVer+" ✓")
		} else {
			printField("ASF module status", "NOT importable (PYTHONPATH issue)")
		}
	} else {
		printField("Engine directory", engineDir+" ✗")
		printField("ASF module status", "not installed")
	}

	editable := findEditableInstall()
	if editable != "" {
		printField("Editable install", editable)
		printField("Remediation", "pip uninstall asf-validator && pip install -e /correct/path")
	}

	printSection("Dependencies")
	checkDep("tesseract", "tesseract --version")
	checkDep("ollama", "ollama --version")
	checkDep("python3", "python3 --version")

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
			cmd := exec.Command(pyPath, "-c", "import asf; print(asf.__file__)")
			cmd.Env = append(os.Environ(), "PYTHONPATH="+asfEngineDir())
			out, err := cmd.Output()
			if err == nil {
				printField("asf module path", strings.TrimSpace(string(out)))
				printField("asf package", "importable ✓")
			} else {
				printField("asf package", "NOT importable")
				// try without PYTHONPATH (for pip-installed engines)
				err2 := exec.Command(pyPath, "-c", "import asf").Run()
				if err2 == nil {
					printField("asf package (pip)", "importable via pip install")
				}
			}
			out2, _ := exec.Command(pyPath, "-m", "pip", "list", "--format=columns").Output()
			for _, line := range strings.Split(string(out2), "\n") {
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

func downloadEngineBundle(version string) error {
	url := engineBundleURL(version)
	fmt.Printf("  Downloading: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	tmpDir, err := os.MkdirTemp("", "asf-engine")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	tarPath := filepath.Join(tmpDir, "engine.tar.gz")
	f, err := os.Create(tarPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return err
	}
	f.Close()

	engineDir := asfEngineDir()
	if err := os.RemoveAll(engineDir); err != nil {
		return fmt.Errorf("remove old engine: %w", err)
	}
	if err := os.MkdirAll(engineDir, 0755); err != nil {
		return err
	}

	cmd := exec.Command("tar", "-xzf", tarPath, "-C", engineDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("extract engine: %w\n%s", err, string(out))
	}

	// Verify the bundled engine path exists
	if !bundledEngineExists() {
		return fmt.Errorf("extracted engine does not contain asf/ package at %s", engineDir)
	}

	return nil
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

	// Check Python engine status
	pyPath := findPython()
	if pyPath == "" {
		fmt.Println("  ✗ Python not found. Install Python 3.9+ and try again.")
		fmt.Println()
		return
	}

	if bundledEngineExists() {
		asfVer := asfModuleVersion(pyPath)
		if asfVer != "" {
			fmt.Println("  ✓ Python ASF engine is already installed (v" + asfVer + ")")
			fmt.Println()
			fmt.Println("  Run 'asf doctor --fix' to reinstall if needed.")
			return
		}
	}

	fmt.Println("  ASF Python engine is missing or broken.")
	fmt.Println()

	version := strings.TrimPrefix(ASFVersion, "v")
	if version == "" {
		version = "1.1.0"
	}

	fmt.Println("  Downloading ASF Python engine v" + version + "...")
	if err := downloadEngineBundle(version); err != nil {
		fmt.Printf("  ✗ Failed: %v\n", err)
		fmt.Println()
		fmt.Println("  Try running install.sh again:")
		fmt.Println("    curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash")
		return
	}
	fmt.Println("  ✓ Engine downloaded and extracted")

	// Verify the import works
	cmd := exec.Command(pyPath, "-c", "import asf; print(asf.__version__)")
	cmd.Env = append(os.Environ(), "PYTHONPATH="+asfEngineDir())
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("  ⚠  Engine extracted but import failed: %v\n", err)
		fmt.Println("  Check that PYTHONPATH includes " + asfEngineDir())
	} else {
		fmt.Printf("  ✓ Python ASF engine v%s importable\n", strings.TrimSpace(string(out)))
	}

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
