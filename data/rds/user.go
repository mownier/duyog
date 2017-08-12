package rds

import (
	"duyog/data/store"
	"duyog/generator"
	"duyog/progerr"
	"time"

	"github.com/garyburd/redigo/redis"
)

type userRepo struct {
	pool   *redis.Pool
	keyGen generator.Key
}

func (r userRepo) Register(u store.User) (store.User, error) {
	var user store.User

	if u.Key == "" {
		return user, progerr.UserInvalidKey
	}

	if u.Email == "" {
		return user, progerr.UserInvalidEmail
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("EXISTS", "user:"+u.Key)

	if err != nil {
		return user, progerr.Internal(err)
	}

	if data.(int64) == 1 {
		return user, progerr.UserAlreadyExists
	}

	_, err = conn.Do("HMSET", "user:"+u.Key,
		"id", u.Key,
		"email", u.Email,
		"avatar", u.Avatar,
		"last_name", u.LastName,
		"first_name", u.FirstName)

	if err != nil {
		return user, progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("SADD", "users", u.Key)
	conn.Do("HMSET", "user:"+u.Key+":timestamp", "created_on", t, "updated_on", t)

	user = u

	return user, nil
}

func (r userRepo) Update(u store.User) (store.User, error) {
	var user store.User

	if u.Key == "" {
		return user, progerr.UserInvalidKey
	}

	if u.Email == "" {
		return user, progerr.UserInvalidEmail
	}

	if u.FirstName == "" && u.LastName == "" && u.Avatar == "" {
		return user, progerr.UserNothingToUpdate
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("HGETALL", "user:"+u.Key))

	if err != nil {
		return user, progerr.Internal(err)
	}

	if len(data) == 0 {
		return user, progerr.UserNotFound
	}

	var us store.User
	err = redis.ScanStruct(data, &us)

	if err != nil {
		return user, progerr.Internal(err)
	}

	var tmp store.User

	if u.Avatar == "" && u.Avatar != us.Avatar {
		_, err = conn.Do("HSET", "user:"+u.Key, "avatar", u.Avatar)

		if err == nil {
			tmp.Avatar = u.Avatar
		}
	}

	if u.LastName == "" && u.LastName != us.LastName {
		_, err = conn.Do("HSET", "user:"+u.Key, "last_name", u.LastName)

		if err == nil {
			tmp.LastName = u.LastName
		}
	}

	if u.FirstName == "" && u.FirstName != us.FirstName {
		_, err = conn.Do("HSET", "user:"+u.Key, "first_name", u.FirstName)

		if err == nil {
			tmp.FirstName = u.FirstName
		}
	}

	if tmp.FirstName == "" && tmp.LastName == "" && tmp.Avatar == "" {
		return user, progerr.UserNothingToUpdate
	}

	conn.Do("HSET", "user:"+u.Key+":timestamp", "updated_on", time.Now().Unix())

	user = tmp

	return user, nil
}

func (r userRepo) AddPlaylists(k store.UserKey, p []store.PlaylistKey) ([]store.PlaylistKey, error) {
	var keys []store.PlaylistKey

	if k == "" {
		return keys, progerr.UserInvalidKey
	}

	if len(p) == 0 {
		return keys, progerr.PlaylistNotAdded
	}

	conn := r.pool.Get()
	defer conn.Close()

	var tmp []store.PlaylistKey

	for _, v := range p {
		data, err := conn.Do("EXISTS", "playlist:"+v)

		if err != nil || data.(int64) == 0 {
			continue
		}

		data, err = conn.Do("SISMEMBER", "user:"+k+":playlists", v)

		if err != nil || data.(int64) == 1 {
			continue
		}

		_, err = conn.Do("SADD", "user:"+k+":playlists", v)

		if err != nil {
			continue
		}

		conn.Do("SADD", "playlist:"+v+":users", k)

		tmp = append(tmp, v)
	}

	if len(tmp) == 0 {
		return keys, progerr.PlaylistNotAdded
	}

	keys = tmp

	return keys, nil
}

func (r userRepo) GetByKey(k store.UserKey) (store.User, error) {
	var user store.User

	if k == "" {
		return user, progerr.UserInvalidKey
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("HGETALL", "user:"+k))

	if err != nil {
		return user, progerr.Internal(err)
	}

	if len(data) == 0 {
		return user, progerr.UserNotRegistered
	}

	var u store.User
	err = redis.ScanStruct(data, &u)

	if err != nil {
		return user, progerr.Internal(err)
	}

	user = u

	return user, nil
}

func (r userRepo) HasPlaylist(uk store.UserKey, pk store.PlaylistKey) error {
	if uk == "" {
		return progerr.UserInvalidKey
	}

	if pk == "" {
		return progerr.PlaylistInvalidKey
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("SISMEMBER", "user:"+uk+":playlists", pk)

	if err != nil {
		return progerr.Internal(err)
	}

	if data.(int64) == 0 {
		return progerr.PlaylistNotFound
	}

	data, err = conn.Do("EXISTS", "playlist:"+pk+":user")

	if err != nil {
		return progerr.Internal(err)
	}

	if data.(int64) == 0 || store.UserKey(data.([]byte)[:]) != uk {
		return progerr.PlaylistNotFound
	}

	return nil
}

// UserRepo method
func UserRepo(p *redis.Pool, g generator.Key) store.UserRepo {
	return userRepo{
		pool:   p,
		keyGen: g,
	}
}
