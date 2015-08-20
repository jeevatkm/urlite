package context

import (
	"encoding/gob"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jeevatkm/urlite/model"
	"github.com/zenazn/goji/web"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
)

type Context struct {
	Config    *Configuration
	DBSession *mgo.Session
	Domains   map[string]*model.Domain
}

type Configuration struct {
	AppName       string `toml:"app_name"`
	RunMode       string `toml:"run_mode"`
	BehindProxy   bool   `toml:"behind_proxy"`
	GatewayHttp   http   `toml:"gateway"`
	DashboardHttp http   `toml:"dashboard"`
	Owner         ownerInfo
	DB            database `toml:"database"`
	Cookie        cookie
	Security      security
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
}

type http struct {
	IP   string `toml:"http_ip"`
	Port string `toml:"http_port"`
}

type security struct {
	PublicPath []string `toml:"public_path"`
}

type database struct {
	Type             string
	Hosts            string
	Name             string `toml:"db_name"`
	User             string `toml:"db_user"`
	Password         string `toml:"db_password"`
	LinkSyncInternal int    `toml:"sync_link_num_interval"`
}

type cookie struct {
	MacSecret string `toml:"mac_secret"`
	Secure    bool
}

func init() {
	gob.Register(bson.ObjectId(""))

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
}

func (c *Context) Init(configFile *string) {
	c.readConfig(configFile)
	c.initLogger()
	c.dailDatabase()

	c.Domains = map[string]*model.Domain{}
	c.loadDomains()
}

func (c *Context) IsDevMode() bool {
	return c.Config.RunMode == "DEV"
}

func (c *Context) IsProdMode() bool {
	return c.Config.RunMode == "PROD"
}

func (c *Context) Close() {
	log.Info("this is parent close method")

	if c.DBSession != nil {
		c.DBSession.Close()
	}
}

func (c *Context) AddDomain(d *model.Domain) {
	c.Domains[d.Name] = &model.Domain{ID: d.ID,
		Name:            d.Name,
		Scheme:          d.Scheme,
		Salt:            d.Salt,
		IsDefault:       d.IsDefault,
		LinkCount:       d.LinkCount,
		CustomLinkCount: d.CustomLinkCount,
		TrackName:       d.TrackName}

	// // Initializing Hash generater for domain
	// hd := hash.NewData()
	// hd.Salt = d.Salt
	// hd.MinLength = 4
	// a.HashGen[d.Name] = hash.NewWithData(hd)

	// // Initial linkstate for domain
	// a.LinkState[d.Name] = false
}

func (c *Context) readConfig(configFile *string) {
	log.Infof("Reading configuration - '%s'", *configFile)
	c.Config = &Configuration{}
	if _, err := toml.DecodeFile(*configFile, &c.Config); err != nil {
		log.Fatalf("Can't read configuration file: %s", err)
		panic(err)
	}
}

func (c *Context) dailDatabase() {
	var err error
	c.DBSession, err = mgo.Dial(c.Config.DB.Hosts)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
		panic(err)
	}
}

func (c *Context) initLogger() {
	if c.IsProdMode() {
		log.SetFormatter(new(log.JSONFormatter))
		log.SetLevel(log.ErrorLevel)
	} else if c.IsDevMode() {
		//log.SetFormatter(new(log.JSONFormatter))
		log.SetLevel(log.DebugLevel)
	}
}

func (c *Context) loadDomains() {
	domains, err := model.GetAllDomain(c.db())
	if err != nil {
		log.Errorf("Error while loading domain details: %q", err)
		return
	}

	for _, v := range domains {
		c.AddDomain(&v)
	}

	log.Infof("%d domains loaded by app context", len(c.Domains))
}

func (c *Context) DB(gc *web.C) *mgo.Database {
	return gc.Env["DB"].(*mgo.Database)
}

func (c *Context) db() *mgo.Database {
	return c.DBSession.DB(c.Config.DB.Name)
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

// //
// // App functions
// //
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

// /*
//  * Private
//  */

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
