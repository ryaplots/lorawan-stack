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
	"crypto/tls"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/web"
	"go.thethings.network/lorawan-stack/v3/pkg/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"google.golang.org/grpc"
)

func NewClusterComponentConnection(ctx context.Context, config Config, role ttnpb.ClusterRole) (*grpc.ClientConn, cluster.Cluster, error) {
	clusterOpts := []cluster.Option{}
	if config.Cluster.TLS {
		tlsConf := config.TLS
		tls := &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: tlsConf.InsecureSkipVerify,
		}
		if err := tlsConf.Client.ApplyTo(tls); err != nil {
			return nil, nil, err
		}
		clusterOpts = append(clusterOpts, cluster.WithTLSConfig(tls))
	}
	c, err := cluster.New(ctx, &config.Cluster, clusterOpts...)
	if err != nil {
		return nil, nil, err
	}
	if err := c.Join(); err != nil {
		return nil, nil, err
	}
	var cc *grpc.ClientConn
	for i := 0; i < maxRetries; i++ {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		cc, err = c.GetPeerConn(ctx, role, nil)
		if err == nil {
			return cc, c, nil
		}
	}
	return nil, nil, err
}

func FetchIdentityServerApplicationSet(ctx context.Context, client ttnpb.ApplicationRegistryClient, clusterAuth grpc.CallOption) (map[string]struct{}, error) {
	applicationIdentityServerSet := make(map[string]struct{})
	pageCounter := uint32(1)
	// Iterate over application list paginated requests and add them to the IS app map.
	for {
		res, err := client.List(ctx, &ttnpb.ListApplicationsRequest{
			Collaborator: nil,
			FieldMask:    &pbtypes.FieldMask{Paths: []string{"ids"}},
			Limit:        limit,
			Page:         pageCounter,
		}, clusterAuth)
		if err != nil {
			return nil, err
		}
		for _, app := range res.Applications {
			applicationIdentityServerSet[unique.ID(ctx, app.GetIds())] = struct{}{}
		}
		if len(res.Applications) < int(limit) {
			break
		}
		pageCounter++
	}
	return applicationIdentityServerSet, nil
}

func FetchIdentityServerEndDeviceSet(ctx context.Context, client ttnpb.EndDeviceRegistryClient, clusterAuth grpc.CallOption) (map[string]struct{}, error) {
	deviceIdentityServerSet := make(map[string]struct{})
	pageCounter := uint32(1)
	// Iterate over application list paginated requests and add them to the IS app map.
	for {
		res, err := client.List(ctx, &ttnpb.ListEndDevicesRequest{
			ApplicationIds: nil,
			FieldMask:      &pbtypes.FieldMask{Paths: []string{"ids"}},
			Limit:          limit,
			Page:           pageCounter,
		}, clusterAuth)
		if err != nil {
			return nil, err
		}
		for _, dev := range res.EndDevices {
			deviceIdentityServerSet[unique.ID(ctx, dev.EndDeviceIdentifiers)] = struct{}{}
		}
		if len(res.EndDevices) < int(limit) {
			break
		}
		pageCounter++
	}
	return deviceIdentityServerSet, nil
}

var (
	// Define limit for pagination (maximum defined in protos).
	limit = uint32(1000)
	// Max retries to join cluster
	maxRetries = 5
	// Delay between retries to join cluster
	delay       = 500
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
			// Initialize AS registry cleaners (together with their local app/dev sets).
			logger.Info("Initiating PubSub client")
			pubsubCleaner, err := NewPubSubCleaner(ctx, &config.Redis)
			if err != nil {
				return err
			}
			webhookCleaner := &web.RegistryCleaner{}
			if config.AS.Webhooks.Target != "" {
				logger.Info("Initiating webhook client")
				webhookCleaner, err = NewWebhookCleaner(ctx, &config.Redis)
				if err != nil {
					return err
				}
			}
			logger.Info("Initiating application packages registry")
			appPackagesCleaner, err := NewPackagesCleaner(ctx, &config.Redis)
			if err != nil {
				return err
			}
			logger.Info("Initiating device registry")
			deviceCleaner, err := NewASDeviceRegistryCleaner(ctx, &config.Redis)
			if err != nil {
				return err
			}
			// Create cluster and grpc connection with identity server.
			conn, cl, err := NewClusterComponentConnection(ctx, *config, ttnpb.ClusterRole_ENTITY_REGISTRY)
			if err != nil {
				return err
			}
			defer func() {
				logger.Debug("Leaving cluster...")
				if err := cl.Leave(); err != nil {
					logger.WithError(err).Error("Could not leave cluster")
				}
				logger.Debug("Left cluster")
			}()
			client := ttnpb.NewApplicationRegistryClient(conn)
			applicationIdentityServerSet, err := FetchIdentityServerApplicationSet(ctx, client, cl.Auth())
			if err != nil {
				return err
			}
			devClient := ttnpb.NewEndDeviceRegistryClient(conn)
			deviceIdentityServerSet, err := FetchIdentityServerEndDeviceSet(ctx, devClient, cl.Auth())
			if err != nil {
				return err
			}
			// Cleanup data from AS registries.
			logger.Info("Cleaning PubSub registry")
			err = pubsubCleaner.CleanData(ctx, applicationIdentityServerSet)
			if err != nil {
				return err
			}
			logger.Info("Cleaning application packages registry")
			err = appPackagesCleaner.CleanData(ctx, deviceIdentityServerSet, applicationIdentityServerSet)
			if err != nil {
				return err
			}
			logger.Info("Cleaning device registry")
			err = deviceCleaner.CleanData(ctx, deviceIdentityServerSet)
			if err != nil {
				return err
			}
			if webhookCleaner.WebRegistry != nil {
				logger.Info("Cleaning webhook registry")
				err = webhookCleaner.CleanData(ctx, applicationIdentityServerSet)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
)

func init() {
	Root.AddCommand(asDBCommand)
	asDBCommand.AddCommand(asDBCleanupCommand)
}
