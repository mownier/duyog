package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
)

var config option

type option struct {
	Version string `json:"version"`

	NetAddr string `json:"network_address"`

	RedisNetAddr     string `json:"redis_network_address"`
	RedisMaxIdle     int    `json:"redis_max_idle"`
	RedisIdleTimeout int    `json:"redis_idle_timeout"`

	AuthURL    string `json:"auth_verifier_url"`
	AuthMethod string `json:"auth_verifier_method"`
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

	if _, err := net.ResolveTCPAddr("tcp", config.RedisNetAddr); err != nil {
		errMessage += fmt.Sprintf("%v[redis_network_address]: %v\n", spacing, err)
	}

	if o.RedisMaxIdle <= 0 {
		errMessage += fmt.Sprintf("%v[redis_max_idle]: %v\n", spacing, `max idle connections must be > 0`)
	}

	if o.RedisIdleTimeout <= 0 {
		errMessage += fmt.Sprintf("%v[redis_idle_timeout]: %v\n", spacing, `idle timeout (seconds) must be > 0`)
	}

	if _, err := url.ParseRequestURI(o.AuthURL); err != nil {
		msg := fmt.Sprintf(`%v`, err)
		errMessage += fmt.Sprintf("%v[auth_verifier_url]: %v\n", spacing, msg)
	}

	if o.AuthMethod == "" {
		errMessage += fmt.Sprintf("%v[auth_verifier_method]: %v\n", spacing, `provide a RPC method for verifying the access token`)
	}

	if errMessage == "" {
		return nil
	}

	const sampleConfig = `
    sample config: 
    {
        "version": "v1",

        "network_address": ":9002",

        "redis_network_address": ":6380",
        "redis_max_idle": 3,
        "redis_idle_timeout": 240,

        "auth_verifier_url": "http://127.0.0.1:9001/v1/rpc/token/verify",
        "auth_verifier_method": "Validator.ValidateAccessToken"
    }
	`
	errMessage += sampleConfig
	errMessage += "\n"
	errMessage = "\n" + errMessage

	return errors.New(errMessage)
}

func (o option) toString() string {
	const optionsFormatString = `
               version : %v
       network address : %v
 redis network address : %v
        redis max idle : %v connection(s)
    redis idle timeout : %v second(s)
              auth url : %v
           auth mehtod : %v
	`
	return fmt.Sprintf(optionsFormatString,
		o.Version,
		o.NetAddr,
		o.RedisNetAddr,
		o.RedisMaxIdle,
		o.RedisIdleTimeout,
		o.AuthURL,
		o.AuthMethod)
}
