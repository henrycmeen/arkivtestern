package cmd

import (
    "fmt"
    "runtime"

    "github.com/spf13/cobra"
)

var (
    // Disse variablene settes via -ldflags under build.
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Vis versjonsinformasjon",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("arkivtestern %s (%s) %s %s\n", version, commit, date, runtime.GOOS+"/"+runtime.GOARCH)
    },
}
