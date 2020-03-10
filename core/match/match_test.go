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

package match

import (
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatch(t *testing.T) {
	b, _ := Match("exact", "a", "a")
	assert.True(t, b)

	Install("notEq", func(v, e string) bool {
		return !(v == e)
	})

	b, _ = Match("notEq", "a", "a")
	assert.False(t, b)
}

func TestMark(t *testing.T) {
	// just call mark() function to escape dead code checking
	// mark() is a pre-committed function designed in go-chassis 2.0
	inv := &invocation.Invocation{}
	mark(inv)
}
