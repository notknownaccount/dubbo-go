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

package config

import (
	"net/url"
	"time"
)

import (
	perrors "github.com/pkg/errors"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
)

// RemoteConfig: usually we need some middleware, including nacos, zookeeper
// this represents an instance of this middleware
// so that other module, like config center, registry could reuse the config
// but now, only metadata report, metadata service, service discovery use this structure
type RemoteConfig struct {
	Protocol   string            `yaml:"protocol"  json:"protocol,omitempty" property:"protocol"`
	Address    string            `yaml:"address" json:"address,omitempty" property:"address"`
	TimeoutStr string            `default:"5s" yaml:"timeout" json:"timeout,omitempty" property:"timeout"`
	Username   string            `yaml:"username" json:"username,omitempty" property:"username"`
	Password   string            `yaml:"password" json:"password,omitempty"  property:"password"`
	Params     map[string]string `yaml:"params" json:"params,omitempty"`
}

// Prefix dubbo.remote.
func (rc *RemoteConfig) Prefix() string {
	return constant.RemotePrefix
}

// Timeout return timeout duration.
// if the configure is invalid, or missing, the default value 5s will be returned
func (rc *RemoteConfig) Timeout() time.Duration {
	if res, err := time.ParseDuration(rc.TimeoutStr); err == nil {
		return res
	}
	return 5 * time.Second
}

// GetParam will return the value of the key. If not found, def will be return;
// def => default value
func (rc *RemoteConfig) GetParam(key string, def string) string {
	param, ok := rc.Params[key]
	if !ok {
		return def
	}
	return param
}

// ToURL config to url
func (rc *RemoteConfig) ToURL() (*common.URL, error) {
	if len(rc.Protocol) == 0 {
		return nil, perrors.Errorf("Must provide protocol in RemoteConfig.")
	}
	return common.NewURL(rc.Address,
		common.WithProtocol(rc.Protocol),
		common.WithUsername(rc.Username),
		common.WithPassword(rc.Password),
		common.WithLocation(rc.Address),
		common.WithParams(rc.getUrlMap()),
	)
}

// getUrlMap get url map
func (rc *RemoteConfig) getUrlMap() url.Values {
	urlMap := url.Values{}
	urlMap.Set(constant.CONFIG_USERNAME_KEY, rc.Username)
	urlMap.Set(constant.CONFIG_PASSWORD_KEY, rc.Password)
	urlMap.Set(constant.CONFIG_TIMEOUT_KEY, rc.TimeoutStr)

	for key, val := range rc.Params {
		urlMap.Set(key, val)
	}
	return urlMap
}
