// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package build

import (
	"fmt"
	"os"
	"strings"

	"launchpad.net/snappy/_integration-tests/testutils"
)

const (
	buildTestCmd      = "go test -c ./_integration-tests/tests"
	buildSnappyCliCmd = "go build -o " + testsBinDir + "snappy ./cmd/snappy"
	buildSnapdCmd     = "go build -o " + testsBinDir + "snapd ./cmd/snapd"

	// IntegrationTestName is the name of the test binary.
	IntegrationTestName = "integration.test"
	defaultGoArm        = "7"

	testsBinDir = "_integration-tests/bin/"
)

var (
	// dependency aliasing
	execCommand      = testutils.ExecCommand
	prepareTargetDir = testutils.PrepareTargetDir
	osRename         = os.Rename
	osSetenv         = os.Setenv
	osGetenv         = os.Getenv
)

// Assets builds the snappy and integration tests binaries for the target
// architecture.
func Assets(useSnappyFromBranch bool, arch string) {
	prepareTargetDir(testsBinDir)

	if useSnappyFromBranch {
		// FIXME We need to build an image that has the snappy from the branch
		// installed. --elopio - 2015-06-25.
		buildSnappyCLI(arch)
		buildSnapd(arch)
	}
	buildTests(arch)
}

func buildSnappyCLI(arch string) {
	fmt.Println("Building snappy CLI...")
	// On the root of the project we have a directory called snappy, so we
	// output the binary for the tests in the tests directory.
	goCall(arch, buildSnappyCliCmd)
}

func buildSnapd(arch string) {
	fmt.Println("Building snapd...")
	// On the root of the project we have a directory called snappy, so we
	// output the binary for the tests in the tests directory.
	goCall(arch, buildSnapdCmd)
}

func buildTests(arch string) {
	fmt.Println("Building tests...")

	goCall(arch, buildTestCmd)
	// XXX Go test 1.3 does not have the output flag, so we move the
	// binaries after they are generated.
	osRename("tests.test", testsBinDir+IntegrationTestName)
}

func goCall(arch string, cmd string) {
	if arch != "" {
		defer osSetenv("GOARCH", osGetenv("GOARCH"))
		osSetenv("GOARCH", arch)
		if arch == "arm" {
			envs := map[string]string{
				"GOARM":       defaultGoArm,
				"CGO_ENABLED": "1",
				"CC":          "arm-linux-gnueabihf-gcc",
			}
			for env, value := range envs {
				defer osSetenv(env, osGetenv(env))
				osSetenv(env, value)
			}
		}
	}
	execCommand(strings.Fields(cmd)...)
}
