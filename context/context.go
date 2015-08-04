package context

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"
	"github.com/jeevatkm/urlite/model"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	VERSION = "1.0-beta"
)

var (
	dmutex sync.RWMutex
)

type App struct {
	Config    *Configuration
	Template  *template.Template
	Store     *sessions.CookieStore
	DBSession *mgo.Session
	Domains   map[string]model.Domain
}

type Configuration struct {
	AppName     string `toml:"app_name"`
	RunMode     string `toml:"run_mode"`
	BehindProxy bool   `toml:"behind_proxy"`
	Http        http
	Owner       ownerInfo
	DB          database `toml:"database"`
	Cookie      cookie
	Csrf        csrf
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
}

type http struct {
	IP   string `toml:"http_ip"`
	Port string `toml:"http_port"`
}

type database struct {
	Type   string
	Hosts  string
	DBName string `toml:"db_name"`
}

type cookie struct {
	MacSecret string `toml:"mac_secret"`
	Secure    bool
}

type csrf struct {
	Key    string
	Cookie string
	Header string
}

func init() {
	gob.Register(bson.ObjectId(""))

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
}

func InitContext(configFile *string) (ac *App) {
	log.Infof("Parsing application config '%s'", *configFile)

	ac = &App{}

	// Parsing configuration
	ac.Config = &Configuration{}
	if _, cerr := toml.DecodeFile(*configFile, &ac.Config); cerr != nil {
		log.Fatalf("Can't read configuration file: %s", cerr)
		panic(cerr)
	}

	// Logger configuration
	if ac.IsProdMode() {
		log.SetFormatter(new(log.JSONFormatter))
		log.SetLevel(log.ErrorLevel)
	}

	// Session store
	hash := sha256.New()
	io.WriteString(hash, ac.Config.Cookie.MacSecret)

	ac.Store = sessions.NewCookieStore(hash.Sum(nil))
	ac.Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600 * 24, // 24 hours
		Secure:   ac.Config.Cookie.Secure,
	}

	// Loading HTML templates
	var templates []string
	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			templates = append(templates, path)
		}
		return nil
	}

	ferr := filepath.Walk("./view", fn)
	if ferr != nil {
		log.Fatalf("Unable to load templates: %v", ferr)
		panic(ferr)
	}

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	ac.Template = template.Must(template.New("").Funcs(funcMap).ParseFiles(templates...))

	// Connecting Database
	var merr error
	ac.DBSession, merr = mgo.Dial(ac.Config.DB.Hosts)
	if merr != nil {
		log.Fatalf("DB connection error: %v", merr)
		panic(merr)
	}
	//app.DBSession.SetMode(mgo.Monotonic, true)

	ac.loadDomains()

	return
}

//
// App functions
//
func (a *App) IsDevMode() bool {
	return a.Config.RunMode == "DEV"
}

func (a *App) IsProdMode() bool {
	return a.Config.RunMode == "PROD"
}

func (a *App) Close() {
	log.Info("Application is cleaning up...")

	if a.DBSession != nil {
		a.DBSession.Close()
	}
}

func (a *App) GetDB(n string) *mgo.Database {
	// dbsession := a.DBSession.Clone()
	// defer dbsession.Close()
	// dbsession.DB(n)

	return a.DBSession.DB(n)
}

func (a *App) DB() *mgo.Database {
	return a.GetDB(a.Config.DB.DBName)
}

func (a *App) Parse(name string, data interface{}) (page string, err error) {
	var doc bytes.Buffer

	err = a.Template.ExecuteTemplate(&doc, name, data)
	if err != nil {
		return
	}

	page = doc.String()
	return
}

func (a *App) ParseF(data interface{}) (page string, err error) {
	page, err = a.Parse("layout/base", data)
	return
}

func (a *App) GetDomainLinkNum(name string) (linkNum int64, err error) {
	dmutex.Lock()
	defer dmutex.Unlock()

	if d, ok := a.Domains[name]; ok {
		linkNum = d.Count
		d.Count++
		return
	}

	err = errors.New("Domain not exists")
	return
}

/*
 * Private
 */
func (a *App) loadDomains() {
	domains, err := model.GetAllDomain(a.DB())
	if err != nil {
		log.Errorf("Error while loading domain details: %q", err)
		return
	}

	a.Domains = map[string]model.Domain{}
	for _, v := range domains {
		a.Domains[v.Name] = v
	}

	log.Infof("%d domains loaded by applciation", len(a.Domains))
}
