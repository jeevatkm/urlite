package main

import (
	"flag"
	"runtime"

	"github.com/jeevatkm/urlite/gateway/context"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"

	log "github.com/Sirupsen/logrus"
	ctr "github.com/jeevatkm/urlite/gateway/controller"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	configFile := flag.String("config", "/etc/urlite/urlite.conf", "Path to the configuration file")
	flag.Parse()

	ctx := context.Init(configFile)
	log.Infof("Urlite gateway context loaded")

	goji.Get("/", ctr.Handle{ctx, ctr.Home})
	goji.Get("/:urlite", ctr.Handle{ctx, ctr.Urlite})

	graceful.PostHook(func() {
		ctx.Close()
	})

	// Assigning Ip and port config
    flag.Set("bind", ctx.Config.GatewayHttp.IP+":"+ctx.Config.GatewayHttp.Port)

	goji.Serve()
	log.Info("Urlite gateway shutdown completed.")
}
