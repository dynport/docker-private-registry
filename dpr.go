package main

import (
	"time"
	"io"
	"encoding/json"
	"strings"
	"fmt"
	"net/http"
	"flag"
	"regexp"
	"path/filepath"
	"io/ioutil"
	"os"
)

var logger = &Logger{}

type Image struct {
	Dir string
}

func (i* Image) Id() (id string) {
	return filepath.Base(i.Dir)
}

func (i* Image) LayerPath() (id string) {
	return i.Dir + "/layer"
}

func (i* Image) Ancestry() (a []string) {
	a = []string{i.Id()}
	current := i
	for {
		atts, err := current.Attributes()
		if err != nil {
			logger.Error(err.Error())
			break
		}
		if atts.Parent != "" {
			a = append(a, atts.Parent)
			current = &Image{filepath.Dir(current.Dir) + "/" + atts.Parent}
		} else {
			break
		}
	}
	return
}

func (i* Image) Attributes() (a *ImageAttributes, err error) {
	a = &ImageAttributes{}
	path := i.Dir + "/json"
	logger.Info("reading attributes from path", path)
	if data, err := ioutil.ReadFile(path); err == nil {
		err = json.Unmarshal(data, a)
	}
	return
}

type ImageAttributes struct {
	Id, Parent, Container string
}

type Repository struct {
	Dir string
}

func (r* Repository) Images() (b []byte, err error) {
	return ioutil.ReadFile(r.Dir + "/images")
}

func writeFile(path string, r io.ReadCloser) (e error) {
	started := time.Now()
	logger.Info("writing to ", path)
	e = os.MkdirAll(filepath.Dir(path), 0755)
	if e != nil { return }

	tmpName := path + ".tmp"
	out, e := os.Create(tmpName)
	if e != nil { return }
	defer out.Close()
	cnt, e := io.Copy(out, r)
	if e != nil { return }
	logger.Info(fmt.Sprintf("Wrote %d bytes in %.06f", cnt, time.Now().Sub(started).Seconds()))
	e = os.Rename(tmpName, path)
	return
}

func (r* Repository) ImagesPath() string {
	return r.Dir + "/images"
}

func (r* Repository) IndexPath() string {
	return r.Dir + "/_index"
}

func (r* Repository) Tags() (m map[string]string) {
	m = make(map[string]string)
	files, err := filepath.Glob(r.Dir + "/tags/*")
	if err != nil {
		return
	}
	for _, path := range files {
		name := filepath.Base(path)
		if data, err := ioutil.ReadFile(path); err == nil {
			m[name] = strings.Replace(string(data), `"`, "", -1)
		}
	}
	return
}

type Handler struct {
	DataDir string
	Mappings []*Mapping
}

func (h* Handler) WriteJsonHeader(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

func (h* Handler) WriteEndpointsHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Docker-Endpoints", r.Host)
}

func (h* Handler) GetPing(w http.ResponseWriter, r *http.Request, p [][]string) {
	w.Header().Add("X-Docker-Registry-Version", "0.0.1")
	w.WriteHeader(200)
	fmt.Fprint(w, "pong")
}

func (h* Handler) GetUsers(w http.ResponseWriter, r *http.Request, p [][]string) {
	w.WriteHeader(201)
	fmt.Fprint(w, "OK")
}

func (h* Handler) GetRepositoryImages(w http.ResponseWriter, r *http.Request, p [][]string) {
	repo := &Repository{h.DataDir + "/repositories/" + p[0][2]}
	if images, err := repo.Images(); err == nil {
		h.WriteJsonHeader(w)
		h.WriteEndpointsHeader(w, r)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(images))
	} else {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	return
}

func (h* Handler) GetImageAncestry(w http.ResponseWriter, r *http.Request, p [][]string) {
	idPrefix := p[0][2]
	if paths, err := filepath.Glob(h.DataDir + "/images/" + idPrefix + "*"); err == nil {
		if len(paths) > 0 {
			image := &Image{paths[0]}
			if out, err := json.Marshal(image.Ancestry()); err == nil {
				h.WriteJsonHeader(w)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, string(out))
				return
			}
		}
	}
	http.NotFound(w, r)
}

func (h* Handler) GetImageLayer(w http.ResponseWriter, r *http.Request, p [][]string) {
	idPrefix := p[0][2]
	if paths, err := filepath.Glob(h.DataDir + "/images/" + idPrefix + "*"); err == nil {
		image := &Image{paths[0]}
		file, err := os.Open(image.LayerPath())
		if err == nil {
			w.Header().Add("Content-Type", "application/x-xz")
			w.WriteHeader(http.StatusOK)
			io.Copy(w, file)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func (h* Handler) GetImageJson(w http.ResponseWriter, r *http.Request, p [][]string) {
	idPrefix := p[0][2]
	if paths, err := filepath.Glob(h.DataDir + "/images/" + idPrefix + "*"); err == nil {
		if len(paths) > 0 {
			image := &Image{paths[0]}
			file, err := os.Open(image.Dir + "/json")
			if err == nil {
				if file, err := os.Open(image.LayerPath()); err == nil {
					if stat, err := file.Stat(); err == nil {
						w.Header().Add("X-Docker-Size", fmt.Sprintf("%d", stat.Size()))
					}
				}
				w.WriteHeader(http.StatusOK)
				io.Copy(w, file)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h* Handler) GetRepositoryTags(w http.ResponseWriter, r *http.Request, p [][]string) {
	repo := &Repository{h.DataDir + "/repositories/" + p[0][2]}
	tagsJson, err := json.Marshal(repo.Tags())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, err.Error())
		return
	}
	h.WriteJsonHeader(w)
	h.WriteEndpointsHeader(w, r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(tagsJson))
	return
}

func (h* Handler) PutImageResource(w http.ResponseWriter, r *http.Request, p [][]string) {
	imageId := p[0][2]
	tagName := p[0][3]

	err := writeFile(h.DataDir + "/images/" + imageId + "/" + tagName, r.Body)
	if err != nil {
		logger.Error(err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h* Handler) PutRepositoryTags(w http.ResponseWriter, r *http.Request, p [][]string) {
	repoName := p[0][2]
	path := h.DataDir + "/repositories/" + repoName + "/tags/" + p[0][3]
	err := writeFile(path, r.Body)
	if err != nil {
		logger.Error(err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h* Handler) PutRepositoryImages(w http.ResponseWriter, r *http.Request, p [][]string) {
	repoName := p[0][2]
	repo := &Repository{h.DataDir + "/repositories/" + repoName}
	err := writeFile(repo.ImagesPath(), r.Body)
	if err != nil {
		logger.Error(err.Error())
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h* Handler) PutRepository(w http.ResponseWriter, r *http.Request, p [][]string) {
	h.WriteJsonHeader(w)
	h.WriteEndpointsHeader(w, r)
	w.Header().Add("WWW-Authenticate" , `Token signature=123abc,repository="dynport/test",access=write`)
	w.Header().Add("X-Docker-Token" , "token")
	w.WriteHeader(http.StatusOK)
	repoName := p[0][2]
	repo := &Repository{h.DataDir + "/repositories/" + repoName}
	err := writeFile(repo.IndexPath(), r.Body)
	if err != nil {
		logger.Error(err.Error())
	}
}

type Mapping struct {
	Method string
	Regexp *regexp.Regexp
	Handler	func(http.ResponseWriter, *http.Request, [][]string)
}

func (h* Handler) Map(t, re string, f func(http.ResponseWriter, *http.Request, [][]string)) {
	if h.Mappings == nil {
		h.Mappings = make([]*Mapping, 0)
	}
	h.Mappings = append(h.Mappings, &Mapping{t, regexp.MustCompile("/v(\\d+)/" + re), f})
}

func (h* Handler) doHandle(w http.ResponseWriter, r *http.Request) (ok bool) {
	for _, mapping := range h.Mappings {
		if r.Method != mapping.Method { continue }
		if res := mapping.Regexp.FindAllStringSubmatch(r.URL.String(), -1); len(res) > 0 {
			mapping.Handler(w, r, res)
			return true
		}
	}
	return false
}

func GenerateUUID() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (h* Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	uuid := GenerateUUID()
	w.Header().Add("X-Request-ID", uuid)
	logger.Info(fmt.Sprintf("%s got request %s %s", uuid, r.Method, r.URL.String()))
	if ok := h.doHandle(w, r); !ok {
		logger.Info("returning 404")
		http.NotFound(w, r)
	}
	logger.Info(fmt.Sprintf("%s finished request in %.06f", uuid, time.Now().Sub(started).Seconds()))
}

func NewHandler(dataDir string) (handler *Handler) {
	handler = &Handler{DataDir: dataDir}

	// dummies
	handler.Map("GET", "_ping",							handler.GetPing)
	handler.Map("GET", "users",							handler.GetUsers)

	// images
	handler.Map("GET", "images/(.*?)/ancestry",			handler.GetImageAncestry)
	handler.Map("GET", "images/(.*?)/layer",			handler.GetImageLayer)
	handler.Map("GET", "images/(.*?)/json",				handler.GetImageJson)
	handler.Map("PUT", "images/(.*?)/(.*)",				handler.PutImageResource)

	// repositories
	handler.Map("GET", "repositories/(.*?)/tags",		handler.GetRepositoryTags)
	handler.Map("GET", "repositories/(.*?)/images",		handler.GetRepositoryImages)
	handler.Map("PUT", "repositories/(.*?)/tags/(.*)",	handler.PutRepositoryTags)
	handler.Map("PUT", "repositories/(.*?)/images",		handler.PutRepositoryImages)
	handler.Map("PUT", "repositories/(.*?)/$",			handler.PutRepository)
	return
}

func startServer(port int, dataDir string) {
	logger.Info("starting server on port ", port)
	logger.Info("using dataDir ", dataDir)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), NewHandler(dataDir)); err != nil {
		logger.Error(err.Error())
	}
}

func main() {
	var port* int
	var dataDir* string
	var doDebug* bool

	port	= flag.Int("p", 80, "Port on which to listen")
	dataDir = flag.String("d", "/data/docker_index", "Directory to store data in")
	doDebug = flag.Bool("D", false, "set log level to debug")
	flag.Parse()

	logger.Level = INFO
	if *doDebug {
		logger.Level = DEBUG
	}
	startServer(*port, *dataDir)
}
