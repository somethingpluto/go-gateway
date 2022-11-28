package initialize

import (
	"github.com/gin-gonic/contrib/sessions"
	"go_gateway/global"
	"log"
)

func InitRedisStore() {
	store, err := sessions.NewRedisStore(10, "tcp", "120.25.255.207:6380", "chx200205173214", []byte("secret"))
	if err != nil {
		log.Fatalln(err)
	}
	store.Options(sessions.Options{
		MaxAge: int(30 * 60),
		Path:   "/",
	})
	global.SessionRedisStore = store
}
