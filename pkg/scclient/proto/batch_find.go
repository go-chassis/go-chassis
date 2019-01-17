// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proto

//FindService specify services you want to fetch
type FindService struct {
	Service *MicroServiceKey `protobuf:"bytes,1,opt,name=service" json:"service"`
	Rev     string           `protobuf:"bytes,2,opt,name=rev" json:"rev,omitempty"`
}

//FindResult is instance list
type FindResult struct {
	Index     int64                   `protobuf:"varint,1,opt,name=index" json:"index"`
	Rev       string                  `protobuf:"bytes,2,opt,name=rev" json:"rev"`
	Instances []*MicroServiceInstance `protobuf:"bytes,3,rep,name=instances" json:"instances,omitempty"`
}

//FindFailedResult represent error message
type FindFailedResult struct {
	Indexes []int64 `protobuf:"varint,1,rep,packed,name=indexes" json:"indexes"`
	Error   *Error  `protobuf:"bytes,2,opt,name=error" json:"error"`
}

//Error is error message
type Error struct {
	Code    int32  `json:"errorCode,string"`
	Message string `json:"errorMessage"`
	Detail  string `json:"detail,omitempty"`
}

//BatchFindResult batch find response
type BatchFindResult struct {
	Failed      []*FindFailedResult `protobuf:"bytes,1,rep,name=failed" json:"failed,omitempty"`
	NotModified []int64             `protobuf:"varint,2,rep,packed,name=notModified" json:"notModified,omitempty"`
	Updated     []*FindResult       `protobuf:"bytes,3,rep,name=updated" json:"updated,omitempty"`
}

//BatchFindInstancesRequest is request body
type BatchFindInstancesRequest struct {
	ConsumerServiceID string         `protobuf:"bytes,1,opt,name=consumerServiceId" json:"consumerServiceId,omitempty"`
	Services          []*FindService `protobuf:"bytes,2,rep,name=services" json:"services,omitempty"`
}

//BatchFindInstancesResponse is response body
type BatchFindInstancesResponse struct {
	Response  *Response        `protobuf:"bytes,1,opt,name=response" json:"response,omitempty"`
	Services  *BatchFindResult `protobuf:"bytes,2,rep,name=services" json:"services,omitempty"`
	Instances *BatchFindResult `protobuf:"bytes,3,rep,name=instances" json:"instances,omitempty"`
}
