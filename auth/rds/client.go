package rds

import (
	"github.com/mownier/duyog/auth/store"
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/progerr"
	"time"

	"github.com/garyburd/redigo/redis"
)

type clientRepo struct {
	keyGen    generator.Key
	apiGen    generator.Token
	secretGen generator.Token

	pool *redis.Pool
}

func (r clientRepo) Create(c store.Client) (store.Client, error) {
	var client store.Client

	if c.Email == "" || c.Name == "" {
		return client, progerr.ClientInvalidInfo
	}

	api, err := generator.GenerateToken(r.apiGen)

	if err != nil {
		return client, progerr.Internal(err)
	}

	secret, err := generator.GenerateToken(r.secretGen)

	if err != nil {
		return client, progerr.Internal(err)
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("EXISTS", "email:"+c.Email+":client")

	if err != nil {
		return client, progerr.Internal(err)
	}

	if data.(int64) == 1 {
		return client, progerr.ClientAlreadyExists
	}

	key := generator.GenerateKey(r.keyGen)
	_, err = conn.Do("HMSET", "client:"+key,
		"id", key,
		"role", c.Role,
		"name", c.Name,
		"email", c.Email,
		"api_key", api.Token,
		"secret_token", secret.Token)

	if err != nil {
		return client, progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("SET", "email:"+c.Email+":client", key)
	conn.Do("SET", "secret_token:"+secret.Token+":client", key)
	conn.Do("SET", "api_key:"+api.Token+":client", key)
	conn.Do("SADD", "clients", key)
	conn.Do("HMSET", "client:"+key+":timestamp", "created_on", t, "updated_on", t)

	client = store.Client{
		Key:         key,
		Role:        c.Role,
		Name:        c.Name,
		Email:       c.Email,
		APIKey:      api.Token,
		SecretToken: secret.Token,
	}

	return client, nil
}

func (r clientRepo) GetByKey(k string) (store.Client, error) {
	var client store.Client

	conn := r.pool.Get()
	defer conn.Close()

	err := r.getByKey(k, conn, &client)

	return client, err
}

func (r clientRepo) GetByAPIKey(k string) (store.Client, error) {
	return r.getByToken(k, "api_key")
}

func (r clientRepo) GetBySecretToken(t string) (store.Client, error) {
	return r.getByToken(t, "secret_token")
}

func (r clientRepo) GetByAccessToken(t string) (store.Client, error) {
	return r.getByToken(t, "access_token")
}

func (r clientRepo) getByKey(k string, conn redis.Conn, client *store.Client) error {
	if k == "" {
		return progerr.ClientInvalidKey
	}

	data, err := redis.Values(conn.Do("HGETALL", "client:"+k))

	if err != nil {
		return progerr.Internal(err)
	}

	c := store.Client{}
	err = redis.ScanStruct(data, &c)

	if err != nil {
		return progerr.Internal(err)
	}

	client = &c

	return nil
}

func (r clientRepo) getByToken(t string, name string) (store.Client, error) {
	var client store.Client

	if t == "" {
		return client, progerr.ClientInvalidAPIKey
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("GET", name+":"+t+":client")

	if err != nil {
		return client, progerr.Internal(err)
	}

	if data == nil || len(data.([]byte)) == 0 {
		return client, progerr.ClientNotFound
	}

	key := string(data.([]byte)[:])
	err = r.getByKey(key, conn, &client)

	return client, err
}

// ClientRepo method
func ClientRepo(k generator.Key, a generator.Token, s generator.Token, p *redis.Pool) store.ClientRepo {
	return clientRepo{
		keyGen:    k,
		apiGen:    a,
		secretGen: s,

		pool: p,
	}
}
