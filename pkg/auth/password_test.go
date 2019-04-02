// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/auth/pbkdf2"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

func init() {
	pbkdf2.SetDefaultIterations(100)
}

func TestHash(t *testing.T) {
	a := assertions.New(t)

	plain := "secret"

	p, err := Hash(plain)
	a.So(err, should.BeNil)

	{
		ok, err := p.Validate(plain)
		a.So(err, should.BeNil)
		a.So(ok, should.BeTrue)
	}

	{
		ok, err := p.Validate("somethingelse")
		a.So(err, should.BeNil)
		a.So(ok, should.BeFalse)
	}

	{
		p := Password("foo")
		ok, err := p.Validate("somethingelse")
		a.So(err, should.NotBeNil)
		a.So(ok, should.BeFalse)
	}

	{
		p := Password("LOL$foo")
		ok, err := p.Validate("somethingelse")
		a.So(err, should.NotBeNil)
		a.So(ok, should.BeFalse)
	}

	{
		p := Password("PBKDF2$foo")
		ok, err := p.Validate("somethingelse")
		a.So(err, should.NotBeNil)
		a.So(ok, should.BeFalse)
	}
}

func TestLegacy(t *testing.T) {
	a := assertions.New(t)

	// this is a pair generated by django
	const plain = "secret"
	const legacy = Password("pbkdf2$sha256$30000$salt$4v3K66vbKbwv3vnwnf32hdzoK8O03GOiBcWFNHul9bo")

	ok, err := legacy.Validate(plain)
	a.So(err, should.BeNil)
	a.So(ok, should.BeTrue)
}