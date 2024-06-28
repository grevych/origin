// Description: This file is the entrypoint for origin.

// Package main implements the main entrypoint for the origin service.
package main

import (
	"context"
	"os"

	"github.com/grevych/gobox/pkg/app"
	"github.com/grevych/gobox/pkg/async"
	"github.com/grevych/gobox/pkg/env"
	"github.com/grevych/gobox/pkg/events"
	"github.com/grevych/gobox/pkg/log"
	"github.com/grevych/gobox/pkg/serviceactivities/automemlimit"
	"github.com/grevych/gobox/pkg/serviceactivities/gomaxprocs"
	"github.com/grevych/gobox/pkg/serviceactivities/shutdown"
	"github.com/grevych/gobox/pkg/trace"
	"github.com/grevych/origin/internal/http"
	"github.com/grevych/origin/internal/origin"
	"github.com/grevych/origin/internal/rpc"
)

// dependencies is a conglomerate struct used for injecting dependencies
// into service activities.
type dependencies struct {
	privateHTTP http.PrivateHTTPDependencies
	publicHTTP  http.PublicHTTPDependencies
	gRPC        rpc.GRPCDependencies
}

// main is the entrypoint for the origin service.
func main() { //nolint: funlen // Why: We can't dwindle this down anymore without adding complexity.
	exitCode := 1
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
		os.Exit(exitCode)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env.ApplyOverrides()
	app.SetName("origin")

	cfg, err := origin.LoadConfig(ctx)
	if err != nil {
		log.Error(ctx, "failed to load config", events.NewErrorInfo(err))
		return
	}

	if err := trace.InitTracer(ctx, "origin"); err != nil {
		log.Error(ctx, "tracing failed to start", events.NewErrorInfo(err))
		return
	}
	defer trace.CloseTracer(ctx)

	// Initialize variable for service activity dependency injection.
	var deps dependencies

	log.Info(ctx, "starting", app.Info(), cfg, log.F{"app.pid": os.Getpid()})

	acts := []async.Runner{
		shutdown.New(),
		gomaxprocs.New(),
		automemlimit.New(),
		http.NewPrivateHTTPService(cfg, &deps.privateHTTP),
		http.NewPublicHTTPService(cfg, &deps.publicHTTP),
		// rpc.NewGRPCService(cfg, &deps.gRPC),
	}

	err = async.RunGroup(acts).Run(ctx)
	if shutdown.HandleShutdownConditions(ctx, err) {
		exitCode = 0
	}
}
