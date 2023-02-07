// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/project-radius/radius/pkg/armrpc/hostoptions"
	"github.com/project-radius/radius/pkg/corerp/backend"
	"github.com/project-radius/radius/pkg/corerp/frontend"
	metricsservice "github.com/project-radius/radius/pkg/telemetry/metrics/service"
	metricshostoptions "github.com/project-radius/radius/pkg/telemetry/metrics/service/hostoptions"
	"github.com/project-radius/radius/pkg/telemetry/trace"

	link_backend "github.com/project-radius/radius/pkg/linkrp/backend"
	link_frontend "github.com/project-radius/radius/pkg/linkrp/frontend"
	"github.com/project-radius/radius/pkg/logging"
	"github.com/project-radius/radius/pkg/ucp/data"
	"github.com/project-radius/radius/pkg/ucp/dataprovider"
	"github.com/project-radius/radius/pkg/ucp/hosting"
	"github.com/project-radius/radius/pkg/ucp/ucplog"
	etcdclient "go.etcd.io/etcd/client/v3"
)

func newLinkHosts(configFile string, enableAsyncWorker bool) ([]hosting.Service, *hostoptions.HostOptions) {
	hostings := []hosting.Service{}
	options, err := hostoptions.NewHostOptionsFromEnvironment(configFile)
	if err != nil {
		log.Fatal(err)
	}
	hostings = append(hostings, link_frontend.NewService(options))
	if enableAsyncWorker {
		hostings = append(hostings, link_backend.NewService(options))
	}

	return hostings, &options
}

func main() {
	var configFile string
	var enableAsyncWorker bool

	var runLink bool
	var linkConfigFile string

	defaultConfig := fmt.Sprintf("radius-%s.yaml", hostoptions.Environment())
	flag.StringVar(&configFile, "config-file", defaultConfig, "The service configuration file.")
	flag.BoolVar(&enableAsyncWorker, "enable-asyncworker", true, "Flag to run async request process worker (for private preview and dev/test purpose).")

	flag.BoolVar(&runLink, "run-link", true, "Flag to run Applications.Link RP (for private preview and dev/test purpose).")
	defaultLinkConfig := fmt.Sprintf("link-%s.yaml", hostoptions.Environment())
	flag.StringVar(&linkConfigFile, "link-config", defaultLinkConfig, "The service configuration file for Applications.Link.")

	if configFile == "" {
		log.Fatal("config-file is empty.")
	}

	flag.Parse()

	options, err := hostoptions.NewHostOptionsFromEnvironment(configFile)
	if err != nil {
		log.Fatal(err)
	}
	hostingSvc := []hosting.Service{frontend.NewService(options)}

	metricOptions := metricshostoptions.NewHostOptionsFromEnvironment(*options.Config)
	if metricOptions.Config.Prometheus.Enabled {
		hostingSvc = append(hostingSvc, metricsservice.NewService(metricOptions))
	}

	logger, flush, err := ucplog.NewLogger(logging.AppCoreLoggerName, &options.Config.Logging)
	if err != nil {
		log.Fatal(err)
	}
	defer flush()

	if enableAsyncWorker {
		logger.Info("Enable AsyncRequestProcessWorker.")
		hostingSvc = append(hostingSvc, backend.NewService(options))
	}

	// Configure Applications.Link to run it with Applications.Core RP.
	var linkOpts *hostoptions.HostOptions
	if runLink && linkConfigFile != "" {
		logger.Info("Run Applications.Link.")
		var linkSvcs []hosting.Service
		linkSvcs, linkOpts = newLinkHosts(linkConfigFile, enableAsyncWorker)
		hostingSvc = append(hostingSvc, linkSvcs...)
	}

	if options.Config.StorageProvider.Provider == dataprovider.TypeETCD &&
		options.Config.StorageProvider.ETCD.InMemory {
		// For in-memory etcd we need to register another service to manage its lifecycle.
		//
		// The client will be initialized asynchronously.
		logger.Info("Enabled in-memory etcd")
		client := hosting.NewAsyncValue[etcdclient.Client]()
		options.Config.StorageProvider.ETCD.Client = client
		options.Config.SecretProvider.ETCD.Client = client
		if linkOpts != nil {
			linkOpts.Config.StorageProvider.ETCD.Client = client
			linkOpts.Config.SecretProvider.ETCD.Client = client
		}
		hostingSvc = append(hostingSvc, data.NewEmbeddedETCDService(data.EmbeddedETCDServiceOptions{ClientConfigSink: client}))
	}

	loggerValues := []any{}
	host := &hosting.Host{
		Services: hostingSvc,

		// Values that will be propagated to all loggers
		LoggerValues: loggerValues,
	}

	ctx, cancel := context.WithCancel(logr.NewContext(context.Background(), logger))

	url := "http://localhost:9411/api/v2/spans"
	shutdown, err := trace.InitServerTracer(url, "appcore-rp")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}

	}()

	stopped, serviceErrors := host.RunAsync(ctx)

	exitCh := make(chan os.Signal, 2)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)

	select {
	// Shutdown triggered
	case <-exitCh:
		logger.Info("Shutting down....")
		cancel()

	// A service terminated with a failure. Shut down
	case <-serviceErrors:
		logger.Info("Error occurred - shutting down....")
		cancel()
	}

	// Finished shutting down. An error returned here is a failure to terminate
	// gracefully, so just crash if that happens.
	err = <-stopped
	if err == nil {
		os.Exit(0)
	} else {
		panic(err)
	}
}
