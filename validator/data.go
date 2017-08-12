package validator

import (
	"bytes"
	"duyog/progerr"
	"net/http"

	"github.com/gorilla/rpc/v2/json2"
)

// DataURL string
type DataURL string

// DataMethod string
type DataMethod string

// DataArgs struct
type DataArgs struct {
	Key string `json:"id"`
}

// DataReply struct
type DataReply struct {
	OK bool `json:"verified"`
}

// AlbumKey string
type AlbumKey string

// ArtistKey string
type ArtistKey string

// PlaylistKey string
type PlaylistKey string

// SongKey string
type SongKey string

// UserKey string
type UserKey string

type data struct {
	url    string
	method string
}

type user struct{ data }
type song struct{ data }
type album struct{ data }
type artist struct{ data }
type playlist struct{ data }

// Album interface
type Album interface {
	Validate(k AlbumKey) DataReply
}

// Artist interface
type Artist interface {
	Validate(k ArtistKey) DataReply
}

// Playlist interface
type Playlist interface {
	Validate(k PlaylistKey) DataReply
}

// Song interface
type Song interface {
	Validate(k SongKey) DataReply
}

// User interface
type User interface {
	Validate(k UserKey) DataReply
}

// Data interfacd
type Data interface {
	Song() Song
	User() User
	Album() Album
	Artist() Artist
	Playlist() Playlist
}

// ValidateAlbum method
func ValidateAlbum(a Album, k AlbumKey) DataReply {
	return a.Validate(k)
}

// ValidateArtist method
func ValidateArtist(a Artist, k ArtistKey) DataReply {
	return a.Validate(k)
}

// ValidatePlaylist method
func ValidatePlaylist(p Playlist, k PlaylistKey) DataReply {
	return p.Validate(k)
}

// ValidateSong method
func ValidateSong(s Song, k SongKey) DataReply {
	return s.Validate(k)
}

// ValidateUser method
func ValidateUser(u User, k UserKey) DataReply {
	return u.Validate(k)
}

func (d data) requestValidation(k string, e error) (DataReply, error) {
	var reply DataReply

	if k == "" {
		return reply, e
	}

	args := DataArgs{
		Key: k,
	}
	msg, err := json2.EncodeClientRequest(d.method, args)

	if err != nil {
		return reply, progerr.Internal(err)
	}

	req, err := http.NewRequest(http.MethodPost, d.url, bytes.NewBuffer(msg))

	if err != nil {
		return reply, progerr.Internal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return reply, progerr.Internal(err)
	}

	var tmp DataReply
	err = json2.DecodeClientResponse(resp.Body, &tmp)

	if err != nil {
		return reply, progerr.Internal(err)
	}

	if tmp.OK == false {
		return reply, e
	}

	reply = tmp

	return reply, nil
}

func (d data) validate(k string, e error) DataReply {
	reply, err := d.requestValidation(k, e)

	if err != nil {
		reply.OK = false
	}

	return reply
}

func (a album) Validate(k AlbumKey) DataReply {
	return a.validate(string(k), progerr.DataAlbumNotVerified)
}

func (a artist) Validate(k ArtistKey) DataReply {
	return a.validate(string(k), progerr.DataArtistNotVerified)
}

func (p playlist) Validate(k PlaylistKey) DataReply {
	return p.validate(string(k), progerr.DataPlaylistNotVerified)
}

func (s song) Validate(k SongKey) DataReply {
	return s.validate(string(k), progerr.DataSongNotVerified)
}

func (u user) Validate(k UserKey) DataReply {
	return u.validate(string(k), progerr.DataUserNotVerified)
}

// RPCAlbum method
func RPCAlbum(u DataURL, m DataMethod) Album {
	return album{
		data: data{
			url:    string(u),
			method: string(m),
		},
	}
}

// RPCArtist method
func RPCArtist(u DataURL, m DataMethod) Artist {
	return artist{
		data: data{
			url:    string(u),
			method: string(m),
		},
	}
}

// RPCPlaylist method
func RPCPlaylist(u DataURL, m DataMethod) Playlist {
	return playlist{
		data: data{
			url:    string(u),
			method: string(m),
		},
	}
}

// RPCSong method
func RPCSong(u DataURL, m DataMethod) Song {
	return song{
		data: data{
			url:    string(u),
			method: string(m),
		},
	}
}

// RPCUser method
func RPCUser(u DataURL, m DataMethod) User {
	return user{
		data: data{
			url:    string(u),
			method: string(m),
		},
	}
}

type rpcData struct {
	song     Song
	user     User
	album    Album
	artist   Artist
	playlist Playlist
}

func (d rpcData) Song() Song         { return d.song }
func (d rpcData) User() User         { return d.user }
func (d rpcData) Album() Album       { return d.album }
func (d rpcData) Artist() Artist     { return d.artist }
func (d rpcData) Playlist() Playlist { return d.playlist }

// RPCData method
func RPCData(url DataURL, song DataMethod, user DataMethod, album DataMethod, artist DataMethod, playlist DataMethod) Data {
	return rpcData{
		song:     RPCSong(url, song),
		user:     RPCUser(url, user),
		album:    RPCAlbum(url, album),
		artist:   RPCArtist(url, artist),
		playlist: RPCPlaylist(url, playlist),
	}
}
