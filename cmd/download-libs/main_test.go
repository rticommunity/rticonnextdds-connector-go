package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func libraryArchive(t *testing.T, platform string) string {
	t.Helper()
	var archive bytes.Buffer
	writer := zip.NewWriter(&archive)
	entry, err := writer.Create("lib/" + platform + "/fixture.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(entry, "test library"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return archive.String()
}

func TestDownloadVersionSelection(t *testing.T) {
	for _, tc := range []struct {
		name     string
		version  string
		override bool
	}{
		{name: "pinned default", version: "v1.4.0"},
		{name: "explicit legacy version", version: "v1.3.1", override: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dest := t.TempDir()
			platform := getPlatform()
			archive := libraryArchive(t, platform)

			oldArgs, oldFlags, oldTransport := os.Args, flag.CommandLine, http.DefaultTransport
			t.Cleanup(func() {
				os.Args, flag.CommandLine, http.DefaultTransport = oldArgs, oldFlags, oldTransport
			})
			os.Args = []string{"download-libs", "-dest", dest}
			if tc.override {
				os.Args = append(os.Args, "-version", tc.version)
			}
			flag.CommandLine = flag.NewFlagSet("download-libs", flag.ContinueOnError)
			var releasePath string
			latestRequested := false
			http.DefaultTransport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
				body := `{"tag_name":"v1.5.0","assets":[{"name":"connector.zip","browser_download_url":"https://downloads.example/connector.zip"}]}`
				if req.URL.Path == "/connector.zip" {
					body = archive
				} else if req.URL.Path == "/repos/"+repoOwner+"/"+repoName+"/releases/latest" {
					latestRequested = true
				} else if strings.HasPrefix(req.URL.Path, "/repos/"+repoOwner+"/"+repoName+"/releases/tags/") {
					releasePath = req.URL.Path
				} else {
					t.Fatalf("unexpected request: %s", req.URL)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			})

			main()

			if latestRequested {
				t.Error("download must not query the latest release")
			}
			if want := "/repos/" + repoOwner + "/" + repoName + "/releases/tags/" + tc.version; releasePath != want {
				t.Errorf("release request = %q, want %q", releasePath, want)
			}
			if _, err := os.Stat(filepath.Join(dest, "rticonnextdds-connector", "lib", platform, "fixture.txt")); err != nil {
				t.Fatalf("library was not installed: %v", err)
			}
		})
	}
}
