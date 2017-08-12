package rds

import (
	"duyog/data/store"

	"github.com/garyburd/redigo/redis"
)

func getSongs(keys []store.SongKey, conn redis.Conn, e error) (songs store.Songs, err error) {
	var tmp store.Songs

	for _, key := range keys {
		data, err := redis.Values(conn.Do("HGETALL", "song:"+key))

		if err != nil || len(data) == 0 {
			continue
		}

		var song store.Song
		err = redis.ScanStruct(data, &song)

		if err != nil {
			continue
		}

		data, err = redis.Values(conn.Do("SMEMBERS", "song:"+key+":artists"))

		if err != nil || len(data) == 0 {
			continue
		}

		tmp.ArtistKeys[key] = []store.ArtistKey{}

		for _, v := range data {
			if len(v.([]byte)) == 0 {
				continue
			}

			artistKey := store.ArtistKey(v.([]byte)[:])
			data, err := redis.Values(conn.Do("HGETALL", "artist:"+artistKey))

			if err != nil || len(data) == 0 {
				continue
			}

			var artist store.Artist
			err = redis.ScanStruct(data, &artist)

			if err != nil {
				continue
			}

			tmp.ArtistKeys[key] = append(tmp.ArtistKeys[key], artistKey)
			tmp.Artists[artistKey] = artist
		}

		if len(tmp.ArtistKeys[key]) == 0 {
			continue
		}

		tmp.Songs[key] = song
		data, err = redis.Values(conn.Do("SMEMBERS", "song:"+key+":albums"))

		if err != nil || len(data) == 0 {
			continue
		}

		tmp.AlbumKeys[key] = []store.AlbumKey{}

		for _, v := range data {
			if len(v.([]byte)) == 0 {
				continue
			}

			albumKey := store.AlbumKey(v.([]byte)[:])
			data, err := redis.Values(conn.Do("HGETALL", "artist:"+albumKey))

			if err != nil || len(data) == 0 {
				continue
			}

			var album store.Album
			err = redis.ScanStruct(data, &album)

			if err != nil {
				continue
			}

			tmp.AlbumKeys[key] = append(tmp.AlbumKeys[key], albumKey)
			tmp.Albums[albumKey] = album
		}
	}

	if len(tmp.Songs) == 0 {
		return songs, e
	}

	songs = tmp

	return songs, nil
}
