package service

import (
	"github.com/mownier/duyog/extractor"
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/validator"
	"github.com/mownier/duyog/writer"

	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// ResourceHandler function
type ResourceHandler func(http.ResponseWriter, *http.Request)

// Resource struct
type Resource struct {
	FormKey string
	MaxMem  int64

	AuthExtractor  extractor.Auth
	AuthValidator  validator.Auth
	DataValidator  validator.Data
	ResponseWriter writer.Response
}

type uploadInput struct {
	dir string

	formKey         string
	sizeLimit       int64
	multipartMaxMem int64
	validExtensions []string

	keyGen func(string) (string, error)

	errExtension    error
	errExceededSize error

	writer writer.Response
}

func writeRespErr(w http.ResponseWriter, r *http.Request, wr writer.Response, e error) {
	var err progerr.Err

	switch e.(type) {
	case progerr.Err:
		err = e.(progerr.Err)

	default:
		err = progerr.Internal(err)
	}

	writer.WriteResponse(wr, w, r, err.HTTPStatus, err.Data())
}

func marshalResp(w http.ResponseWriter, r *http.Request, wr writer.Response, v interface{}) bool {
	data, err := json.Marshal(v)

	if err != nil {
		writeRespErr(w, r, wr, err)
		return false
	}

	writer.WriteResponse(wr, w, r, http.StatusOK, data)

	return true
}

func serveFile(w http.ResponseWriter, r *http.Request, wr writer.Response, dir string) bool {
	path := filepath.Join(dir, mux.Vars(r)["name"])
	stat, err := os.Stat(path)

	if err != nil {
		writeRespErr(w, r, wr, err)
		return false
	}

	if stat.IsDir() {
		writeRespErr(w, r, wr, progerr.FileNotFound)
		return false
	}

	http.ServeFile(w, r, path)

	return true
}

func extractAuth(a extractor.Auth, r *http.Request) validator.AccessToken {
	switch r.Method {
	case http.MethodGet:
		return validator.AccessToken(mux.Vars(r)["token"])

	default:
		return validator.AccessToken(extractor.ExtractAuth(a, r.Header.Get("Authorization")))
	}
}

func validate(w http.ResponseWriter, r *http.Request, res Resource, c pathCode, e error) bool {
	token := extractAuth(res.AuthExtractor, r)
	auth, err := validator.ValidateAuth(res.AuthValidator, token)

	if err != nil {
		writeRespErr(w, r, res.ResponseWriter, err)
		return false
	}

	if isForbidden(c) && strings.ToLower(auth.Role) != "admin" {
		writeRespErr(w, r, res.ResponseWriter, progerr.RequestPathForbidden)
		return false
	}

	var reply validator.DataReply
	id := mux.Vars(r)["id"]

	switch pathResource(c) {
	case "album":
		key := validator.AlbumKey(id)
		reply = validator.ValidateAlbum(res.DataValidator.Album(), key)

	case "atist":
		key := validator.ArtistKey(id)
		reply = validator.ValidateArtist(res.DataValidator.Artist(), key)

	case "playlist":
		key := validator.PlaylistKey(id)
		reply = validator.ValidatePlaylist(res.DataValidator.Playlist(), key)

	case "song":
		key := validator.SongKey(id)
		reply = validator.ValidateSong(res.DataValidator.Song(), key)

	case "user":
		key := validator.UserKey(id)
		reply = validator.ValidateUser(res.DataValidator.User(), key)

	default:
		reply = validator.DataReply{
			OK: false,
		}
	}

	if reply.OK == false {
		writeRespErr(w, r, res.ResponseWriter, e)
		return false
	}

	return true
}

func containsExtension(ext string, exts []string) bool {
	if ext == "" || len(exts) == 0 {
		return false
	}

	for _, v := range exts {
		if v == ext {
			return true
		}
	}

	return false
}

func upload(w http.ResponseWriter, r *http.Request, i uploadInput) bool {
	err := r.ParseMultipartForm(i.multipartMaxMem)

	if err != nil {
		writeRespErr(w, r, i.writer, err)
		return false
	}

	file, header, err := r.FormFile(i.formKey)
	defer file.Close()

	if err != nil {
		writeRespErr(w, r, i.writer, err)
		return false
	}

	size, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)

	if err != nil {
		writeRespErr(w, r, i.writer, err)
		return false
	}

	if size > i.sizeLimit {
		writeRespErr(w, r, i.writer, i.errExceededSize)
		return false
	}

	ext := filepath.Ext(header.Filename)

	if containsExtension(ext, i.validExtensions) == false {
		writeRespErr(w, r, i.writer, i.errExtension)
		return false
	}

	fileKey, err := i.keyGen(mux.Vars(r)["id"])

	if err != nil {
		writeRespErr(w, r, i.writer, err)
		return false
	}

	filename := fileKey + ext
	dst, err := os.Create(filepath.Join(i.dir, filename))
	defer dst.Close()

	if err != nil {
		writeRespErr(w, r, i.writer, err)
		return false
	}

	_, err = io.Copy(dst, file)

	if err != nil {
		writeRespErr(w, r, i.writer, err)
		return false
	}

	var scheme = "http"

	if r.URL.Scheme != "" {
		scheme = r.URL.Scheme
	}

	url := url.URL{
		Scheme: scheme,
		Host:   r.Host,
		Path:   filepath.Join(filepath.Dir(r.URL.Path), filename),
	}

	resp := map[string]string{
		"id":  fileKey,
		"url": url.String(),
	}

	return marshalResp(w, r, i.writer, resp)
}

func download(w http.ResponseWriter, r *http.Request, wr writer.Response, dir string, verifyFunc func(string, string) error) bool {
	vars := mux.Vars(r)
	fname := strings.Replace(vars["name"], filepath.Ext(vars["name"]), "", -1)
	err := verifyFunc(vars["id"], fname)

	if err != nil {
		writeRespErr(w, r, wr, err)
		return false
	}

	serveFile(w, r, wr, dir)
	return true
}
