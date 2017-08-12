package service

type pathCode int

const (
	// pathCodeFileUploadMusic code
	pathCodeFileUploadMusic pathCode = 1000

	// pathCodeFileUploadAlbumPhoto code
	pathCodeFileUploadAlbumPhoto pathCode = 1001

	// pathCodeFileUploadUserAvatar code
	pathCodeFileUploadUserAvatar pathCode = 1002

	// pathCodeFileUploadArtistPhoto code
	pathCodeFileUploadArtistPhoto pathCode = 1003

	// pathCodeFileUploadPlaylistPhoto code
	pathCodeFileUploadPlaylistPhoto pathCode = 1004

	// pathCodeFileGetMusic code
	pathCodeFileGetMusic pathCode = 1005

	// pathCodeFileGetAlbumPhoto code
	pathCodeFileGetAlbumPhoto pathCode = 1006

	// pathCodeFileGetUserAvatar code
	pathCodeFileGetUserAvatar pathCode = 1007

	// pathCodeFileGetArtistPhoto code
	pathCodeFileGetArtistPhoto pathCode = 1008

	// pathCodeFileGetPlaylistPhoto code
	pathCodeFileGetPlaylistPhoto pathCode = 1009
)

func pathResource(c pathCode) string {
	switch c {
	case pathCodeFileGetMusic:
		fallthrough
	case pathCodeFileUploadMusic:
		return "song"

	case pathCodeFileGetAlbumPhoto:
		fallthrough
	case pathCodeFileUploadAlbumPhoto:
		return "album"

	case pathCodeFileGetUserAvatar:
		fallthrough
	case pathCodeFileUploadUserAvatar:
		return "user"

	case pathCodeFileGetArtistPhoto:
		fallthrough
	case pathCodeFileUploadArtistPhoto:
		return "artist"

	case pathCodeFileGetPlaylistPhoto:
		fallthrough
	case pathCodeFileUploadPlaylistPhoto:
		return "playlist"

	default:
		return ""
	}
}

var forbiddenPath = map[pathCode]bool{
	pathCodeFileUploadMusic:         false,
	pathCodeFileUploadAlbumPhoto:    true,
	pathCodeFileUploadUserAvatar:    true,
	pathCodeFileUploadArtistPhoto:   true,
	pathCodeFileUploadPlaylistPhoto: true,

	pathCodeFileGetMusic:         true,
	pathCodeFileGetAlbumPhoto:    true,
	pathCodeFileGetUserAvatar:    true,
	pathCodeFileGetArtistPhoto:   true,
	pathCodeFileGetPlaylistPhoto: true,
}

func isForbidden(c pathCode) bool {
	return forbiddenPath[c]
}
