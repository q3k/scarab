package scarab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/golang/glog"
	"github.com/q3k/scarab/templates"
)

func (s *Service) RunHTTPServer(ctx context.Context, bind string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.httpRoot)
	mux.HandleFunc("/job/type/", s.httpJobType)
	mux.HandleFunc("/json/job/definition/", s.httpJsonJobDefinition)

	srv := http.Server{
		Addr:    bind,
		Handler: mux,
	}

	lisErr := make(chan error)

	go func() {
		err := srv.ListenAndServe()
		lisErr <- err
	}()

	var err error
	select {
	case <-ctx.Done():
		glog.Infof("Stopping HTTP...")
		srv.Close()
		err = <-lisErr
	case err = <-lisErr:
	}
	if err != http.ErrServerClosed {
		return err
	}
	return ctx.Err()
}
func loadTemplate(names ...string) *template.Template {
	var t *template.Template
	for _, n := range names {
		asset, err := templates.Asset(n)
		if err != nil {
			panic(fmt.Sprintf("unknown template %q", n))
		}

		if t == nil {
			t = template.New(n)
		} else {
			t = t.New(n)
		}

		_, err = t.Parse(string(asset))
		if err != nil {
			panic(fmt.Sprintf("template %q parse failed: %v", n, err))
		}
	}
	return t
}

func renderTemplate(w http.ResponseWriter, t *template.Template, data interface{}) {
	err := t.Execute(w, data)
	if err == nil {
		return
	}

	glog.Errorf("Error executing template %q: %v", t.Name(), err)
}

func getArg(r *http.Request, root ...string) string {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != len(root)+1 {
		return ""
	}

	for i, v := range root {
		if v != parts[i] {
			return ""
		}
	}

	return parts[len(root)]
}

func (s *Service) httpRoot(w http.ResponseWriter, r *http.Request) {
	templateRoot := loadTemplate("root.html", "base.html")
	data := struct {
		renderData
		RenderSelectedJobType string
	}{
		renderData:            renderData{Service: s, RenderSubtitle: "Home"},
		RenderSelectedJobType: "",
	}
	renderTemplate(w, templateRoot, data)
}

type renderData struct {
	*Service
	RenderSubtitle string
}

func (s *Service) httpJobType(w http.ResponseWriter, r *http.Request) {
	t := getArg(r, "job", "type")
	if t == "" {
		http.NotFound(w, r)
		return
	}
	def, ok := s.Definitions[t]
	if !ok {
		http.NotFound(w, r)
		return
	}

	templateJob := loadTemplate("job.html", "base.html")
	data := struct {
		renderData
		RenderSelectedJobType string
		RenderJobs            []*RunningJob
	}{
		renderData: renderData{
			Service:        s,
			RenderSubtitle: def.Description,
		},
		RenderSelectedJobType: t,
		RenderJobs:            []*RunningJob{},
	}

	for _, rj := range s.Jobs {
		if rj.definition.Name != t {
			continue
		}
		data.RenderJobs = append(data.RenderJobs, rj)
	}

	renderTemplate(w, templateJob, data)
}

func (s *Service) httpJsonJobDefinition(w http.ResponseWriter, r *http.Request) {
	t := getArg(r, "json", "job", "definition")
	if t == "" {
		http.NotFound(w, r)
		return
	}

	def, ok := s.Definitions[t]
	if !ok {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(def)
}
