package main

import (
	"github.com/mownier/duyog/validator"
	"fmt"
)

const optionsFormatString = `

	version    : %v

	port       : %v

	redis port         : %v
	redis max idle     : %v
	redis idle timeout : %v
	
	auth url    : %v
	auth method : %v

	data url        : %v
	user method     : %v
	song method     : %v
	album method    : %v
	artist method   : %v
	playlist method : %v

	upload form key      : %v
	image size limit     : %v
	audio size limit     : %v
	multipart max memory : %v

	user directory     : %v
	song directory     : %v
	album directory    : %v
	artist directory   : %v
	playlist directory : %v

	valid image extensions : %v
	valid audio extensions : %v
`

type option struct {
	Version string `json:"version"`

	NetAddr string `json:"network_address"`

	RedisNetAddr     string `json:"redis_network_address"`
	RedisMaxIdle     int    `json:"redis_max_idle"`
	RedisIdleTimeout int    `json:"redis_idle_timeout"`

	AuthURL    validator.AuthURL    `json:"auth_verifier_url"`
	AuthMethod validator.AuthMethod `json:"auth_verifier_method"`

	DataURL        validator.DataURL    `json:"data_verifier_url"`
	UserMethod     validator.DataMethod `json:"user_verifier_method"`
	SongMethod     validator.DataMethod `json:"song_verifier_method"`
	AlbumMethod    validator.DataMethod `json:"album_verifier_method"`
	ArtistMethod   validator.DataMethod `json:"artist_verifier_method"`
	PlaylistMethod validator.DataMethod `json:"playlist_verifier_method"`

	UploadFormKey   string `json:"file_form_key"`
	ImageSizeLimit  int64  `json:"image_size_limit"`
	AudioSizeLimit  int64  `json:"audio_size_limit"`
	MultipartMaxMem int64  `json:"multipart_max_mem"`

	UserDir     string `json:"user_dir"`
	SongDir     string `json:"song_dir"`
	AlbumDir    string `json:"album_dir"`
	ArtistDrir  string `json:"artist_dir"`
	PlaylistDir string `json:"playlist_dir"`

	ValidImageExtensions []string `json:"valid_image_extensions"`
	ValidAudioExtensions []string `json:"valid_audio_extensions"`
}

func (o option) toString() string {
	const optionsFormatString = `
    version : %v

    network address : %v

    redis max idle        : %v connection(s)
    redis idle timeout    : %v second(s)
    redis network address : %v

    auth url    : %v
    auth method : %v

    data url        : %v
    user method     : %v
    song method     : %v
    album method    : %v
    artist method   : %v
    playlist method : %v

    upload form key      : %v
    image size limit     : %v MB
    audio size limit     : %v MB
    multipart max memory : %v MB

    user directory     : %v
    song directory     : %v
    album directory    : %v
    artist directory   : %v
    playlist directory : %v

    valid image extensions : %v
    valid audio extensions : %v
	`
	return fmt.Sprintf(optionsFormatString,
		o.Version,

		o.NetAddr,

		o.RedisMaxIdle,
		o.RedisIdleTimeout,
		o.RedisNetAddr,

		o.AuthURL,
		o.AuthMethod,

		o.DataURL,
		o.UserMethod,
		o.SongMethod,
		o.AlbumMethod,
		o.ArtistMethod,
		o.PlaylistMethod,

		o.UploadFormKey,
		o.ImageSizeLimit/1048576,
		o.AudioSizeLimit/1048576,
		o.MultipartMaxMem/1048576,

		o.UserDir,
		o.SongDir,
		o.AlbumDir,
		o.ArtistDrir,
		o.PlaylistDir,

		o.ValidImageExtensions,
		o.ValidAudioExtensions)
}
