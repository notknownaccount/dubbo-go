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

package cluster_impl

import (
	"context"
	"fmt"
	"testing"
)

import (
	"github.com/golang/mock/gomock"

	perrors "github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

import (
	"dubbo.apache.org/dubbo-go/v3/cluster/directory"
	"dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	"dubbo.apache.org/dubbo-go/v3/protocol/mock"
)

var failsafeUrl, _ = common.NewURL(
	fmt.Sprintf("dubbo://%s:%d/com.ikurento.user.UserProvider", constant.LOCAL_HOST_VALUE, constant.DEFAULT_PORT))

// registerFailsafe register failsafeCluster to cluster extension.
func registerFailsafe(invoker *mock.MockInvoker) protocol.Invoker {
	extension.SetLoadbalance("random", loadbalance.NewRandomLoadBalance)
	failsafeCluster := NewFailsafeCluster()

	invokers := []protocol.Invoker{}
	invokers = append(invokers, invoker)
	invoker.EXPECT().IsAvailable().Return(true).AnyTimes()

	invoker.EXPECT().GetUrl().Return(failbackUrl)

	staticDir := directory.NewStaticDirectory(invokers)
	clusterInvoker := failsafeCluster.Join(staticDir)
	return clusterInvoker
}

func TestFailSafeInvokeSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invoker := mock.NewMockInvoker(ctrl)
	clusterInvoker := registerFailsafe(invoker)

	invoker.EXPECT().IsAvailable().Return(true).AnyTimes()

	invoker.EXPECT().GetUrl().Return(failsafeUrl).AnyTimes()

	mockResult := &protocol.RPCResult{Rest: rest{tried: 0, success: true}}

	invoker.EXPECT().Invoke(gomock.Any()).Return(mockResult)
	result := clusterInvoker.Invoke(context.Background(), &invocation.RPCInvocation{})

	assert.NoError(t, result.Error())
	res := result.Result().(rest)
	assert.True(t, res.success)
}

func TestFailSafeInvokeFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invoker := mock.NewMockInvoker(ctrl)
	clusterInvoker := registerFailsafe(invoker)
	invoker.EXPECT().IsAvailable().Return(true).AnyTimes()

	invoker.EXPECT().GetUrl().Return(failsafeUrl).AnyTimes()

	mockResult := &protocol.RPCResult{Err: perrors.New("error")}

	invoker.EXPECT().Invoke(gomock.Any()).Return(mockResult)
	result := clusterInvoker.Invoke(context.Background(), &invocation.RPCInvocation{})

	assert.NoError(t, result.Error())
	assert.Nil(t, result.Result())
}
