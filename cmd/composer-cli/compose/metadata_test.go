// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

func TestCmdComposeMetadata(t *testing.T) {
	// Test the "compose metadata" command
	mc := root.SetupCmdTest(func(request *http.Request) (*http.Response, error) {
		log := `This is a poor approximation of a logfile.
But it has multiple lines.
And should do the job.`

		tar, err := root.MakeTarBytes("b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7.json", log)
		require.Nil(t, err)

		resp := http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(tar)),
			Header:     http.Header{},
		}
		resp.Header.Set("Content-Disposition", "attachment; filename=b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7-metadata.tar")
		resp.Header.Set("Content-Type", "application/x-tar")

		return &resp, nil
	})

	// Change to a temporary directory for the file to be saved in
	dir, err := ioutil.TempDir("", "test-metadata-*")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	prevDir, _ := os.Getwd()
	err = os.Chdir(dir)
	require.Nil(t, err)
	//nolint:errcheck
	defer os.Chdir(prevDir)

	// Get the logs
	cmd, out, err := root.ExecuteTest("compose", "metadata", "b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7")
	defer out.Close()
	require.Nil(t, err)
	require.NotNil(t, out.Stdout)
	require.NotNil(t, out.Stderr)
	require.NotNil(t, cmd)
	assert.Equal(t, cmd, metadataCmd)
	stdout, err := ioutil.ReadAll(out.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stdout)
	stderr, err := ioutil.ReadAll(out.Stderr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(""), stderr)
	assert.Equal(t, "GET", mc.Req.Method)
	assert.Equal(t, "/api/v1/compose/metadata/b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7", mc.Req.URL.Path)

	_, err = os.Stat("b27c5a7b-d1f6-4c8c-8526-6d6de464f1c7-metadata.tar")
	assert.Nil(t, err)
}
