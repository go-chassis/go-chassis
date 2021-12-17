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

package quota_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/pkg/backends/quota"
	"github.com/stretchr/testify/assert"
	"testing"
)

type inMemory struct {
}

func (im *inMemory) SetLimit(domain, project, resourceType string, limit int64) error {
	panic("implement me")
}

func (im *inMemory) IncreaseUsed(domain, project, resource string, used int64) error {
	return nil
}
func (im *inMemory) DecreaseUsed(domain, project, resource string, used int64) error {
	return nil
}
func (im *inMemory) GetQuota(domain, project, resource string) (*quota.Quota, error) {
	return &quota.Quota{ResourceType: "cpu", Used: 10, Limit: 20}, nil
}
func (im *inMemory) GetQuotas(domain, project string) ([]*quota.Quota, error) {
	return []*quota.Quota{
		{ResourceType: "cpu", Used: 10, Limit: 20}, {ResourceType: "mem", Used: 10, Limit: 256},
	}, nil
}
func TestInit(t *testing.T) {
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.quota.plugin", "mock")
	t.Run("circuit qms plugin", func(t *testing.T) {
		err := quota.Init(quota.Options{
			Endpoint: "",
			Plugin:   "",
		})
		assert.NoError(t, err)
		err = quota.PreCreate("", "", "some", 1)
		assert.NoError(t, err)
	})
	t.Run("install and init", func(t *testing.T) {
		quota.Install("mock", func(options quota.Options) (quota.Manager, error) {
			return &inMemory{}, nil
		})
		err := quota.Init(quota.Options{
			Endpoint: "",
			Plugin:   archaius.GetString("servicecomb.quota.plugin", "mock"),
		})
		assert.NoError(t, err)
	})
	t.Run("pre create,should success", func(t *testing.T) {
		err := quota.PreCreate("", "", "cpu", 2)
		assert.NoError(t, err)
	})
	t.Run("pre create reached maximum,should success", func(t *testing.T) {
		err := quota.PreCreate("", "", "cpu", 12)
		assert.Error(t, err)
	})
	t.Run("no limits", func(t *testing.T) {
		err := quota.PreCreate("", "", "other", 12)
		assert.NoError(t, err)
	})

}
