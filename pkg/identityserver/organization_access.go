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

package identityserver

import (
	"context"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/email"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/emails"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
)

var (
	evtCreateOrganizationAPIKey = events.Define(
		"organization.api-key.create", "create organization API key",
		events.WithVisibility(ttnpb.RIGHT_ORGANIZATION_SETTINGS_API_KEYS),
		events.WithAuthFromContext(),
		events.WithClientInfoFromContext(),
	)
	evtUpdateOrganizationAPIKey = events.Define(
		"organization.api-key.update", "update organization API key",
		events.WithVisibility(ttnpb.RIGHT_ORGANIZATION_SETTINGS_API_KEYS),
		events.WithAuthFromContext(),
		events.WithClientInfoFromContext(),
	)
	evtDeleteOrganizationAPIKey = events.Define(
		"organization.api-key.delete", "delete organization API key",
		events.WithVisibility(ttnpb.RIGHT_ORGANIZATION_SETTINGS_API_KEYS),
		events.WithAuthFromContext(),
		events.WithClientInfoFromContext(),
	)
	evtUpdateOrganizationCollaborator = events.Define(
		"organization.collaborator.update", "update organization collaborator",
		events.WithVisibility(
			ttnpb.RIGHT_ORGANIZATION_SETTINGS_MEMBERS,
			ttnpb.RIGHT_USER_ORGANIZATIONS_LIST,
		),
		events.WithAuthFromContext(),
		events.WithClientInfoFromContext(),
	)
	evtDeleteOrganizationCollaborator = events.Define(
		"organization.collaborator.delete", "delete organization collaborator",
		events.WithVisibility(
			ttnpb.RIGHT_ORGANIZATION_SETTINGS_MEMBERS,
			ttnpb.RIGHT_USER_ORGANIZATIONS_LIST,
		),
		events.WithAuthFromContext(),
		events.WithClientInfoFromContext(),
	)
)

func (is *IdentityServer) listOrganizationRights(ctx context.Context, ids *ttnpb.OrganizationIdentifiers) (*ttnpb.Rights, error) {
	orgRights, err := rights.ListOrganization(ctx, *ids)
	if err != nil {
		return nil, err
	}
	return orgRights.Intersect(ttnpb.AllEntityRights.Union(ttnpb.AllOrganizationRights)), nil
}

func (is *IdentityServer) createOrganizationAPIKey(ctx context.Context, req *ttnpb.CreateOrganizationAPIKeyRequest) (key *ttnpb.APIKey, err error) {
	// Require that caller has rights to manage API keys.
	if err = rights.RequireOrganization(ctx, *req.GetOrganizationIds(), ttnpb.RIGHT_ORGANIZATION_SETTINGS_API_KEYS); err != nil {
		return nil, err
	}
	// Require that caller has at least the rights of the API key.
	if err = rights.RequireOrganization(ctx, *req.GetOrganizationIds(), req.Rights...); err != nil {
		return nil, err
	}
	key, token, err := GenerateAPIKey(ctx, req.Name, req.ExpiresAt, req.Rights...)
	if err != nil {
		return nil, err
	}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		key, err = store.GetAPIKeyStore(db).CreateAPIKey(ctx, req.GetOrganizationIds().GetEntityIdentifiers(), key)
		return err
	})
	if err != nil {
		return nil, err
	}
	key.Key = token
	events.Publish(evtCreateOrganizationAPIKey.NewWithIdentifiersAndData(ctx, req.GetOrganizationIds(), nil))
	err = is.SendContactsEmail(ctx, req, func(data emails.Data) email.MessageData {
		data.SetEntity(req)
		return &emails.APIKeyCreated{Data: data, Key: key, Rights: key.Rights}
	})
	if err != nil {
		log.FromContext(ctx).WithError(err).Error("Could not send API key creation notification email")
	}
	return key, nil
}

func (is *IdentityServer) listOrganizationAPIKeys(ctx context.Context, req *ttnpb.ListOrganizationAPIKeysRequest) (keys *ttnpb.APIKeys, err error) {
	if err = rights.RequireOrganization(ctx, *req.GetOrganizationIds(), ttnpb.RIGHT_ORGANIZATION_SETTINGS_API_KEYS); err != nil {
		return nil, err
	}
	var total uint64
	ctx = store.WithPagination(ctx, req.Limit, req.Page, &total)
	defer func() {
		if err == nil {
			setTotalHeader(ctx, total)
		}
	}()
	keys = &ttnpb.APIKeys{}
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		keys.ApiKeys, err = store.GetAPIKeyStore(db).FindAPIKeys(ctx, req.GetOrganizationIds().GetEntityIdentifiers())
		return err
	})
	if err != nil {
		return nil, err
	}
	for _, key := range keys.ApiKeys {
		key.Key = ""
	}
	return keys, nil
}

func (is *IdentityServer) getOrganizationAPIKey(ctx context.Context, req *ttnpb.GetOrganizationAPIKeyRequest) (key *ttnpb.APIKey, err error) {
	if err = rights.RequireOrganization(ctx, *req.GetOrganizationIds(), ttnpb.RIGHT_ORGANIZATION_SETTINGS_API_KEYS); err != nil {
		return nil, err
	}

	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		_, key, err = store.GetAPIKeyStore(db).GetAPIKey(ctx, req.KeyId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	key.Key = ""
	return key, nil
}

func (is *IdentityServer) updateOrganizationAPIKey(ctx context.Context, req *ttnpb.UpdateOrganizationAPIKeyRequest) (key *ttnpb.APIKey, err error) {
	// Require that caller has rights to manage API keys.
	if err = rights.RequireOrganization(ctx, *req.GetOrganizationIds(), ttnpb.RIGHT_ORGANIZATION_SETTINGS_API_KEYS); err != nil {
		return nil, err
	}

	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		if len(req.APIKey.Rights) > 0 {
			_, key, err := store.GetAPIKeyStore(db).GetAPIKey(ctx, req.APIKey.Id)
			if err != nil {
				return err
			}

			newRights := ttnpb.RightsFrom(req.APIKey.Rights...)
			existingRights := ttnpb.RightsFrom(key.Rights...)

			// Require the caller to have all added rights.
			if err := rights.RequireOrganization(ctx, *req.GetOrganizationIds(), newRights.Sub(existingRights).GetRights()...); err != nil {
				return err
			}
			// Require the caller to have all removed rights.
			if err := rights.RequireOrganization(ctx, *req.GetOrganizationIds(), existingRights.Sub(newRights).GetRights()...); err != nil {
				return err
			}
		}

		key, err = store.GetAPIKeyStore(db).UpdateAPIKey(ctx, req.GetOrganizationIds().GetEntityIdentifiers(), &req.APIKey, req.FieldMask)
		return err
	})
	if err != nil {
		return nil, err
	}
	if key == nil { // API key was deleted.
		events.Publish(evtDeleteOrganizationAPIKey.NewWithIdentifiersAndData(ctx, req.GetOrganizationIds(), nil))
		return &ttnpb.APIKey{}, nil
	}
	key.Key = ""
	events.Publish(evtUpdateOrganizationAPIKey.NewWithIdentifiersAndData(ctx, req.GetOrganizationIds(), nil))
	err = is.SendContactsEmail(ctx, req, func(data emails.Data) email.MessageData {
		data.SetEntity(req)
		return &emails.APIKeyChanged{Data: data, Key: key, Rights: key.Rights}
	})
	if err != nil {
		log.FromContext(ctx).WithError(err).Error("Could not send API key update notification email")
	}

	return key, nil
}

func (is *IdentityServer) getOrganizationCollaborator(ctx context.Context, req *ttnpb.GetOrganizationCollaboratorRequest) (*ttnpb.GetCollaboratorResponse, error) {
	if err := rights.RequireOrganization(ctx, *req.GetOrganizationIds(), ttnpb.RIGHT_ORGANIZATION_SETTINGS_MEMBERS); err != nil {
		return nil, err
	}
	res := &ttnpb.GetCollaboratorResponse{
		OrganizationOrUserIdentifiers: *req.GetCollaborator(),
	}
	err := is.withDatabase(ctx, func(db *gorm.DB) error {
		rights, err := is.getMembershipStore(ctx, db).GetMember(
			ctx,
			req.GetCollaborator(),
			req.GetOrganizationIds().GetEntityIdentifiers(),
		)
		if err != nil {
			return err
		}
		res.Rights = rights.GetRights()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

var errOrganizationNeedsCollaborator = errors.DefineFailedPrecondition("organization_needs_collaborator", "every organization needs at least one collaborator with all rights")

func (is *IdentityServer) setOrganizationCollaborator(ctx context.Context, req *ttnpb.SetOrganizationCollaboratorRequest) (*pbtypes.Empty, error) {
	// Require that caller has rights to manage collaborators.
	if err := rights.RequireOrganization(ctx, *req.GetOrganizationIds(), ttnpb.RIGHT_ORGANIZATION_SETTINGS_MEMBERS); err != nil {
		return nil, err
	}

	err := is.withDatabase(ctx, func(db *gorm.DB) error {
		store := is.getMembershipStore(ctx, db)

		existingRights, err := store.GetMember(
			ctx,
			&req.GetCollaborator().OrganizationOrUserIdentifiers,
			req.GetOrganizationIds().GetEntityIdentifiers(),
		)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		existingRights = existingRights.Implied()
		newRights := ttnpb.RightsFrom(req.GetCollaborator().GetRights()...).Implied()
		addedRights := newRights.Sub(existingRights)
		removedRights := existingRights.Sub(newRights)

		// Require the caller to have all added rights.
		if len(addedRights.GetRights()) > 0 {
			if err := rights.RequireOrganization(ctx, *req.GetOrganizationIds(), addedRights.GetRights()...); err != nil {
				return err
			}
		}

		// Unless we're deleting the collaborator, require the caller to have all removed rights.
		if len(newRights.GetRights()) > 0 && len(removedRights.GetRights()) > 0 {
			if err := rights.RequireOrganization(ctx, *req.GetOrganizationIds(), removedRights.GetRights()...); err != nil {
				return err
			}
		}

		if removedRights.IncludesAll(ttnpb.RIGHT_ORGANIZATION_ALL) {
			memberRights, err := is.getMembershipStore(ctx, db).FindMembers(ctx, req.GetOrganizationIds().GetEntityIdentifiers())
			if err != nil {
				return err
			}
			var hasOtherOwner bool
			for member, rights := range memberRights {
				if unique.ID(ctx, member) == unique.ID(ctx, &req.GetCollaborator().OrganizationOrUserIdentifiers) {
					continue
				}
				if rights.Implied().IncludesAll(ttnpb.RIGHT_ORGANIZATION_ALL) {
					hasOtherOwner = true
					break
				}
			}
			if !hasOtherOwner {
				return errOrganizationNeedsCollaborator.New()
			}
		}

		return store.SetMember(
			ctx,
			&req.GetCollaborator().OrganizationOrUserIdentifiers,
			req.GetOrganizationIds().GetEntityIdentifiers(),
			ttnpb.RightsFrom(req.GetCollaborator().GetRights()...),
		)
	})
	if err != nil {
		return nil, err
	}
	if len(req.GetCollaborator().GetRights()) > 0 {
		events.Publish(evtUpdateOrganizationCollaborator.New(ctx, events.WithIdentifiers(req.GetOrganizationIds(), &req.GetCollaborator().OrganizationOrUserIdentifiers)))
		err = is.SendContactsEmail(ctx, req, func(data emails.Data) email.MessageData {
			data.SetEntity(req)
			return &emails.CollaboratorChanged{Data: data, Collaborator: *req.GetCollaborator()}
		})
		if err != nil {
			log.FromContext(ctx).WithError(err).Error("Could not send collaborator updated notification email")
		}
	} else {
		events.Publish(evtDeleteOrganizationCollaborator.New(ctx, events.WithIdentifiers(req.GetOrganizationIds(), &req.GetCollaborator().OrganizationOrUserIdentifiers)))
	}
	return ttnpb.Empty, nil
}

func (is *IdentityServer) listOrganizationCollaborators(ctx context.Context, req *ttnpb.ListOrganizationCollaboratorsRequest) (collaborators *ttnpb.Collaborators, err error) {
	if err = is.RequireAuthenticated(ctx); err != nil {
		return nil, err
	}
	if err = rights.RequireOrganization(ctx, *req.GetOrganizationIds(), ttnpb.RIGHT_ORGANIZATION_SETTINGS_MEMBERS); err != nil {
		defer func() { collaborators = collaborators.PublicSafe() }()
	}
	var total uint64
	ctx = store.WithPagination(ctx, req.Limit, req.Page, &total)
	defer func() {
		if err == nil {
			setTotalHeader(ctx, total)
		}
	}()
	err = is.withDatabase(ctx, func(db *gorm.DB) (err error) {
		memberRights, err := is.getMembershipStore(ctx, db).FindMembers(ctx, req.GetOrganizationIds().GetEntityIdentifiers())
		if err != nil {
			return err
		}
		collaborators = &ttnpb.Collaborators{}
		for member, rights := range memberRights {
			collaborators.Collaborators = append(collaborators.Collaborators, &ttnpb.Collaborator{
				OrganizationOrUserIdentifiers: *member,
				Rights:                        rights.GetRights(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return collaborators, nil
}

type organizationAccess struct {
	*IdentityServer
}

func (oa *organizationAccess) ListRights(ctx context.Context, req *ttnpb.OrganizationIdentifiers) (*ttnpb.Rights, error) {
	return oa.listOrganizationRights(ctx, req)
}

func (oa *organizationAccess) CreateAPIKey(ctx context.Context, req *ttnpb.CreateOrganizationAPIKeyRequest) (*ttnpb.APIKey, error) {
	return oa.createOrganizationAPIKey(ctx, req)
}

func (oa *organizationAccess) ListAPIKeys(ctx context.Context, req *ttnpb.ListOrganizationAPIKeysRequest) (*ttnpb.APIKeys, error) {
	return oa.listOrganizationAPIKeys(ctx, req)
}

func (oa *organizationAccess) GetAPIKey(ctx context.Context, req *ttnpb.GetOrganizationAPIKeyRequest) (*ttnpb.APIKey, error) {
	return oa.getOrganizationAPIKey(ctx, req)
}

func (oa *organizationAccess) UpdateAPIKey(ctx context.Context, req *ttnpb.UpdateOrganizationAPIKeyRequest) (*ttnpb.APIKey, error) {
	return oa.updateOrganizationAPIKey(ctx, req)
}

func (oa *organizationAccess) GetCollaborator(ctx context.Context, req *ttnpb.GetOrganizationCollaboratorRequest) (*ttnpb.GetCollaboratorResponse, error) {
	return oa.getOrganizationCollaborator(ctx, req)
}

func (oa *organizationAccess) SetCollaborator(ctx context.Context, req *ttnpb.SetOrganizationCollaboratorRequest) (*pbtypes.Empty, error) {
	return oa.setOrganizationCollaborator(ctx, req)
}

func (oa *organizationAccess) ListCollaborators(ctx context.Context, req *ttnpb.ListOrganizationCollaboratorsRequest) (*ttnpb.Collaborators, error) {
	return oa.listOrganizationCollaborators(ctx, req)
}
