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

package secret

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

//GenRSAPrivateKey generate a rsa private key
func GenRSAPrivateKey(bits int) ([]byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	bufferPrivate := new(bytes.Buffer)
	err = pem.Encode(bufferPrivate, block)
	if err != nil {
		return nil, err
	}
	b := bufferPrivate.Bytes()
	return b, nil
}

//GenRSAKeyPair create rsa key pair
func GenRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	public := &private.PublicKey
	return private, public, nil
}

//RSAPrivate2Bytes expose bytes of private key
func RSAPrivate2Bytes(privateKey *rsa.PrivateKey) ([]byte, error) {
	k := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: k,
	}
	bufferPrivate := new(bytes.Buffer)
	err := pem.Encode(bufferPrivate, block)
	if err != nil {
		return nil, err
	}
	b := bufferPrivate.Bytes()
	return b, nil
}

//RSAPublicKey2Bytes expose bytes of public key
func RSAPublicKey2Bytes(publicKey *rsa.PublicKey) ([]byte, error) {
	k, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: k,
	}

	b := pem.EncodeToMemory(block)
	return b, nil
}

//ParseRSAPrivateKey convert string to private key
func ParseRSAPrivateKey(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("failed to parse private key")
	}

	p, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return p, nil
}

//ParseRSAPPublicKey convert string to pub key
func ParseRSAPPublicKey(key string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("failed to parse public key")
	}

	p, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub, ok := p.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not RSA")
	}
	return pub, nil
}
