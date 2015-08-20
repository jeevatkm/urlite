package context

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jeevatkm/urlite/context"
)

type Context struct {
	context.Context
}

func Init(configFile *string) (c *Context) {
	c = &Context{}

	c.Context.Init(configFile)

	return
}

func (c *Context) Close() {
	log.Info("Cleaning up...")

	c.Context.Close()
}
