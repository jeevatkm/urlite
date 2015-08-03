package main

import (
	"flag"
	"net/http"

	"github.com/jeevatkm/urlite/context"
	"github.com/jeevatkm/urlite/controller/api"
	"github.com/jeevatkm/urlite/controller/web"
	"github.com/jeevatkm/urlite/middleware"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"

	log "github.com/Sirupsen/logrus"
	groc "github.com/gorilla/context"
	ctr "github.com/jeevatkm/urlite/controller"
	gw "github.com/zenazn/goji/web"
	gm "github.com/zenazn/goji/web/middleware"
)

func main() {
	configFile := flag.String("config", "/etc/urlite/urlite.conf", "Path to the configuration file")
	flag.Parse()

	context := context.InitContext(configFile)
	log.Infof("Application context loaded")

	// Middleware
	goji.Use(groc.ClearHandler) // Gorilla session clear
	goji.Use(middleware.AppInfo(context))
	goji.Use(middleware.Session(context))
	goji.Use(middleware.Auth(context))

	/*
	 * Root level routes
	 */
	goji.Get("/", ctr.Handle{context, web.Home})

	// Login routes
	goji.Get("/login", ctr.Handle{context, web.Login})
	goji.Post("/login", ctr.Handle{context, web.LoginPost})

	// Logout route
	goji.Get("/logout", ctr.Handle{context, web.Logout})

	/*
	 * Admin Rroutes
	 */
	ar := gw.New()
	ar.Use(gm.SubRouter)
	ar.Use(middleware.AdminAuth(context))
	ar.Get("/", ctr.Handle{context, web.Dashboard})

	goji.Handle("/admin/*", ar)
	goji.Get("/admin", http.RedirectHandler("/admin/", 301))

	/*
	 * API routes
	 */
	apirt := gw.New()
	apirt.Use(gm.SubRouter)
	apirt.Use(middleware.ApiAuth(context))
	apirt.Post("/shorten", ctr.Handle{context, api.Shorten})

	goji.Handle("/api/*", apirt)

	/*
	 * Static Files handling
	 */
	goji.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	goji.Get("/robots.txt", http.FileServer(http.Dir("./static")))
	goji.Get("/favicon.ico", http.FileServer(http.Dir("./static/images")))

	graceful.PostHook(func() {
		context.Close()
	})

	// Assigning Ip and port config
	flag.Set("bind", context.Config.HttpIp+":"+context.Config.HttpPort)

	goji.Serve()
	log.Info("Application shutdown completed.")
}
