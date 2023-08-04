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

package codes

type Code interface {
	Code() string
	Error() string
	Msg() string
}

// 通用错误
const (
	OK              code = 200 // ok
	NeedRedirect    code = 301
	InvalidArgs     code = 400 // invalid request params
	Unauthorized    code = 401 // token expired or error
	Forbidden       code = 403
	NotFound        code = 404
	Conflict        code = 409
	TooManyRequests code = 429
	ResultError     code = 500
	ThirdPartyError code = 510
	DatabaseError   code = 598
	EMAILError      code = 597
	CSRFDetected    code = 599
)

// special error
const (
	ErrorcodeExample code = 5000
	NotFindAccount   code = 5107
	SigninWrongInfo  code = 5100 // account or password error
	SigninFailed     code = 5101 // login failed
	SigninBlocked    code = 5102 // too many login failed and the account is blocked temporarily
	InvalidToken     code = 5103 // token, refresh_token expired or error
	OverQuota        code = 5104
	OpIsNotConfirmed code = 5105 // need to check password
	SessionError     code = 5106 // session error
)

var codeHumanize = map[code]string{
	OK:              "ok",
	NeedRedirect:    "need redirect",
	InvalidArgs:     "invalid args",
	Unauthorized:    "unauthorized",
	Forbidden:       "forbidden",
	NotFound:        "not found",
	Conflict:        "entry exist",
	TooManyRequests: "too many requests",
	ResultError:     "response result error",
	ThirdPartyError: "third party interface response error",
	DatabaseError:   "database err",
	EMAILError:      "email send err",
	CSRFDetected:    "csrf attack detected",
	NotFindAccount:  "not find account",
}

type code int

func (c code) Code() int {
	return int(c)
}

func (c code) Error() int {
	return c.Code()
}

func (c code) Msg() string {
	return codeHumanize[c]
}
