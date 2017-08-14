package main

import (
	"github.com/mownier/duyog/extractor"
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/logger"
	"github.com/mownier/duyog/storage/rds"
	"github.com/mownier/duyog/storage/service"
	"github.com/mownier/duyog/validator"
	"github.com/mownier/duyog/writer"

	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

var config option

func init() {
	configPath := flag.String("config", "./config.json", "Path of the configuration named 'config.json'")
	flag.Parse()

	data, err := ioutil.ReadFile(*configPath)

	if err != nil {
		log.Println(err)
		os.Exit(0)
	}

	err = json.Unmarshal(data, &config)

	if err != nil {
		log.Println("error: can not parse config file to json")
		os.Exit(0)
	}

	log.Println(config.toString())
}

func main() {
	pool := newPool()
	res := newResource()
	keyGen := generator.XIDKey()

	songRepo := rds.SongRepo(keyGen, pool)
	userRepo := rds.UserRepo(keyGen, pool)
	albumRepo := rds.AlbumRepo(keyGen, pool)
	artistRepo := rds.ArtistRepo(keyGen, pool)
	playlistRepo := rds.PlaylistRepo(keyGen, pool)

	songRes := service.SongResource(res, songRepo, config.SongDir, config.AudioSizeLimit, config.ValidAudioExtensions)
	userRes := service.UserResource(res, userRepo, config.UserDir, config.ImageSizeLimit, config.ValidImageExtensions)
	albumRes := service.AlbumResource(res, albumRepo, config.AlbumDir, config.ImageSizeLimit, config.ValidImageExtensions)
	artistRes := service.ArtistResource(res, artistRepo, config.ArtistDrir, config.ImageSizeLimit, config.ValidImageExtensions)
	playlistRes := service.PlaylistResource(res, playlistRepo, config.PlaylistDir, config.ImageSizeLimit, config.ValidImageExtensions)

	r := mux.NewRouter()

	setupSongRoutes(r, songRes)
	setupUserRoutes(r, userRes)
	setupAlbumRoutes(r, albumRes)
	setupArtistRoutes(r, artistRes)
	setupPlaylistRoutes(r, playlistRes)

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
	respWriter := writer.ServiceResponse(logger.ResponseLog())
	authValidator := validator.RPCAuth(config.AuthURL, config.AuthMethod)
	authExtractor := extractor.AccessToken()
	dataValidator := validator.RPCData(config.DataURL, config.SongMethod, config.UserMethod, config.AlbumMethod, config.ArtistMethod, config.PlaylistMethod)

	return service.Resource{
		FormKey: config.UploadFormKey,
		MaxMem:  config.MultipartMaxMem,

		AuthValidator:  authValidator,
		AuthExtractor:  authExtractor,
		DataValidator:  dataValidator,
		ResponseWriter: respWriter,
	}
}
