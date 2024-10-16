package main

import (
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/mholt/archiver"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestDownloadFileUnreachableURL(t *testing.T) {
	dummyFile, err := os.Open("/dev/null")
	require.NoError(t, err)

	err = downloadFile("http://unreachable.invalid", dummyFile)
	require.Error(t, err)
}

func TestDownloadFileHTTPError(t *testing.T) {
	dummyFile, err := os.Open("/dev/null")
	require.NoError(t, err)

	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()

	err = downloadFile(ts.URL, dummyFile)
	require.Error(t, err)
}

func TestDownloadFileSuccessfully(t *testing.T) {
	dummyFile, err := os.CreateTemp("", "test.txt")
	require.NoError(t, err)

	ts := httptest.NewServer(http.FileServer(http.Dir(os.TempDir())))
	defer ts.Close()

	err = downloadFile(ts.URL, dummyFile)
	require.NoError(t, err)

	err = dummyFile.Close()
	require.NoError(t, err)

	err = os.Remove(dummyFile.Name())
	require.NoError(t, err)
}

func TestDownloadAnUnTarXZSuccessfully(t *testing.T) {
	archiveDir, err := os.MkdirTemp("", "testdir")
	require.NoError(t, err)

	dummyFile, err := os.CreateTemp("", "test")
	require.NoError(t, err)

	const archiveFileName = "test.tar.xz"
	archiveFile := path.Join(archiveDir, archiveFileName)

	err = archiver.Archive([]string{dummyFile.Name()}, archiveFile)
	require.NoError(t, err)

	ts := httptest.NewServer(http.FileServer(http.Dir(archiveDir)))
	defer ts.Close()

	destinationDir, err := os.MkdirTemp("", "destination")
	require.NoError(t, err)

	err = downloadAndUnTarXZ(ts.URL+"/"+archiveFileName, destinationDir)
	require.NoError(t, err)
	pathExists, err := pathutil.IsPathExists(dummyFile.Name())
	require.NoError(t, err)
	require.True(t, pathExists)

	require.NoError(t, os.RemoveAll(archiveDir))
	require.NoError(t, os.RemoveAll(destinationDir))
	require.NoError(t, os.RemoveAll(dummyFile.Name()))
}

func TestNormalizePre117SemanticVersion(t *testing.T) {
	require.Equal(t, "v1.16.9", normalizeFlutterVersion("1.16.9"))
}

func TestNormalizePost117SemanticVersion(t *testing.T) {
	require.Equal(t, "1.17.0", normalizeFlutterVersion("1.17.0"))
}

func TestNormalizeNonSemanticVersion(t *testing.T) {
	require.Equal(t, "master", normalizeFlutterVersion("master"))
}
