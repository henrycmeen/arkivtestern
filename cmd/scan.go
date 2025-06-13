package cmd

import (
    "fmt"
    "path/filepath"
    "encoding/json"

    "github.com/spf13/cobra"
    "go.uber.org/zap"

    "github.com/henrycmeen/arkivtestern/internal/sfwrap"
)

var scanCmd = &cobra.Command{
    Use:   "scan [path]",
    Short: "Skann en fil eller mappe (stub)",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        target := args[0]
        abs, err := filepath.Abs(target)
        if err != nil {
            return err
        }
        zap.L().Info("Scanning", zap.String("path", abs))

        // Siegfried-integrasjon via Go-bibliotek
        id, err := sfwrap.IdentifyFile(abs)
        if err != nil {
            return fmt.Errorf("identifikasjon feilet: %w", err)
        }
        b, err := json.MarshalIndent(id, "", "  ")
        if err != nil {
            return fmt.Errorf("marshal feilet: %w", err)
        }
        fmt.Println(string(b))
        return nil
    },
}
