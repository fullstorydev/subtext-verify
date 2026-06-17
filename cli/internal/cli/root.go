package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/fullstorydev/subtext-verify/cli/internal/config"
)

const exitUsage = 2

// globalFlags holds the persistent flag values bound to rootCmd. They are read
// by the tunnel commands (auth resolution + flag forwarding to the daemon).
var globalFlags struct {
	apiKey     string
	region     string
	endpoint   string
	format     string
	configPath string
}

// globalConfig holds the loaded config file, available after cobra's OnInitialize fires.
var globalConfig config.File

var rootCmd = &cobra.Command{
	Use:   "subtext",
	Short: "Subtext Verify tunnel client",
	Long: `subtext is the Subtext Verify reverse-tunnel client.

It exposes the tunnel that lets the Subtext hosted browser reach a local dev
server — both as direct commands (connect / disconnect / status) and as an MCP
stdio server:

  subtext tunnel mcp`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute is the entry point called from main.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}
	return nil
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&globalFlags.apiKey, "api-key", "", "API key (or \"-\" to read from stdin)")
	pf.StringVar(&globalFlags.region, "region", "", "Region: na1 (default), eu1")
	pf.StringVar(&globalFlags.endpoint, "endpoint", "", "Override MCP endpoint URL")
	pf.StringVar(&globalFlags.format, "format", "json", "Output format: json (default) or text")
	pf.StringVar(&globalFlags.configPath, "config", "", "Config file path (default: ~/.config/subtext/config.yaml)")

	cobra.OnInitialize(loadConfig)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(tunnelCmd)
}

func loadConfig() {
	cfg, err := config.Load(globalFlags.configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(exitUsage)
	}
	globalConfig = cfg
}

// Version is set at build time via ldflags:
//
//	-X github.com/fullstorydev/subtext-verify/cli/internal/cli.Version=v1.2.3
var Version string

func binaryVersion() string {
	if Version != "" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" && info.Main.Version != "" {
		return info.Main.Version
	}
	return "dev"
}
