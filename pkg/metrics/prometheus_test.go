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

package metrics_test

import (
	"github.com/go-chassis/go-chassis/v2/pkg/metrics"
	"testing"
)

func TestSplit(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 string
	}{
		{"key is namespace_subsystem_name, should return namespace, subsystem",
			args{"namespace_subsystem_name"}, "namespace", "subsystem", "name"},
		{"key is namespace_subsystem_name, should return namespace, subsystem",
			args{"namespace_subsystem_name_sub"}, "namespace", "subsystem", "name_sub"},
		{"key is namespace_name, should return namespace",
			args{"namespace_name"}, "namespace", "", "name"},
		{"key is name, should only return name",
			args{"name"}, "", "", "name"},
		{"key is empty, should be empty",
			args{""}, "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := metrics.Split(tt.args.key)
			if got != tt.want {
				t.Errorf("Split() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Split() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("Split() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
