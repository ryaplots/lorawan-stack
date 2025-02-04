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

package jsonpb

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/TheThingsIndustries/protoc-gen-go-json/jsonplugin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var jsonpluginFeatureFlag = func() bool {
	b, _ := strconv.ParseBool(os.Getenv("TTN_LW_EXP_JSONPLUGIN"))
	return b
}()

// TTN returns the default TTN JSONPb marshaler.
func TTN() *TTNMarshaler {
	return &TTNMarshaler{
		GoGoJSONPb: &GoGoJSONPb{
			OrigName:    true,
			EnumsAsInts: true,
		},
	}
}

type TTNMarshaler struct {
	*GoGoJSONPb
}

func (*TTNMarshaler) ContentType() string { return "application/json" }

func (m *TTNMarshaler) Marshal(v interface{}) ([]byte, error) {
	if jsonpluginFeatureFlag {
		if marshaler, ok := v.(jsonplugin.Marshaler); ok {
			b, err := jsonplugin.MarshalerConfig{
				EnumsAsInts: true,
			}.Marshal(marshaler)
			if err != nil {
				return nil, err
			}
			if m.GoGoJSONPb.Indent == "" {
				return b, nil
			} else {
				var buf bytes.Buffer
				json.Indent(&buf, b, "", m.GoGoJSONPb.Indent)
				return buf.Bytes(), nil
			}
		}
	}
	return m.GoGoJSONPb.Marshal(v)
}

func (m *TTNMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return &TTNEncoder{w: w, gogo: m.GoGoJSONPb}
}

type TTNEncoder struct {
	w    io.Writer
	gogo *GoGoJSONPb
}

func (e *TTNEncoder) Encode(v interface{}) error {
	if jsonpluginFeatureFlag {
		if marshaler, ok := v.(jsonplugin.Marshaler); ok {
			b, err := jsonplugin.MarshalerConfig{
				EnumsAsInts: true,
			}.Marshal(marshaler)
			if err != nil {
				return err
			}
			if e.gogo.Indent == "" {
				_, err = e.w.Write(b)
			} else {
				var buf bytes.Buffer
				json.Indent(&buf, b, "", e.gogo.Indent)
				io.Copy(e.w, &buf)
			}
			return err
		}
	}
	return e.gogo.NewEncoder(e.w).Encode(v)
}

func (m *TTNMarshaler) Unmarshal(data []byte, v interface{}) error {
	if jsonpluginFeatureFlag {
		if unmarshaler, ok := v.(jsonplugin.Unmarshaler); ok {
			return jsonplugin.UnmarshalerConfig{}.Unmarshal(data, unmarshaler)
		}
	}
	return m.GoGoJSONPb.Unmarshal(data, v)
}

func (m *TTNMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return &NewDecoder{r: r, gogo: m.GoGoJSONPb}
}

type NewDecoder struct {
	r    io.Reader
	gogo *GoGoJSONPb
}

func (d *NewDecoder) Decode(v interface{}) error {
	if jsonpluginFeatureFlag {
		if unmarshaler, ok := v.(jsonplugin.Unmarshaler); ok {
			var data json.RawMessage
			err := json.NewDecoder(d.r).Decode(&data)
			if err != nil {
				return err
			}
			return jsonplugin.UnmarshalerConfig{}.Unmarshal(data, unmarshaler)
		}
	}
	return d.gogo.NewDecoder(d.r).Decode(v)
}

// TTNEventStream returns a TTN JsonPb marshaler with double newlines for
// text/event-stream compatibility.
func TTNEventStream() runtime.Marshaler {
	return &ttnEventStream{TTNMarshaler: TTN()}
}

type ttnEventStream struct {
	*TTNMarshaler
}

func (s *ttnEventStream) ContentType() string { return "text/event-stream" }

func (s *ttnEventStream) Delimiter() []byte { return []byte{'\n', '\n'} }
