package main

import (
	"flag"
	"net/http"
	"runtime"

	"github.com/jeevatkm/urlite/dashboard/context"
	"github.com/jeevatkm/urlite/dashboard/controller/api"
	"github.com/jeevatkm/urlite/dashboard/controller/web"
	"github.com/jeevatkm/urlite/dashboard/middleware"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"

	log "github.com/Sirupsen/logrus"
	groc "github.com/gorilla/context"
	jm "github.com/jeevatkm/middleware"
	ctr "github.com/jeevatkm/urlite/dashboard/controller"
	gw "github.com/zenazn/goji/web"
	gm "github.com/zenazn/goji/web/middleware"
)

const (
	VERSION = "0.1.0"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	configFile := flag.String("config", "/etc/urlite/urlite.conf", "Path to the configuration file")
	flag.Parse()

	ctx := context.Init(configFile)
	ctx.Version = VERSION
	log.Infof("Dashboard context loaded")

	goji.Abandon(gm.AutomaticOptions)
	goji.Use(middleware.AutomaticOptions)

	// Middleware
	goji.Use(groc.ClearHandler) // Gorilla session clear
	goji.Use(jm.Minify)
	goji.Use(middleware.AppInfo(ctx))
	goji.Use(middleware.Database(ctx))
	goji.Use(middleware.Session(ctx))
	goji.Use(middleware.Auth(ctx))

	/*
	 * Root level routes
	 */
	goji.Get("/", ctr.Handle{ctx, web.Home})
	goji.Post("/urlite", ctr.Handle{ctx, web.Urlite})
	goji.Get("/profile", ctr.Handle{ctx, web.Profile})
	goji.Post("/profile", ctr.Handle{ctx, web.ProfilePost})
	goji.Get("/urlites", ctr.Handle{ctx, web.Urlites})
	goji.Get("/urlites/data", ctr.Handle{ctx, web.UrlitesData})

	// Login routes
	goji.Get("/login", ctr.Handle{ctx, web.Login})
	goji.Post("/login", ctr.Handle{ctx, web.LoginPost})

	// Logout route
	goji.Get("/logout", ctr.Handle{ctx, web.Logout})

	/*
	 * Admin Rroutes
	 */
	ar := gw.New()
	ar.Use(gm.SubRouter)
	ar.Use(middleware.AutomaticOptions)
	ar.Use(middleware.AdminAuth(ctx))
	ar.Get("/", http.RedirectHandler("/admin/urlites", 301))
	ar.Get("/urlites", ctr.Handle{ctx, web.Urlites})
	ar.Get("/urlites/data", ctr.Handle{ctx, web.UrlitesData})
	ar.Get("/domains", ctr.Handle{ctx, web.Domains})
	ar.Get("/domains/validate", ctr.Handle{ctx, web.DomainsValidate})
	ar.Post("/domains", ctr.Handle{ctx, web.DomainsPost})
	ar.Get("/users", ctr.Handle{ctx, web.Users})
	ar.Get("/users/validate", ctr.Handle{ctx, web.UsersValidate})
	ar.Get("/users/data", ctr.Handle{ctx, web.UsersData})
	ar.Post("/users", ctr.Handle{ctx, web.UsersPost})

	goji.Handle("/admin/*", ar)
	goji.Get("/admin", http.RedirectHandler("/admin/urlites", 301))

	/*
	 * API routes
	 */
	apirt := gw.New()
	apirt.Use(gm.SubRouter)
	apirt.Use(middleware.RESTAutomaticOptions)
	apirt.Use(middleware.MediaTypeCheck)
	apirt.Use(middleware.ApiAuth(ctx))
	apirt.Post("/urlite", ctr.Handle{ctx, api.Urlite})
	apirt.Get("/stats", ctr.Handle{ctx, api.Stats})
	apirt.Get("/stats/:name", ctr.Handle{ctx, api.Stats})
	apirt.Get("/domains", ctr.Handle{ctx, api.Domains})
	apirt.Get("/domains/:name", ctr.Handle{ctx, api.Domains})

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
		ctx.Close()
	})

	// Assigning Ip and port config
	flag.Set("bind", ctx.Config.DashboardHttp.IP+":"+ctx.Config.DashboardHttp.Port)

	goji.Serve()
	log.Info("Application shutdown completed.")
}
