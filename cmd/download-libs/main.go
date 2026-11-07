// RTI Connector Library Downloader
// This tool helps users download RTI Connector libraries when using go get
package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
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

	// Check if libraries already exist (unless force is specified)
	libDir := filepath.Join(*destination, "rticonnextdds-connector")
	if !*force {
		if _, err := os.Stat(libDir); err == nil {
			fmt.Printf("✅ Libraries already exist at %s\n", libDir)
			fmt.Printf("Use -force flag to overwrite, or -current to check installation\n")
			return
		}
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
		if targetVersion == "" {
			fmt.Printf("Error: Latest version is empty\n")
			os.Exit(1)
		}
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
		case "arm":
			// linux-arm (32-bit) was removed in v1.4.0
			fmt.Printf("⚠️  Warning: linux-arm (32-bit) support was removed in v1.4.0\n")
			fmt.Printf("   Use v1.3.1 or earlier for 32-bit ARM support\n")
			return "linux-arm"
		default:
			fmt.Printf("Unsupported Linux architecture: %s\n", runtime.GOARCH)
			os.Exit(1)
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			// osx-x64 (Intel Mac) was removed in v1.4.0
			fmt.Printf("⚠️  Warning: osx-x64 (Intel Mac) support was removed in v1.4.0\n")
			fmt.Printf("   Use v1.3.1 or earlier for Intel Mac support\n")
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

// getPlatform returns the platform identifier for the current system
func getPlatform() string {
	return detectPlatform()
}

// createHTTPClient creates an HTTP client with optional GitHub token authentication
func createHTTPClient() *http.Client {
	return &http.Client{}
}

// createRequest creates an HTTP request with optional GitHub token authentication
func createRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Add GitHub token if available
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req, nil
}

func listVersions() {
	fmt.Println("📋 Available Versions:")

	req, err := createRequest("GET", baseURL+"/releases")
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	client := createHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching versions: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: API request failed with status %s\n", resp.Status)
		if resp.StatusCode == http.StatusForbidden {
			fmt.Println("Hint: You may have hit GitHub API rate limits. Set GITHUB_TOKEN environment variable.")
		}
		return
	}

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
	req, err := createRequest("GET", baseURL+"/releases/latest")
	if err != nil {
		return "", err
	}

	client := createHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("API request failed: %s", resp.Status)
		if resp.StatusCode == http.StatusForbidden {
			errMsg += "\nHint: You may have hit GitHub API rate limits. Set GITHUB_TOKEN environment variable."
		}
		return "", errors.New(errMsg)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	if release.TagName == "" {
		return "", fmt.Errorf("tag_name is empty in API response")
	}

	return release.TagName, nil
}

func getDownloadURL(version string) (string, string, error) {
	releaseURL := fmt.Sprintf("%s/releases/tags/%s", baseURL, version)

	req, err := createRequest("GET", releaseURL)
	if err != nil {
		return "", "", err
	}

	client := createHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("release %s not found (status: %s)", version, resp.Status)
		if resp.StatusCode == http.StatusForbidden {
			errMsg += "\nHint: You may have hit GitHub API rate limits. Set GITHUB_TOKEN environment variable."
		}
		return "", "", errors.New(errMsg)
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
	platform := getPlatform()
	libPath := filepath.Join(dest, "rticonnextdds-connector", "lib", platform)

	fmt.Println("📋 Current Installation:")
	fmt.Printf("  Platform: %s\n", platform)
	fmt.Printf("  Library path: %s\n", libPath)

	// Check for version information
	version := getInstalledVersion(dest)
	if version != "" {
		fmt.Printf("  Version: %s\n", version)
	}

	// Check if libraries exist
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		fmt.Printf("  Status: ❌ No libraries found\n")
		fmt.Printf("  Run: go run github.com/rticommunity/rticonnextdds-connector-go/cmd/download-libs@latest\n")
		return
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

// getInstalledVersion attempts to detect the installed version of RTI Connector libraries
func getInstalledVersion(dest string) string {
	return detectVersionFromLibraries(dest)
} // detectVersionFromLibraries tries to determine version based on library characteristics
func detectVersionFromLibraries(dest string) string {
	platform := getPlatform()
	libPath := filepath.Join(dest, "rticonnextdds-connector", "lib", platform)

	// Check if libraries exist
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		return ""
	}

	// Try to extract version from the connector library binary
	connectorLib := filepath.Join(libPath, "librtiddsconnector")

	// Add appropriate file extension based on platform
	switch runtime.GOOS {
	case "linux":
		connectorLib += ".so"
	case "darwin":
		connectorLib += ".dylib"
	case "windows":
		connectorLib += ".dll"
	}

	if version := extractVersionFromBinary(connectorLib); version != "" {
		return version
	}

	return "unknown (installed before version tracking)"
}

// extractVersionFromBinary attempts to extract version information from the connector library
func extractVersionFromBinary(libPath string) string {
	// Check if file exists
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		return ""
	}

	// Use strings command to extract version information
	cmd := exec.Command("strings", libPath)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Look for RTICONNECTOR_BUILD pattern in the output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "RTICONNECTOR_BUILD_") {
			// Extract version from line like "RTICONNECTOR_BUILD_7.6.0.0_20250912T000000Z_RTI_REL"
			// Remove the prefix and split by underscore
			withoutPrefix := strings.TrimPrefix(line, "RTICONNECTOR_BUILD_")
			parts := strings.Split(withoutPrefix, "_")
			if len(parts) >= 1 {
				version := parts[0] // This should be the version number like "7.6.0.0"
				return fmt.Sprintf("RTI Connext DDS %s", version)
			}
		}
	}

	return ""
}

func downloadLibraries(version, dest string, force bool) error {
	platform := getPlatform()
	libDir := filepath.Join(dest, "rticonnextdds-connector")

	fmt.Printf("🌐 Downloading RTI Connector %s...\n", version)
	fmt.Printf("📱 Target platform: %s\n", platform)

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

	// Verify platform-specific libraries exist
	platformLibDir := filepath.Join(libDir, "lib", platform)
	if _, err := os.Stat(platformLibDir); os.IsNotExist(err) {
		// List available platforms for debugging
		libBase := filepath.Join(libDir, "lib")
		entries, readErr := os.ReadDir(libBase)
		if readErr == nil {
			var availablePlatforms []string
			for _, e := range entries {
				if e.IsDir() {
					availablePlatforms = append(availablePlatforms, e.Name())
				}
			}
			return fmt.Errorf("platform '%s' not found in %s. Available platforms: %v. "+
				"This version may not support your platform. Try: -version v1.3.1",
				platform, version, availablePlatforms)
		}
		return fmt.Errorf("platform '%s' libraries not found in downloaded archive", platform)
	}

	fmt.Printf("✅ Platform libraries verified: %s\n", platformLibDir)
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

	// Detect archive structure:
	// - v1.3.1 and earlier: lib/* (need to prepend connectorDir)
	// - v1.4.0 and later: rticonnextdds-connector/lib/* (already has prefix)
	hasConnectorPrefix := false
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "rticonnextdds-connector/") {
			hasConnectorPrefix = true
			break
		}
	}

	// Extract files
	for _, f := range r.File {
		var path string
		if hasConnectorPrefix {
			// v1.4.0+: Extract directly to dest, the archive already has rticonnextdds-connector/ prefix
			path = filepath.Join(dest, f.Name)
		} else {
			// v1.3.1 and earlier: Prepend connectorDir to create rticonnextdds-connector/lib/*
			path = filepath.Join(connectorDir, f.Name)
		}

		// Check for ZipSlip vulnerability - ensure path is within the destination tree
		// Use absolute paths to avoid issues with relative path comparisons
		absDest, err := filepath.Abs(dest)
		if err != nil {
			return fmt.Errorf("getting absolute path for dest: %v", err)
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("getting absolute path for file: %v", err)
		}
		if !strings.HasPrefix(absPath, absDest+string(os.PathSeparator)) && absPath != absDest {
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
	platform := getPlatform()
	libPath := filepath.Join(dest, "rticonnextdds-connector", "lib", platform)

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
