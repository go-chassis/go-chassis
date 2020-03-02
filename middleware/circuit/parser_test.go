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

package circuit_test

import (
	"testing"

	"github.com/go-chassis/go-chassis/middleware/circuit"
	"github.com/stretchr/testify/assert"
)

func TestExtractSchemaAndOperation(t *testing.T) {
	s, sch, op, m := circuit.ExtractServiceSchemaOperationMetrics("service.schema.op.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "schema", sch)
	assert.Equal(t, "op", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = circuit.ExtractServiceSchemaOperationMetrics("service.schema.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "schema", sch)
	assert.Equal(t, "", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = circuit.ExtractServiceSchemaOperationMetrics("service.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "", sch)
	assert.Equal(t, "", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = circuit.ExtractServiceSchemaOperationMetrics("service.schema.op.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "schema", sch)
	assert.Equal(t, "op", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = circuit.ExtractServiceSchemaOperationMetrics("ErrServer.rest./sayhimessage.metrics")
	assert.Equal(t, "ErrServer", s)
	assert.Equal(t, "rest", sch)
	assert.Equal(t, "/sayhimessage", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = circuit.ExtractServiceSchemaOperationMetrics("service.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "", sch)
	assert.Equal(t, "", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = circuit.ExtractServiceSchemaOperationMetrics("service.schema.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "schema", sch)
	assert.Equal(t, "", op)
	assert.Equal(t, "metrics", m)
}

func TestExtractMetricKey(t *testing.T) {
	key, target, sch, op := circuit.ParseCircuitCMD("Consumer.ErrServer.rest./sayhimessage.rejects")
	assert.Equal(t, "ErrServer", target)
	assert.Equal(t, "rest", sch)
	assert.Equal(t, "/sayhimessage", op)
	assert.Equal(t, "Consumer.rejects", key)
}
func TestGetEventType(t *testing.T) {
	m := circuit.GetEventType("Consumer.ErrServer.rest./sayhimessage.rejects")
	assert.Equal(t, "rejects", m)
}

func TestGetMetricsName(t *testing.T) {
	m := circuit.GetMetricsName("Consumer.ErrServer.rest./sayhimessage.rejects")
	assert.Equal(t, "Consumer.rejects", m)
}
