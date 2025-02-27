// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

// +build integration

package weldr

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProjects(t *testing.T) {
	projects, r, err := testState.client.ListProjects("")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, projects)
	assert.GreaterOrEqual(t, len(projects), 2)
}

func TestListProjectsDistro(t *testing.T) {
	projects, r, err := testState.client.ListProjects(testState.distros[0])
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, projects)
	assert.GreaterOrEqual(t, len(projects), 2)
}

func TestProjectsInfo(t *testing.T) {
	projects, r, err := testState.client.ProjectsInfo([]string{"bash"}, "")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, projects)
	assert.Equal(t, 1, len(projects))
}

func TestProjectsInfoDistro(t *testing.T) {
	projects, r, err := testState.client.ProjectsInfo([]string{"bash"}, testState.distros[0])
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, projects)
	assert.Equal(t, 1, len(projects))
}

func TestProjectsInfoMultiple(t *testing.T) {
	projects, r, err := testState.client.ProjectsInfo([]string{"bash", "filesystem", "tmux"}, "")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, projects)
	assert.Equal(t, 3, len(projects))
}

func TestProjectsInfoMultipleDistro(t *testing.T) {
	projects, r, err := testState.client.ProjectsInfo([]string{"bash", "filesystem", "tmux"}, testState.distros[0])
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, projects)
	assert.Equal(t, 3, len(projects))
}

func TestProjectsInfoOneError(t *testing.T) {
	projects, r, err := testState.client.ProjectsInfo([]string{"bart"}, "")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Nil(t, projects)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, "UnknownProject", r.Errors[0].ID)
	assert.Equal(t, "No packages have been found.", r.Errors[0].Msg)
}

func TestProjectsInfoOneErrorDistro(t *testing.T) {
	projects, r, err := testState.client.ProjectsInfo([]string{"bart"}, testState.distros[0])
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Nil(t, projects)
	assert.Equal(t, false, r.Status)
	assert.Equal(t, 1, len(r.Errors))
	assert.Equal(t, "UnknownProject", r.Errors[0].ID)
	assert.Equal(t, "No packages have been found.", r.Errors[0].Msg)
}

func TestProjectsInfoMultipleOneError(t *testing.T) {
	projects, r, err := testState.client.ProjectsInfo([]string{"bash", "filesystem", "bart"}, "")
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, projects)
	assert.Equal(t, 2, len(projects))
}

func TestProjectsInfoMultipleOneErrorDistro(t *testing.T) {
	projects, r, err := testState.client.ProjectsInfo([]string{"bash", "filesystem", "bart"}, testState.distros[0])
	require.Nil(t, err)
	require.Nil(t, r)
	require.NotNil(t, projects)
	assert.Equal(t, 2, len(projects))
}

func TestDepsolveProjects(t *testing.T) {
	deps, errors, err := testState.client.DepsolveProjects([]string{"bash"}, "")
	require.Nil(t, err)
	require.Nil(t, errors)
	require.NotNil(t, deps)
	require.GreaterOrEqual(t, len(deps), 1)

	// Decode a bit of the response for testing
	type projectDeps struct {
		Name string
	}

	// Encode it using json
	data := new(bytes.Buffer)
	err = json.NewEncoder(data).Encode(deps)
	require.Nil(t, err)

	// Decode the parts we care about
	var projects []projectDeps
	err = json.Unmarshal(data.Bytes(), &projects)
	require.Nil(t, err)

	// Do not depend on exact version numbers for dependencies, just check some package names
	var names []string
	for _, p := range projects {
		names = append(names, p.Name)
	}
	assert.Contains(t, names, "bash")
	assert.Contains(t, names, "filesystem")
}

func TestDepsolveProjectsDistro(t *testing.T) {
	deps, errors, err := testState.client.DepsolveProjects([]string{"bash"}, testState.distros[0])
	require.Nil(t, err)
	require.Nil(t, errors)
	require.NotNil(t, deps)
	require.GreaterOrEqual(t, len(deps), 1)

	// Decode a bit of the response for testing
	type projectDeps struct {
		Name string
	}

	// Encode it using json
	data := new(bytes.Buffer)
	err = json.NewEncoder(data).Encode(deps)
	require.Nil(t, err)

	// Decode the parts we care about
	var projects []projectDeps
	err = json.Unmarshal(data.Bytes(), &projects)
	require.Nil(t, err)

	// Do not depend on exact version numbers for dependencies, just check some package names
	var names []string
	for _, p := range projects {
		names = append(names, p.Name)
	}
	assert.Contains(t, names, "bash")
	assert.Contains(t, names, "filesystem")
}
