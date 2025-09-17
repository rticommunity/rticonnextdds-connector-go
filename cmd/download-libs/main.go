// RTI Connector Library Downloader
// This tool helps users download RTI Connector libraries when using go get
package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	repoOwner = "rticommunity"
	repoName  = "rticonnextdds-connector"
	baseURL   = "https://api.github.com/repos/" + repoOwner + "/" + repoName
)

type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func main() {
	var (
		version     = flag.String("version", "", "Specific version to download (e.g., v1.3.1)")
		list        = flag.Bool("list", false, "List available versions")
		current     = flag.Bool("current", false, "Show current installation info")
		force       = flag.Bool("force", false, "Force download even if libraries exist")
		destination = flag.String("dest", ".", "Destination directory for libraries")
	)
	flag.Parse()

	if *list {
		listVersions()
		return
	}

	if *current {
		showCurrent(*destination)
		return
	}

	targetVersion := *version
	if targetVersion == "" {
		var err error
		targetVersion, err = getLatestVersion()
		if err != nil {
			fmt.Printf("Error getting latest version: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Latest version: %s\n", targetVersion)
	}

	err := downloadLibraries(targetVersion, *destination, *force)
	if err != nil {
		fmt.Printf("Error downloading libraries: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Libraries downloaded successfully!")
	showSetupInstructions(*destination)
}

func detectPlatform() string {
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return "linux-x64"
		case "arm64":
			return "linux-arm64"
		default:
			fmt.Printf("Unsupported Linux architecture: %s\n", runtime.GOARCH)
			os.Exit(1)
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return "osx-x64"
		case "arm64":
			return "osx-arm64"
		default:
			fmt.Printf("Unsupported macOS architecture: %s\n", runtime.GOARCH)
			os.Exit(1)
		}
	case "windows":
		return "win-x64"
	default:
		fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
		os.Exit(1)
	}
	return ""
}

// getPlatformWithFallback returns the preferred platform and a fallback if needed
func getPlatformWithFallback() (string, string) {
	preferred := detectPlatform()

	// Define fallbacks for platforms that might not be available in older releases
	fallbacks := map[string]string{
		"osx-arm64":   "osx-x64",   // Apple Silicon can run Intel binaries
		"linux-arm64": "linux-x64", // Some ARM64 systems can run x64 with emulation
	}

	fallback := fallbacks[preferred]
	return preferred, fallback
}

func listVersions() {
	fmt.Println("📋 Available Versions:")
	resp, err := http.Get(baseURL + "/releases")
	if err != nil {
		fmt.Printf("Error fetching versions: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		fmt.Printf("Error parsing releases: %v\n", err)
		return
	}

	for i, release := range releases {
		if i >= 10 { // Show latest 10 versions
			break
		}
		fmt.Printf("  %s\n", release.TagName)
	}
}

func getLatestVersion() (string, error) {
	resp, err := http.Get(baseURL + "/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func getDownloadURL(version string) (string, string, error) {
	releaseURL := fmt.Sprintf("%s/releases/tags/%s", baseURL, version)
	resp, err := http.Get(releaseURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("release %s not found", version)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", err
	}

	// Find the ZIP asset
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".zip") {
			return asset.BrowserDownloadURL, asset.Name, nil
		}
	}

	return "", "", fmt.Errorf("no ZIP asset found in release %s", version)
}

func showCurrent(dest string) {
	preferred, fallback := getPlatformWithFallback()
	libPath := filepath.Join(dest, "rticonnextdds-connector", "lib", preferred)

	fmt.Println("📋 Current Installation:")
	fmt.Printf("  Platform: %s", preferred)
	if fallback != "" {
		fmt.Printf(" (fallback: %s)", fallback)
	}
	fmt.Printf("\n")
	fmt.Printf("  Library path: %s\n", libPath)

	// Check preferred platform first
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		// Try fallback if preferred doesn't exist
		if fallback != "" {
			fallbackPath := filepath.Join(dest, "rticonnextdds-connector", "lib", fallback)
			if _, err := os.Stat(fallbackPath); err == nil {
				fmt.Printf("  Status: ✅ Libraries installed (using fallback %s)\n", fallback)
				fmt.Printf("  Fallback path: %s\n", fallbackPath)
				libPath = fallbackPath // Use fallback path for listing files
			} else {
				fmt.Printf("  Status: ❌ No libraries found\n")
				fmt.Printf("  Run: go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest\n")
				return
			}
		} else {
			fmt.Printf("  Status: ❌ No libraries found\n")
			fmt.Printf("  Run: go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest\n")
			return
		}
	} else {
		fmt.Printf("  Status: ✅ Libraries installed\n")
	}

	// List library files
	entries, err := os.ReadDir(libPath)
	if err != nil {
		fmt.Printf("  Error reading directory: %v\n", err)
		return
	}

	fmt.Printf("  Libraries:\n")
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".so") ||
			strings.HasSuffix(entry.Name(), ".dylib") ||
			strings.HasSuffix(entry.Name(), ".dll")) {

			info, err := entry.Info()
			if err == nil {
				fmt.Printf("    %s (%d bytes)\n", entry.Name(), info.Size())
			} else {
				fmt.Printf("    %s\n", entry.Name())
			}
		}
	}
}

func downloadLibraries(version, dest string, force bool) error {
	preferred, fallback := getPlatformWithFallback()
	libDir := filepath.Join(dest, "rticonnextdds-connector")

	// Check if libraries already exist
	if !force {
		if _, err := os.Stat(libDir); err == nil {
			fmt.Printf("⚠️  Libraries already exist at %s\n", libDir)
			fmt.Printf("Use -force flag to overwrite, or -current to check installation\n")
			return nil
		}
	}

	fmt.Printf("🌐 Downloading RTI Connector %s...\n", version)
	fmt.Printf("📱 Target platform: %s", preferred)
	if fallback != "" {
		fmt.Printf(" (fallback available: %s)", fallback)
	}
	fmt.Printf("\n")

	// Get the actual download URL from GitHub API
	downloadURL, archiveName, err := getDownloadURL(version)
	if err != nil {
		return fmt.Errorf("finding download URL: %v", err)
	}

	fmt.Printf("📦 Downloading: %s\n", archiveName)

	// Create temporary file
	tmpFile, err := os.CreateTemp("", archiveName)
	if err != nil {
		return fmt.Errorf("creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Download file
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("downloading file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s (check if version %s exists)", resp.Status, version)
	}

	// Copy to temp file with progress
	fmt.Printf("⬇️  Downloading...")
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("saving file: %v", err)
	}
	fmt.Printf(" Done!\n")

	// Extract archive
	fmt.Printf("📂 Extracting archive...\n")
	err = extractZip(tmpFile.Name(), dest, libDir)
	if err != nil {
		return fmt.Errorf("extracting archive: %v", err)
	}

	fmt.Printf("✅ Libraries installed to: %s\n", libDir)
	return nil
}

func extractZip(src, dest, connectorDir string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create the rticonnextdds-connector directory
	os.MkdirAll(connectorDir, 0755)

	// Extract files
	for _, f := range r.File {
		// Map the extracted path to include the connector directory
		// The ZIP contains lib/* which we want to extract to rticonnextdds-connector/lib/*
		path := filepath.Join(connectorDir, f.Name)

		// Check for ZipSlip vulnerability - ensure path is within the destination tree
		if !strings.HasPrefix(path, filepath.Clean(connectorDir)+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.FileInfo().Mode())
			continue
		}

		// Create directories if needed
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// Extract file
		rc, err := f.Open()
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func showSetupInstructions(dest string) {
	preferred, fallback := getPlatformWithFallback()
	libPath := filepath.Join(dest, "rticonnextdds-connector", "lib", preferred)

	// Check if preferred platform exists, otherwise use fallback
	if _, err := os.Stat(libPath); os.IsNotExist(err) && fallback != "" {
		fallbackPath := filepath.Join(dest, "rticonnextdds-connector", "lib", fallback)
		if _, err := os.Stat(fallbackPath); err == nil {
			libPath = fallbackPath
			fmt.Printf("ℹ️  Using %s libraries (fallback for %s)\n\n", fallback, preferred)
		}
	}

	fmt.Println("🔧 Setup Instructions:")
	fmt.Println("Add the following to your environment:")

	// Determine OS for environment setup
	switch runtime.GOOS {
	case "linux":
		fmt.Printf("export LD_LIBRARY_PATH=%s:$LD_LIBRARY_PATH\n", libPath)
	case "darwin":
		fmt.Printf("export DYLD_LIBRARY_PATH=%s:$DYLD_LIBRARY_PATH\n", libPath)
	case "windows":
		fmt.Printf("set PATH=%s;%%PATH%%\n", libPath)
	}

	fmt.Println("\n📝 Example usage:")
	fmt.Println("  go run your_program.go")
}
