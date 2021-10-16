package html

import (
	"embed"
	"github.com/torbenschinke/picd/internal/dashboard"
	"github.com/torbenschinke/picd/pkg/logging"
	"github.com/torbenschinke/picd/pkg/server"
	"html/template"
	"net/http"
)

//go:embed *.gohtml
var templates embed.FS

type Handler struct {
	service *dashboard.Service
	tpls    *template.Template
}

func NewHandler(service *dashboard.Service) *Handler {
	tpls, err := template.ParseFS(templates, "*")
	if err != nil {
		panic(err) //must parse
	}

	return &Handler{service: service, tpls: tpls}
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	cams,err := h.service.ListCameras()
	if err != nil{
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.tpls.ExecuteTemplate(w, "dashboard", cams); err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadCam(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	params := server.ParamsFromContext(r.Context())
	id := params.ByName("id")
	if err := h.service.LoadCameraImage(id, w); err != nil {
		logger.Println(err)
	}

}

func (h *Handler) Configure(r server.Router) {
	r.HandlerFunc(http.MethodGet, "/", h.index)
	r.HandlerFunc(http.MethodGet, "/camera/:id/photo", h.loadCam)
}
