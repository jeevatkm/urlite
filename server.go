package main

import (
	"flag"
	"net/http"
	"runtime"

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

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	configFile := flag.String("config", "/etc/urlite/urlite.conf", "Path to the configuration file")
	flag.Parse()

	context := context.InitContext(configFile)
	log.Infof("Application context loaded")

	goji.Abandon(gm.AutomaticOptions)
	goji.Use(middleware.AutomaticOptions)

	// Middleware
	goji.Use(groc.ClearHandler) // Gorilla session clear
	goji.Use(middleware.AppInfo(context))
	goji.Use(middleware.Session(context))
	goji.Use(middleware.Auth(context))

	/*
	 * Root level routes
	 */
	goji.Get("/", ctr.Handle{context, web.Home})
	goji.Get("/profile", ctr.Handle{context, web.Profile})
	goji.Post("/profile", ctr.Handle{context, web.ProfilePost})
	goji.Get("/dashboard", ctr.Handle{context, web.Dashboard})

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
	ar.Use(middleware.AutomaticOptions)
	ar.Use(middleware.AdminAuth(context))
	ar.Get("/", http.RedirectHandler("/admin/dashboard", 301))
	ar.Get("/dashboard", ctr.Handle{context, web.Dashboard})
	ar.Get("/domains", ctr.Handle{context, web.Domains})
	ar.Post("/domains", ctr.Handle{context, web.DomainsPost})
	ar.Get("/users", ctr.Handle{context, web.Users}) 
	ar.Get("/users/data", ctr.Handle{context, web.UsersData})
	ar.Post("/users", ctr.Handle{context, web.UsersPost})

	goji.Handle("/admin/*", ar)
	goji.Get("/admin", http.RedirectHandler("/admin/dashboard", 301))

	/*
	 * API routes
	 */
	apirt := gw.New()
	apirt.Use(gm.SubRouter)
	apirt.Use(middleware.RESTAutomaticOptions)
	apirt.Use(middleware.MediaTypeCheck)
	apirt.Use(middleware.ApiAuth(context))
	apirt.Post("/urlite", ctr.Handle{context, api.Urlite})
	apirt.Get("/domains", ctr.Handle{context, api.Domains})
	apirt.Get("/domains/:name", ctr.Handle{context, api.Domains})

	goji.Handle("/api/*", apirt)

	/*
	 * Static Files handling
	 */
	goji.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	goji.Get("/robots.txt", http.FileServer(http.Dir("./static")))
	goji.Get("/favicon.ico", http.FileServer(http.Dir("./static/images")))

	// goji.NotFound(func(w http.ResponseWriter, r *http.Request) {
	// 	http.Error(w, "Umm... have you tried turning it off and on again?", 404)
	// })

	graceful.PostHook(func() {
		context.Close()
	})

	// Assigning Ip and port config
	flag.Set("bind", context.Config.Http.IP+":"+context.Config.Http.Port)

	goji.Serve()
	log.Info("Application shutdown completed.")
}
