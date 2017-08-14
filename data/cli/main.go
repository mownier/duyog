package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mownier/duyog/data/rds"
	"github.com/mownier/duyog/data/service"
	"github.com/mownier/duyog/extractor"
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/logger"
	"github.com/mownier/duyog/validator"
	"github.com/mownier/duyog/writer"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

func main() {
	pool := newPool()
	res := newResource()
	keyGen := generator.XIDKey()

	songRepo := rds.SongRepo(pool, keyGen)
	userRepo := rds.UserRepo(pool, keyGen)
	albumRepo := rds.AlbumRepo(pool, keyGen)
	artistRepo := rds.ArtistRepo(pool, keyGen)
	playlistRepo := rds.PlaylistRepo(pool, keyGen)

	meRes := service.MeResource(res, userRepo, playlistRepo)
	songRes := service.SongResource(res, songRepo)
	userRes := service.UserResource(res, userRepo, playlistRepo)
	albumRes := service.AlbumResource(res, albumRepo)
	artistRes := service.ArtistResource(res, artistRepo)
	playlistRes := service.PlaylistResource(res, playlistRepo)

	verifier := service.NewVerifier(userRepo, songRepo, albumRepo, artistRepo, playlistRepo, logger.RequestLog())

	r := mux.NewRouter().StrictSlash(true)

	setupMeRoutes(r, meRes)
	setupSongRoutes(r, songRes)
	setupUserRoutes(r, userRes)
	setupAlbumRoutes(r, albumRes)
	setupArtistRoutes(r, artistRes)
	setupPlaylistRoutes(r, playlistRes)

	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterService(verifier, "Verifier")
	r.Handle("/"+config.Version+"/data/verify", s)

	fmt.Println(config.toString())

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
	authURL := validator.AuthURL(config.AuthURL)
	authMethod := validator.AuthMethod(config.AuthMethod)

	return service.Resource{
		AuthExtractor:    extractor.AccessToken(),
		AuthValidator:    validator.RPCAuth(authURL, authMethod),
		RequestValidator: validator.JSONRequest(),
		ResponseWriter:   writer,
	}
}
