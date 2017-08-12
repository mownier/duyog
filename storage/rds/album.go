package rds

import (
	"duyog/generator"
	"duyog/progerr"
	"duyog/storage/store"
	"time"

	"github.com/garyburd/redigo/redis"
)

type albumRepo struct {
	keyGen generator.Key

	pool *redis.Pool
}

func (r albumRepo) AddPhoto(k string) (string, error) {
	if k == "" {
		return "", progerr.DataAlbumNotVerified
	}

	conn := r.pool.Get()
	defer conn.Close()

	key := generator.GenerateKey(r.keyGen)
	_, err := conn.Do("HSET", "album:"+k, "photo", key)

	if err != nil {
		return "", progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("HSET", "album:"+k+":timestamp", "updated_on", t)
	conn.Do("SADD", "album:"+k+":files", key)
	conn.Do("SADD", "album:files", k)
	conn.Do("SADD", "files", key)

	return key, nil
}

func (r albumRepo) Verify(k, fk string) error {
	if fk == "" {
		return progerr.FileNotFound
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("SISMEMBER", "album:"+k+":files", fk)

	if err != nil {
		return progerr.Internal(err)
	}

	if data.(int64) == 0 {
		return progerr.FileNotFound
	}

	return nil
}

// AlbumRepo method
func AlbumRepo(k generator.Key, p *redis.Pool) store.AlbumRepo {
	return albumRepo{
		keyGen: k,

		pool: p,
	}
}
