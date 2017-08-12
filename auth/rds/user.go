package rds

import (
	"duyog/auth/store"
	"duyog/generator"
	"duyog/progerr"
	"time"

	"github.com/garyburd/redigo/redis"
)

type userRepo struct {
	keyGen  generator.Key
	passGen generator.Pass

	pool *redis.Pool
}

func (r userRepo) Create(u store.User) (store.User, error) {
	var user store.User

	if u.Email == "" {
		return user, progerr.UserInvalidEmail
	}

	if u.Password == "" {
		return user, progerr.UserEmptyPass
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("EXISTS", "email:"+u.Email+":user")

	if err != nil {
		return user, progerr.Internal(err)
	}

	if data.(int64) == 1 {
		return user, progerr.UserAlreadyExists
	}

	key := generator.GenerateKey(r.keyGen)
	pass := generator.GeneratePass(r.passGen, u.Password)
	_, err = conn.Do("HMSET", "user:"+key,
		"id", key,
		"email", u.Email,
		"password", pass)

	if err != nil {
		return user, progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("SET", "email:"+u.Email+":user", key)
	conn.Do("SADD", "users", key)
	conn.Do("HMSET", "user:"+key+":timestamp", "created_on", t, "updated_on", t)

	user = store.User{
		Key:   key,
		Email: u.Email,
	}

	return user, nil
}

func (r userRepo) ChangePass(i store.ChangePassInput) error {
	if i.Email == "" {
		return progerr.UserInvalidEmail
	}

	if i.CurPass == i.NewPass || i.NewPass == "" {
		return progerr.UserInvalidPassword
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("EXISTS", "email:"+i.Email+":user")

	if err != nil {
		return progerr.Internal(err)
	}

	if data.(int64) == 0 {
		return progerr.UserNotFound
	}

	key := string(data.([]byte)[:])
	data, err = conn.Do("HGET", "user:"+key, "password")

	if err != nil {
		return progerr.Internal(err)
	}

	curPass := string(data.([]byte)[:])

	if generator.GeneratePass(r.passGen, i.CurPass) != curPass {
		return progerr.UserMismatchedCurrentPass
	}

	newPass := generator.GeneratePass(r.passGen, i.NewPass)
	_, err = conn.Do("HSET", "user:"+key, "password", newPass)

	if err != nil {
		return progerr.Internal(err)
	}

	conn.Do("HSET", "user:"+key+":timestamp", "updated_on", time.Now().Unix())

	return nil
}

func (r userRepo) ValidateCredential(u store.User) (store.User, error) {
	var user store.User

	if u.Email == "" {
		return user, progerr.UserInvalidEmail
	}

	if u.Password == "" {
		return user, progerr.UserEmptyPass
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("GET", "email:"+u.Email+":user")

	if err != nil {
		return user, progerr.Internal(err)
	}

	if data == nil || len(data.([]byte)) == 0 {
		return user, progerr.UserNotFound
	}

	key := string(data.([]byte)[:])
	data, err = conn.Do("HGET", "user:"+key, "password")

	if err != nil {
		return user, progerr.Internal(err)
	}

	if generator.GeneratePass(r.passGen, u.Password) != string(data.([]byte)[:]) {
		return user, progerr.UserMismatchedPass
	}

	user = store.User{
		Key:   key,
		Email: u.Email,
	}

	return user, nil
}

func (r userRepo) GetByKey(k string) (store.User, error) {
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

	if data == nil || len(data) == 0 {
		return user, progerr.UserNotFound
	}

	u := store.User{}
	err = redis.ScanStruct(data, &u)

	if err != nil {
		return user, progerr.Internal(err)
	}

	user = store.User{
		Key:   u.Key,
		Email: u.Email,
	}

	return user, nil
}

// UserRepo method
func UserRepo(k generator.Key, g generator.Pass, p *redis.Pool) store.UserRepo {
	return userRepo{
		keyGen:  k,
		passGen: g,

		pool: p,
	}
}
