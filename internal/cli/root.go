package cli

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var (
    flagConfigDir  string
    flagPluginsDir string
    flagVerbose    bool
    flagNoBackup   bool
)

// RootCmd is the singleton root command used by subcommands to attach themselves.
var RootCmd = &cobra.Command{
    Use:   "palettesmith",
    Short: "PaletteSmith - plugin-based TUI/CLI theme editor",
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // placeholder for global setup
        return nil
    },
}

func init() {
    RootCmd.PersistentFlags().StringVar(&flagConfigDir, "config", "", "config directory (default: XDG config)")
    RootCmd.PersistentFlags().StringVar(&flagPluginsDir, "plugins", "", "additional plugin directory")
    RootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "enable verbose logging")
    RootCmd.PersistentFlags().BoolVar(&flagNoBackup, "no-backup", false, "skip backups when applying changes")
}

// Execute runs the CLI and exits on error.
func Execute() {
    if err := RootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

