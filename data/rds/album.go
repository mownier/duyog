package rds

import (
	"duyog/data/store"
	"duyog/generator"
	"duyog/progerr"
	"time"

	"github.com/garyburd/redigo/redis"
)

type albumRepo struct {
	pool   *redis.Pool
	keyGen generator.Key
}

func (r albumRepo) Create(a store.Album) (store.Album, error) {
	var album store.Album

	if a.Year <= 0 || a.Photo == "" || a.Title == "" {
		return album, progerr.AlbumInvalidInfo
	}

	conn := r.pool.Get()
	defer conn.Close()

	key := generator.GenerateKey(r.keyGen)

	_, err := conn.Do("HMSET", "album:"+key,
		"id", key,
		"year", a.Year,
		"photo", a.Photo,
		"title", a.Title,
		"description", a.Desc)

	if err != nil {
		return album, progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("SADD", "albums", key)
	conn.Do("HMSET", "album:"+key+":timestamp", "created_on", t, "updated_on", t)

	album = a
	album.Key = store.AlbumKey(key)

	return album, nil
}

func (r albumRepo) Update(a store.Album) (store.Album, error) {
	var album store.Album

	if a.Key == "" {
		return album, progerr.AlbumInvalidKey
	}

	if a.Year <= 0 && a.Photo == "" && a.Title == "" && a.Desc == "" {
		return album, progerr.AlbumNothingToUpdate
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("HGETALL", "album:"+a.Key))

	if err != nil {
		return album, progerr.Internal(err)
	}

	if len(data) == 0 {
		return album, progerr.AlbumNotFound
	}

	var al store.Album
	err = redis.ScanStruct(data, &al)

	if err != nil {
		return album, progerr.Internal(err)
	}

	tmp := store.Album{
		Key: a.Key,
	}

	if a.Year > 0 && a.Year != al.Year {
		_, err = conn.Do("HSET", "album:"+a.Key, "year", a.Year)

		if err == nil {
			tmp.Year = a.Year
		}
	}

	if a.Photo == "" && a.Photo != al.Photo {
		_, err = conn.Do("HSET", "album:"+a.Key, "photo", a.Photo)

		if err == nil {
			tmp.Photo = a.Photo
		}
	}

	if a.Title == "" && a.Title != al.Title {
		_, err := conn.Do("HSET", "album:"+a.Key, "title", a.Title)

		if err == nil {
			tmp.Title = a.Title
		}
	}

	if a.Desc == "" && a.Desc != al.Desc {
		_, err = conn.Do("HSET", "artist:"+a.Key, "description", a.Desc)

		if err == nil {
			tmp.Desc = a.Desc
		}
	}

	if tmp.Year <= 0 && tmp.Photo == "" && tmp.Title == "" && tmp.Desc == "" {
		return album, progerr.AlbumNothingToUpdate
	}

	conn.Do("HSET", "album:"+a.Key+":timestamp", "updated_on", time.Now().Unix())

	album = tmp

	return album, nil
}

func (r albumRepo) GetByKey(k store.AlbumKey) (store.Album, error) {
	var album store.Album

	if k == "" {
		return album, progerr.AlbumInvalidKey
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("HGETALL", "album:"+k))

	if err != nil {
		return album, progerr.Internal(err)
	}

	if len(data) == 0 {
		return album, progerr.AlbumNotFound
	}

	var al store.Album
	err = redis.ScanStruct(data, &al)

	if err != nil {
		return album, progerr.Internal(err)
	}

	return album, nil
}

func (r albumRepo) GetSongs(k store.AlbumKey) (store.Songs, error) {
	var songs store.Songs

	if k == "" {
		return songs, progerr.AlbumInvalidKey
	}

	_, err := store.GetAlbumByKey(r, k)

	if err != nil {
		return songs, err
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("SMEMBERS", "album:"+k+":songs"))

	if err != nil {
		return songs, progerr.Internal(err)
	}

	if len(data) == 0 {
		return songs, progerr.AlbumHasNoSongs
	}

	var keys []store.SongKey

	for _, v := range data {
		if len(v.([]byte)) == 0 {
			continue
		}

		keys = append(keys, store.SongKey(v.([]byte)[:]))
	}

	if len(keys) == 0 {
		return songs, progerr.AlbumHasNoSongs
	}

	tmp, err := getSongs(keys, conn, progerr.AlbumHasNoSongs)

	if err != nil {
		return songs, err
	}

	songs = tmp

	return songs, nil
}

func (r albumRepo) GetArtists(k store.AlbumKey) (store.Artists, error) {
	var artists store.Artists

	if k == "" {
		return artists, progerr.AlbumInvalidKey
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("SMEMBERS", "album:"+k+":artists"))

	if err != nil {
		return artists, progerr.Internal(err)
	}

	if len(data) == 0 {
		return artists, progerr.AlbumHasNoArtists
	}

	var tmp store.Artists

	for _, v := range data {
		if len(v.([]byte)) == 0 {
			continue
		}

		key := store.ArtistKey(string(v.([]byte)[:]))

		if key == "" {
			continue
		}

		data, err := redis.Values(conn.Do("HGETALL", "artist:"+key))

		if err != nil || len(data) == 0 {
			continue
		}

		var artist store.Artist
		err = redis.ScanStruct(data, &artist)

		if err != nil {
			continue
		}

		tmp[key] = artist
	}

	if len(tmp) == 0 {
		return artists, progerr.AlbumHasNoArtists
	}

	artists = tmp

	return artists, nil
}

// AlbumRepo method
func AlbumRepo(p *redis.Pool, g generator.Key) store.AlbumRepo {
	return albumRepo{
		pool:   p,
		keyGen: g,
	}
}
