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
	Set *jet.Set
}

func New(loader jet.Loader, opts ...jet.Option) *jetRenderer {
	return &jetRenderer{jet.NewSet(loader, opts...)}
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

func (r *renderInstance) render(w http.ResponseWriter) error {
	if strings.ContainsRune(r.name, '#') {
		splitted := strings.Split(r.name, "#")
		if len(splitted) != 2 {
			return errors.New("only one '#' is allowed in the template name as a template fragment")
		}

		templatePath := splitted[0]
		blockName := splitted[1]

		expr := fmt.Sprintf(`{{ import "%s" }} {{ yield %s() . }}`, templatePath, blockName)
		t, err := r.Set.Parse(templatePath, expr)
		if err != nil {
			return fmt.Errorf("failed to parse expr '%s': %w", expr, err)
		}

		return t.Execute(w, nil, r.data)
	}

	t, err := r.Set.GetTemplate(r.name)
	if err != nil {
		return err
	}

	return t.Execute(w, nil, r.data)
}

func (r *renderInstance) Render(w http.ResponseWriter) error {
	if err := r.render(w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (*renderInstance) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
