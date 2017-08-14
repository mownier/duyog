package rds

import (
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/storage/store"
	"time"

	"github.com/garyburd/redigo/redis"
)

type playlistRepo struct {
	keyGen generator.Key

	pool *redis.Pool
}

func (r playlistRepo) AddPhoto(k string) (string, error) {
	if k == "" {
		return "", progerr.DataPlaylistNotVerified
	}

	conn := r.pool.Get()
	defer conn.Close()

	key := generator.GenerateKey(r.keyGen)
	_, err := conn.Do("HSET", "playlist:"+k, "photo", key)

	if err != nil {
		return "", progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("HSET", "playlist:"+k+":timestamp", "updated_on", t)
	conn.Do("SADD", "playlist:"+k+":files", key)
	conn.Do("SADD", "playlist:files", key)
	conn.Do("SADD", "files", key)

	return key, nil
}

func (r playlistRepo) Verify(k, fk string) error {
	if fk == "" {
		return progerr.FileNotFound
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("SISMEMBER", "playlist:"+k+":files", fk)

	if err != nil {
		return progerr.Internal(err)
	}

	if data.(int64) == 0 {
		return progerr.FileNotFound
	}

	return nil
}

// PlaylistRepo method
func PlaylistRepo(k generator.Key, p *redis.Pool) store.PlaylistRepo {
	return playlistRepo{
		keyGen: k,

		pool: p,
	}
}
