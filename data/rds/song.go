package rds

import (
	"time"

	"github.com/mownier/duyog/data/store"
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/progerr"

	"github.com/garyburd/redigo/redis"
)

type songRepo struct {
	pool   *redis.Pool
	keyGen generator.Key
}

func (r songRepo) Create(s store.Song, ar []store.ArtistKey, al []store.AlbumKey) (store.Songs, error) {
	songs := store.NewSongs()

	if s.Title == "" || s.Duration <= 0 || s.Year <= 0 {
		return songs, progerr.SongInvalidInfo
	}

	if len(ar) == 0 {
		return songs, progerr.SongHasNoArtist
	}

	conn := r.pool.Get()
	defer conn.Close()

	key := generator.GenerateKey(r.keyGen)

	_, err := conn.Do("HMSET", "song:"+key,
		"id", key,
		"genre", s.Genre,
		"audio_url", s.AudioURL,
		"title", s.Title,
		"duration", s.Duration,
		"year", s.Year)

	if err != nil {
		return songs, progerr.Internal(err)
	}

	var artistKeys []store.ArtistKey
	artists := store.Artists{}

	for _, v := range ar {
		if v == "" {
			continue
		}

		data, err := redis.Values(conn.Do("HGETALL", "artist:"+v))

		if err != nil || len(data) == 0 {
			continue
		}

		var artist store.Artist
		err = redis.ScanStruct(data, &artist)

		if err != nil {
			continue
		}

		_, err = conn.Do("SADD", "song:"+key+":artists", v)

		if err != nil {
			continue
		}

		conn.Do("SADD", "artist:"+v+":songs", key)

		artistKeys = append(artistKeys, v)
		artists[v] = artist
	}

	if len(artistKeys) == 0 {
		return songs, progerr.SongHasNoArtist
	}

	var albumKeys []store.AlbumKey
	albums := store.Albums{}

	for _, v := range al {
		if v == "" {
			continue
		}

		data, err := redis.Values(conn.Do("HGETALL", "album:"+v))

		if err != nil || len(data) == 0 {
			continue
		}

		var album store.Album
		err = redis.ScanStruct(data, &album)

		if err != nil {
			continue
		}

		_, err = conn.Do("SADD", "song:"+key+":albums", v)

		if err != nil {
			continue
		}

		conn.Do("SADD", "album:"+v+":songs", key)

		albumKeys = append(albumKeys, v)
		albums[v] = album
	}

	for _, albumKey := range albumKeys {
		for _, artistKey := range artistKeys {
			conn.Do("SADD", "album:"+albumKey+":artists", artistKey)
		}
	}

	for _, artistKey := range artistKeys {
		for _, albumKey := range albumKeys {
			conn.Do("SADD", "artist:"+artistKey+":albums", albumKey)
		}
	}

	t := time.Now().Unix()
	conn.Do("HMSET", "song:"+key+":timestamp", "created_on", t, "upated_on", t)
	conn.Do("SADD", "songs", key)

	song := s
	song.Key = store.SongKey(key)

	songs.Songs[song.Key] = song
	songs.AlbumKeys[song.Key] = albumKeys
	songs.ArtistKeys[song.Key] = artistKeys
	songs.Albums = albums
	songs.Artists = artists

	return songs, nil
}

func (r songRepo) Update(s store.Song) (store.Song, error) {
	var song store.Song

	if s.Key == "" {
		return song, progerr.SongInvalidKey
	}

	if s.Duration <= 0 && s.Year <= 0 && s.Genre == "" && s.AudioURL == "" && s.Title == "" {
		return song, progerr.SongNothingToUpdate
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("HGETALL", "song:"+s.Key))

	if err != nil {
		return song, progerr.Internal(err)
	}

	so := store.Song{}
	err = redis.ScanStruct(data, &so)

	if err != nil {
		return song, progerr.Internal(err)
	}

	tmp := store.Song{
		Key: s.Key,
	}

	if s.Duration > 0 && s.Duration != so.Duration {
		_, err = conn.Do("HSET", "song:"+s.Key, "duration", s.Duration)

		if err == nil {
			tmp.Duration = s.Duration
		}
	}

	if s.Year > 0 && s.Year != so.Year {
		_, err = conn.Do("HSET", "song:"+s.Key, "year", s.Year)

		if err == nil {
			tmp.Year = s.Year
		}
	}

	if s.Genre == "" && s.Genre != so.Genre {
		_, err = conn.Do("HSET", "song:"+s.Key, "genre", s.Genre)

		if err == nil {
			tmp.Genre = s.Genre
		}
	}

	if s.Title != "" && s.Title != so.Title {
		_, err = conn.Do("HSET", "song:"+s.Key, "title", s.Title)

		if err == nil {
			tmp.Title = s.Title
		}
	}

	if s.AudioURL != "" && s.AudioURL != so.AudioURL {
		_, err = conn.Do("HSET", "song:"+s.Key, "audio_url", s.AudioURL)

		if err == nil {
			tmp.AudioURL = s.AudioURL
		}
	}

	if tmp.Duration <= 0 && tmp.Year <= 0 && tmp.Genre == "" && tmp.Title == "" && tmp.AudioURL == "" {
		return song, progerr.SongNothingToUpdate
	}

	song = tmp

	return song, nil
}

func (r songRepo) UpdateAlbums(k store.SongKey, a []store.AlbumKey) ([]store.AlbumKey, error) {
	var keys []store.AlbumKey
	var strings []string

	for _, key := range a {
		strings = append(strings, string(key))
	}

	tmp, err := r.update("album", string(k), strings)

	if err != nil {
		return keys, err
	}

	for _, key := range tmp {
		keys = append(keys, store.AlbumKey(key))
	}

	return keys, nil
}

func (r songRepo) UpdateArtists(k store.SongKey, a []store.ArtistKey) ([]store.ArtistKey, error) {
	var keys []store.ArtistKey
	var strings []string

	for _, key := range a {
		strings = append(strings, string(key))
	}

	tmp, err := r.update("artist", string(k), strings)

	if err != nil {
		return keys, err
	}

	for _, key := range tmp {
		keys = append(keys, store.ArtistKey(key))
	}

	return keys, nil
}

func (r songRepo) GetByKey(k store.SongKey) (store.Songs, error) {
	var songs store.Songs

	if k == "" {
		return songs, progerr.SongInvalidKey
	}

	conn := r.pool.Get()
	defer conn.Close()

	keys := []store.SongKey{k}
	tmp, err := getSongs(keys, conn, progerr.SongNotFound)

	if err != nil {
		return songs, err
	}

	songs = tmp

	return songs, nil
}

func (r songRepo) update(name string, k string, a []string) ([]string, error) {
	var keys []string

	if k == "" {
		return keys, progerr.SongInvalidKey
	}

	if len(a) == 0 {
		return keys, progerr.SongNothingToUpdate
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := conn.Do("EXISTS", "song:"+k)

	if err != nil {
		return keys, progerr.Internal(err)
	}

	if data.(int64) == 0 {
		return keys, progerr.SongNotFound
	}

	var akeys []string

	for _, key := range a {
		data, err = conn.Do("EXISTS", name+":"+key)

		if err != nil || data.(int64) == 0 {
			continue
		}

		data, err = conn.Do("SISMEMBER", "song:"+k+":"+name+"s", key)

		if err != nil || data.(int64) == 1 {
			continue
		}

		akeys = append(akeys, key)
	}

	if len(akeys) == 0 {
		return keys, progerr.SongNothingToUpdate
	}

	_, err = conn.Do("DEL", "song:"+k+":"+name+"s")

	if err != nil {
		return keys, progerr.Internal(err)
	}

	var tmp []string

	for _, key := range akeys {
		_, err = conn.Do("SADD", "song:"+k+":"+name+"s", key)

		if err != nil {
			continue
		}

		conn.Do("SADD", name+":"+key+":songs", k)
		tmp = append(tmp, key)
	}

	if len(tmp) == 0 {
		return keys, progerr.SongNothingToUpdate
	}

	keys = tmp

	return keys, nil
}

// SongRepo method
func SongRepo(p *redis.Pool, g generator.Key) store.SongRepo {
	return songRepo{
		pool:   p,
		keyGen: g,
	}
}
