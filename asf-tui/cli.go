package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	indent1 = "  "
	indent2 = "    "
)

func printUsage() {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("ASF0 v%s — Security Assumption Framework\n", ASFVersion))
	b.WriteString("\nUsage:\n")
	b.WriteString(indent1 + "asf                         Launch the TUI\n")
	b.WriteString(indent1 + "asf analyze <file> [flags]   Run analysis (JSON output)\n")
	b.WriteString(indent1 + "asf export <file> [flags]    Export analysis results\n")
	b.WriteString(indent1 + "asf doctor [flags]           System diagnostics\n")
	b.WriteString(indent1 + "asf config [subcommand]      View/edit configuration\n")
	b.WriteString(indent1 + "asf completion [shell]       Generate shell completion\n")
	b.WriteString(indent1 + "asf --version, -v            Show version\n")
	b.WriteString(indent1 + "asf --license                Show license status\n")
	b.WriteString(indent1 + "asf --version-check          Check for newer version\n")
	b.WriteString(indent1 + "asf --help, -h               Show this help\n")
	b.WriteString("\nAnalyze Flags:\n")
	b.WriteString(indent2 + "-e, --evidence <path>        Evidence files/directories\n")
	b.WriteString(indent2 + "--graph                      Include dependency graph\n")
	b.WriteString(indent2 + "--report-type <type>         Report pack filter\n")
	b.WriteString("\nExport Flags:\n")
	b.WriteString(indent2 + "-f, --format <format>        Output format (json|markdown|csv|html|pdf|jsonl)\n")
	b.WriteString(indent2 + "-o, --output <dir>           Output directory\n")
	b.WriteString(indent2 + "--narrative                  Include narrative output\n")
	b.WriteString(indent2 + "--trust                      Include trust chain output\n")
	b.WriteString("\nDoctor Flags:\n")
	b.WriteString(indent2 + "--verbose, -v                Detailed diagnostics\n")
	b.WriteString(indent2 + "--fix                        Clean stale binaries\n")
	b.WriteString("\nConfig Subcommands:\n")
	b.WriteString(indent2 + "asf config                   Show current config\n")
	b.WriteString(indent2 + "asf config set <key> <val>   Set a config value\n")
	b.WriteString(indent2 + "asf config path              Show config file path\n")
	b.WriteString("\nCompletion:\n")
	b.WriteString(indent2 + "asf completion bash          Generate bash completion\n")
	b.WriteString(indent2 + "asf completion zsh           Generate zsh completion\n")
	b.WriteString(indent2 + "asf completion fish          Generate fish completion\n")
	b.WriteString("\nConfiguration:\n")
	b.WriteString(fmt.Sprintf("  Config:  %s\n", asfConfigPath()))
	b.WriteString(fmt.Sprintf("  Cache:   %s\n", asfCacheDir()))
	b.WriteString(fmt.Sprintf("  License: %s\n", asfLicensePath()))
	b.WriteString("\nDocumentation: https://github.com/moksh5936-2/asfassumption\n")

	os.Stdout.WriteString(b.String())
}

func runConfigCLI(args []string) {
	if len(args) == 0 {
		cfg, err := LoadConfig(ConfigPath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(ExitGeneralError)
		}
		fmt.Printf("Config file: %s\n\n", ConfigPath())
		fmt.Printf("Theme:         %s\n", cfg.Appearance.Theme)
		fmt.Printf("Fox Style:     %s\n", cfg.Appearance.FoxStyle)
		fmt.Printf("Analysis Depth: %s\n", cfg.Analysis.Depth)
		fmt.Printf("Risk Threshold: %s\n", cfg.Analysis.RiskThreshold)
		fmt.Printf("STRIDE:        %v\n", cfg.Analysis.Stride)
		fmt.Printf("Controls:      %v\n", cfg.Analysis.Controls)
		fmt.Printf("Export Format: %s\n", cfg.Output.Default)
		fmt.Printf("Export Dir:    %s\n", cfg.Output.Directory)
		fmt.Printf("AI Enabled:    %v\n", cfg.AI.Enabled)
		fmt.Printf("Debug:         %v\n", cfg.General.Debug)
		fmt.Printf("Telemetry:     %v\n", cfg.Telemetry.OptIn)
		return
	}

	switch args[0] {
	case "set":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: asf config set <key> <value>\n")
			os.Exit(ExitInvalidCmd)
		}
		setConfigValue(args[1], args[2])
	case "path":
		fmt.Println(ConfigPath())
	case "--help", "-h":
		fmt.Println("Usage: asf config [subcommand]")
		fmt.Println()
		fmt.Println("Subcommands:")
		fmt.Println("  (no subcommand)    Show current configuration")
		fmt.Println("  set <key> <value>  Set a configuration value")
		fmt.Println("  path               Show config file path")
	default:
		fmt.Fprintf(os.Stderr, "Unknown config subcommand: %s\n", args[0])
		os.Exit(ExitInvalidCmd)
	}
}

func setConfigValue(key, value string) {
	cfg, err := LoadConfig(ConfigPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(ExitGeneralError)
	}

	switch key {
	case "theme":
		cfg.Appearance.Theme = value
		cfg.General.Theme = value
	case "fox_style", "fox-style":
		cfg.Appearance.FoxStyle = value
		cfg.General.FoxStyle = value
	case "depth":
		cfg.Analysis.Depth = value
	case "risk_threshold", "risk-threshold":
		cfg.Analysis.RiskThreshold = value
	case "stride":
		cfg.Analysis.Stride = value == "true" || value == "yes" || value == "1"
	case "controls":
		cfg.Analysis.Controls = value == "true" || value == "yes" || value == "1"
	case "export_format", "export-format":
		cfg.Output.Default = value
	case "export_dir", "export-dir":
		cfg.Output.Directory = value
	case "ai_enabled", "ai-enabled":
		cfg.AI.Enabled = value == "true" || value == "yes" || value == "1"
	case "debug":
		cfg.General.Debug = value == "true" || value == "yes" || value == "1"
	case "telemetry":
		cfg.Telemetry.OptIn = value == "true" || value == "yes" || value == "1"
	case "strict_cli", "strict-cli":
		cfg.StrictCLI = value == "true" || value == "yes" || value == "1"
	default:
		fmt.Fprintf(os.Stderr, "Unknown config key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Valid keys: theme, fox-style, depth, risk-threshold, stride, controls, export-format, export-dir, ai-enabled, debug, telemetry, strict-cli\n")
		os.Exit(ExitInvalidCmd)
	}

	if err := cfg.Save(ConfigPath()); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(ExitGeneralError)
	}
	fmt.Printf("Set %s = %s\n", key, value)
}

func runCompletionCLI(args []string) {
	shell := "bash"
	if len(args) > 0 {
		shell = args[0]
	}

	switch shell {
	case "bash":
		fmt.Print(bashCompletion)
	case "zsh":
		fmt.Print(zshCompletion)
	case "fish":
		fmt.Print(fishCompletion)
	default:
		fmt.Fprintf(os.Stderr, "Unknown shell: %s\n", shell)
		fmt.Fprintf(os.Stderr, "Supported shells: bash, zsh, fish\n")
		os.Exit(ExitInvalidCmd)
	}
}

const bashCompletion = `# asf bash completion
_asf() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="analyze export doctor config completion --version --license --version-check --help"

    if [[ ${cur} == -* ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi

    case "${prev}" in
        analyze)
            COMPREPLY=( $(compgen -f -- ${cur}) )
            ;;
        export)
            COMPREPLY=( $(compgen -f -- ${cur}) )
            ;;
        config)
            COMPREPLY=( $(compgen -W "set path" -- ${cur}) )
            ;;
        completion)
            COMPREPLY=( $(compgen -W "bash zsh fish" -- ${cur}) )
            ;;
        -f|--format)
            COMPREPLY=( $(compgen -W "json markdown csv html pdf jsonl" -- ${cur}) )
            ;;
        --report-type)
            COMPREPLY=( $(compgen -W "board executive technical architect-narrative executive-summary technical-summary" -- ${cur}) )
            ;;
        *)
            if [[ ${cur} == -* ]] ; then
                COMPREPLY=( $(compgen -W "-e --evidence --graph --report-type -f --format -o --output --narrative --trust --verbose --fix -v --help" -- ${cur}) )
            fi
            ;;
    esac
}
complete -F _asf asf
`

const zshCompletion = `#compdef asf
_asf() {
    local -a subcmds
    subcmds=(
        'analyze:Run analysis on architecture documents'
        'export:Export analysis results'
        'doctor:System diagnostics'
        'config:View/edit configuration'
        'completion:Generate shell completion'
    )
    _arguments \
        '--version[Show version]' \
        '--license[Show license status]' \
        '--version-check[Check for newer version]' \
        '--help[Show help]' \
        '-v[Show version]' \
        '-h[Show help]' \
        '*:: :->subcmd'
    case $state in
        subcmd)
            _describe 'command' subcmds
            ;;
    esac
}
_asf "$@"
`

const fishCompletion = `# asf fish completion
complete -c asf -f
complete -c asf -n "not __fish_seen_subcommand_from analyze export doctor config completion" -a analyze -d "Run analysis on architecture documents"
complete -c asf -n "not __fish_seen_subcommand_from analyze export doctor config completion" -a export -d "Export analysis results"
complete -c asf -n "not __fish_seen_subcommand_from analyze export doctor config completion" -a doctor -d "System diagnostics"
complete -c asf -n "not __fish_seen_subcommand_from analyze export doctor config completion" -a config -d "View/edit configuration"
complete -c asf -n "not __fish_seen_subcommand_from analyze export doctor config completion" -a completion -d "Generate shell completion"
complete -c asf -l version -d "Show version"
complete -c asf -l license -d "Show license status"
complete -c asf -l version-check -d "Check for newer version"
complete -c asf -l help -d "Show help"
complete -c asf -s v -d "Show version"
complete -c asf -s h -d "Show help"
`
