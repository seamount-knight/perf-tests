/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	Local = "local"
	KS3 = "ks3"
)

// ClusterLoaderConfig represents all flags used by CLusterLoader
type ClusterLoaderConfig struct {
	ClusterConfig     ClusterConfig `json: clusterConfig`
	ReportDir         string        `json: reportDir`
	TestConfigPath    string        `json: testConfigPath`
	TestOverridesPath []string      `json: testOverrides`
	TestJob 		  string 		`json: testJob`
	BuildNumber 	  string  		`json: buildNumber`
	ReportType 		  string `json: reportType`
	AccessKeyID string `json: accessKeyID`
	AccessKeySecret string `json: accessKeySecret`
	Bucket string `json: bucket`
	Region string `json: region`
	Config string
}

// ClusterConfig is a structure that represents cluster description.
type ClusterConfig struct {
	KubeConfigPath string `json: kubeConfigPath`
	Nodes          int    `json: nodes`
	Provider       string `json: provider`
	// TODO(krzysied): Add support for HA cluster with more than one master.
	MasterIP   string `json: masterIP`
	MasterName string `json: masterName`
}

func (config *ClusterLoaderConfig) ParseClusterLoaderConfig()  error {
	b, err := ioutil.ReadFile(config.Config)
	if err != nil {
		return  err
	}
	if err := yaml.NewYAMLOrJSONDecoder(bytes.NewBuffer(b), 4096).Decode(&config); err != nil {
		return fmt.Errorf("decoding failed: %v", err)
	}
	return nil
}
