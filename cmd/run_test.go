// Copyright (c) Alex Ellis 2017-2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types/mount"
)

var validRequest = TaskRequest{
	Image:         "input-image",
	Networks:      []string{"net1", "net2"},
	Constraints:   []string{"node.id=2ivku8v2gvtg4", "engine.labels.operatingsystem==ubuntu 14.04"},
	EnvVars:       []string{"ev1=val1", "ev2=val2"},
	Mounts:        []string{"hostVol1=taskVol1", "hostVol2=taskVol2"},
	EnvFiles:      []string{".env1", ".env2"},
	Secrets:       []string{"secret1", "secret1"},
	ShowLogs:      true,
	Timeout:       "12",
	RemoveService: true,
	RegistryAuth:  "true",
	Command:       "echo 'some output'",
}

func TestMakeServiceSpecValid(t *testing.T) {
	spec := makeServiceSpec(validRequest)
	if spec.TaskTemplate.ContainerSpec.Image != "input-image" {
		t.Errorf("Container spec image should be %s, was %s", validRequest.Image, spec.TaskTemplate.ContainerSpec.Image)
	}

	if !reflect.DeepEqual(spec.Networks, []string{"net1", "net2"}) {
		t.Errorf("Container spec networks should be %s, was %s", validRequest.Networks, spec.Networks)
	}

	if !reflect.DeepEqual(spec.TaskTemplate.ContainerSpec.Env, []string{"ev1=val1", "ev2=val2"}) {
		t.Errorf("Container spec env should be %s, was %s", validRequest.EnvVars, spec.TaskTemplate.ContainerSpec.Env)
	}

	expectedMounts := []mount.Mount{
		{Source: "hostVol1", Target: "taskVol1"},
		{Source: "hostVol2", Target: "taskVol2"},
	}
	if !reflect.DeepEqual(spec.TaskTemplate.ContainerSpec.Mounts, expectedMounts) {

		t.Error("Container spec mounts should include:")
		for _, m := range expectedMounts {
			t.Errorf("{Source: %s, Target: %s}", m.Source, m.Target)
		}
		t.Error("But contained instead:")
		for _, m := range spec.TaskTemplate.ContainerSpec.Mounts {
			t.Errorf("{Source: %s, Target: %s}", m.Source, m.Target)
		}
	}
}
