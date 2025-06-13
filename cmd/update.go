package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"net/http"
	"io"
	"regexp"
	"strconv"
	"github.com/spf13/cobra"
	"github.com/henrycmeen/arkivtestern/internal/sfwrap"
)

var updateCmd = &cobra.Command{
	Use:   "update-signatures",
	Short: "Last ned og bygg inn nyeste PRONOM/DROID-signaturer",
	Long: `Laster ned siste DROID og container signature XML fra The National Archives, kompilerer til Siegfried-bundle og genererer Go-embed-fil for statisk binary. Kjør deretter 'go build' for å ta i bruk nye signaturer.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// URLs for latest signature files (kan evt. parameteriseres)
		// Finn siste DROID signature file automatisk
		fmt.Println("Henter URL til siste DROID signature file...")
		latestDroidURL, latestDroidFile, err := findLatestDroidSignature()
		if err != nil {
			return fmt.Errorf("klarte ikke å finne siste DROID signature file: %w", err)
		}
		signatureDir := filepath.Join("internal", "sfwrap", "signatures")
		if err := os.MkdirAll(signatureDir, 0755); err != nil {
			return fmt.Errorf("kunne ikke opprette signaturmappe: %w", err)
		}
		containerURL := "https://www.nationalarchives.gov.uk/documents/container-signature-20240715.xml"
		// Filnavn
		droidURL := latestDroidURL
		droidFile := filepath.Join(signatureDir, latestDroidFile)
		containerFile := filepath.Join(signatureDir, "container-signature.xml")
		bundleFile := filepath.Join(signatureDir, "default.sig")
		// Last ned DROID signature file
		fmt.Println("Laster ned DROID signature file...")
		err = downloadFile(droidURL, droidFile)
		if err != nil {
			return fmt.Errorf("feil ved nedlasting av DROID signature: %w", err)
		}
		// Last ned container signature file
		fmt.Println("Laster ned container signature file...")
		err = downloadFile(containerURL, containerFile)
		if err != nil {
			return fmt.Errorf("feil ved nedlasting av container signature: %w", err)
		}
		// Kjør Go-funksjon for å lage embed-fil fra bundle
		fmt.Println("Genererer Go-embed-fil fra default.sig...")
		embedPath := filepath.Join("internal", "sfwrap", "sfembed.go")
		err = sfwrap.GenerateSfEmbedGo(bundleFile, embedPath)
		if err != nil {
			return fmt.Errorf("feil ved generering av sfembed.go: %w", err)
		}
		fmt.Println("Ferdig! Kjør nå 'go build' for å bygge binary med oppdaterte signaturer.")
		return nil
	},
}

// findLatestDroidSignature finner URL og filnavn til siste DROID signature file
func findLatestDroidSignature() (string, string, error) {
	const droidListURL = "https://www.nationalarchives.gov.uk/aboutapps/pronom/droid-signature-files.htm"
	resp, err := http.Get(droidListURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("feil ved nedlasting av DROID signature list: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	// Finn alle forekomster av DROID_SignatureFile_Vxxx.xml
	// og velg høyeste versjon
	re := regexp.MustCompile(`https://cdn\.nationalarchives\.gov\.uk/documents/DROID_SignatureFile_V(\d+)\.xml`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	maxVer := -1
	var bestURL, bestFile string
	for _, m := range matches {
		if len(m) < 2 { continue }
		v, _ := strconv.Atoi(m[1])
		if v > maxVer {
			maxVer = v
			bestURL = m[0]
			parts := regexp.MustCompile(`/`).Split(bestURL, -1)
			bestFile = parts[len(parts)-1]
		}
	}
	if bestURL == "" {
		return "", "", fmt.Errorf("fant ingen DROID signature files på %s", droidListURL)
	}
	return bestURL, bestFile, nil
}

func downloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("feil ved nedlasting: %s", resp.Status)
	}
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
