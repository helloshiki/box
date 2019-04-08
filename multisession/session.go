package multisession

import (
	"fmt"
	"github.com/boj/redistore"
	"github.com/gin-gonic/gin"
	"github.com/helloshiki/box/ginhelp"
	"github.com/pkg/errors"
	"net/http"
)

type sessionBuilder struct {
	size     int
	network  string
	address  string
	password string
	db       string
	keyPairs [][]byte

	maxAge int
}

func NewSessionBuilder() *sessionBuilder {
	return &sessionBuilder{
		size:     10,
		network:  "tcp",
		address:  "127.0.0.1:6379",
		password: "",
		db:       "0",
		keyPairs: [][]byte{[]byte("key")},
	}
}

func (opts *sessionBuilder) SetSize(size int) *sessionBuilder {
	opts.size = size
	return opts
}

func (opts *sessionBuilder) SetAddress(addr string) *sessionBuilder {
	opts.address = addr
	return opts
}

func (opts *sessionBuilder) SetPassword(password string) *sessionBuilder {
	opts.password = password
	return opts
}

func (opts *sessionBuilder) SetDB(db string) *sessionBuilder {
	opts.db = db
	return opts
}

func (opts *sessionBuilder) SetMaxAge(maxAge int) *sessionBuilder {
	opts.maxAge = maxAge
	return opts
}

func (opts *sessionBuilder) SetKeyPairs(keyPairs [][]byte) *sessionBuilder {
	opts.keyPairs = keyPairs
	return opts
}

func (opts *sessionBuilder) Build() (*sessionHandler, error) {
	store, err := redistore.NewRediStoreWithDB(opts.size, opts.network, opts.address, opts.password, opts.db, opts.keyPairs...)
	if err != nil {
		return nil, errors.WithMessagef(err, "NewRediStoreWithDB %+v", *opts)
	}

	store.DefaultMaxAge = opts.maxAge
	store.SetSerializer(redistore.JSONSerializer{})
	return &sessionHandler{
		RediStore:   store,
		sessionName: "login",
	}, nil
}

type sessionHandler struct {
	*redistore.RediStore
	sessionName string
}

func (store *sessionHandler) PreLoginHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, store.sessionName)
		if err != nil {
			ginhelp.GinAbort(c, http.StatusInternalServerError, http.StatusInternalServerError, err.Error())
			return
		}

		if !session.IsNew {
			ginhelp.GinAbort(c, http.StatusOK, http.StatusOK, "Online")
			return
		}

		c.Next()
	}
}

func (store *sessionHandler) PostLoginHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		values, ok := c.Get("session")
		if ok {
			session, _ := store.Get(c.Request, store.sessionName)
			session.Values = values.(map[interface{}]interface{})
			session.Options.MaxAge = store.DefaultMaxAge
			_ = session.Save(c.Request, c.Writer)
		}
	}
}

func (store *sessionHandler) CheckLoginHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, store.sessionName)
		if err != nil {
			ginhelp.GinAbort(c, http.StatusInternalServerError, http.StatusInternalServerError, err.Error())
			return
		}

		if session.IsNew {
			ginhelp.GinAbort(c, http.StatusUnauthorized, http.StatusInternalServerError, "Offline")
			return
		}

		c.Set("session", session.Values)
		c.Next()
	}
}

func (store *sessionHandler) LogoutHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println(22, c.GetBool("deleteSession"))
		if c.GetBool("deleteSession") {
			session, _ := store.Get(c.Request, store.sessionName)

			session.Options.MaxAge = -1
			_ = session.Save(c.Request, c.Writer)
		}
	}
}
