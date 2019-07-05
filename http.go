package scarab

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/golang/glog"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	"github.com/q3k/scarab/js"
	cpb "github.com/q3k/scarab/proto/common"
	"github.com/q3k/scarab/templates"
)

type grpcManage struct {
	s *Service
}

func (s *grpcManage) Definitions(ctx context.Context, req *cpb.DefinitionsRequest) (*cpb.DefinitionsResponse, error) {
	res := &cpb.DefinitionsResponse{
		Jobs: make([]*cpb.JobDefinition, len(s.s.Definitions)),
	}
	i := 0
	for _, job := range s.s.Definitions {
		res.Jobs[i] = &cpb.JobDefinition{
			Name:        job.Name,
			Description: job.Description,
			Arguments:   job.ArgsDescriptor,
			Steps:       make([]*cpb.StepDefinition, len(job.Steps)),
		}
		for j, step := range job.Steps {
			res.Jobs[i].Steps[j] = &cpb.StepDefinition{
				Name:        step.Name,
				Description: step.Description,
			}
		}
		i += 1
	}
	return res, nil
}

func (s *Service) RunHTTPServer(ctx context.Context, bind string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.httpRoot)
	mux.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		ruri := strings.TrimPrefix(r.RequestURI, "/js/")
		data, ok := js.Data[ruri]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/javascript")
		w.Write(data)
	})

	grpcServer := grpc.NewServer()
	manage := &grpcManage{s}
	cpb.RegisterManageServer(grpcServer, manage)
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	srv := http.Server{
		Addr: bind,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if wrappedGrpc.IsGrpcWebRequest(r) {
				wrappedGrpc.ServeHTTP(w, r)
				return
			}
			mux.ServeHTTP(w, r)
		}),
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
		asset, ok := templates.Data["templates/"+n]
		if !ok {
			panic(fmt.Sprintf("unknown template %q", n))
		}

		if t == nil {
			t = template.New(n)
		} else {
			t = t.New(n)
		}

		_, err := t.Parse(string(asset))
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
