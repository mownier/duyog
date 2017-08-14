package rds

import (
	"github.com/mownier/duyog/auth/store"
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/progerr"

	"github.com/garyburd/redigo/redis"
)

type tokenRepo struct {
	keyGen     generator.Key
	accessGen  generator.Access
	refreshGen generator.Token

	pool *redis.Pool
}

func (r tokenRepo) Create(i generator.AccessInput) (store.Token, error) {
	var token store.Token

	if i.ClientKey == "" {
		return token, progerr.ClientInvalidKey
	}

	if i.UserKey == "" {
		return token, progerr.UserInvalidKey
	}

	if i.SecretToken == "" {
		return token, progerr.ClientInvalidSecretToken
	}

	if i.Expiry <= 0 {
		return token, progerr.TokenInvalidExpiry
	}

	access, err := generator.GenerateAccess(r.accessGen, i)

	if err != nil {
		return token, err
	}

	refresh, err := generator.GenerateToken(r.refreshGen)

	if err != nil {
		return token, err
	}

	conn := r.pool.Get()
	defer conn.Close()

	key := generator.GenerateKey(r.keyGen)
	_, err = conn.Do("HMSET", "token:"+key,
		"id", key,
		"expiry", i.Expiry,
		"access_token", access.Token,
		"refresh_token", refresh.Token)

	if err != nil {
		return token, progerr.Internal(err)
	}

	conn.Do("SET", "token:"+key+":user", i.UserKey)
	conn.Do("SET", "token:"+key+":client", i.ClientKey)
	conn.Do("SET", "refresh:"+refresh.Token+":token", key)
	conn.Do("SET", "access:"+access.Token+":token", key)
	conn.Do("SADD", "access_tokens", access.Token)
	conn.Do("SADD", "refresh_tokens", refresh.Token)
	conn.Do("SADD", "user:"+i.UserKey+":tokens", key)
	conn.Do("SADD", "client:"+i.ClientKey+":tokens", key)
	conn.Do("SADD", "tokens", key)
	conn.Do("HMSET", "token:"+key+":timestamp", "created_on", access.CreatedOn, "updated_on", access.CreatedOn)

	token = store.Token{
		Key:     key,
		Access:  access.Token,
		Refresh: refresh.Token,
		Expiry:  i.Expiry,
	}

	return token, nil
}

func (r tokenRepo) Refresh(i store.RefreshTokenInput) (store.Token, error) {
	var token store.Token

	if i.Expiry <= 0 {
		return token, progerr.TokenInvalidExpiry
	}

	if i.RefreshToken == "" {
		return token, progerr.TokenInvalidRefresh
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("GET", "refresh:"+i.RefreshToken+":token")

	if err != nil {
		return token, progerr.Internal(err)
	}

	if data == nil || len(data.([]byte)) == 0 {
		return token, progerr.TokenInvalidRefresh
	}

	data, err = redis.Values(conn.Do("HGETALL", "token:"+string(data.([]byte)[:])))

	if err != nil {
		return token, progerr.Internal(err)
	}

	t := store.Token{}
	err = redis.ScanStruct(data.([]interface{}), &t)

	if err != nil {
		return token, progerr.Internal(err)
	}

	data, err = conn.Do("GET", "token:"+t.Key+":user")

	if err != nil {
		return token, progerr.Internal(err)
	}

	if data == nil || len(data.([]byte)) == 0 {
		return token, progerr.TokenInvalidRefresh
	}

	userKey := string(data.([]byte)[:])
	data, err = conn.Do("GET", "token:"+t.Key+":client")

	if err != nil {
		return token, progerr.Internal(err)
	}

	if data == nil || len(data.([]byte)) == 0 {
		return token, progerr.TokenInvalidRefresh
	}

	clientKey := string(data.([]byte)[:])
	data, err = conn.Do("HGET", "client:"+clientKey, "secret_token")

	if err != nil {
		return token, progerr.Internal(err)
	}

	if data == nil || len(data.([]byte)) == 0 {
		return token, progerr.TokenInvalidRefresh
	}

	input := generator.AccessInput{
		ClientKey:   clientKey,
		UserKey:     userKey,
		Expiry:      i.Expiry,
		SecretToken: string(data.([]byte)[:]),
	}

	access, err := generator.GenerateAccess(r.accessGen, input)

	if err != nil {
		return token, err
	}

	_, err = conn.Do("HMSET", "token:"+t.Key,
		"access_token", access.Token,
		"updated_timestamp", access.CreatedOn,
		"expiry", i.Expiry)

	if err != nil {
		return token, progerr.Internal(err)
	}

	conn.Do("DEL", "access:"+t.Access+":token")
	conn.Do("SREM", "access_tokens", t.Access)
	conn.Do("SET", "access:"+access.Token+":token", t.Key)
	conn.Do("SADD", "access_tokens", access.Token)
	conn.Do("HSET", "token:"+t.Key+":timestamp", "updated_on", access.CreatedOn)

	token = store.Token{
		Access:  access.Token,
		Refresh: i.RefreshToken,
		Expiry:  i.Expiry,
	}

	return token, nil
}

// TokenRepo method
func TokenRepo(k generator.Key, a generator.Access, r generator.Token, p *redis.Pool) store.TokenRepo {
	return tokenRepo{
		keyGen:     k,
		accessGen:  a,
		refreshGen: r,

		pool: p,
	}
}
