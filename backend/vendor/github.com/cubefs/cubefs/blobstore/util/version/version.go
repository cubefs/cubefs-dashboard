// Copyright 2022 The CubeFS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package version

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	version string      = ""
	fPerm   os.FileMode = 0o600
)

func init() {
	if len(os.Args) > 1 && os.Args[1] == "-version" {
		fmt.Println("version:", version)
		os.Exit(0)
	}
	writeFile(".version", version)
}

func Version() string {
	return version
}

func writeFile(fname, field string) {
	if field != "" {
		ioutil.WriteFile(fname, []byte(field), fPerm)
	}
}
