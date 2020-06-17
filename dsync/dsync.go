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

//Package dsync solves mutex problem in distributed system.
//in another hand go standard pkg sync solves mutex problem in one process.
//typically, you have a lot of micorservice instance doing same action, only allow one action at a time,
//then you need dsync
package dsync

import (
	"sync"
)

var globalMux sync.Mutex

type Mutex interface {
	ID() string
	Lock(key string, wait bool) (err error)
	Unlock() (err error)
}

//Lock lock a mutex resource
func Lock(wait bool, opts ...LockOption) (*Mutex, error) {
	//TDOD
	return nil, nil
}
