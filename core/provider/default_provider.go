package provider

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
)

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type operation struct {
	sync.Mutex // protects counters
	method     reflect.Method
	In         []reflect.Type
	Out        []reflect.Type
}

func (o *operation) Method() reflect.Method {
	return o.method
}
func (o *operation) Args() []reflect.Type {
	return o.In
}
func (o *operation) Reply() []reflect.Type {
	return o.Out
}

// Schema struct is having schema name, receiver, and registered methods
type Schema struct {
	name    string                // name of schema
	rcvr    reflect.Value         // receiver of methods for the schema
	typ     reflect.Type          // type of the receiver
	methods map[string]*operation // registered methods
}

// DefaultProvider default provider
type DefaultProvider struct {
	mu               sync.RWMutex // protects the schemaMap
	MicroServiceName string
	SchemaMap        map[string]*Schema //string=schemaID
	OperationMap     map[string]*operation
}

// NewProvider returns the object of DefaultProvider
func NewProvider(microserviceName string) Provider {
	return &DefaultProvider{MicroServiceName: microserviceName}
}

// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
//	- exported method of exported type
//	- two arguments, both of exported type
//	- the second argument is a pointer
//	- one return value, of type error
// It returns an error if the receiver is not an exported type or has
// no suitable methods. It also logs the error using package log.
// The client accesses each method using a string of the form "Type.Method",
// where Type is the receiver's concrete type.
func (p *DefaultProvider) Register(schema interface{}) (string, error) {
	return p.register(schema, "", false)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (p *DefaultProvider) RegisterName(name string, rcvr interface{}) error {
	_, err := p.register(rcvr, name, true)
	return err
}
func (p *DefaultProvider) register(schema interface{}, name string, useName bool) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.SchemaMap == nil {
		p.SchemaMap = make(map[string]*Schema)
	}
	s := new(Schema)
	s.typ = reflect.TypeOf(schema)
	s.rcvr = reflect.ValueOf(schema)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if useName {
		sname = name
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + s.typ.String()
		return "", errors.New(s)
	}
	if !isExported(sname) && !useName {
		s := "rpc.Register: type " + sname + " is not exported"
		return "", errors.New(s)
	}
	if _, present := p.SchemaMap[sname]; present {
		return "", errors.New("rpc: service already defined: " + sname)
	}
	s.name = sname

	// Install the methods
	s.methods = suitableMethods(s.typ, true)

	if len(s.methods) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(reflect.PtrTo(s.typ), false)
		if len(method) != 0 {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
		}
		return "", errors.New(str)
	}
	p.SchemaMap[s.name] = s
	return sname, nil
}

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using log if reportErr is true.
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*operation {
	methods := make(map[string]*operation)
	for m := 0; m < typ.NumMethod(); m++ {

		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name

		// Method must be exported.
		if method.PkgPath != "" {
			if reportErr {
				lager.Logger.Warnf(nil, "Method must be exported")
			}
			continue
		}
		// Method needs three ins: receiver, *anyArg, *request.
		if mtype.NumIn() != 3 {
			lager.Logger.Warnf(nil, "method has wrong number of ins, method:%s, nujm:%d", mname, mtype.NumIn())
			continue
		}

		// second arg need not be a pointer.
		any := mtype.In(1)
		if !isExportedOrBuiltinType(any) {
			if reportErr {
				lager.Logger.Warnf(nil, "argument type not exported, method:%s, nujm:%s", mname, any)
			}
			continue
		}

		// Second arg must be a pointer.
		requestType := mtype.In(2)
		if requestType.Kind() != reflect.Ptr {
			if reportErr {
				lager.Logger.Warnf(nil, "method reply type not a pointer, method:%s, requestType:%s", mname, requestType)
			}
			continue
		}
		// request type must be exported.
		if !isExportedOrBuiltinType(requestType) {
			if reportErr {
				lager.Logger.Warnf(nil, "method reply type not exported, method:%s, requestType:%s", mname, requestType)
			}
			continue
		}
		var in = []reflect.Type{any, requestType}
		// Method needs 2 out.
		// response must be a pointer.
		if mtype.NumOut() != 2 {
			lager.Logger.Warnf(nil, "method has wrong number of outs, method:%s, requestType:%d", mname, mtype.NumOut())
			continue
		}
		reponseType := mtype.Out(0)
		if reponseType.Kind() != reflect.Ptr {
			lager.Logger.Warnf(nil, "method reply type not a pointe, method:%s, reponseType:%s", mname, reponseType)
			continue
		}

		// The second return type of the method must be error.
		returnType := mtype.Out(1)
		if returnType != typeOfError {
			if reportErr {
				lager.Logger.Warnf(nil, "method returns method:%s, returnType.String():%s", mname, returnType.String())
			}
			continue
		}
		var out = []reflect.Type{reponseType, returnType}
		methods[mname] = &operation{method: method, In: in, Out: out}
	}

	return methods
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// Invoke is for to invoke the methods of defaultprovider
func (p *DefaultProvider) Invoke(inv *invocation.Invocation) (interface{}, error) {
	schema := p.SchemaMap[inv.SchemaID]
	op := schema.methods[inv.OperationID]
	var err error

	defer func() {
		if r := recover(); r != nil {
			lager.Logger.Errorf(nil, "Invoke returns error:%s", error.Error(r.(error)))
			err = r.(error)
		}
	}()

	function := op.method.Func
	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{schema.rcvr, reflect.Indirect(reflect.New(op.In[0])), reflect.ValueOf(inv.Args)})
	// The return value for the method is an error.
	errInter := returnValues[1].Interface()

	if errInter != nil {
		err = errInter.(error)
	}
	return returnValues[0].Interface(), err

}

// GetOperation get operation
func (p *DefaultProvider) GetOperation(schemaID string, operationID string) (Operation, error) {
	s := p.SchemaMap[schemaID]
	if s == nil {
		return nil, fmt.Errorf("Schema [%s] doesn't exist ", schemaID)
	}
	if s.methods[operationID] == nil {
		return nil, fmt.Errorf("Schema [%s] doesn't exist ", schemaID)
	}
	return s.methods[operationID], nil
}

// Exist check the schema, operation is present or not
func (p *DefaultProvider) Exist(schemaID string, operationID string) bool {
	op, err := p.GetOperation(schemaID, operationID)
	if err != nil {
		return false
	}
	if op != nil {
		return true
	}
	return false
}

func init() {
	InstallProviderPlugin("default", NewProvider)
}
