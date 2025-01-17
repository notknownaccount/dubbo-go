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

package impl

import (
	"fmt"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
)

var (
	serializers = make(map[string]Serializer)
	nameMaps    = make(map[byte]string)
)

func init() {
	nameMaps = map[byte]string{
		constant.S_Hessian2: constant.HESSIAN2_SERIALIZATION,
		constant.S_Proto:    constant.PROTOBUF_SERIALIZATION,
	}
}

func SetSerializer(name string, serializer Serializer) {
	serializers[name] = serializer
}

func GetSerializerById(id byte) (Serializer, error) {
	name, ok := nameMaps[id]
	if !ok {
		panic(fmt.Sprintf("serialId %d not found", id))
	}
	serializer, ok := serializers[name]
	if !ok {
		panic(fmt.Sprintf("serialization %s not found", name))
	}
	return serializer, nil
}
