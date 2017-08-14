package rds

import (
	"github.com/mownier/duyog/data/store"

	"github.com/garyburd/redigo/redis"
)

func getSongs(keys []store.SongKey, conn redis.Conn, e error) (store.Songs, error) {
	songs := store.NewSongs()

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

		songs.ArtistKeys[key] = []store.ArtistKey{}

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

			songs.ArtistKeys[key] = append(songs.ArtistKeys[key], artistKey)
			songs.Artists[artistKey] = artist
		}

		if len(songs.ArtistKeys[key]) == 0 {
			continue
		}

		songs.Songs[key] = song
		data, err = redis.Values(conn.Do("SMEMBERS", "song:"+key+":albums"))

		if err != nil || len(data) == 0 {
			continue
		}

		songs.AlbumKeys[key] = []store.AlbumKey{}

		for _, v := range data {
			if len(v.([]byte)) == 0 {
				continue
			}

			albumKey := store.AlbumKey(v.([]byte)[:])
			data, err := redis.Values(conn.Do("HGETALL", "album:"+albumKey))

			if err != nil || len(data) == 0 {
				continue
			}

			var album store.Album
			err = redis.ScanStruct(data, &album)

			if err != nil {
				continue
			}

			songs.AlbumKeys[key] = append(songs.AlbumKeys[key], albumKey)
			songs.Albums[albumKey] = album
		}
	}

	if len(songs.Songs) == 0 {
		return songs, e
	}

	return songs, nil
}
