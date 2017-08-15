package store

// AlbumKey is the key type of Album
type AlbumKey string

// ArtistKey is the key type of Artist
type ArtistKey string

// PlaylistKey is the key type of Playlist
type PlaylistKey string

// SongKey is the key type of Song
type SongKey string

// UserKey is the key type of User
type UserKey string

// Album struct
type Album struct {
	Key   AlbumKey `json:"id,omitempty" redis:"id"`
	Year  int64    `json:"year,omitempty" redis:"year"`
	Desc  string   `json:"description,omitempty" redis:"description"`
	Photo string   `json:"photo,omitempty" redis:"photo"`
	Title string   `json:"title,omitempty" redis:"title"`
}

// Artist struct
type Artist struct {
	Key    ArtistKey `json:"id,omitempty" redis:"id"`
	Bio    string    `json:"bio,omitempty" redis:"bio"`
	Name   string    `json:"name,omitempty" redis:"name"`
	Genre  string    `json:"genre,omitempty" redis:"genre"`
	Avatar string    `json:"avatar,omitempty" redis:"avatar"`
}

// Playlist struct
type Playlist struct {
	Key   PlaylistKey `json:"id,omitempty" redis:"id"`
	Name  string      `json:"name,omitempty" redis:"name"`
	Desc  string      `json:"description,omitempty" redis:"description"`
	Photo string      `json:"photo,omitempty" redis:"photo"`
}

// Song data struct
type Song struct {
	Key      SongKey `json:"id,omitempty" redis:"id"`
	Year     int64   `json:"year,omitempty" redis:"year"`
	Genre    string  `json:"genre,omitempty" redis:"genre"`
	Title    string  `json:"title,omitempty" redis:"title"`
	Duration float64 `json:"duration,omitempty" redis:"duration"`
	AudioURL string  `json:"audio_url,omitempty" redis:"audio_url"`
}

// User struct
type User struct {
	Key       UserKey `json:"id,omitempty" redis:"id"`
	Email     string  `json:"email,omitempty" redis:"email"`
	Avatar    string  `json:"avatar,omitempty" redis:"avatar"`
	LastName  string  `json:"last_name,omitempty" redis:"last_name"`
	FirstName string  `json:"first_name,omitempty" redis:"first_name"`
}

// Albums map
type Albums map[AlbumKey]Album

// Artists map
type Artists map[ArtistKey]Artist

// Users map
type Users map[UserKey]User

// Playlists struct
type Playlists struct {
	Creators  Users                    `json:"creators,omitempty"`
	Playlists map[PlaylistKey]Playlist `json:"playlists"`
}

// Songs struct
type Songs struct {
	Songs   map[SongKey]Song `json:"songs"`
	Albums  Albums           `json:"albums,omitempty"`
	Artists Artists          `json:"artists,omitempty"`

	AlbumKeys  map[SongKey][]AlbumKey  `json:"album_keys,omitempty"`
	ArtistKeys map[SongKey][]ArtistKey `json:"artist_keys,omitempty"`
}

// NewSongs method
func NewSongs() Songs {
	return Songs{
		Songs:   map[SongKey]Song{},
		Albums:  Albums{},
		Artists: Artists{},

		AlbumKeys:  map[SongKey][]AlbumKey{},
		ArtistKeys: map[SongKey][]ArtistKey{},
	}
}
