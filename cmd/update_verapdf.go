package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"github.com/henrycmeen/arkivtestern/internal/verapdf"
	"github.com/spf13/cobra"
)

var updateVeraPDFCmd = &cobra.Command{
	Use:   "update-verapdf",
	Short: "Last ned og bygg inn siste veraPDF JAR",
	Long:  `Laster ned siste veraPDF CLI JAR fra GitHub, og genererer Go-embed-fil for statisk binary. Kjør deretter 'go build' for å ta i bruk oppdatert veraPDF.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Henter URL til siste veraPDF release...")
		latestURL, jarFile, err := findLatestVeraPDFJar()
		if err != nil {
			return fmt.Errorf("klarte ikke å finne siste veraPDF JAR: %w", err)
		}
		verapdfDir := filepath.Join("internal", "verapdf")
		if err := os.MkdirAll(verapdfDir, 0755); err != nil {
			return fmt.Errorf("kunne ikke opprette veraPDF-mappe: %w", err)
		}
		jarPath := filepath.Join(verapdfDir, jarFile)
		fmt.Println("Laster ned veraPDF JAR...")
		err = downloadFile(latestURL, jarPath)
		if err != nil {
			return fmt.Errorf("feil ved nedlasting av veraPDF JAR: %w", err)
		}
		fmt.Println("Genererer Go-embed-fil...")
		embedPath := filepath.Join(verapdfDir, "verapdfembed.go")
		err = verapdf.GenerateVeraPDFEmbedGo(jarPath, embedPath)
		if err != nil {
			return fmt.Errorf("feil ved generering av verapdfembed.go: %w", err)
		}
		fmt.Println("Ferdig! Kjør nå 'go build' for å bygge binary med oppdatert veraPDF.")
		return nil
	},
}


func findLatestVeraPDFJar() (string, string, error) {
	resp, err := http.Get("https://api.github.com/repos/veraPDF/veraPDF-apps/releases/latest")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var release struct {
		Assets []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", err
	}
	for _, asset := range release.Assets {
		matched, _ := regexp.MatchString(`verapdf-verapdf-.*-cli\.jar$`, asset.Name)
		if matched {
			return asset.BrowserDownloadURL, asset.Name, nil
		}
	}
	return "", "", fmt.Errorf("fant ikke veraPDF CLI JAR i siste release")
}



func init() {
	rootCmd.AddCommand(updateVeraPDFCmd)
}
