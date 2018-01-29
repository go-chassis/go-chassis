package tcp

import (
	"context"
	"errors"
	"io"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/ServiceComb/go-chassis/client/highway/pb"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/provider"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/core/util/metadata"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/codec"
	microServer "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	"log"
)

const (
	//Name is a variable of type string which says about the protocol used
	Name = "highway"
)

// constants for request and login
const (
	Request = 0
	Login   = 1
)

var remoteLogin = true

type highwayServer struct {
	tr   transport.Transport
	opts microServer.Options
	exit chan chan error
	sync.RWMutex
}

/*
NOTICE:

SOFTWARE: github.com/micro/go-micro
VERSION: 0.2.0

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   Copyright 2015 Asim Aslam.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

func (s *highwayServer) Init(opts ...microServer.Option) error {
	s.Lock()
	for _, o := range opts {
		o(&s.opts)
	}
	lager.Logger.Debugf("server init,transport:%s", s.opts.Transport.String())
	s.Unlock()
	return nil
}
func (s *highwayServer) Options() microServer.Options {
	s.RLock()
	opts := s.opts
	s.RUnlock()
	return opts
}
func (s *highwayServer) Register(schema interface{}, options ...microServer.RegisterOption) (string, error) {
	opts := microServer.RegisterOptions{}
	var pn string
	for _, o := range options {
		o(&opts)
	}
	if opts.MicroServiceName == "" {
		opts.MicroServiceName = config.SelfServiceName
	}
	mc := config.MicroserviceDefinition
	if mc == nil {
		pn = common.DefaultProvider
	}
	if mc == nil || mc.Provider == "" {
		pn = common.DefaultProvider
	} else {
		if mc.Provider == "" {
			pn = common.DefaultProvider
		} else {
			pn = mc.Provider
		}

	}
	provider.RegisterProvider(pn, opts.MicroServiceName)
	if opts.SchemaID != "" {
		err := provider.RegisterSchemaWithName(opts.MicroServiceName, opts.SchemaID, schema)
		return opts.SchemaID, err
	}
	schemaID, err := provider.RegisterSchema(opts.MicroServiceName, schema)
	return schemaID, err

}

func (s *highwayServer) Start() error {
	opts := s.Options()
	//TODO lot of options
	l, err := opts.Transport.Listen(opts.Address)
	if err != nil {
		return err
	}
	lager.Logger.Warnf(nil, "Highway server listening on: %s", l.Addr())
	s.Lock()
	s.opts.Address = l.Addr()
	s.Unlock()

	go l.Accept(s.accept)

	go func() {
		ch := <-s.exit
		ch <- l.Close()
	}()
	return nil
}
func (s *highwayServer) Stop() error {
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}
func (s *highwayServer) serveSocket(sock transport.Socket, Header []byte, Body []byte, metadata map[string]string, ID int) {
	rpcRequest := &highway.RequestHeader{}
	rpcResponse := &highway.ResponseHeader{}

	codecFunc := codec.NewPBCodec()
	defer func() {
		if r := recover(); r != nil {
			err := r.(string)
			writeError(rpcResponse, codecFunc, sock, ID, errors.New(err))
		}
	}()

	//默认为OK状态
	rpcResponse.StatusCode = 200
	rpcResponse.Reason = ""
	rpcResponse.Flags = 0

	//解码 请求头
	err := codecFunc.Unmarshal(Header, rpcRequest)
	if err != nil {
		writeError(rpcResponse, codecFunc, sock, ID, err)
		return
	}

	//TODO 请求头是Login
	switch rpcRequest.MsgType {
	case Login:
		err = s.loginHandler(sock, rpcResponse, Body, ID)
		if err != nil {
			lager.Logger.Errorf(err, "highway server deal with login message failed")
			return
		}

	case Request:
		err = s.messageHandler(sock, Body, rpcRequest, rpcResponse, ID, s.opts.ChainName)
		if err != nil {
			lager.Logger.Errorf(err, "highway server deal with request message failed")
			return
		}
	default:
		lager.Logger.Errorf(err, "highway server receive an unknow  message type")
		//TODO 异常请求不断链是否OK
		return
	}
}

func (s *highwayServer) accept(sock transport.Socket) {
	defer func() {
		// close socket
		sock.Close()
		if r := recover(); r != nil {
			lager.Logger.Warnf(nil, string(debug.Stack()), r)
		}
	}()

	for {
		Header, Body, md, ID, err := sock.Recv()
		if err != nil {
			if err != io.EOF {
				lager.Logger.Errorf(err, "Server Receive Err")
			}
			return
		}
		s.serveSocket(sock, Header, Body, md, ID)
	}
}

func newHighwayServer(opts ...microServer.Option) microServer.Server {
	return &highwayServer{
		opts: newOptions(opts...),
		exit: make(chan chan error),
	}
}
func newOptions(opt ...microServer.Option) microServer.Options {
	opts := microServer.Options{
		Metadata: map[string]string{},
	}
	if opts.Codecs == nil {
		opts.Codecs = codec.GetCodecMap()
	}
	for _, o := range opt {
		o(&opts)
	}

	return opts
}
func (s *highwayServer) String() string {
	return Name
}
func init() {
	server.InstallPlugin(Name, newHighwayServer)
}

func writeError(rpcResponse *highway.ResponseHeader,
	codec codec.Codec,
	sock transport.Socket,
	id int,
	err error) {
	lager.Logger.Errorf(err, "highway server socket  error")
	rpcResponse.Reason = err.Error()
	//TODO 505 定义为服务端异常
	rpcResponse.StatusCode = 505

	respBytes, err := codec.Marshal(rpcResponse)
	if err != nil {
		lager.Logger.Errorf(err, "server marshal err")
		return
	}

	//TODO 只发送响应
	err = sock.Send(respBytes, nil, nil, id)
	if err != nil {
		lager.Logger.Errorf(err, "sock send err")
	}
	return
}

func (s *highwayServer) loginHandler(sock transport.Socket, rpcResponse *highway.ResponseHeader, Body []byte, ID int) error {
	codecFunc := codec.NewPBCodec()
	loginRequest := &highway.LoginRequest{}

	//TODO 解码请求头
	err := codecFunc.Unmarshal(Body, loginRequest)
	if err != nil {
		lager.Logger.Errorf(err, "highway server unmarshal loginRequest failed")
		return err
	}

	//header:ResponseHeader
	//Body  :LoginResponse

	if loginRequest.UseProtobufMapCodec == remoteLogin {
		loginHeaderBytes, err := codecFunc.Marshal(rpcResponse)
		if err != nil {
			lager.Logger.Errorf(err, "login server marshal loginRequest failed")
			return err
		}
		//TODO 设置服务端编码支持新的编码方式
		loginResponse := &highway.LoginResponse{
			Protocol:            "highway",
			ZipName:             "z",
			UseProtobufMapCodec: remoteLogin,
		}

		loginResponseBytes, err := codecFunc.Marshal(loginResponse)
		if err != nil {
			lager.Logger.Errorf(err, "login server marshal loginResponse failed")
			return err
		}

		err = sock.Send(loginHeaderBytes, loginResponseBytes, nil, ID)
		if err != nil {
			lager.Logger.Errorf(err, "sock send err")
		}
		return err
	}
	return nil
}

func (s *highwayServer) messageHandler(sock transport.Socket, Body []byte, rpcRequest *highway.RequestHeader, rpcResponse *highway.ResponseHeader, ID int, chainName string) error {

	codecFunc := codec.NewPBCodec()
	op, err := provider.GetOperation(rpcRequest.DestMicroservice, rpcRequest.SchemaID, rpcRequest.OperationName)
	if err != nil {
		writeError(rpcResponse, codecFunc, sock, ID, err)
		return err
	}

	i := &invocation.Invocation{}
	if op != nil && op.Args() != nil && len(op.Args()) > 0 {
		if op.Args()[1].Kind() != reflect.Ptr {
			err = errors.New("second arg not ptr")
			writeError(rpcResponse, codecFunc, sock, ID, err)
			return err
		}

		argv := reflect.New(op.Args()[1].Elem())
		err = codecFunc.Unmarshal(Body, argv.Interface())
		if err != nil {
			writeError(rpcResponse, codecFunc, sock, ID, err)
			return err
		}
		i.Args = argv.Interface()
	}

	i.MicroServiceName = rpcRequest.DestMicroservice
	i.SchemaID = rpcRequest.SchemaID
	i.OperationID = rpcRequest.OperationName
	if rpcRequest.GetContext() != nil {
		i.SourceMicroService = rpcRequest.GetContext()[common.HeaderSourceName]
	}
	i.Ctx = metadata.NewContext(context.Background(), rpcRequest.Context)
	i.Protocol = common.ProtocolHighway

	c, err := handler.GetChain(common.Provider, chainName)
	if err != nil {
		lager.Logger.Errorf(err, "Handler chain init err")
	}
	c.Next(i, func(ir *invocation.InvocationResponse) error {
		if ir.Err != nil {
			writeError(rpcResponse, codecFunc, sock, ID, ir.Err)
			return ir.Err
		}
		p, err := provider.GetProvider(i.MicroServiceName)
		if err != nil {
			return err
		}
		r, err := p.Invoke(i)
		if err != nil {
			return err
		}
		log.Println(r)
		result, err := codecFunc.Marshal(r)
		if err != nil {
			lager.Logger.Errorf(err, "Marshal result error")
			writeError(rpcResponse, codecFunc, sock, ID, err)
			return err
		}

		//Todo Context 带什么信息
		respBytes, err := codecFunc.Marshal(rpcResponse)
		if err != nil {
			writeError(rpcResponse, codecFunc, sock, ID, err)
			return err
		}

		//存放响应头
		err = sock.Send(respBytes, result, nil, ID)
		if err != nil {
			lager.Logger.Errorf(err, "sock send err")
			return err
		}
		return err
	})
	return nil
}
