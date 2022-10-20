// Copyright 2022 SAP SE or an SAP affiliate company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// ControllerTestEnv is a convenience interface to be used by tests to access controller-runtime testEnv.
type ControllerTestEnv interface {
	// GetClient provides access to the kubernetes client.Client to access the Kube ApiServer.
	GetClient() client.Client
	// GetConfig provides access to *rest.Config.
	GetConfig() *rest.Config
	// Delete deletes the resources created as part of testEnv.
	Delete()
}

type controllerTestEnv struct {
	client     client.Client
	restConfig *rest.Config
	testEnv    *envtest.Environment
}

// CreateControllerTestEnv creates a controller-runtime testEnv and provides access to the convenience interface to interact with it.
func CreateControllerTestEnv() (ControllerTestEnv, error) {
	testEnv := &envtest.Environment{}
	cfg, err := testEnv.Start()
	if err != nil {
		log.Fatalf("error in starting testEnv: %v", err)
	}
	if cfg == nil {
		log.Fatalf("Got nil config from testEnv")
	}
	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to create new client %w", err)
	}
	return &controllerTestEnv{
		client:     k8sClient,
		restConfig: cfg,
		testEnv:    testEnv,
	}, nil
}

func (te *controllerTestEnv) GetClient() client.Client {
	return te.client
}

func (te *controllerTestEnv) GetConfig() *rest.Config {
	return te.restConfig
}

func (te *controllerTestEnv) Delete() {
	err := te.testEnv.Stop()
	if err != nil {
		log.Printf("failed to cleanly stop controller test environment %v", err)
		return
	}
}
