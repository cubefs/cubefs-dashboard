// Copyright 2023 The CubeFS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func Md5Encryption(s string) string {
	data := []byte(s)
	sum := md5.Sum(data)
	md5String := fmt.Sprintf("%x", sum)
	return md5String
}

func Encrypt(src []byte) (string, error) {
	return EncryptByKey(src, GetDefaultKey())
}

func EncryptByKey(src []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	src = PKCS7Padding(src, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	return base64.StdEncoding.EncodeToString(dst), nil
}

func Decrypt(str string) (string, error) {
	return DecryptByKey(str, GetDefaultKey())
}

func DecryptByKey(str string, key []byte) (string, error) {
	if len(str) == 0 {
		return "", nil
	}
	src, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	return string(PKCS7UnPadding(dst)), nil
}

func PKCS7Padding(cipher []byte, blockSize int) []byte {
	padSize := blockSize - len(cipher)%blockSize
	padText := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(cipher, padText...)
}

func PKCS7UnPadding(src []byte) []byte {
	length := len(src)
	size := int(src[length-1])
	return src[:(length - size)]
}

func GetDefaultKey() []byte {
	return []byte{
		54, 38, 98, 35, 89, 102, 99, 38, 57, 52, 106, 71, 76, 112, 99, 52, 40, 66, 68,
		113, 76, 120, 50, 93, 56, 46, 79, 52, 80, 76, 101, 46,
	}
}

func GetKey(key string) []byte {
	return []byte(key)
}

func GenKey(n int) []byte {
	if n < 1 {
		return nil
	}
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)
	return b
}
