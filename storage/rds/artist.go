package rds

import (
	"duyog/generator"
	"duyog/progerr"
	"duyog/storage/store"
	"time"

	"github.com/garyburd/redigo/redis"
)

type artistRepo struct {
	keyGen generator.Key

	pool *redis.Pool
}

func (r artistRepo) AddPhoto(k string) (string, error) {
	if k == "" {
		return "", progerr.DataArtistNotVerified
	}

	conn := r.pool.Get()
	defer conn.Close()

	key := generator.GenerateKey(r.keyGen)
	_, err := conn.Do("HSET", "artist:"+k, "photo", key)

	if err != nil {
		return "", progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("HSET", "artist:"+k+":timestamp", "updated_on", t)
	conn.Do("SADD", "artist:"+k+":files", key)
	conn.Do("SADD", "artist:files", key)
	conn.Do("SADD", "files", key)

	return key, nil
}

func (r artistRepo) Verify(k, fk string) error {
	if fk == "" {
		return progerr.FileNotFound
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("SISMEMBER", "artist:"+k+":files", fk)

	if err != nil {
		return progerr.Internal(err)
	}

	if data.(int64) == 0 {
		return progerr.FileNotFound
	}

	return nil
}

// ArtistRepo method
func ArtistRepo(k generator.Key, p *redis.Pool) store.ArtistRepo {
	return artistRepo{
		keyGen: k,

		pool: p,
	}
}
