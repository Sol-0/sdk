// Copyright (c) 2020 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package setid_test

import (
	"context"
	"testing"

	"github.com/networkservicemesh/api/pkg/api/registry"
	"github.com/stretchr/testify/require"
)

type emptyBulkRegisterNSEServer struct {
	registry.NetworkServiceRegistry_BulkRegisterNSEServer
}

func (s *emptyBulkRegisterNSEServer) Send(*registry.NSERegistration) error {
	return nil
}

func (s *emptyBulkRegisterNSEServer) Context() context.Context {
	return context.Background()
}

type assertServer struct {
	*testing.T
	registry.NetworkServiceRegistryServer
}

func (n *assertServer) BulkRegisterNSE(s registry.NetworkServiceRegistry_BulkRegisterNSEServer) error {
	nse := &registry.NetworkServiceEndpoint{
		NetworkServiceName: "ns-1",
		Payload:            "IP",
	}
	registration := &registry.NSERegistration{
		NetworkService: &registry.NetworkService{
			Name:    "ns-1",
			Payload: "IP",
		},
		NetworkServiceEndpoint: nse,
	}
	require.Nil(n, s.Send(registration))
	require.NotEmpty(n, registration.NetworkServiceEndpoint.Name)
	return nil
}