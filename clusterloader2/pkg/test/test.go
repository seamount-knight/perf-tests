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

package test

import (
	"fmt"
	"path/filepath"

	"k8s.io/perf-tests/clusterloader2/pkg/config"
	"k8s.io/perf-tests/clusterloader2/pkg/errors"
	"k8s.io/perf-tests/clusterloader2/pkg/framework"
	"k8s.io/perf-tests/clusterloader2/pkg/state"
)

var (
	// CreateContext global function for creating context.
	// This function should be set by Context implementation.
	CreateContext = createSimpleContext

	// Test is a singleton for test execution object.
	// This object should be set by TestExecutor implementation.
	Test = createSimpleTestExecutor()
)

// RunTest runs test based on provided test configuration.
func RunTest(f *framework.Framework, clusterLoaderConfig *config.ClusterLoaderConfig) *errors.ErrorList {
	if f == nil {
		return errors.NewErrorList(fmt.Errorf("framework must be provided"))
	}
	if clusterLoaderConfig == nil {
		return errors.NewErrorList(fmt.Errorf("cluster loader config must be provided"))
	}
	if CreateContext == nil {
		return errors.NewErrorList(fmt.Errorf("no CreateContext function installed"))
	}
	if Test == nil {
		return errors.NewErrorList(fmt.Errorf("no Test installed"))
	}

	ctx := CreateContext(clusterLoaderConfig, f, state.NewState())
	testConfigFilename := filepath.Base(clusterLoaderConfig.TestConfigPath)

	mapping, err := config.GetOverridesMapping(clusterLoaderConfig.TestOverridesPath)
	if err != nil {
		return errors.NewErrorList(fmt.Errorf("mapping creation error: %v", err))
	}
	mapping["Nodes"] = clusterLoaderConfig.ClusterConfig.Nodes
	testConfig, err := ctx.GetTemplateProvider().TemplateToConfig(testConfigFilename, mapping)
	fmt.Printf("step len: %d\n\n", len(testConfig.Steps))
	//for index, step := range testConfig.Steps {
	//	fmt.Printf("index: %d, name: %v\n", index, step.Name)
	//	fmt.Printf("measu: %d\n", len(step.Measurements))
	//	for _, m := range step.Measurements {
	//		fmt.Println(m.Identifier, m.Method, m.Params)
	//	}
	//	fmt.Printf("phases: %d\n", len(step.Phases))
	//	for _, p := range step.Phases {
	//		fmt.Println(p.TuningSet, p.NamespaceRange, p.ObjectBundle)
	//	}
	//	fmt.Println("------------------------------------------\n")
	//}
	//return nil
	if err != nil {
		return errors.NewErrorList(fmt.Errorf("config reading error: %v", err))
	}
	return Test.ExecuteTest(ctx, testConfig)
}
