package service

type pathCode int

const (
	// pathCodeSongCreate code
	pathCodeSongCreate pathCode = 1000

	// pathCodeSongUpdate code
	pathCodeSongUpdate pathCode = 1001

	// pathCodeSongGetInfo code
	pathCodeSongGetInfo pathCode = 1002

	// pathCodeSongUpdateAlbums code
	pathCodeSongUpdateAlbums pathCode = 1003

	// pathCodeSongUpdateArtists code
	pathCodeSongUpdateArtists pathCode = 1004

	// pathCodeAlbumCreate code
	pathCodeAlbumCreate pathCode = 2000

	// pathCodeAlbumUpdate code
	pathCodeAlbumUpdate pathCode = 2001

	// pathCodeAlbumGetInfo code
	pathCodeAlbumGetInfo pathCode = 2002

	// pathCodeAlbumGetSongs code
	pathCodeAlbumGetSongs pathCode = 2003

	// pathCodeAlbumGetArtists code
	pathCodeAlbumGetArtists pathCode = 2004

	// pathCodeArtistCreate code
	pathCodeArtistCreate pathCode = 3000

	// pathCodeArtistUpdate code
	pathCodeArtistUpdate pathCode = 3001

	// pathCodeArtistGetInfo code
	pathCodeArtistGetInfo pathCode = 3002

	// pathCodeArtistGetSongs code
	pathCodeArtistGetSongs pathCode = 3003

	// pathCodeArtistGetAlbums code
	pathCodeArtistGetAlbums pathCode = 3004

	// pathCodePlaylistCreate code
	pathCodePlaylistCreate pathCode = 4000

	// pathCodePlaylistUpdate code
	pathCodePlaylistUpdate pathCode = 4001

	// pathCodePlaylistGetInfo code
	pathCodePlaylistGetInfo pathCode = 4002

	// pathCodePlaylistGetSongs code
	pathCodePlaylistGetSongs pathCode = 4003

	// pathCodePlaylistAddSongs code
	pathCodePlaylistAddSongs pathCode = 4004

	// pathCodeMeUpdate code
	pathCodeMeUpdate pathCode = 5000

	// pathCodeMeGetPlaylists code
	pathCodeMeGetPlaylists pathCode = 5001

	// pathCodeMeAddPlaylists code
	pathCodeMeAddPlaylists pathCode = 5002

	// pathCodeMeUpdatePlaylist code
	pathCodeMeUpdatePlaylist pathCode = 5003

	// pathCodeMeAddPlaylistSongs code
	pathCodeMeAddPlaylistSongs pathCode = 5004

	// pathCodeUserRegister code
	pathCodeUserRegister pathCode = 6000

	// pathCodeUserGetPlaylists code
	pathCodeUserGetPlaylists pathCode = 6001

	// pathCodeUserGetInfo code
	pathCodeUserGetInfo pathCode = 6002
)

var forbiddenPath = map[pathCode]bool{
	pathCodeSongCreate:        true,
	pathCodeSongUpdate:        true,
	pathCodeSongGetInfo:       false,
	pathCodeSongUpdateAlbums:  true,
	pathCodeSongUpdateArtists: true,

	pathCodeAlbumCreate:     true,
	pathCodeAlbumUpdate:     true,
	pathCodeAlbumGetInfo:    false,
	pathCodeAlbumGetSongs:   false,
	pathCodeAlbumGetArtists: false,

	pathCodeArtistCreate:    true,
	pathCodeArtistUpdate:    true,
	pathCodeArtistGetInfo:   false,
	pathCodeArtistGetSongs:  false,
	pathCodeArtistGetAlbums: false,

	pathCodePlaylistCreate:   false,
	pathCodePlaylistUpdate:   true,
	pathCodePlaylistGetInfo:  false,
	pathCodePlaylistGetSongs: false,
	pathCodePlaylistAddSongs: true,

	pathCodeMeUpdate:           false,
	pathCodeMeGetPlaylists:     false,
	pathCodeMeAddPlaylists:     false,
	pathCodeMeUpdatePlaylist:   false,
	pathCodeMeAddPlaylistSongs: false,

	pathCodeUserRegister:     false,
	pathCodeUserGetPlaylists: false,
	pathCodeUserGetInfo:      false,
}

func isForbidden(c pathCode) bool {
	return forbiddenPath[c]
}
