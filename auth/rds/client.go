// Package rds provides the authentication server
// a redis database access on clients, tokens and users.
package rds

import (
	"time"

	"github.com/mownier/duyog/auth/store"
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/progerr"

	"github.com/garyburd/redigo/redis"
)

// clientRepo implements the ClientRepo interface in
// package store of the authentication server.
type clientRepo struct {
	// keyGen generates a unique key upon creating a client
	keyGen generator.Key

	// apiGen generates API key upon creating a client
	apiGen generator.Token

	// secretGen generates secret token upon creating a client
	secretGen generator.Token

	// pool maintains a pool of connections for the redis database
	pool *redis.Pool
}

// Create creates and stores a new client in the redis database.
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

// GetByKey retrieves a client information from the
// redis database using a key of a client.
func (r clientRepo) GetByKey(k string) (store.Client, error) {
	var client store.Client

	conn := r.pool.Get()
	defer conn.Close()

	tmp, err := r.getByKey(k, conn)

	if err != nil {
		return client, err
	}

	client = tmp

	return client, nil
}

// GetByAPIKey retrieves a client information from the
// redis database using an API key assigned to the client.
func (r clientRepo) GetByAPIKey(k string) (store.Client, error) {
	return r.getByToken(k, "api_key")
}

// GetBySecretToken retrieves a client information from the
// redis database using a secret token assigned to the client.
func (r clientRepo) GetBySecretToken(t string) (store.Client, error) {
	return r.getByToken(t, "secret_token")
}

// GetByAccessToken retrieves a client information from the
// redis database using an access token issued by the client.
func (r clientRepo) GetByAccessToken(t string) (store.Client, error) {
	var client store.Client

	if t == "" {
		return client, progerr.TokenInvalidAccess
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("GET", "access:"+t+":token")

	if err != nil {
		return client, progerr.TokenInvalidAccess
	}

	if len(data.([]byte)) == 0 {
		return client, progerr.TokenInvalidAccess
	}

	tokenKey := string(data.([]byte)[:])

	data, err = conn.Do("GET", "token:"+tokenKey+":client")

	if err != nil {
		return client, progerr.TokenInvalidAccess
	}

	if len(data.([]byte)) == 0 {
		return client, progerr.TokenInvalidAccess
	}

	clientKey := string(data.([]byte)[:])
	tmp, err := r.getByKey(clientKey, conn)

	if err != nil {
		return client, progerr.TokenInvalidAccess
	}

	client = tmp

	return client, nil
}

// getByKey is a reusable function that retrieves a client
// information from the redis database using a client key.
func (r clientRepo) getByKey(k string, conn redis.Conn) (store.Client, error) {
	var client store.Client

	if k == "" {
		return client, progerr.ClientInvalidKey
	}

	data, err := redis.Values(conn.Do("HGETALL", "client:"+k))

	if err != nil {
		return client, progerr.Internal(err)
	}

	c := store.Client{}
	err = redis.ScanStruct(data, &c)

	if err != nil {
		return client, progerr.Internal(err)
	}

	client = c

	return client, nil
}

// getByToken is a reusable functon that retrieves a client
// information from the redis database using the assigned
// tokens which are the API key and secret token.
// The parameter name is either "api_key" or "secret_token".
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
	tmp, err := r.getByKey(key, conn)

	if err != nil {
		return client, err
	}

	client = tmp

	return client, nil
}

// ClientRepo returns an implementation of the ClientRepo
// interface in the package store of the authentication server.
func ClientRepo(k generator.Key, a generator.Token, s generator.Token, p *redis.Pool) store.ClientRepo {
	return clientRepo{
		keyGen:    k,
		apiGen:    a,
		secretGen: s,

		pool: p,
	}
}
