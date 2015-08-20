package context

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/jeevatkm/urlite/context"

	"github.com/gorilla/sessions"
	"github.com/jeevatkm/urlite/dashboard/tpl"
	"github.com/jeevatkm/urlite/model"

	log "github.com/Sirupsen/logrus"
	hash "github.com/speps/go-hashids"
)

var (
	dmutex sync.RWMutex
	dbSync *time.Ticker
)

type Context struct {
	context.Context
	Store          *sessions.CookieStore
	HashGen        map[string]*hash.HashID
	LinkState      map[string]bool
	NoSessionRoute []string
}

func Init(configFile *string) (c *Context) {
	c = &Context{}

	c.Context.Init(configFile)

	c.NoSessionRoute = append(c.Config.Security.PublicPath, "/api")

	// Session store
	chash := sha256.New()
	io.WriteString(chash, c.Config.Cookie.MacSecret)

	c.Store = sessions.NewCookieStore(chash.Sum(nil))
	c.Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600 * 24, // 24 hours
		Secure:   c.Config.Cookie.Secure,
	}

	// Loading HTML templates
	c.LoadTemplates("./view", &tpl.FuncMap)

	c.HashGen = map[string]*hash.HashID{}
	c.LinkState = map[string]bool{}
	c.registerHash()

	c.startSyncTasks()

	return
}

func (c *Context) AddDomain(d *model.Domain) {
	c.Context.AddDomain(d)

	// Initializing Hash generater for domain
	hd := hash.NewData()
	hd.Salt = d.Salt
	hd.MinLength = 4
	c.HashGen[d.Name] = hash.NewWithData(hd)

	// Initial linkstate for domain
	c.LinkState[d.Name] = false
}

func (c *Context) Close() {
	c.syncDomainCountToDB()
	c.stopSyncTasks()

	log.Info("Cleaning up...")
	c.Context.Close()
}

func (c *Context) Parse(name string, data interface{}) (page string, err error) {
	if c.IsDevMode() {
		c.LoadTemplates("./view", &tpl.FuncMap)
	}

	var doc bytes.Buffer
	err = c.Template.ExecuteTemplate(&doc, name, data)
	if err != nil {
		return
	}

	page = doc.String()
	return
}

func (c *Context) ParseF(data interface{}) (page string, err error) {
	page, err = c.Parse("layout/base", data)
	return
}

func (c *Context) CheckDomainTrackName(name string) bool {
	for _, v := range c.Domains {
		if v.TrackName == name {
			return true
		}
	}

	return false
}

func (c *Context) GetDomainDetail(name string) (*model.Domain, error) {
	if d, ok := c.Domains[name]; ok {
		return d, nil
	}
	return nil, errors.New("Domain not exists")
}

func (c *Context) GetDomainLinkNum(name string) (n int64) {
	dmutex.Lock()
	defer dmutex.Unlock()

	n = c.Domains[name].LinkCount
	log.Debugf("New generated link count: %d", n)
	c.Domains[name].LinkCount++
	c.LinkState[name] = true
	log.Debugf("Next new generate link count: %d", c.Domains[name].LinkCount)
	return
}

func (c *Context) IncDomainCustomLink(name string) {
	dmutex.Lock()
	defer dmutex.Unlock()

	log.Debugf("Current Custom link count: %d", c.Domains[name].CustomLinkCount)
	c.Domains[name].CustomLinkCount++
	c.LinkState[name] = true
	log.Debugf("New custom link count: %d", c.Domains[name].CustomLinkCount)
}

func (c *Context) GetUrliteID(dn string, n int64) (ul string, err error) {
	ul, err = c.HashGen[dn].EncodeInt64([]int64{n})
	return
}

func (c *Context) db() *mgo.Database {
	return c.DBSession.DB(c.Config.DB.Name)
}

func (c *Context) registerHash() {
	for _, d := range c.Domains {
		// Initializing Hash generater for domain
		hd := hash.NewData()
		hd.Salt = d.Salt
		hd.MinLength = 4
		c.HashGen[d.Name] = hash.NewWithData(hd)

		// Initial linkstate for domain
		c.LinkState[d.Name] = false
	}

}

func (c *Context) startSyncTasks() {
	dbSync = time.NewTicker(time.Minute * time.Duration(c.Config.DB.LinkSyncInternal))

	// Listening for channel input
	go func(cc *Context) {
		for {
			select {
			case <-dbSync.C:
				cc.syncDomainCountToDB()
			}
		}
	}(c)
}

func (c *Context) stopSyncTasks() {
	dbSync.Stop()
}

func (c *Context) syncDomainCountToDB() {
	//log.Infof("Starting Domain link count sync at %v", time.Now())

	for k, v := range c.LinkState {
		if v {
			log.Infof("Syncing domain: %v", k)
			err := model.UpdateDomainLinkCount(c.db(), c.Domains[k])
			if err != nil {
				log.Errorf("Error occurred while syncing domain urlite count for %v and error is %q", k, err)
			}

			c.LinkState[k] = false
		}
	}

	//log.Infof("Completed Domain link count sync at %v", time.Now())
}

// func InitContext(configFile *string) (ac *App) {
// 	log.Infof("Parsing application config '%s'", *configFile)

// 	ac = &App{}

// 	ac.parseConfig(configFile)

// 	// Parsing configuration
// 	// ac.Config = &Configuration{}
// 	// if _, cerr := toml.DecodeFile(*configFile, &ac.Config); cerr != nil {
// 	// 	log.Fatalf("Can't read configuration file: %s", cerr)
// 	// 	panic(cerr)
// 	// }

// 	// Logger configuration
// 	if ac.IsProdMode() {
// 		log.SetFormatter(new(log.JSONFormatter))
// 		log.SetLevel(log.ErrorLevel)
// 	}

// 	// No Session Routes
// 	ac.NoSessionRoute = append(ac.Config.Security.PublicPath, "/api")

// 	// Session store
// 	chash := sha256.New()
// 	io.WriteString(chash, ac.Config.Cookie.MacSecret)

// 	ac.Store = sessions.NewCookieStore(chash.Sum(nil))
// 	ac.Store.Options = &sessions.Options{
// 		Path:     "/",
// 		HttpOnly: true,
// 		MaxAge:   3600 * 24, // 24 hours
// 		Secure:   ac.Config.Cookie.Secure,
// 	}

// 	// Loading HTML templates
// 	ac.loadTemplates()

// 	// Connecting Database
// 	ac.initDatabase()
// 	// var merr error
// 	// ac.DBSession, merr = mgo.Dial(ac.Config.DB.Hosts)
// 	// if merr != nil {
// 	// 	log.Fatalf("DB connection error: %v", merr)
// 	// 	panic(merr)
// 	// }
// 	//app.DBSession.SetMode(mgo.Monotonic, true)

// 	ac.Domains = map[string]*model.Domain{}
// 	ac.HashGen = map[string]*hash.HashID{}
// 	ac.LinkState = map[string]bool{}
// 	ac.loadDomains()

// 	ac.startSyncTasks()

// 	return
// }

//
// App functions
//
// func (a *App) IsDevMode() bool {
// 	return a.Config.RunMode == "DEV"
// }

// func (a *App) IsProdMode() bool {
// 	return a.Config.RunMode == "PROD"
// }

// func (a *App) GatewayClose() {
// 	log.Info("Application is cleaning up...")

// 	if a.DBSession != nil {
// 		a.DBSession.Close()
// 	}
// }

// func (a *App) Close() {
// 	a.syncDomainCountToDB()

// 	a.stopSyncTasks()

// 	log.Info("Application is cleaning up...")

// 	if a.DBSession != nil {
// 		a.DBSession.Close()
// 	}
// }

// func (a *App) DB(c *web.C) *mgo.Database {
// 	return c.Env["DB"].(*mgo.Database)
// }

// func (a *App) db() *mgo.Database {
// 	return a.DBSession.DB(a.Config.DB.DBName)
// }

// func (a *App) Parse(name string, data interface{}) (page string, err error) {
// 	if a.IsDevMode() {
// 		a.loadTemplates()
// 	}

// 	var doc bytes.Buffer
// 	err = a.Template.ExecuteTemplate(&doc, name, data)
// 	if err != nil {
// 		return
// 	}

// 	page = doc.String()
// 	return
// }

// func (a *App) ParseF(data interface{}) (page string, err error) {
// 	page, err = a.Parse("layout/base", data)
// 	return
// }

// func (a *App) GetDomainDetail(name string) (*model.Domain, error) {
// 	if d, ok := a.Domains[name]; ok {
// 		return d, nil
// 	}
// 	return nil, errors.New("Domain not exists")
// }

// func (a *App) GetDomainLinkNum(name string) (n int64) {
// 	dmutex.Lock()
// 	defer dmutex.Unlock()

// 	n = a.Domains[name].LinkCount
// 	log.Debugf("New generated link count: %d", n)
// 	a.Domains[name].LinkCount++
// 	a.LinkState[name] = true
// 	log.Debugf("Next new generate link count: %d", a.Domains[name].LinkCount)
// 	return
// }

// func (a *App) IncDomainCustomLink(name string) {
// 	dmutex.Lock()
// 	defer dmutex.Unlock()

// 	log.Debugf("Current Custom link count: %d", a.Domains[name].CustomLinkCount)
// 	a.Domains[name].CustomLinkCount++
// 	a.LinkState[name] = true
// 	log.Debugf("New custom link count: %d", a.Domains[name].CustomLinkCount)
// }

// func (a *App) GetUrliteID(dn string, n int64) (ul string, err error) {
// 	ul, err = a.HashGen[dn].EncodeInt64([]int64{n})
// 	return
// }

// func (a *App) AddDomain(d *model.Domain) {
// 	a.Domains[d.Name] = &model.Domain{ID: d.ID,
// 		Name:            d.Name,
// 		Scheme:          d.Scheme,
// 		Salt:            d.Salt,
// 		IsDefault:       d.IsDefault,
// 		LinkCount:       d.LinkCount,
// 		CustomLinkCount: d.CustomLinkCount,
// 		TrackName:       d.TrackName}

// 	// Initializing Hash generater for domain
// 	hd := hash.NewData()
// 	hd.Salt = d.Salt
// 	hd.MinLength = 4
// 	a.HashGen[d.Name] = hash.NewWithData(hd)

// 	// Initial linkstate for domain
// 	a.LinkState[d.Name] = false
// }

// func (a *App) CheckDomainTrackName(name string) bool {
// 	//dTrack := "urlite_" + name
// 	for _, v := range a.Domains {
// 		if v.TrackName == name {
// 			return true
// 		}
// 	}

// 	return false
// }

/*
 * Private
 */

// // Parsing 'urlite.conf' configuration
// func (a *App) parseConfig(configFile *string) {
// 	a.Config = &Configuration{}
// 	if _, cerr := toml.DecodeFile(*configFile, &a.Config); cerr != nil {
// 		log.Fatalf("Can't read configuration file: %s", cerr)
// 		panic(cerr)
// 	}
// }

// // Connecting Database
// func (a *App) initDatabase() {
// 	var merr error
// 	a.DBSession, merr = mgo.Dial(a.Config.DB.Hosts)
// 	if merr != nil {
// 		log.Fatalf("DB connection error: %v", merr)
// 		panic(merr)
// 	}
// }

// func (a *App) loadTemplates() {
// 	var templates []string
// 	fn := func(path string, f os.FileInfo, err error) error {
// 		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
// 			templates = append(templates, path)
// 		}
// 		return nil
// 	}

// 	ferr := filepath.Walk("./view", fn)
// 	if ferr != nil {
// 		log.Fatalf("Unable to load templates: %v", ferr)
// 		panic(ferr)
// 	}

// 	a.Template = template.Must(template.New("").Funcs(tpl.FuncMap).ParseFiles(templates...))
// }

// func (a *App) loadDomains() {
// 	domains, err := model.GetAllDomain(a.db())
// 	if err != nil {
// 		log.Errorf("Error while loading domain details: %q", err)
// 		return
// 	}

// 	for _, v := range domains {
// 		a.AddDomain(&v)
// 	}

// 	log.Infof("%d domains loaded and it's hash generater have been initialized", len(a.Domains))
// }

// func (a *App) startSyncTasks() {
// 	dbSync = time.NewTicker(time.Minute * time.Duration(a.Config.DB.LinkSyncInternal))

// 	// Listening for channel input
// 	go func(aa *App) {
// 		for {
// 			select {
// 			case <-dbSync.C:
// 				aa.syncDomainCountToDB()
// 			}
// 		}
// 	}(a)
// }

// func (a *App) stopSyncTasks() {
// 	dbSync.Stop()
// }

// func (a *App) syncDomainCountToDB() {
// 	//log.Infof("Starting Domain link count sync at %v", time.Now())

// 	for k, v := range a.LinkState {
// 		if v {
// 			log.Infof("Syncing domain: %v", k)
// 			err := model.UpdateDomainLinkCount(a.db(), a.Domains[k])
// 			if err != nil {
// 				log.Errorf("Error occurred while syncing domain urlite count for %v and error is %q", k, err)
// 			}

// 			a.LinkState[k] = false
// 		}
// 	}

// 	//log.Infof("Completed Domain link count sync at %v", time.Now())
// }
