package rds

import (
	"time"

	"github.com/mownier/duyog/data/store"
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/progerr"

	"github.com/garyburd/redigo/redis"
)

type artistRepo struct {
	pool   *redis.Pool
	keyGen generator.Key
}

func (r artistRepo) Create(a store.Artist) (store.Artist, error) {
	var artist store.Artist

	if a.Bio == "" || a.Genre == "" || a.Name == "" {
		return artist, progerr.ArtistInvalidInfo
	}

	conn := r.pool.Get()
	defer conn.Close()

	key := generator.GenerateKey(r.keyGen)

	_, err := conn.Do("HMSET", "artist:"+key,
		"id", a.Key,
		"bio", a.Bio,
		"name", a.Name,
		"genre", a.Genre,
		"avatar", a.Avatar)

	if err != nil {
		return artist, progerr.Internal(err)
	}

	t := time.Now().Unix()
	conn.Do("HMSET", "artist:"+key+":timestamp", "created_on", t, "updated_on", t)
	conn.Do("SADD", "artists", key)

	artist = a
	artist.Key = store.ArtistKey(key)

	return artist, nil
}

func (r artistRepo) Update(a store.Artist) (store.Artist, error) {
	var artist store.Artist

	if a.Key == "" {
		return artist, progerr.ArtistInvalidKey
	}

	if a.Bio == "" && a.Genre == "" && a.Name == "" && a.Avatar == "" {
		return artist, progerr.ArtistNothingToUpdate
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("HGETALL", "artist:"+a.Key))

	if err != nil {
		return artist, progerr.Internal(err)
	}

	if len(data) == 0 {
		return artist, progerr.ArtistNotFound
	}

	var ar store.Artist
	err = redis.ScanStruct(data, &ar)

	if err != nil {
		return artist, progerr.Internal(err)
	}

	tmp := store.Artist{
		Key: a.Key,
	}

	if a.Bio != "" && a.Bio != ar.Bio {
		_, err = conn.Do("HSET", "artist:"+a.Key, "bio", a.Bio)

		if err == nil {
			tmp.Bio = a.Bio
		}
	}

	if a.Avatar != "" && a.Avatar != ar.Avatar {
		_, err = conn.Do("HSET", "artist:"+a.Key, "avatar", a.Avatar)

		if err == nil {
			tmp.Avatar = a.Avatar
		}
	}

	if a.Genre != "" && a.Genre != ar.Genre {
		_, err = conn.Do("HSET", "artist:"+a.Key, "genre", a.Genre)

		if err == nil {
			tmp.Genre = a.Genre
		}
	}

	if a.Name != "" && a.Name != ar.Name {
		_, err = conn.Do("HSET", "artist:"+a.Key, "name", a.Name)

		if err == nil {
			tmp.Name = a.Name
		}
	}

	if tmp.Bio == "" && tmp.Avatar == "" && tmp.Genre == "" && tmp.Name == "" {
		return artist, progerr.ArtistNothingToUpdate
	}

	conn.Do("HSET", "artist:"+a.Key+":timestamp", "updated_on", time.Now().Unix())

	artist = tmp

	return artist, nil
}

func (r artistRepo) GetByKey(k store.ArtistKey) (store.Artist, error) {
	var artist store.Artist

	if k == "" {
		return artist, progerr.ArtistInvalidKey
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("HGETALL", "artist:"+k))

	if err != nil {
		return artist, progerr.Internal(err)
	}

	if len(data) == 0 {
		return artist, progerr.ArtistNotFound
	}

	var tmp store.Artist
	err = redis.ScanStruct(data, &tmp)

	if err != nil {
		return artist, progerr.Internal(err)
	}

	artist = tmp

	return artist, nil
}

func (r artistRepo) GetSongs(k store.ArtistKey) (store.Songs, error) {
	var songs store.Songs

	if k == "" {
		return songs, progerr.ArtistInvalidKey
	}

	_, err := store.GetArtistByKey(r, k)

	if err != nil {
		return songs, err
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("SMEMBERS", "artist:"+k+":songs"))

	if err != nil {
		return songs, progerr.Internal(err)
	}

	if len(data) == 0 {
		return songs, progerr.ArtistHasNoSongs
	}

	var keys []store.SongKey

	for _, v := range data {
		if len(v.([]byte)) == 0 {
			continue
		}

		keys = append(keys, store.SongKey(v.([]byte)[:]))
	}

	if len(keys) == 0 {
		return songs, progerr.ArtistHasNoSongs
	}

	tmp, err := getSongs(keys, conn, progerr.ArtistHasNoSongs)

	if err != nil {
		return songs, err
	}

	songs = tmp

	return songs, nil
}

func (r artistRepo) GetAlbums(k store.ArtistKey) (store.Albums, error) {
	albums := store.Albums{}

	if k == "" {
		return albums, progerr.ArtistInvalidKey
	}

	conn := r.pool.Get()
	defer conn.Close()

	data, err := redis.Values(conn.Do("SMEMBERS", "artist:"+k+":albums"))

	if err != nil {
		return albums, progerr.Internal(err)
	}

	if len(data) == 0 {
		return albums, progerr.ArtistHasNoAlbums
	}

	for _, v := range data {
		if len(v.([]byte)) == 0 {
			continue
		}

		key := store.AlbumKey(v.([]byte)[:])

		if key == "" {
			continue
		}

		data, err := redis.Values(conn.Do("HGETALL", "album:"+key))

		if err != nil || len(data) == 0 {
			continue
		}

		var album store.Album
		err = redis.ScanStruct(data, &album)

		if err != nil {
			continue
		}

		albums[key] = album
	}

	if len(albums) == 0 {
		return albums, progerr.ArtistHasNoAlbums
	}

	return albums, nil
}

// ArtistRepo method
func ArtistRepo(p *redis.Pool, g generator.Key) store.ArtistRepo {
	return artistRepo{
		pool:   p,
		keyGen: g,
	}
}
