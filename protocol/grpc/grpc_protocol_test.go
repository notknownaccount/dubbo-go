/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grpc

import (
	"dubbo.apache.org/dubbo-go/v3/config/application"
	"dubbo.apache.org/dubbo-go/v3/config/base"
	"dubbo.apache.org/dubbo-go/v3/config/instance"
	"dubbo.apache.org/dubbo-go/v3/config/method"
	"dubbo.apache.org/dubbo-go/v3/config/provider"
	"dubbo.apache.org/dubbo-go/v3/config/service"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/protocol/grpc/internal/helloworld"
)

func doInitProvider() {
	providerConfig := provider.Config{
		BaseConfig: base.Config{
			ApplicationConfig: &application.Config{
				Organization: "dubbo_org",
				Name:         "BDTService",
				Module:       "module",
				Version:      "0.0.1",
				Owner:        "dubbo",
				Environment:  "test",
			},
		},
		Services: map[string]*service.Config{
			"GrpcGreeterImpl": {
				InterfaceName: "io.grpc.examples.helloworld.GreeterGrpc$IGreeter",
				Protocol:      "grpc",
				Registry:      "shanghai_reg1,shanghai_reg2,hangzhou_reg1,hangzhou_reg2,hangzhou_service_discovery_reg",
				Cluster:       "failover",
				Loadbalance:   "random",
				Retries:       "3",
				Methods: []*method.Config{
					{
						Name:        "SayHello",
						Retries:     "2",
						LoadBalance: "random",
						Weight:      200,
					},
				},
			},
		},
	}
	provider.SetProviderConfig(providerConfig)
}

func TestGrpcProtocolExport(t *testing.T) {
	// Export
	instance.SetProviderService(helloworld.NewService())
	doInitProvider()

	url, err := common.NewURL(helloworldURL)
	assert.NoError(t, err)

	proto := GetProtocol()
	exporter := proto.Export(protocol.NewBaseInvoker(url))
	time.Sleep(time.Second)

	// make sure url
	eq := exporter.GetInvoker().GetURL().URLEqual(url)
	assert.True(t, eq)

	// make sure exporterMap after 'Unexport'
	_, ok := proto.(*GrpcProtocol).ExporterMap().Load(url.ServiceKey())
	assert.True(t, ok)
	exporter.Unexport()
	_, ok = proto.(*GrpcProtocol).ExporterMap().Load(url.ServiceKey())
	assert.False(t, ok)

	// make sure serverMap after 'Destroy'
	_, ok = proto.(*GrpcProtocol).serverMap[url.Location]
	assert.True(t, ok)
	proto.Destroy()
	_, ok = proto.(*GrpcProtocol).serverMap[url.Location]
	assert.False(t, ok)
}

func TestGrpcProtocolRefer(t *testing.T) {
	server, err := helloworld.NewServer("127.0.0.1:30000")
	assert.NoError(t, err)
	go server.Start()
	defer server.Stop()

	url, err := common.NewURL(helloworldURL)
	assert.NoError(t, err)

	proto := GetProtocol()
	invoker := proto.Refer(url)

	// make sure url
	eq := invoker.GetURL().URLEqual(url)
	assert.True(t, eq)

	// make sure invokers after 'Destroy'
	invokersLen := len(proto.(*GrpcProtocol).Invokers())
	assert.Equal(t, 1, invokersLen)
	proto.Destroy()
	invokersLen = len(proto.(*GrpcProtocol).Invokers())
	assert.Equal(t, 0, invokersLen)
}
