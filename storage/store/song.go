package store

// SongRepo interface
type SongRepo interface {
	AddAudio(k string) (string, error)
	Verify(k, fk string) error
}

// AddSongAudio method
func AddSongAudio(r SongRepo, k string) (string, error) {
	return r.AddAudio(k)
}

// VerifySong method
func VerifySong(r SongRepo, k, fk string) error {
	return r.Verify(k, fk)
}
