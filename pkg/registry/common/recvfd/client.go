// Copyright (c) 2021 Cisco and/or its affiliates.
//
// Copyright (c) 2021 Doc.ai and/or its affiliates.
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

// +build linux

package recvfd

import (
	"context"
	"net/url"
	"os"

	"github.com/edwarnicke/grpcfd"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/registry"
	"google.golang.org/grpc"

	"github.com/networkservicemesh/sdk/pkg/registry/core/next"
)

type recvfdNSEClient struct {
	fileMaps perEndpointFileMapMap
}

func (n *recvfdNSEClient) Register(ctx context.Context, in *registry.NetworkServiceEndpoint, opts ...grpc.CallOption) (*registry.NetworkServiceEndpoint, error) {
	return next.NetworkServiceEndpointRegistryClient(ctx).Register(ctx, in, opts...)
}

func (n *recvfdNSEClient) Find(ctx context.Context, in *registry.NetworkServiceEndpointQuery, opts ...grpc.CallOption) (registry.NetworkServiceEndpointRegistry_FindClient, error) {
	rpcCredentials := grpcfd.PerRPCCredentials(grpcfd.PerRPCCredentialsFromCallOptions(opts...))
	opts = append(opts, grpc.PerRPCCredentials(rpcCredentials))
	recv, _ := grpcfd.FromPerRPCCredentials(rpcCredentials)
	resp, err := next.NetworkServiceEndpointRegistryClient(ctx).Find(ctx, in, opts...)
	if err != nil {
		return nil, err
	}
	return &recvfdNSEFindClient{
		transceiver: recv,
		NetworkServiceEndpointRegistry_FindClient: resp,
		fileMaps: &n.fileMaps,
	}, nil
}

func (n *recvfdNSEClient) Unregister(ctx context.Context, in *registry.NetworkServiceEndpoint, opts ...grpc.CallOption) (*empty.Empty, error) {
	return next.NetworkServiceEndpointRegistryClient(ctx).Unregister(ctx, in, opts...)
}

// NewNetworkServiceEndpointRegistryClient - returns a new null client that does nothing but call next.NetworkServiceEndpointRegistryClient(ctx).
func NewNetworkServiceEndpointRegistryClient() registry.NetworkServiceEndpointRegistryClient {
	return new(recvfdNSEClient)
}

type recvfdNSEFindClient struct {
	registry.NetworkServiceEndpointRegistry_FindClient
	transceiver grpcfd.FDTransceiver
	fileMaps    *perEndpointFileMapMap
}

func (x *recvfdNSEFindClient) Recv() (*registry.NetworkServiceEndpointResponse, error) {
	nseResp, err := x.NetworkServiceEndpointRegistry_FindClient.Recv()
	if err != nil {
		return nil, err
	}
	if x.transceiver != nil {
		// Get the fileMap
		fileMap := &perEndpointFileMap{
			filesByInodeURL:    make(map[string]*os.File),
			inodeURLbyFilename: make(map[string]*url.URL),
		}
		endpointName := nseResp.GetNetworkServiceEndpoint().GetName()
		// If name is specified, let's use it, since it could be heal/update request
		if endpointName != "" {
			fileMap, _ = x.fileMaps.LoadOrStore(nseResp.NetworkServiceEndpoint.GetName(), fileMap)
		}

		// Recv the FD and swap theInode to File in the Parameters for the returned connection mechanism
		err = recvFDAndSwapInodeToUnix(x.Context(), fileMap, nseResp.GetNetworkServiceEndpoint(), x.transceiver)
	}
	return nseResp, err
}
