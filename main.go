package scarab

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	spb "github.com/q3k/scarab/proto/state"
	"github.com/q3k/scarab/templates"
)

type JobDefinition struct {
	Name        string
	Description string

	ArgsDescriptor []*spb.ArgumentDefinition

	Steps []StepDefinition
}

type StepDefinition struct {
	Name        string
	Description string
	DependsStep string
}

type RunningJob struct {
	definition *JobDefinition

	Args  proto.Message
	State proto.Message
}

type RunningStep struct {
	definition *StepDefinition
	job        *RunningJob
}

type Service struct {
	Definitions map[string]*JobDefinition
	Jobs        []*RunningJob
}

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

var (
//templateRoot = loadTemplate("root.html", "base.html")
)

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
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "job" || parts[1] != "type" {
		http.NotFound(w, r)
		return
	}

	def, ok := s.Definitions[parts[2]]
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
		RenderSelectedJobType: parts[2],
		RenderJobs:            []*RunningJob{},
	}

	for _, rj := range s.Jobs {
		if rj.definition.Name != parts[2] {
			continue
		}
		data.RenderJobs = append(data.RenderJobs, rj)
	}

	renderTemplate(w, templateJob, data)
}

func (s *Service) httpJsonJobDefinition(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[0] != "json" || parts[1] != "job" || parts[2] != "definition" {
		http.NotFound(w, r)
		return
	}

	def, ok := s.Definitions[parts[3]]
	if !ok {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(def)
}
