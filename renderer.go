package jetr

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin/render"
)

type jetRenderer struct {
	*jet.Set
}

var _ render.HTMLRender = &jetRenderer{}

func (r *jetRenderer) Instance(name string, data interface{}) render.Render {
	return &jetRender{}
}

type jetRender struct{}

var _ render.Render = &jetRender{}

func (r *jetRender) Render(wr http.ResponseWriter) error {
	panic("unimplemented")
}

func (*jetRender) WriteContentType(w http.ResponseWriter) {
	panic("unimplemented")
}
