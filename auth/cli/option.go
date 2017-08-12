package main

import (
	"duyog/auth/store"
	"errors"
	"fmt"
	"net"
)

type option struct {
	Version string `json:"version"`

	NetAddr string `json:"network_address"`

	TokenExpiry int64 `json:"token_expiry"`

	RedisNetAddr     string `json:"redis_network_address"`
	RedisMaxIdle     int    `json:"redis_max_idle"`
	RedisIdleTimeout int    `json:"redis_idle_timeout"`
}

func (o option) valid() error {
	const spacing = "    "

	var errMessage string

	if o.Version == "" {
		errMessage += fmt.Sprintf("%v[version]: %v\n", spacing, `provide a version string`)
	}

	if _, err := net.ResolveTCPAddr("tcp", config.NetAddr); err != nil {
		errMessage += fmt.Sprintf("%v[network_address]: %v\n", spacing, err)
	}

	if o.TokenExpiry <= 0 {
		errMessage += fmt.Sprintf("%v[token_expiry]: %v\n", spacing, `expiry (seconds) must be > 0`)
	}

	if _, err := net.ResolveTCPAddr("tcp", config.RedisNetAddr); err != nil {
		errMessage += fmt.Sprintf("%v[redis_network_address]: %v\n", spacing, err)
	}

	if o.RedisMaxIdle <= 0 {
		errMessage += fmt.Sprintf("%v[redis_max_idle]: %v\n", spacing, `max idle connections must be > 0`)
	}

	if o.RedisIdleTimeout <= 0 {
		errMessage += fmt.Sprintf("%v[redis_idle_timeout]: %v\n", spacing, `idle timeout (seconds) must be > 0`)
	}

	if errMessage == "" {
		return nil
	}

	const sampleConfig = `
    sample config: 
    {
        "version": "v1",

        "network_address": ":9001",
		
        "token_expiry": 3600,

        "redis_network_address": ":6379",
        "redis_max_idle": 3,
        "redis_idle_timeout": 240
    }
	`
	errMessage += sampleConfig
	errMessage += "\n"
	errMessage = "\n" + errMessage

	return errors.New(errMessage)
}

func toString(v interface{}) string {
	switch v.(type) {
	case store.Client:
		const clientFormatString = `
                 key : %v
                name : %v
                role : %v
               email : %v
             api_key : %v
        secret_token : %v
        `
		return fmt.Sprintf(clientFormatString,
			v.(store.Client).Key,
			v.(store.Client).Name,
			v.(store.Client).Role,
			v.(store.Client).Email,
			v.(store.Client).APIKey,
			v.(store.Client).SecretToken)

	case option:
		const optionsFormatString = `
                   version : %v
           network address : %v
              token expiry : %v second(s)
     redis network address : %v
            redis max idle : %v connection(s)
        redis idle timeout : %v seconds(s)
        `
		return fmt.Sprintf(optionsFormatString,
			v.(option).Version,
			v.(option).NetAddr,
			v.(option).TokenExpiry,
			v.(option).RedisNetAddr,
			v.(option).RedisMaxIdle,
			v.(option).RedisIdleTimeout)

	default:
		return fmt.Sprintf("%v", v)
	}
}
