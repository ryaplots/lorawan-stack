// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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

package commands

import (
	"context"

	as "go.thethings.network/lorawan-stack/v3/pkg/applicationserver"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages"
	packageredis "go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/pubsub"
	pubsubredis "go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/pubsub/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/web"
	webredis "go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/web/redis"
	asredis "go.thethings.network/lorawan-stack/v3/pkg/applicationserver/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/redis"
)

func NewPubSubCleaner(ctx context.Context, config *redis.Config) (*pubsub.RegistryCleaner, error) {
	cleaner := &pubsub.RegistryCleaner{
		PubSubRegistry: &pubsubredis.PubSubRegistry{
			Redis: redis.New(config.WithNamespace("as", "io", "pubsub")),
		},
	}
	err := cleaner.RangeToLocalSet(ctx)
	if err != nil {
		return nil, err
	}
	return cleaner, nil
}

func NewPackagesCleaner(ctx context.Context, config *redis.Config) (*packages.RegistryCleaner, error) {
	cleaner := &packages.RegistryCleaner{
		ApplicationPackagesRegistry: &packageredis.ApplicationPackagesRegistry{
			Redis: redis.New(config.WithNamespace("as", "io", "applicationpackages")),
		},
	}
	err := cleaner.RangeToLocalSet(ctx)
	if err != nil {
		return nil, err
	}
	return cleaner, nil
}

func NewASDeviceRegistryCleaner(ctx context.Context, config *redis.Config) (*as.RegistryCleaner, error) {
	cleaner := &as.RegistryCleaner{
		DevRegistry: &asredis.DeviceRegistry{
			Redis: redis.New(config.WithNamespace("as", "devices")),
		},
		AppUpsRegistry: &asredis.ApplicationUplinkRegistry{
			Redis: redis.New(config.WithNamespace("as", "applicationups")),
		},
	}
	err := cleaner.RangeToLocalSet(ctx)
	if err != nil {
		return nil, err
	}
	return cleaner, nil
}

func NewWebhookCleaner(ctx context.Context, config *redis.Config) (*web.RegistryCleaner, error) {
	cleaner := &web.RegistryCleaner{
		WebRegistry: &webredis.WebhookRegistry{
			Redis: redis.New(config.WithNamespace("as", "io", "webhooks")),
		},
	}
	err := cleaner.RangeToLocalSet(ctx)
	if err != nil {
		return nil, err
	}
	return cleaner, nil
}
