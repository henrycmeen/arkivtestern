package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "go.uber.org/zap"
)

var (
    cfgFile   string
    verbose   bool
    quiet     bool
    rootCmd   = &cobra.Command{
        Use:   "arkivtestern",
        Short: "Arkivtestern CLI – AI-drevet arkivtest for oss",
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            if quiet {
                // force zap to production logger with silence
                zap.ReplaceGlobals(zap.NewNop())
                return nil
            }
            var logger *zap.Logger
            var err error
            if verbose {
                logger, err = zap.NewDevelopment()
            } else {
                logger, err = zap.NewProduction()
            }
            if err != nil {
                return err
            }
            zap.ReplaceGlobals(logger)
            return nil
        },
    }
)

// Execute is the entry point from main.go
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    // Global flags
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.arkivtestern.yaml)")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
    rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress all non-error output")

    // Bind to viper
    _ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
    _ = viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))

    // Sub-commands
    rootCmd.AddCommand(scanCmd)
    rootCmd.AddCommand(versionCmd)
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        if err == nil {
            viper.AddConfigPath(home)
        }
        viper.SetConfigName(".arkivtestern")
    }
    viper.AutomaticEnv()
    _ = viper.ReadInConfig() // ignore error if not found
}
