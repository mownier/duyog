package rds

import (
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/storage/store"
	"time"

	"github.com/garyburd/redigo/redis"
)

type songRepo struct {
	keyGen generator.Key

	pool *redis.Pool
}

func (r songRepo) AddAudio(k string) (string, error) {
	if k == "" {
		return "", progerr.DataSongNotVerified
	}

	conn := r.pool.Get()
	defer conn.Close()

	key := generator.GenerateKey(r.keyGen)
	_, err := conn.Do("HSET", "song:"+k, "audio", key)

	if err != nil {
		return "", progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("HSET", "song:"+k+":timestamp", "updated_on", t)
	conn.Do("SADD", "song:"+k+":files", key)
	conn.Do("SADD", "song:files", key)
	conn.Do("SADD", "files", key)

	return key, nil
}

func (r songRepo) Verify(k, fk string) error {
	if fk == "" {
		return progerr.FileNotFound
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("SISMEMBER", "song:"+k+":files", fk)

	if err != nil {
		return progerr.Internal(err)
	}

	if data.(int64) == 0 {
		return progerr.FileNotFound
	}

	return nil
}

// SongRepo method
func SongRepo(k generator.Key, p *redis.Pool) store.SongRepo {
	return songRepo{
		keyGen: k,

		pool: p,
	}
}
