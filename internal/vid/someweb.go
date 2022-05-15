package vid

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/torbenschinke/picd/internal/pic"
	"github.com/torbenschinke/picd/internal/sensor"
	"github.com/torbenschinke/picd/pkg/server"
	"html/template"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//go:embed index.gohtml
var pageIndex string

//go:embed cluster.gohtml
var pageCluster string

type idxModel struct {
	ImageCount   int
	LastImageMod string
	GBInUse      string
	Cluster      []idxCluster
	Degree       string
	Hum          string
}

type idxCluster struct {
	First      os.FileInfo
	Last       os.FileInfo
	Files      []os.DirEntry
	ImageCount int
}

type clusterDetails struct {
	Cluster   idxCluster
	JsonNames template.JS
}

func (c idxCluster) From() string {
	return c.First.ModTime().In(globalLoc).Format("Mon, 02 Jan 2006 15:04:05")
}

func (c idxCluster) To() string {
	return c.Last.ModTime().In(globalLoc).Format("Mon, 02 Jan 2006 15:04:05")
}

type Controller struct {
	sensors            *sensor.SenseService
	camera             *pic.CameraService
	app                *Application
	idxTpl, detailsTpl *template.Template
	loc                *time.Location
}

var globalLoc *time.Location

func NewController(app *Application, sensors *sensor.SenseService, camera *pic.CameraService) *Controller {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	globalLoc = loc
	idx := template.Must(template.New("index.gohtml").Parse(pageIndex))
	details := template.Must(template.New("cluster.gohtml").Parse(pageCluster))
	return &Controller{app: app, sensors: sensors, camera: camera, idxTpl: idx, loc: loc, detailsTpl: details}
}

func (c *Controller) index(w http.ResponseWriter, r *http.Request) {
	m := c.model()
	if err := c.idxTpl.Execute(w, m); err != nil {
		fmt.Println(err)
	}
}

func (c *Controller) cluster(w http.ResponseWriter, r *http.Request) {
	params := server.ParamsFromContext(r.Context())
	sidx := params.ByName("idx")
	idx, _ := strconv.Atoi(sidx)
	m := c.model()
	if !(idx >= 0 && idx < len(m.Cluster)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid cluster index"))
		return
	}

	details := clusterDetails{Cluster: m.Cluster[idx]}
	var names []string
	for _, file := range details.Cluster.Files {
		names = append(names, file.Name())
	}

	buf, _ := json.Marshal(names)
	details.JsonNames = template.JS(strconv.Quote(string(buf)))
	if err := c.detailsTpl.Execute(w, details); err != nil {
		fmt.Println(err)
	}
}

func (c *Controller) model() idxModel {
	var m idxModel
	files, _ := os.ReadDir(c.app.timelapseDir)
	m.ImageCount = len(files)
	m.Degree = fmt.Sprintf("%d°C", int(math.Round(c.sensors.T().Celsius())))
	m.Hum = fmt.Sprintf("%d %%", int(math.Round(c.sensors.RH().Humidity())))
	if len(files) > 0 {
		if info, err := files[len(files)-1].Info(); err == nil {
			m.LastImageMod = info.ModTime().In(c.loc).Format("Mon, 02 Jan 2006 15:04:05")
		}

	}

	b := int64(0)
	var cluster idxCluster
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		cluster.Files = append(cluster.Files, file)

		b += info.Size()
		cluster.ImageCount++
		if cluster.First == nil {
			cluster.First = info
		}

		if cluster.ImageCount > 60*60 {
			cluster.Last = info
			m.Cluster = append(m.Cluster, cluster)
			cluster = idxCluster{}
		}

	}

	if cluster.Last == nil && len(files) > 0 {
		info, err := files[len(files)-1].Info()
		if err != nil {
			fmt.Println(err)
		}
		cluster.Last = info
		m.Cluster = append(m.Cluster, cluster)
	}

	m.GBInUse = fmt.Sprintf("%.2f GiB", float64(b)/1024/1024/1024)

	return m
}

func (c *Controller) current(w http.ResponseWriter, r *http.Request) {
	fname := filepath.Join(c.app.timelapseDir, "current.jpg")
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer file.Close()

	io.Copy(w, file)
}

func (c *Controller) image(w http.ResponseWriter, r *http.Request) {
	params := server.ParamsFromContext(r.Context())
	name := params.ByName("name")
	if strings.Contains(name, "..") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fname := filepath.Join(c.app.timelapseDir, name)
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer file.Close()

	io.Copy(w, file)
}

func (c *Controller) Configure(r server.Router) {
	r.HandlerFunc(http.MethodGet, "/", c.index)
	r.HandlerFunc(http.MethodGet, "/index.html", c.index)
	r.HandlerFunc(http.MethodGet, "/current.jpg", c.current)
	r.HandlerFunc(http.MethodGet, "/cluster/:idx", c.cluster)
	r.HandlerFunc(http.MethodGet, "/image/:name", c.image)
}
