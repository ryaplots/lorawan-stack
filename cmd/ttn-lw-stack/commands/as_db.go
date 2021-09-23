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
	"time"

	"github.com/spf13/cobra"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages"
	asioapredis "go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/pubsub"
	asiopsredis "go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/pubsub/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/web"
	asiowebredis "go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/web/redis"
	asredis "go.thethings.network/lorawan-stack/v3/pkg/applicationserver/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/redis"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
)

func toApplicationSet(list []*ttnpb.ApplicationIdentifiers) (set map[string]struct{}) {
	set = make(map[string]struct{})
	for _, ids := range list {
		set[unique.ID(ctx, ids)] = struct{}{}
	}
	return set
}

func toDeviceSet(list []*ttnpb.EndDeviceIdentifiers) (set map[string]struct{}) {
	set = make(map[string]struct{})
	for _, ids := range list {
		set[unique.ID(ctx, ids)] = struct{}{}
	}
	return set
}

var (
	asDBCommand = &cobra.Command{
		Use:   "as-db",
		Short: "Manage Application Server database",
	}
	asDBCleanupCommand = &cobra.Command{
		Use:   "cleanup",
		Short: "Clean stale Application Server application data",
		RunE: func(cmd *cobra.Command, args []string) error {
			if config.Redis.IsZero() {
				panic("Only Redis is supported by this command")
			}

			logger.Info("Connecting to Identity Server database...")
			db, err := store.Open(ctx, config.IS.DatabaseURI)
			if err != nil {
				return err
			}
			defer db.Close()

			logger.Info("Fetching Identity Server application set")
			appIds, err := store.GetApplicationStore(db).FindAllApplications(ctx)
			if err != nil {
				return err
			}
			isApplicationSet := toApplicationSet(appIds)
			logger.Info("Fetching Identity Server device set")
			devIds, err := store.GetEndDeviceStore(db).FindAllEndDevices(ctx)
			if err != nil {
				return err
			}
			logger.Info("Get end device set")
			isDeviceSet := toDeviceSet(devIds)
			logger.Info("Cleaning up pubsub registry")
			pubsubCleaner := &pubsub.RegistryCleaner{
				PubSubRegistry: &asiopsredis.PubSubRegistry{
					Redis: redis.New(config.Redis.WithNamespace("as", "io", "pubsub")),
				},
			}
			err = pubsubCleaner.CleanData(ctx, isApplicationSet)
			if err != nil {
				return err
			}
			logger.Info("Cleaning up webhook registry")
			if config.AS.Webhooks.Target != "" {
				webhookCleaner := &web.RegistryCleaner{
					WebRegistry: &asiowebredis.WebhookRegistry{
						Redis: redis.New(config.Redis.WithNamespace("as", "io", "webhooks")),
					},
				}
				err = webhookCleaner.CleanData(ctx, isApplicationSet)
				if err != nil {
					return err
				}
			}
			logger.Info("Cleaning up application packages registry")
			appPackagesCleaner := &packages.RegistryCleaner{
				ApplicationPackagesRegistry: &asioapredis.ApplicationPackagesRegistry{
					Redis:   redis.New(config.Redis.WithNamespace("as", "io", "applicationpackages")),
					LockTTL: 10 * time.Second,
				},
			}
			err = appPackagesCleaner.CleanData(ctx, isDeviceSet, isApplicationSet)
			if err != nil {
				return err
			}
			logger.Info("Cleaning up device data")
			deviceCleaner := &applicationserver.RegistryCleaner{
				DevRegistry: &asredis.DeviceRegistry{
					Redis: NewComponentDeviceRegistryRedis(*config, "as"),
				},
				AppUpsRegistry: &asredis.ApplicationUplinkRegistry{
					Redis: redis.New(config.Redis.WithNamespace("as", "applicationups")),
					Limit: config.AS.UplinkStorage.Limit,
				},
			}
			err = deviceCleaner.CleanData(ctx, isDeviceSet)
			if err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	Root.AddCommand(asDBCommand)
	asDBCommand.AddCommand(asDBCleanupCommand)
}
