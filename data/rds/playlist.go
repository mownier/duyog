package rds

import (
	"duyog/data/store"
	"duyog/generator"
	"duyog/progerr"
	"time"

	"github.com/garyburd/redigo/redis"
)

type playlistRepo struct {
	pool   *redis.Pool
	keyGen generator.Key
}

func (r playlistRepo) Create(k store.UserKey, p store.Playlist) (store.Playlist, error) {
	var playlist store.Playlist

	if k == "" {
		return playlist, progerr.UserInvalidKey
	}

	if p.Name == "" {
		return playlist, progerr.PlaylistInvalidName
	}

	key := generator.GenerateKey(r.keyGen)

	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("HMSET", "playlist:"+key,
		"id", key,
		"photo", p.Photo,
		"description", p.Desc,
		"name", p.Name)

	if err != nil {
		return playlist, progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("HMSET", "playlist:"+key+":timestamp", "created_on", t, "updated_on", t)
	conn.Do("SADD", "playlists", key)
	conn.Do("SADD", "user:"+k+":playlists", key)
	conn.Do("SADD", "playlist:"+key+":users", k)
	conn.Do("SET", "playlist:"+key+":user", k)

	playlist = p
	playlist.Key = store.PlaylistKey(key)

	return playlist, nil
}

func (r playlistRepo) Update(p store.Playlist) (store.Playlist, error) {
	var playlist store.Playlist

	_, err := r.getCreator(p.Key)

	if err != nil {
		return playlist, err
	}

	if p.Desc == "" && p.Name == "" && p.Photo == "" {
		return playlist, progerr.PlaylistNothingToUpdate
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("HGETALL", "playlist:"+p.Key))

	if err != nil {
		return playlist, progerr.Internal(err)
	}

	if len(data) == 0 {
		return playlist, progerr.PlaylistNotFound
	}

	var pl store.Playlist
	err = redis.ScanStruct(data, &pl)

	if err != nil {
		return playlist, progerr.Internal(err)
	}

	tmp := store.Playlist{
		Key: p.Key,
	}

	if p.Name == "" && p.Name != pl.Name {
		_, err = conn.Do("HSET", "playlist:"+p.Key, "name", p.Name)

		if err == nil {
			tmp.Name = p.Name
		}
	}

	if p.Photo == "" && p.Photo != pl.Photo {
		_, err = conn.Do("HSET", "playlist:"+p.Key, "photo", p.Photo)

		if err == nil {
			tmp.Photo = p.Photo
		}
	}

	if p.Desc == "" && p.Desc != pl.Desc {
		_, err = conn.Do("HSET", "playlist:"+p.Key, "description", p.Desc)

		if err == nil {
			tmp.Desc = p.Desc
		}
	}

	if tmp.Name == "" && tmp.Photo == "" && tmp.Desc == "" {
		return playlist, progerr.PlaylistNothingToUpdate
	}

	conn.Do("HSET", "playlist:"+p.Key+":timestamp", "updated_on", time.Now().Unix())

	playlist = tmp

	return playlist, nil
}

func (r playlistRepo) AddSongs(pk store.PlaylistKey, sk []store.SongKey) ([]store.SongKey, error) {
	var songKeys []store.SongKey

	_, err := r.getCreator(pk)

	if err != nil {
		return songKeys, err
	}

	if len(sk) == 0 {
		return songKeys, progerr.PlaylistNoSongsAdded
	}

	conn := r.pool.Get()
	defer conn.Close()

	var tmp []store.SongKey

	for _, key := range sk {
		data, err := conn.Do("EXISTS", "song:"+key)

		if err != nil || data.(int64) == 0 {
			continue
		}

		data, err = conn.Do("SISMEMBER", "playlist:"+pk+":songs", key)

		if err != nil || data.(int64) == 1 {
			continue
		}

		_, err = conn.Do("SADD", "playlist:"+pk+":songs", key)

		if err != nil {
			continue
		}

		conn.Do("SADD", "song:"+key+":playlists", pk)

		tmp = append(tmp, key)
	}

	if len(tmp) == 0 {
		return songKeys, progerr.PlaylistNoSongsAdded
	}

	songKeys = tmp

	return songKeys, nil
}

func (r playlistRepo) GetSongs(k store.PlaylistKey) (store.Songs, error) {
	var songs store.Songs

	_, err := r.getCreator(k)

	if err != nil {
		return songs, err
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("SMEMBERS", "playlist:"+k+":songs"))

	if err != nil {
		return songs, progerr.Internal(err)
	}

	if len(data) == 0 {
		return songs, progerr.PlaylistHasNoSongs
	}

	var keys []store.SongKey

	for _, v := range data {
		if len(v.([]byte)) == 0 {
			continue
		}

		keys = append(keys, store.SongKey(v.([]byte)[:]))
	}

	if len(keys) == 0 {
		return songs, progerr.PlaylistHasNoSongs
	}

	tmp, err := getSongs(keys, conn, progerr.PlaylistHasNoSongs)

	if err != nil {
		return songs, err
	}

	songs = tmp

	return songs, nil
}

func (r playlistRepo) GetByKey(k store.PlaylistKey) (store.Playlist, error) {
	var playlist store.Playlist

	_, err := r.getCreator(k)

	if err != nil {
		return playlist, err
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("HGETALL", "playlist:"+k))

	if err != nil {
		return playlist, progerr.Internal(err)
	}

	if len(data) == 0 {
		return playlist, progerr.PlaylistNotFound
	}

	var pl store.Playlist
	err = redis.ScanStruct(data, &pl)

	if err != nil {
		return playlist, progerr.Internal(err)
	}

	playlist = pl

	return playlist, nil
}

func (r playlistRepo) GetByUser(k store.UserKey) (store.Playlists, error) {
	var playlists store.Playlists

	if k == "" {
		return playlists, progerr.UserInvalidKey
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("SMEMBERS", "user:"+k+":playlists"))

	if err != nil {
		return playlists, progerr.Internal(err)
	}

	if len(data) == 0 {
		return playlists, progerr.UserHasNoPlaylists
	}

	var tmp store.Playlists

	for _, v := range data {
		if len(v.([]byte)) == 0 {
			continue
		}

		key := store.PlaylistKey(v.([]byte)[:])

		if key == "" {
			continue
		}

		data, err := redis.Values(conn.Do("HGETALL", "playlist:"+key))

		if err != nil || len(data) == 0 {
			continue
		}

		var p store.Playlist
		err = redis.ScanStruct(data, &p)

		if err != nil {
			continue
		}

		tmp.Playlists[key] = p
	}

	return playlists, nil
}

func (r playlistRepo) getCreator(k store.PlaylistKey) (store.UserKey, error) {
	var key store.UserKey

	if k == "" {
		return key, progerr.PlaylistInvalidKey
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("GET", "playlist:"+k+":user")

	if err != nil {
		return key, progerr.Internal(err)
	}

	if data == nil || len(data.([]byte)) == 0 {
		return key, progerr.PlaylistHasNoCreator
	}

	key = store.UserKey(data.([]byte)[:])

	return key, nil
}

// PlaylistRepo method
func PlaylistRepo(p *redis.Pool, g generator.Key) store.PlaylistRepo {
	return playlistRepo{
		pool:   p,
		keyGen: g,
	}
}
