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

package web_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/formatters"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/mock"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/web"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/web/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	componenttest "go.thethings.network/lorawan-stack/v3/pkg/component/test"
	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

type mockComponent struct{}

func (mockComponent) StartTask(conf *component.TaskConfig) {
	component.DefaultStartTask(conf)
}

func (mockComponent) FromRequestContext(ctx context.Context) context.Context {
	return ctx
}

func createdPooledSink(ctx context.Context, t *testing.T, sink web.Sink) web.Sink {
	return web.NewPooledSink(ctx, mockComponent{}, sink, 1, 4)
}

func TestWebhooks(t *testing.T) {
	_, ctx := test.New(t)

	redisClient, flush := test.NewRedis(ctx, "web_test")
	defer flush()
	defer redisClient.Close()
	downlinks := web.DownlinksConfig{
		PublicAddress: "https://example.com/api/v3",
	}
	registry := &redis.WebhookRegistry{
		Redis: redisClient,
	}
	ids := ttnpb.ApplicationWebhookIdentifiers{
		ApplicationIdentifiers: registeredApplicationID,
		WebhookId:              registeredWebhookID,
	}
	for _, tc := range []struct {
		prefix string
		suffix string
	}{
		{
			prefix: "",
			suffix: "",
		},
		{
			prefix: "",
			suffix: "/",
		},
		{
			prefix: "/",
			suffix: "",
		},
		{
			prefix: "/",
			suffix: "/",
		},
	} {
		t.Run(fmt.Sprintf("Prefix%q/Suffix%q", tc.prefix, tc.suffix), func(t *testing.T) {
			_, ctx := test.New(t)
			_, err := registry.Set(ctx, ids, nil, func(_ *ttnpb.ApplicationWebhook) (*ttnpb.ApplicationWebhook, []string, error) {
				return &ttnpb.ApplicationWebhook{
						ApplicationWebhookIdentifiers: ids,
						BaseUrl:                       "https://myapp.com/api/ttn/v3{/appID,devID}" + tc.suffix,
						Headers: map[string]string{
							"Authorization": "key secret",
						},
						DownlinkApiKey: "foo.secret",
						Format:         "json",
						UplinkMessage: &ttnpb.ApplicationWebhook_Message{
							Path: tc.prefix + "up{?devEUI}",
						},
						JoinAccept: &ttnpb.ApplicationWebhook_Message{
							Path: tc.prefix + "join{?joinEUI}",
						},
						DownlinkAck: &ttnpb.ApplicationWebhook_Message{
							Path: tc.prefix + "down/ack",
						},
						DownlinkNack: &ttnpb.ApplicationWebhook_Message{
							Path: tc.prefix + "down/nack",
						},
						DownlinkSent: &ttnpb.ApplicationWebhook_Message{
							Path: tc.prefix + "down/sent",
						},
						DownlinkQueued: &ttnpb.ApplicationWebhook_Message{
							Path: tc.prefix + "down/queued",
						},
						DownlinkQueueInvalidated: &ttnpb.ApplicationWebhook_Message{
							Path: tc.prefix + "down/invalidated",
						},
						DownlinkFailed: &ttnpb.ApplicationWebhook_Message{
							Path: tc.prefix + "down/failed",
						},
						LocationSolved: &ttnpb.ApplicationWebhook_Message{
							Path: tc.prefix + "location",
						},
						ServiceData: &ttnpb.ApplicationWebhook_Message{
							Path: tc.prefix + "service/data",
						},
					},
					[]string{
						"base_url",
						"downlink_api_key",
						"downlink_ack",
						"downlink_failed",
						"downlink_nack",
						"downlink_queued",
						"downlink_queue_invalidated",
						"downlink_sent",
						"format",
						"headers",
						"ids",
						"service_data",
						"join_accept",
						"location_solved",
						"uplink_message",
					}, nil
			})
			if err != nil {
				t.Fatalf("Failed to set webhook in registry: %s", err)
			}

			t.Run("Upstream", func(t *testing.T) {
				baseURL := fmt.Sprintf("https://myapp.com/api/ttn/v3/%s/%s", registeredApplicationID.ApplicationId, registeredDeviceID.DeviceId)
				testSink := &mockSink{
					ch: make(chan *http.Request, 1),
				}
				for _, sink := range []web.Sink{
					testSink,
					createdPooledSink(ctx, t, testSink),
					createdPooledSink(ctx, t,
						createdPooledSink(ctx, t, testSink),
					),
				} {
					t.Run(fmt.Sprintf("%T", sink), func(t *testing.T) {
						ctx, cancel := context.WithCancel(ctx)
						defer cancel()
						c := componenttest.NewComponent(t, &component.Config{})
						as := mock.NewServer(c)
						_, err := web.NewWebhooks(ctx, as, registry, sink, downlinks)
						if err != nil {
							t.Fatalf("Unexpected error %v", err)
						}
						for _, tc := range []struct {
							Name    string
							Message *ttnpb.ApplicationUp
							OK      bool
							URL     string
						}{
							{
								Name: "UplinkMessage/RegisteredDevice",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: registeredDeviceID,
									Up: &ttnpb.ApplicationUp_UplinkMessage{
										UplinkMessage: &ttnpb.ApplicationUplink{
											SessionKeyId: []byte{0x11},
											FPort:        42,
											FCnt:         42,
											FrmPayload:   []byte{0x1, 0x2, 0x3},
										},
									},
								},
								OK:  true,
								URL: fmt.Sprintf("%s/up?devEUI=%s", baseURL, registeredDeviceID.DevEui),
							},
							{
								Name: "UplinkMessage/UnregisteredDevice",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: unregisteredDeviceID,
									Up: &ttnpb.ApplicationUp_UplinkMessage{
										UplinkMessage: &ttnpb.ApplicationUplink{
											SessionKeyId: []byte{0x22},
											FPort:        42,
											FCnt:         42,
											FrmPayload:   []byte{0x1, 0x2, 0x3},
										},
									},
								},
								OK: false,
							},
							{
								Name: "JoinAccept",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: registeredDeviceID,
									Up: &ttnpb.ApplicationUp_JoinAccept{
										JoinAccept: &ttnpb.ApplicationJoinAccept{
											SessionKeyId: []byte{0x22},
										},
									},
								},
								OK:  true,
								URL: fmt.Sprintf("%s/join?joinEUI=%s", baseURL, registeredDeviceID.JoinEui),
							},
							{
								Name: "DownlinkMessage/Ack",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: registeredDeviceID,
									Up: &ttnpb.ApplicationUp_DownlinkAck{
										DownlinkAck: &ttnpb.ApplicationDownlink{
											SessionKeyId: []byte{0x22},
											FCnt:         42,
											FPort:        42,
											FrmPayload:   []byte{0x1, 0x2, 0x3},
										},
									},
								},
								OK:  true,
								URL: fmt.Sprintf("%s/down/ack", baseURL),
							},
							{
								Name: "DownlinkMessage/Nack",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: registeredDeviceID,
									Up: &ttnpb.ApplicationUp_DownlinkNack{
										DownlinkNack: &ttnpb.ApplicationDownlink{
											SessionKeyId: []byte{0x22},
											FCnt:         42,
											FPort:        42,
											FrmPayload:   []byte{0x1, 0x2, 0x3},
										},
									},
								},
								OK:  true,
								URL: fmt.Sprintf("%s/down/nack", baseURL),
							},
							{
								Name: "DownlinkMessage/Sent",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: registeredDeviceID,
									Up: &ttnpb.ApplicationUp_DownlinkSent{
										DownlinkSent: &ttnpb.ApplicationDownlink{
											SessionKeyId: []byte{0x22},
											FCnt:         42,
											FPort:        42,
											FrmPayload:   []byte{0x1, 0x2, 0x3},
										},
									},
								},
								OK:  true,
								URL: fmt.Sprintf("%s/down/sent", baseURL),
							},
							{
								Name: "DownlinkMessage/Queued",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: registeredDeviceID,
									Up: &ttnpb.ApplicationUp_DownlinkQueued{
										DownlinkQueued: &ttnpb.ApplicationDownlink{
											SessionKeyId: []byte{0x22},
											FCnt:         42,
											FPort:        42,
											FrmPayload:   []byte{0x1, 0x2, 0x3},
										},
									},
								},
								OK:  true,
								URL: fmt.Sprintf("%s/down/queued", baseURL),
							},
							{
								Name: "DownlinkMessage/QueueInvalidated",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: registeredDeviceID,
									Up: &ttnpb.ApplicationUp_DownlinkQueueInvalidated{
										DownlinkQueueInvalidated: &ttnpb.ApplicationInvalidatedDownlinks{
											Downlinks: []*ttnpb.ApplicationDownlink{
												{
													SessionKeyId: []byte{0x22},
													FCnt:         42,
													FPort:        42,
													FrmPayload:   []byte{0x1, 0x2, 0x3},
												},
											},
											LastFCntDown: 42,
											SessionKeyId: []byte{0x22},
										},
									},
								},
								OK:  true,
								URL: fmt.Sprintf("%s/down/invalidated", baseURL),
							},
							{
								Name: "DownlinkMessage/Failed",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: registeredDeviceID,
									Up: &ttnpb.ApplicationUp_DownlinkFailed{
										DownlinkFailed: &ttnpb.ApplicationDownlinkFailed{
											ApplicationDownlink: ttnpb.ApplicationDownlink{
												SessionKeyId: []byte{0x22},
												FCnt:         42,
												FPort:        42,
												FrmPayload:   []byte{0x1, 0x2, 0x3},
											},
											Error: ttnpb.ErrorDetails{
												Name: "test",
											},
										},
									},
								},
								OK:  true,
								URL: fmt.Sprintf("%s/down/failed", baseURL),
							},
							{
								Name: "LocationSolved",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: registeredDeviceID,
									Up: &ttnpb.ApplicationUp_LocationSolved{
										LocationSolved: &ttnpb.ApplicationLocation{
											Location: ttnpb.Location{
												Latitude:  10,
												Longitude: 20,
												Altitude:  30,
											},
											Service: "test",
										},
									},
								},
								OK:  true,
								URL: fmt.Sprintf("%s/location", baseURL),
							},
							{
								Name: "ServiceData",
								Message: &ttnpb.ApplicationUp{
									EndDeviceIdentifiers: registeredDeviceID,
									Up: &ttnpb.ApplicationUp_ServiceData{
										ServiceData: &ttnpb.ApplicationServiceData{
											Data: &types.Struct{
												Fields: map[string]*types.Value{
													"battery": {
														Kind: &types.Value_NumberValue{
															NumberValue: 42.0,
														},
													},
												},
											},
											Service: "test",
										},
									},
								},
								OK:  true,
								URL: fmt.Sprintf("%s/service/data", baseURL),
							},
						} {
							t.Run(tc.Name, func(t *testing.T) {
								a := assertions.New(t)
								err := as.Publish(ctx, tc.Message)
								if !a.So(err, should.BeNil) {
									t.FailNow()
								}
								var req *http.Request
								select {
								case req = <-testSink.ch:
									if !tc.OK {
										t.Fatalf("Did not expect message but received: %v", req)
									}
								case <-time.After(timeout):
									if tc.OK {
										t.Fatal("Expected message but nothing received")
									} else {
										return
									}
								}
								a.So(req.URL.String(), should.Equal, tc.URL)
								a.So(req.Header.Get("Authorization"), should.Equal, "key secret")
								a.So(req.Header.Get("Content-Type"), should.Equal, "application/json")
								a.So(req.Header.Get("X-Downlink-Apikey"), should.Equal, "foo.secret")
								a.So(req.Header.Get("X-Downlink-Push"), should.Equal,
									"https://example.com/api/v3/as/applications/foo-app/webhooks/foo-hook/devices/foo-device/down/push")
								a.So(req.Header.Get("X-Downlink-Replace"), should.Equal,
									"https://example.com/api/v3/as/applications/foo-app/webhooks/foo-hook/devices/foo-device/down/replace")
								a.So(req.Header.Get("X-Tts-Domain"), should.Equal, "example.com")
								actualBody, err := ioutil.ReadAll(req.Body)
								if !a.So(err, should.BeNil) {
									t.FailNow()
								}
								expectedBody, err := formatters.JSON.FromUp(tc.Message)
								if !a.So(err, should.BeNil) {
									t.FailNow()
								}
								a.So(actualBody, should.Resemble, expectedBody)
							})
						}
					})
				}
			})
		})
	}

	t.Run("Downstream", func(t *testing.T) {
		is, isAddr := startMockIS(ctx)
		is.add(ctx, registeredApplicationID, registeredApplicationKey)
		httpAddress := "0.0.0.0:8098"
		conf := &component.Config{
			ServiceBase: config.ServiceBase{
				GRPC: config.GRPC{
					Listen:                      ":0",
					AllowInsecureForCredentials: true,
				},
				Cluster: cluster.Config{
					IdentityServer: isAddr,
				},
				HTTP: config.HTTP{
					Listen: httpAddress,
				},
			},
		}
		c := componenttest.NewComponent(t, conf)
		io := mock.NewServer(c)
		testSink := &mockSink{
			Component: c,
			Server:    io,
		}
		w, err := web.NewWebhooks(ctx, testSink.Server, registry, testSink, downlinks)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		c.RegisterWeb(w)
		componenttest.StartComponent(t, c)
		defer c.Close()

		mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)

		t.Run("Authorization", func(t *testing.T) {
			for _, tc := range []struct {
				Name       string
				ID         ttnpb.ApplicationIdentifiers
				Key        string
				ExpectCode int
			}{
				{
					Name:       "Valid",
					ID:         registeredApplicationID,
					Key:        registeredApplicationKey,
					ExpectCode: http.StatusOK,
				},
				{
					Name:       "InvalidKey",
					ID:         registeredApplicationID,
					Key:        "invalid key",
					ExpectCode: http.StatusForbidden,
				},
				{
					Name:       "InvalidIDAndKey",
					ID:         ttnpb.ApplicationIdentifiers{ApplicationId: "--invalid-id"},
					Key:        "invalid key",
					ExpectCode: http.StatusBadRequest,
				},
			} {
				t.Run(tc.Name, func(t *testing.T) {
					a := assertions.New(t)
					url := fmt.Sprintf("http://%s/api/v3/as/applications/%s/webhooks/%s/devices/%s/down/replace",
						httpAddress, tc.ID.ApplicationId, registeredWebhookID, registeredDeviceID.DeviceId,
					)
					body := bytes.NewReader([]byte(`{"downlinks":[]}`))
					req, err := http.NewRequest(http.MethodPost, url, body)
					if !a.So(err, should.BeNil) {
						t.FailNow()
					}
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.Key))
					res, err := http.DefaultClient.Do(req)
					if !a.So(err, should.BeNil) {
						t.FailNow()
					}
					a.So(res.StatusCode, should.Equal, tc.ExpectCode)
					downlinks, err := io.DownlinkQueueList(ctx, registeredDeviceID)
					if !a.So(err, should.BeNil) {
						t.FailNow()
					}
					a.So(downlinks, should.Resemble, []*ttnpb.ApplicationDownlink{})
				})
			}
		})
	})
}

type mockSink struct {
	Component *component.Component
	Server    io.Server
	ch        chan *http.Request
}

func (s *mockSink) FillContext(ctx context.Context) context.Context {
	return s.Component.FillContext(ctx)
}

func (s *mockSink) Process(req *http.Request) error {
	s.ch <- req
	return nil
}
