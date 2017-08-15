package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/mownier/duyog/validator"
)

var config option

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

	if _, err := url.ParseRequestURI(string(o.AuthURL)); err != nil {
		msg := fmt.Sprintf(`%v`, err)
		errMessage += fmt.Sprintf("%v[auth_verifier_url]: %v\n", spacing, msg)
	}

	if o.AuthMethod == "" {
		errMessage += fmt.Sprintf("%v[auth_verifier_method]: %v\n", spacing, `provide a RPC method for verifying the access token`)
	}

	if _, err := url.ParseRequestURI(string(o.DataURL)); err != nil {
		msg := fmt.Sprintf(`%v`, err)
		errMessage += fmt.Sprintf("%v[data_verifier_url]: %v\n", spacing, msg)
	}

	if o.UserMethod == "" {
		errMessage += fmt.Sprintf("%v[user_verifier_method]: %v\n", spacing, `provide a RPC method for verifying the user`)
	}

	if o.SongMethod == "" {
		errMessage += fmt.Sprintf("%v[song_verifier_method]: %v\n", spacing, `provide a RPC method for verifying the song`)
	}

	if o.AlbumMethod == "" {
		errMessage += fmt.Sprintf("%v[album_verifier_method]: %v\n", spacing, `provide a RPC method for verifying the album`)
	}

	if o.ArtistMethod == "" {
		errMessage += fmt.Sprintf("%v[artist_verifier_method]: %v\n", spacing, `provide a RPC method for verifying the artist`)
	}

	if o.PlaylistMethod == "" {
		errMessage += fmt.Sprintf("%v[playlist_verifier_method]: %v\n", spacing, `provide a RPC method for verifying the playlist`)
	}

	if o.UploadFormKey == "" {
		errMessage += fmt.Sprintf("%v[file_form_key]: %v\n", spacing, `provide an upload form key`)
	}

	if o.ImageSizeLimit <= 0 {
		errMessage += fmt.Sprintf("%v[image_size_limit]: %v\n", spacing, `must be > 0`)
	}

	if o.AudioSizeLimit <= 0 {
		errMessage += fmt.Sprintf("%v[audio_size_limit]: %v\n", spacing, `must be > 0`)
	}

	if o.MultipartMaxMem <= 0 {
		errMessage += fmt.Sprintf("%v[multipart_max_mem]: %v\n", spacing, `must be > 0`)
	}

	if o.UserDir == "" {
		errMessage += fmt.Sprintf("%v[user_dir]: %v\n", spacing, `provide a directory to store user-related files`)
	}

	if o.SongDir == "" {
		errMessage += fmt.Sprintf("%v[song_dir]: %v\n", spacing, `provide a directory to store song-related files`)
	}

	if o.AlbumDir == "" {
		errMessage += fmt.Sprintf("%v[album_dir]: %v\n", spacing, `provide a directory to store album-related files`)
	}

	if o.ArtistDrir == "" {
		errMessage += fmt.Sprintf("%v[artist_dir]: %v\n", spacing, `provide a directory to store artist-related files`)
	}

	if o.PlaylistDir == "" {
		errMessage += fmt.Sprintf("%v[playlist_dir]: %v\n", spacing, `provide a directory to store playlist-related files`)
	}

	if len(o.ValidImageExtensions) == 0 {
		errMessage += fmt.Sprintf("%v[valid_image_extensions]: %v\n", spacing, `provide valid image extensions`)
	}

	if len(o.ValidAudioExtensions) == 0 {
		errMessage += fmt.Sprintf("%v[valid_audio_extensions]: %v\n", spacing, `provide valid audio extensions`)
	}

	dirs := map[string]bool{}
	dirs[strings.ToLower(o.SongDir)] = true
	dirs[strings.ToLower(o.UserDir)] = true
	dirs[strings.ToLower(o.AlbumDir)] = true
	dirs[strings.ToLower(o.ArtistDrir)] = true
	dirs[strings.ToLower(o.PlaylistDir)] = true

	if len(dirs) != 5 {
		errMessage += fmt.Sprintf("%v[directory]: %v\n", spacing, `provide unique directory names (not case sensitive)`)
	}

	return nil
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
