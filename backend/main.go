// Copyright 2023 The CubeFS Authors.
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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cubefs/cubefs-dashboard/backend/config"
	"github.com/cubefs/cubefs-dashboard/backend/helper"
	"github.com/cubefs/cubefs-dashboard/backend/model/migrate"
	"github.com/cubefs/cubefs-dashboard/backend/model/mysql"
	"github.com/cubefs/cubefs-dashboard/backend/router"
	"github.com/cubefs/cubefs-dashboard/backend/helper/crypt"
)

var confPath = flag.String("c", "", "please input your conf file path")
var enc = flag.String("e", "", "please input your encrypted string")

func main() {
	flag.Parse()
	if *enc != "" {
		str, err := crypt.Encrypt([]byte(*enc))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(str)
		os.Exit(0)
	}

	helper.Must(config.Init(*confPath))
	helper.Must(mysql.Init())
	helper.Must(migrate.Init())
	router.RunHTTPServer()
}
