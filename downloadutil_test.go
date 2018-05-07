package main

import (
	"testing"
	"os"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"github.com/mholt/archiver"
	"path"
	"github.com/bitrise-io/go-utils/pathutil"
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
	dummyFile, err := ioutil.TempFile("", "test.txt")
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
	archiveDir, err := ioutil.TempDir("", "testdir")
	require.NoError(t, err)

	dummyFile, err := ioutil.TempFile("", "test")
	require.NoError(t, err)

	const archiveFileName = "test.tar.xz"
	archiveFile := path.Join(archiveDir, archiveFileName)

	err = archiver.TarXZ.Make(archiveFile, []string{dummyFile.Name()})
	require.NoError(t, err)

	ts := httptest.NewServer(http.FileServer(http.Dir(archiveDir)))
	defer ts.Close()

	destinationDir, err := ioutil.TempDir("", "destination")
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
