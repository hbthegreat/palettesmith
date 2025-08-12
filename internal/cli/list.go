package cli

import (
    "encoding/json"
    "fmt"
    "sort"

    "github.com/spf13/cobra"

    "palettesmith/internal/plugin"
)

func newListCmd() *cobra.Command {
    var flagJSON bool
    cmd := &cobra.Command{
        Use:   "list",
        Short: "List detected applications",
        RunE: func(cmd *cobra.Command, args []string) error {
            reg := plugin.NewRegistry()
            // builtins then external, external overrides on conflict
            for id, p := range plugin.LoadBuiltinPlugins() {
                reg.Register(id, p)
            }
            if ext, err := plugin.LoadExternalPlugins(); err == nil {
                for id, p := range ext { reg.Register(id, p) }
            }
            detected := reg.DetectAll()
            // sort output by name
            type row struct { ID, Name string; Detected bool }
            rows := make([]row, 0, len(reg.List()))
            for id, p := range reg.List() {
                rows = append(rows, row{ID: id, Name: p.GetManifest().Metadata.Name, Detected: detected[id]})
            }
            sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
            if flagJSON {
                enc := json.NewEncoder(cmd.OutOrStdout())
                enc.SetIndent("", "  ")
                return enc.Encode(rows)
            }
            for _, r := range rows {
                status := "✗"
                if r.Detected { status = "✓" }
                fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", r.Name, status)
            }
            return nil
        },
    }
    cmd.Flags().BoolVar(&flagJSON, "json", false, "output JSON")
    return cmd
}

func init() {
    // Attach to root
    RootCmd.AddCommand(newListCmd())
}

