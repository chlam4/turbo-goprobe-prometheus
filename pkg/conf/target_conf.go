package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/golang/glog"
)


// Configuration Parameters to connect to a Prometheus Target
type PrometheusTargetConf struct {
	Address       string
}

func NewPrometheusTargetConf(targetConfigFilePath string) (*PrometheusTargetConf, error) {

	glog.Infof("Read configuration from %s\n", targetConfigFilePath)
	metaConfig := readConfig(targetConfigFilePath)

	return metaConfig, nil
}

// Get the config from file.
func readConfig(path string) *PrometheusTargetConf {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		glog.Infof("File error: %v\n", e)
		os.Exit(1)
	}
	glog.Infoln(string(file))

	var config PrometheusTargetConf
	err := json.Unmarshal(file, &config)

	if err != nil {
		glog.Errorf("Unmarshall error :%v\n", err)
	}
	glog.Infof("Results: %+v\n", config)

	return &config
}
