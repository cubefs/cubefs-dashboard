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

package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/util/log"

	"github.com/cubefs/cubefs-dashboard/backend/config"
	"github.com/cubefs/cubefs-dashboard/backend/helper/enums"
	"github.com/cubefs/cubefs-dashboard/backend/helper/middleware"
)

func RunHTTPServer() {
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Conf.Server.Port),
		Handler:        getHandler(),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		// service connections
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func getHandler() *gin.Engine {
	// set server mode
	gin.SetMode(enums.GetGinMode(config.Conf.Server.Mode))
	// set log mode
	log.SetOutputLevel(getLoggerLevel(config.Conf.Server.Mode))
	r := gin.New()
	// Global middleware
	r.Use(
		gin.Recovery(),
		middleware.Cors(),
		middleware.Logger(),
		middleware.InitSession(),
		middleware.Authorization,
		middleware.Default(),
		middleware.Session(),
		middleware.RecordOpLog(),
	)
	// init router
	Register(r)
	//helper.Must(auth.InitAuth())
	return r
}

func getLoggerLevel(mode string) log.Level {
	switch mode {
	case "dev":
		return log.Ldebug
	default:
		return log.Linfo
	}
}
