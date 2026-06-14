package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ExitSuccess      = 0
	ExitGeneralError = 1
	ExitInvalidCmd   = 2
	ExitAnalysisErr  = 4
	ExitExportErr    = 6
	ExitLicenseErr   = 7
)

func main() {
	if err := initLogger(); err != nil {
		asfLog = log.New(io.Discard, "[asf] ", log.Ldate|log.Ltime|log.Lshortfile)
		debugLog = log.New(io.Discard, "[asf-debug] ", log.Ltime|log.Lshortfile)
	}
	asfLog.Printf("ASF0 v%s starting", ASFVersion)

	args := os.Args[1:]
	if len(args) == 0 {
		launchTUI(nil)
		return
	}

	mainCLI(args)
}

func isStrictCLI() bool {
	return os.Getenv("ASF_STRICT") == "1" || os.Getenv("ASF_NONINTERACTIVE") == "1"
}

func mainCLI(args []string) {
	strict := false
	var positional []string
	for _, a := range args {
		if a == "--strict" {
			strict = true
			os.Setenv("ASF_STRICT", "1")
		} else {
			positional = append(positional, a)
		}
	}
	args = positional

	if len(args) == 0 {
		if strict {
			fmt.Fprintf(os.Stderr, "Error: no command in strict mode\n")
			os.Exit(ExitInvalidCmd)
		}
		launchTUI(nil)
		return
	}

	switch args[0] {
	case "--version", "-v":
		fmt.Printf("ASF0 v%s\n", ASFVersion)
		if msg := VersionCheckMessage(); msg != "" {
			fmt.Println(msg)
		}
	case "--version-check":
		if msg := VersionCheckMessage(); msg != "" {
			fmt.Println(msg)
		} else {
			fmt.Printf("ASF0 v%s is up to date.\n", ASFVersion)
		}
	case "--license":
		l := LoadLicense()
		if l != nil && l.Valid {
			fmt.Printf("License: %s\n", l.Message)
			return
		}
		fmt.Println("No valid license found.")
		fmt.Printf("Place your license key in %s\n", asfLicensePath())
		os.Exit(ExitLicenseErr)
	case "doctor", "--doctor", "diagnose":
		verbose := false
		fix := false
		for _, a := range args[1:] {
			switch a {
			case "--verbose", "-v":
				verbose = true
			case "--fix":
				fix = true
			}
		}
		if fix {
			runDoctorFix()
		} else {
			runDoctor(verbose)
		}
	case "analyze":
		strict = true
		os.Setenv("ASF_STRICT", "1")
		runAnalyzeCLI(args[1:])
	case "export":
		runExportCLI(args[1:])
	case "config":
		runConfigCLI(args[1:])
	case "completion":
		runCompletionCLI(args[1:])
	case "--help", "-h":
		printUsage()
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", args[0])
		fmt.Fprintf(os.Stderr, "Run 'asf --help' for usage.\n")
		os.Exit(ExitInvalidCmd)
	}
}

func launchTUI(cfg *Config) {
	if cfg == nil {
		var err error
		cfg, err = LoadConfig(ConfigPath())
		if err != nil {
			def := DefaultConfig()
			cfg = &def
		}
	}
	asfLog.Printf("config path: %s", ConfigPath())

	if err := ensureRuntimeDirs(); err != nil {
		debugLog.Printf("runtime dirs: %v", err)
	}

	if isStrictCLI() {
		fmt.Fprintf(os.Stderr, "Error: cannot launch TUI in strict CLI mode\n")
		fmt.Fprintf(os.Stderr, "Use 'asf analyze' or 'asf export' for non-interactive use.\n")
		os.Exit(ExitInvalidCmd)
	}

	m := newMainModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		asfLog.Printf("received signal %v, shutting down", sig)
		p.Quit()
	}()

	defer func() {
		if cfg != nil {
			if err := cfg.Save(ConfigPath()); err != nil {
				debugLog.Printf("config save on exit: %v", err)
			}
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitGeneralError)
	}
}

func pluralize(n int, s string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, s)
	}
	return fmt.Sprintf("%d %ss", n, s)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
