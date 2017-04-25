package main

import (
	"flag"
	"os"
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
	"github.com/chlam4/turbo-goprobe-prometheus/pkg/registration"
	"github.com/chlam4/turbo-goprobe-prometheus/pkg/discovery"
)

func main() {
	flag.Parse()

	targetConf := "target-conf.json"
	turboCommConf := "turbo-server-conf.json"

	communicator, err := service.ParseTurboCommunicationConfig(turboCommConf)
	if err != nil {
		glog.Infof("Error while parsing the turbo communicator config file %v: %v\n", turboCommConf, err)
		os.Exit(1)
	}

	// Example Probe Registration Client
	registrationClient := &registration.PrometheusRegistrationClient{}
	// Example Probe Registration Client
	discoveryClient, err := discovery.NewDiscoveryClient(targetConf)
	if err != nil {
		glog.Infof("Error while instantiating a discovery client at %v with config %v: %v\n", turboCommConf, targetConf, err)
		os.Exit(1)
	}

	tapService, err := service.NewTAPServiceBuilder().
		WithTurboCommunicator(communicator).
		WithTurboProbe(probe.NewProbeBuilder(registration.TargetType, registration.ProbeCategory).
			RegisteredBy(registrationClient).
			DiscoversTarget(discoveryClient.ClientConf.Address, discoveryClient)).Create()

	if err != nil {
		glog.Infof("Error while building turbo tap service on target %v: %v\n", discoveryClient.ClientConf.Address, err)
		os.Exit(1)
	}

	// Connect to the Turbo server
	tapService.ConnectToTurbo()

	select {}
}
