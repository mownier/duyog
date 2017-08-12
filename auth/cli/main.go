package main

import (
	"duyog/auth/rds"
	"duyog/auth/service"
	"duyog/generator"
	"duyog/logger"
	"duyog/writer"
	"fmt"
	"time"

	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

func main() {

	pool := newPool()
	res := newResource()

	keyGen := generator.XIDKey()
	apiGen := generator.APIToken(tokenLen)
	passGen := generator.HashedPass()
	accessGen := generator.AccessToken()
	secretGen := generator.SecretToken(tokenLen)
	refreshGen := generator.RefreshToken(tokenLen)

	userRepo := rds.UserRepo(keyGen, passGen, pool)
	tokenRepo := rds.TokenRepo(keyGen, accessGen, refreshGen, pool)
	clientRepo := rds.ClientRepo(keyGen, apiGen, secretGen, pool)

	userRes := service.UserResource(res, userRepo, tokenRepo, clientRepo, config.TokenExpiry)
	tokenRes := service.TokenResource(res, userRepo, tokenRepo, clientRepo)
	validator := service.AuthValidator(userRepo, clientRepo, logger.RequestLog())

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/"+config.Version+"/user/signin", service.SignInUser(userRes))
	r.HandleFunc("/"+config.Version+"/user/signup", service.SignUpUser(userRes))
	r.HandleFunc("/"+config.Version+"/user/changepass", service.ChangeUserPass(userRes))

	r.HandleFunc("/"+config.Version+"/token/refresh", service.RefreshToken(tokenRes))

	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterService(validator, "AuthValidator")
	r.Handle("/"+config.Version+"/rpc/token/verify", s)

	fmt.Println(toString(config))

	log.Println("starting server at", config.NetAddr)
	log.Fatal(http.ListenAndServe(config.NetAddr, r))
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     config.RedisMaxIdle,
		IdleTimeout: time.Duration(config.RedisIdleTimeout) * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisNetAddr)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func newResource() service.Resource {
	logger := logger.ResponseLog()
	writer := writer.ServiceResponse(logger)

	return service.Resource{
		ResponseWriter: writer,
	}
}
