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
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"
	"github.com/jeevatkm/urlite/model"
	"github.com/jeevatkm/urlite/tpl"

	log "github.com/Sirupsen/logrus"
	hash "github.com/speps/go-hashids"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	VERSION = "0.1.0"
)

var (
	dmutex sync.RWMutex
	dbSync *time.Ticker
)

type App struct {
	Config    *Configuration
	Template  *template.Template
	Store     *sessions.CookieStore
	DBSession *mgo.Session
	Domains   map[string]*model.Domain
	HashGen   map[string]*hash.HashID
	LinkState map[string]bool
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
	Type             string
	Hosts            string
	DBName           string `toml:"db_name"`
	LinkSyncInternal int    `toml:"sync_link_num_interval"`
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
	chash := sha256.New()
	io.WriteString(chash, ac.Config.Cookie.MacSecret)

	ac.Store = sessions.NewCookieStore(chash.Sum(nil))
	ac.Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600 * 24, // 24 hours
		Secure:   ac.Config.Cookie.Secure,
	}

	// Loading HTML templates
	// var templates []string
	// fn := func(path string, f os.FileInfo, err error) error {
	// 	if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
	// 		templates = append(templates, path)
	// 	}
	// 	return nil
	// }

	// ferr := filepath.Walk("./view", fn)
	// if ferr != nil {
	// 	log.Fatalf("Unable to load templates: %v", ferr)
	// 	panic(ferr)
	// }

	// funcMap := template.FuncMap{
	// 	"safeHTML": func(s string) template.HTML {
	// 		return template.HTML(s)
	// 	},
	// 	"isset": func(a interface{}, key interface{}) bool {
	// 		av := reflect.ValueOf(a)
	// 		kv := reflect.ValueOf(key)

	// 		switch av.Kind() {
	// 		case reflect.Array, reflect.Chan, reflect.Slice:
	// 			if int64(av.Len()) > kv.Int() {
	// 				return true
	// 			}
	// 		case reflect.Map:
	// 			if kv.Type() == av.Type().Key() {
	// 				return av.MapIndex(kv).IsValid()
	// 			}
	// 		}

	// 		return false
	// 	},
	// 	"safeCSS": func(text string) template.CSS {
	// 		return template.CSS(text)
	// 	},
	// 	"safeURL": func(text string) template.URL {
	// 		return template.URL(text)
	// 	},
	// 	"getEnv": func() string {
	// 		return ac.Config.RunMode
	// 	},
	// 	"frdlyDateTime": func(t time.Time) string {
	// 		if t.IsZero() {
	// 			return ""
	// 		}
	// 		if t.Year() == time.Now().Year() {
	// 			return t.Format("Jan 2, 3:04:05 pm")
	// 		}
	// 		return t.Format("Jan 2, 2006, 3:04:05 pm")
	// 	},
	// 	"toCommaStr": func(v []string) string {
	// 		return strings.Join(v, ", ")
	// 	},
	// 	"add": func(a, b interface{}) (interface{}, error) {
	// 		av := reflect.ValueOf(a)
	// 		bv := reflect.ValueOf(b)

	// 		if av.Kind() != bv.Kind() {
	// 			return nil, errors.New("Different kinds, can't add them.")
	// 		}

	// 		switch av.Kind() {
	// 		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// 			return av.Int() + bv.Int(), nil
	// 		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	// 			return av.Uint() + bv.Uint(), nil
	// 		case reflect.Float32, reflect.Float64:
	// 			return av.Float() + bv.Float(), nil
	// 		// case reflect.String:
	// 		//     return av.String() + bv.String(), nil
	// 		default:
	// 			return nil, errors.New("Type does not support addition.")
	// 		}
	// 	},
	// 	"num2str": func(n int64, sep rune) string {
	// 	    s := strconv.FormatInt(n, 10)

	// 	    startOffset := 0
	// 	    var buff bytes.Buffer
	// 	    if n < 0 {
	// 	        startOffset = 1
	// 	        buff.WriteByte('-')
	// 	    }

	// 	    l := len(s)
	// 	    commaIndex := 3 - ((l - startOffset) % 3)
	// 	    if (commaIndex == 3) {
	// 	        commaIndex = 0
	// 	    }

	// 	    for i := startOffset; i < l; i++ {
	// 	        if (commaIndex == 3) {
	// 	            buff.WriteRune(sep)
	// 	            commaIndex = 0
	// 	        }
	// 	        commaIndex++

	// 	        buff.WriteByte(s[i])
	// 	    }

	// 	    return buff.String()
	// 	},
	// }

	// ac.Template = template.Must(template.New("").Funcs(funcMap).ParseFiles(templates...))

	ac.loadTemplates()

	// Connecting Database
	var merr error
	ac.DBSession, merr = mgo.Dial(ac.Config.DB.Hosts)
	if merr != nil {
		log.Fatalf("DB connection error: %v", merr)
		panic(merr)
	}
	//app.DBSession.SetMode(mgo.Monotonic, true)

	ac.Domains = map[string]*model.Domain{}
	ac.HashGen = map[string]*hash.HashID{}
	ac.LinkState = map[string]bool{}
	ac.loadDomains()

	ac.startSyncTasks()

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
	a.syncDomainCountToDB()

	a.stopSyncTasks()

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
	if a.IsDevMode() {
		a.loadTemplates()
	}

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

func (a *App) GetDomainDetail(name string) (*model.Domain, error) {
	if d, ok := a.Domains[name]; ok {
		return d, nil
	}
	return nil, errors.New("Domain not exists")
}

func (a *App) GetDomainLinkNum(name string) (n int64) {
	dmutex.Lock()
	defer dmutex.Unlock()

	n = a.Domains[name].LinkCount
	log.Debugf("New generated link count: %d", n)
	a.Domains[name].LinkCount++
	a.LinkState[name] = true
	log.Debugf("Next new generate link count: %d", a.Domains[name].LinkCount)
	return
}

func (a *App) IncDomainCustomLink(name string) {
	dmutex.Lock()
	defer dmutex.Unlock()

	log.Debugf("Current Custom link count: %d", a.Domains[name].CustomLinkCount)
	a.Domains[name].CustomLinkCount++
	a.LinkState[name] = true
	log.Debugf("New custom link count: %d", a.Domains[name].CustomLinkCount)
}

func (a *App) GetUrliteID(dn string, n int64) (ul string, err error) {
	ul, err = a.HashGen[dn].EncodeInt64([]int64{n})
	return
}

func (a *App) AddDomain(d *model.Domain) {
	a.Domains[d.Name] = &model.Domain{ID: d.ID,
		Name:            d.Name,
		Scheme:          d.Scheme,
		Salt:            d.Salt,
		LinkCount:       d.LinkCount,
		CustomLinkCount: d.CustomLinkCount,
		CollName:        d.CollName,
		StatsCollName:   d.StatsCollName}

	// Initializing Hash generater for domain
	hd := hash.NewData()
	hd.Salt = d.Salt
	hd.MinLength = 5
	a.HashGen[d.Name] = hash.NewWithData(hd)

	// Initial linkstate for domain
	a.LinkState[d.Name] = false
}

func (a *App) AllLinkCount() (al int64) {
	for _, v := range a.Domains {
		al += v.LinkCount + v.CustomLinkCount
	}
	return
}

/*
 * Private
 */
func (a *App) loadTemplates() {
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

	a.Template = template.Must(template.New("").Funcs(tpl.FuncMap).ParseFiles(templates...))
}

func (a *App) loadDomains() {
	domains, err := model.GetAllDomain(a.DB())
	if err != nil {
		log.Errorf("Error while loading domain details: %q", err)
		return
	}

	for _, v := range domains {
		a.AddDomain(&v)
	}

	log.Infof("%d domains loaded and it's hash generater have been initialized", len(a.Domains))
}

func (a *App) startSyncTasks() {
	dbSync = time.NewTicker(time.Minute * time.Duration(a.Config.DB.LinkSyncInternal))

	// Listening for channel input
	go func(aa *App) {
		for {
			select {
			case <-dbSync.C:
				aa.syncDomainCountToDB()
			}
		}
	}(a)
}

func (a *App) stopSyncTasks() {
	dbSync.Stop()
}

func (a *App) syncDomainCountToDB() {
	log.Infof("Starting Domain link count sync at %v", time.Now())

	for k, v := range a.LinkState {
		if v {
			log.Infof("Syncing domain: %v", k)
			err := model.UpdateDomainLinkCount(a.DB(), a.Domains[k])
			if err != nil {
				log.Errorf("Error occurred while syncing domain count for %v and error is %q", k, err)
			}

			a.LinkState[k] = false
		} else {
			log.Infof("Skipping sync domain: %v", k)
		}
	}

	log.Infof("Completed Domain link count sync at %v", time.Now())
}
