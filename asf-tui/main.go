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

func printUsage() {
	fmt.Printf("ASF v%s — Architecture Security Framework\n", ASFVersion)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  asf                        Launch the TUI")
	fmt.Println("  asf --version, -v          Show version")
	fmt.Println("  asf --license              Show license status")
	fmt.Println("  asf analyze <file>         Run native analysis (JSON output)")
	fmt.Println("  asf analyze <file> -e <ev> ...   With evidence files/dirs")
	fmt.Println("  asf analyze <file> --graph Include graph in JSON output")
	fmt.Println("  asf doctor                 Run system diagnostics")
	fmt.Println("  asf doctor --verbose       Detailed diagnostics")
	fmt.Println("  asf doctor --fix           Clean stale binaries")
	fmt.Println("  asf --help, -h             Show this help")
	fmt.Println("  asf --version-check       Check for newer version")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Printf("  Config:  %s\n", asfConfigPath())
	fmt.Printf("  Cache:   %s\n", asfCacheDir())
	fmt.Printf("  License: %s\n", asfLicensePath())
	fmt.Println()
	fmt.Println("Documentation: https://github.com/moksh5936-2/asfassumption")
}

func main() {
	if err := initLogger(); err != nil {
		asfLog = log.New(io.Discard, "[asf] ", log.Ldate|log.Ltime|log.Lshortfile)
		debugLog = log.New(io.Discard, "[asf-debug] ", log.Ltime|log.Lshortfile)
	}
	asfLog.Printf("ASF v%s starting", ASFVersion)

	args := os.Args[1:]

	if len(args) > 0 {
		switch args[0] {
		case "--version", "-v":
			fmt.Printf("ASF v%s\n", ASFVersion)
			if msg := VersionCheckMessage(); msg != "" {
				fmt.Println(msg)
			}
			os.Exit(ExitSuccess)
		case "--version-check":
			if msg := VersionCheckMessage(); msg != "" {
				fmt.Println(msg)
			} else {
				fmt.Printf("ASF v%s is up to date.\n", ASFVersion)
			}
			os.Exit(ExitSuccess)
		case "--license":
			l := LoadLicense()
			if l != nil && l.Valid {
				fmt.Printf("License: %s\n", l.Message)
				os.Exit(ExitSuccess)
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
			os.Exit(ExitSuccess)
		case "analyze":
			runAnalyzeCLI(args[1:])
			os.Exit(ExitSuccess)
		case "--help", "-h":
			printUsage()
			os.Exit(ExitSuccess)
		default:
			fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", args[0])
			fmt.Fprintf(os.Stderr, "Run 'asf --help' for usage.\n")
			os.Exit(ExitInvalidCmd)
		}
	}

	cfg, err := LoadConfig(asfConfigPath())
	if err != nil {
		def := DefaultConfig()
		cfg = &def
	}
	asfLog.Printf("config path: %s", asfConfigPath())

	if err := ensureRuntimeDirs(); err != nil {
		debugLog.Printf("runtime dirs: %v", err)
	}

	m := newMainModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	go func() {
		<-sigCh
		p.Quit()
	}()

	defer func() {
		if cfg != nil {
			if err := cfg.Save(asfConfigPath()); err != nil {
				debugLog.Printf("config save on exit: %v", err)
			}
		}
	}()

	if _, err := p.Run(); err != nil {
		os.Exit(ExitGeneralError)
	}
}
