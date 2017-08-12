package rds

import (
	"duyog/generator"
	"duyog/progerr"
	"duyog/storage/store"
	"time"

	"github.com/garyburd/redigo/redis"
)

type userRepo struct {
	keyGen generator.Key

	pool *redis.Pool
}

func (r userRepo) AddAvatar(k string) (string, error) {
	if k == "" {
		return "", progerr.DataUserNotVerified
	}

	conn := r.pool.Get()
	defer conn.Close()

	key := generator.GenerateKey(r.keyGen)
	_, err := conn.Do("HSET", "user:"+k, "avatar", key)

	if err != nil {
		return "", progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("HSET", "user:"+k+":timestamp", "updated_on", t)
	conn.Do("SADD", "user:"+k+":files", key)
	conn.Do("SADD", "user:files", key)
	conn.Do("SADD", "files", key)

	return key, nil
}

func (r userRepo) Verify(k, fk string) error {
	if fk == "" {
		return progerr.FileNotFound
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("SISMEMBER", "user:"+k+":files", fk)

	if err != nil {
		return progerr.Internal(err)
	}

	if data.(int64) == 0 {
		return progerr.FileNotFound
	}

	return nil
}

// UserRepo method
func UserRepo(k generator.Key, p *redis.Pool) store.UserRepo {
	return userRepo{
		keyGen: k,

		pool: p,
	}
}
