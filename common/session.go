package common

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/boj/redistore"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/sessions"
)

// SessionFactory let you create session with config.
type SessionFactory interface {
	Get(r *http.Request) (sess Session)
}

type sf struct {
	store sessions.Store
	name  string
}

// BuildSession factory
func BuildSession(cfg Config) (fac SessionFactory, err error) {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg["RedisAddr"])
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	store, err := redistore.NewRediStoreWithPool(pool, []byte(cfg["SessSecret"]))
	if err != nil {
		return
	}

	fac = &sf{store, cfg["SessName"]}
	return
}

func (f *sf) Get(r *http.Request) (session Session) {
	s, err := f.store.Get(r, f.name)
	return &sess{err, s}
}

// Session is wrapper for github.com/gorilla/sessions
//
// It does not return error in method calls. Instead, you have to call Err() explicitly to
// see if anything goes wrong in previous method calls.
//
//     sess.Set("key", "val")
//     sess.Save(r, w)
//     if err := sess.Err(); err != nil { handling(err) }
type Session interface {
	Get(key string) (value string)
	Set(key, value string)
	Unset(key string)
	AddFlash(value string)
	Flashes() []string
	Save(r *http.Request, w http.ResponseWriter)
	Err() (err error)
}

type sess struct {
	err error
	s   *sessions.Session
}

func (s *sess) Get(key string) (value string) {
	if s.err != nil {
		return
	}
	switch v := s.s.Values[key].(type) {
	case string:
		value = v
	case fmt.Stringer:
		value = v.String()
	default:
		s.err = fmt.Errorf("The value of %s in session store is %T, not a string", key, v)
	}
	return
}

func (s *sess) Set(key, value string) {
	if s.err == nil {
		s.s.Values[key] = value
	}
}

func (s *sess) Unset(key string) {
	if s.err == nil {
		delete(s.s.Values, key)
	}
}

func (s *sess) AddFlash(value string) {
	if s.err == nil {
		s.s.AddFlash(value)
	}
}

func (s *sess) Flashes() []string {
	if s.err != nil {
		return make([]string, 0)
	}
	f := s.s.Flashes()
	values := make([]string, len(f))
	for k, v := range f {
		val := ""
		switch t := v.(type) {
		case string:
			val = t
		case fmt.Stringer:
			val = t.String()
		default:
			s.err = errors.New("The flash value in session store is not a string")
		}
		values[k] = val
	}
	return values
}

func (s *sess) Save(r *http.Request, w http.ResponseWriter) {
	if s.err != nil {
		return
	}
	if err := s.s.Save(r, w); err != nil {
		s.err = err
	}
}

func (s *sess) Err() (err error) {
	return s.err
}
