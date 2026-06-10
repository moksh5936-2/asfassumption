package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func printUsage() {
	fmt.Printf("ASF v%s — Architecture Security Framework\n", ASFVersion)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  asf                        Launch the TUI")
	fmt.Println("  asf --version, -v          Show version")
	fmt.Println("  asf --license              Show license status")
	fmt.Println("  asf doctor                 Run system diagnostics")
	fmt.Println("  asf --help, -h             Show this help")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Printf("  Config:  %s\n", asfConfigPath())
	fmt.Printf("  Cache:   %s\n", asfCacheDir())
	fmt.Printf("  License: %s\n", asfLicensePath())
	fmt.Println()
	fmt.Println("Documentation: https://github.com/moksh5936-2/asfassumption")
}

func main() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--version", "-v":
			fmt.Printf("ASF v%s\n", ASFVersion)
			os.Exit(0)
		case "--license":
			l := LoadLicense()
			if l != nil && l.Valid {
				fmt.Printf("License: %s\n", l.Message)
				os.Exit(0)
			}
			fmt.Println("No valid license found.")
			fmt.Printf("Place your license key in %s\n", asfLicensePath())
			os.Exit(1)
		case "doctor", "--doctor", "diagnose":
			runDoctor()
			os.Exit(0)
		case "--help", "-h":
			printUsage()
			os.Exit(0)
		}
	}

	helpFlags := map[string]bool{"--help": true, "-h": true}
	if len(os.Args) > 1 && !helpFlags[os.Args[1]] {
		fmt.Printf("ASF v%s — Architecture Security Framework\n", ASFVersion)
		fmt.Println()
		printUsage()
		os.Exit(0)
	}

	cfg, err := LoadConfig(asfConfigPath())
	if err != nil {
		def := DefaultConfig()
		cfg = &def
	}

	m := newMainModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}

	if cfg != nil {
		cfg.Save(asfConfigPath())
	}
}
