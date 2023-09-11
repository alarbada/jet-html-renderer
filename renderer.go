package jetr

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin/render"
)

type jetRenderer struct {
	set *jet.Set
}

var _ render.HTMLRender = &jetRenderer{}

func (r *jetRenderer) Instance(name string, data any) render.Render {
	return &renderInstance{
		jetRenderer: r,
		name:        name,
		data:        data,
	}
}

type renderInstance struct {
	*jetRenderer
	name string
	data any
}

var _ render.Render = &renderInstance{}

func (r *renderInstance) Render(w http.ResponseWriter) error {
	if strings.ContainsRune(r.name, '#') {
		splitted := strings.Split(r.name, "#")
		if len(splitted) != 2 {
			return errors.New("only one # is allowed in the template name as a template fragment")
		}

		expr := fmt.Sprintf(`{{ import "%s" }} {{ block %s() . }} {{ end }}`, splitted[0], splitted[1])
		t, err := r.set.Parse(r.name, expr)
		if err != nil {
			return err
		}

		return t.Execute(w, nil, r.data)
	}

	t, err := r.set.GetTemplate(r.name)
	if err != nil {
		return err
	}

	return t.Execute(w, nil, r.data)
}

func (*renderInstance) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
