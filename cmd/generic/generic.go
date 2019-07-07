package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/q3k/scarab"
	cpb "github.com/q3k/scarab/proto/common"
	gpb "github.com/q3k/scarab/proto/generic"
)

type flags struct {
	configuration string
	databasePath  string
	httpBind      string
}

func init() {
	flag.Set("logtostderr", "true")
}

func ValidateProtoJob(jd *cpb.JobDefinition) error {
	if jd.Name == "" {
		return fmt.Errorf("name must be set")
	}
	for i, arg := range jd.Arguments {
		if err := ValidateProtoArgument(arg); err != nil {
			return fmt.Errorf("argument %d: %v", i, err)
		}
	}
	for i, step := range jd.Steps {
		if err := ValidateProtoStep(step); err != nil {
			return fmt.Errorf("step %d: %v", i, err)
		}
	}
	return nil
}

func ValidateProtoArgument(a *cpb.ArgumentDefinition) error {
	// TODO(q3k): Implement
	return nil
}

func ValidateProtoStep(sd *cpb.StepDefinition) error {
	if sd.Name == "" {
		return fmt.Errorf("name must be set")
	}
	return nil
}

func main() {
	f := flags{
		configuration: "configuration.proto.text",
		httpBind:      "127.0.0.1:2137",
		databasePath:  "scarab.db",
	}

	flag.StringVar(&f.configuration, "configuration", f.configuration, "Location of Scarab instance configuration. If ends in .text, will be unmarshaled as protobuf text, otherwise as binary protobuf.")
	flag.StringVar(&f.httpBind, "http_bind", f.httpBind, "Address to bind HTTP server to. If empty, no HTTP server will be started.")
	flag.StringVar(&f.databasePath, "db_path", f.databasePath, "Local path of storage database.")
	flag.Parse()

	config := gpb.Configuration{}

	data, err := ioutil.ReadFile(f.configuration)
	if err != nil {
		glog.Exitf("Could not open config file %q: %v", f.configuration, err)
	}

	if strings.HasSuffix(f.configuration, ".text") {
		if err := proto.UnmarshalText(string(data), &config); err != nil {
			glog.Exitf("Could not parse text config: %v", err)
		}
	} else {
		if err := proto.Unmarshal(data, &config); err != nil {
			glog.Exitf("Could not parse config: %v", err)
		}
	}

	definitions := make([]*scarab.JobDefinition, len(config.Job))
	for i, job := range config.Job {
		if err := ValidateProtoJob(job); err != nil {
			glog.Exitf("Configuration validation failed: job %d: %v", i, err)
		}

		definitions[i] = &scarab.JobDefinition{
			Name:        job.Name,
			Description: job.Description,
			Arguments:   job.Arguments,
		}
	}

	storage, err := scarab.NewLevelDBStorage(f.databasePath)
	if err != nil {
		glog.Exitf("Storage failed: %v", err)
	}

	s, err := scarab.NewService(definitions, storage)
	if err != nil {
		glog.Exitf("Creating service failed failed: %v", err)
	}

	glog.Infof("Service loaded!")

	ctx := context.Background()
	ctxC, cancel := context.WithCancel(ctx)

	if f.httpBind != "" {
		glog.Infof("Starting HTTP at %s...", f.httpBind)
		go func() {
			err := s.RunHTTPServer(ctxC, f.httpBind)
			if err != ctxC.Err() {
				glog.Exitf("Could not run HTTP server: %v", err)
			}
		}()
	} else {
		glog.Infof("Not starting HTTP")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancel()
	time.Sleep(100 * time.Millisecond)
}
