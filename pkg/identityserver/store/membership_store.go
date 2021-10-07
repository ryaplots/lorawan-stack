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

package store

import (
	"context"
	"fmt"
	"runtime/trace"

	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// GetMembershipStore returns an MembershipStore on the given db (or transaction).
func GetMembershipStore(db *gorm.DB) MembershipStore {
	return &membershipStore{store: newStore(db)}
}

type membershipStore struct {
	*store
}

func (s *membershipStore) queryWithIndirectMemberships(ctx context.Context, entityType string, entityIDs ...string) *gorm.DB {
	idColumnName := fmt.Sprintf(`"%[1]ss"."%[1]s_id"`, entityType)
	if entityType == "organization" {
		idColumnName = `"organization_accounts"."uid"`
	}
	query := s.query(ctx, modelForEntityType(entityType)).Select([]string{
		`"indirect_accounts"."account_type" "indirect_account_type"`,
		`"indirect_accounts"."uid" "indirect_account_friendly_id"`,
		`"indirect_memberships"."rights" "indirect_account_rights"`,
		`"direct_accounts"."account_type" "direct_account_type"`,
		`"direct_accounts"."uid" "direct_account_friendly_id"`,
		`"direct_memberships"."rights" "direct_account_rights"`,
		`"direct_memberships"."entity_type" "entity_type"`,
		idColumnName + ` "entity_friendly_id"`,
	})
	if len(entityIDs) > 0 {
		query = query.Where(idColumnName+" IN (?)", entityIDs)
	}
	if entityType == "organization" {
		query = query.Joins(`JOIN "accounts" "organization_accounts" ON "organization_accounts"."account_id" = "organizations"."id" AND "organization_accounts"."account_type" = 'organization'`)
	}
	query = query.Joins(fmt.Sprintf(`JOIN "memberships" "direct_memberships" ON "direct_memberships"."entity_id" = "%[1]ss"."id" AND "direct_memberships"."entity_type" = '%[1]s'`, entityType)).
		Joins(`JOIN "accounts" "direct_accounts" ON "direct_accounts"."id" = "direct_memberships"."account_id"`).
		Joins(`LEFT JOIN "memberships" "indirect_memberships" ON "indirect_memberships"."entity_id" = "direct_accounts"."account_id" AND "indirect_memberships"."entity_type" = "direct_accounts"."account_type"`).
		Joins(`LEFT JOIN "accounts" "indirect_accounts" ON "indirect_accounts"."id" = "indirect_memberships"."account_id"`)
	query = query.Where(`"direct_accounts"."deleted_at" IS NULL AND "indirect_accounts"."deleted_at" IS NULL`)
	return query
}

func (s *membershipStore) queryWithDirectMemberships(ctx context.Context, entityType string, entityIDs ...string) *gorm.DB {
	idColumnName := fmt.Sprintf(`"%[1]ss"."%[1]s_id"`, entityType)
	if entityType == "organization" {
		idColumnName = `"organization_accounts"."uid"`
	}
	query := s.query(ctx, modelForEntityType(entityType)).Select([]string{
		`"direct_accounts"."account_type" "direct_account_type"`,
		`"direct_accounts"."uid" "direct_account_friendly_id"`,
		`"direct_memberships"."rights" "direct_account_rights"`,
		`"direct_memberships"."entity_type" "entity_type"`,
		idColumnName + ` "entity_friendly_id"`,
	})
	if len(entityIDs) > 0 {
		query = query.Where(idColumnName+" IN (?)", entityIDs)
	}
	if entityType == "organization" {
		query = query.Joins(`JOIN "accounts" "organization_accounts" ON "organization_accounts"."account_id" = "organizations"."id" AND "organization_accounts"."account_type" = 'organization'`)
	}
	query = query.Joins(fmt.Sprintf(`JOIN "memberships" "direct_memberships" ON "direct_memberships"."entity_id" = "%[1]ss"."id" AND "direct_memberships"."entity_type" = '%[1]s'`, entityType)).
		Joins(`JOIN "accounts" "direct_accounts" ON "direct_accounts"."id" = "direct_memberships"."account_id"`)
	query = query.Where(`"direct_accounts"."deleted_at" IS NULL`)
	return query
}

func (s *membershipStore) queryMemberships(ctx context.Context, accountID *ttnpb.OrganizationOrUserIdentifiers, entityType string, entityIDs ...string) *gorm.DB {
	switch ids := accountID.Ids.(type) {
	case *ttnpb.OrganizationOrUserIdentifiers_OrganizationIds:
		return s.queryWithDirectMemberships(ctx, entityType, entityIDs...).Where(
			`"direct_accounts"."account_type" = 'organization' AND "direct_accounts"."uid" = ?`,
			ids.OrganizationIds.GetOrganizationId(),
		)
	case *ttnpb.OrganizationOrUserIdentifiers_UserIds:
		return s.queryWithIndirectMemberships(ctx, entityType, entityIDs...).Where(
			`("direct_accounts"."account_type" = 'user' AND "direct_accounts"."uid" = ?) OR ("indirect_accounts"."account_type" = 'user' AND "indirect_accounts"."uid" = ?)`,
			ids.UserIds.GetUserId(), ids.UserIds.GetUserId(),
		)
	default:
		panic("missed oneof type in queryMemberships")
	}
}

func (s *membershipStore) FindMemberships(ctx context.Context, accountID *ttnpb.OrganizationOrUserIdentifiers, entityType string) ([]*ttnpb.EntityIdentifiers, error) {
	defer trace.StartRegion(ctx, fmt.Sprintf("find %s memberships of %s", entityType, accountID.IDString())).End()

	membershipsQuery := s.queryMemberships(ctx, accountID, entityType).Select(`"direct_memberships"."entity_id"`).QueryExpr()
	query := s.query(ctx, modelForEntityType(entityType)).Where(fmt.Sprintf(`"%[1]ss"."id" IN (?)`, entityType), membershipsQuery)
	switch entityType {
	case "organization":
		query = query.
			Joins(`JOIN "accounts" ON "accounts"."account_type" = 'organization' AND "accounts"."account_id" = "organizations"."id"`).
			Select(`"accounts"."uid" AS "friendly_id"`)
	default:
		query = query.
			Select(fmt.Sprintf(`"%[1]ss"."%[1]s_id" AS "friendly_id"`, entityType))
	}

	query = query.Order(orderFromContext(ctx, fmt.Sprintf("%[1]ss", entityType), "friendly_id", "ASC"))
	page := query
	if limit, offset := limitAndOffsetFromContext(ctx); limit != 0 {
		page = query.Limit(limit).Offset(offset)
	}
	var results []struct {
		FriendlyID string
	}
	if err := page.Scan(&results).Error; err != nil {
		return nil, err
	}
	if limit, offset := limitAndOffsetFromContext(ctx); limit != 0 && (offset > 0 || len(results) == int(limit)) {
		countTotal(ctx, query)
	} else {
		setTotal(ctx, uint64(len(results)))
	}
	identifiers := make([]*ttnpb.EntityIdentifiers, len(results))
	for i, result := range results {
		identifiers[i] = buildIdentifiers(entityType, result.FriendlyID)
	}
	return identifiers, nil
}

type membershipChain struct {
	IndirectAccountType       string
	IndirectAccountFriendlyID string
	IndirectAccountRights     Rights
	DirectAccountType         string
	DirectAccountFriendlyID   string
	DirectAccountRights       Rights
	EntityType                string
	EntityFriendlyID          string
}

func (m membershipChain) GetMembershipChain() *MembershipChain {
	indirectAccountRights := ttnpb.Rights(m.IndirectAccountRights)
	directAccountRights := ttnpb.Rights(m.DirectAccountRights)
	c := &MembershipChain{
		RightsOnEntity:    &directAccountRights,
		EntityIdentifiers: buildIdentifiers(m.EntityType, m.EntityFriendlyID),
	}
	switch m.DirectAccountType {
	case "user":
		c.UserIdentifiers = &ttnpb.UserIdentifiers{UserId: m.DirectAccountFriendlyID}
	case "organization":
		if m.IndirectAccountType == "user" {
			c.UserIdentifiers = &ttnpb.UserIdentifiers{UserId: m.IndirectAccountFriendlyID}
			c.RightsOnOrganization = &indirectAccountRights
		}
		c.OrganizationIdentifiers = &ttnpb.OrganizationIdentifiers{OrganizationId: m.DirectAccountFriendlyID}
	}
	return c
}

// MembershipChain is a User -> (Membership -> Organization) -> Membership -> Entity chain.
type MembershipChain struct {
	UserIdentifiers         *ttnpb.UserIdentifiers
	RightsOnOrganization    *ttnpb.Rights
	OrganizationIdentifiers *ttnpb.OrganizationIdentifiers
	RightsOnEntity          *ttnpb.Rights
	EntityIdentifiers       *ttnpb.EntityIdentifiers
}

// GetRights returns the intersected rights.
func (m *MembershipChain) GetRights() *ttnpb.Rights {
	if m.RightsOnOrganization == nil {
		return m.RightsOnEntity.Implied()
	}
	return m.RightsOnEntity.Implied().Intersect(m.RightsOnOrganization.Implied())
}

func (s *membershipStore) FindAccountMembershipChains(ctx context.Context, accountID *ttnpb.OrganizationOrUserIdentifiers, entityType string, entityIDs ...string) ([]*MembershipChain, error) {
	defer trace.StartRegion(ctx, fmt.Sprintf("find membership chains of user on %ss", entityType)).End()
	query := s.queryMemberships(ctx, accountID, entityType, entityIDs...)
	var results []membershipChain
	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}
	chains := make([]*MembershipChain, len(results))
	for i, result := range results {
		chains[i] = result.GetMembershipChain()
	}
	return chains, nil
}

func (s *membershipStore) FindMembers(ctx context.Context, entityID *ttnpb.EntityIdentifiers) (map[*ttnpb.OrganizationOrUserIdentifiers]*ttnpb.Rights, error) {
	defer trace.StartRegion(ctx, fmt.Sprintf("find members of %s", entityID.EntityType())).End()
	query := s.queryWithDirectMemberships(ctx, entityID.EntityType(), entityID.IDString()).Order("direct_account_friendly_id")
	page := query
	if limit, offset := limitAndOffsetFromContext(ctx); limit != 0 {
		page = query.Limit(limit).Offset(offset)
	}
	var results []membershipChain
	if err := page.Scan(&results).Error; err != nil {
		return nil, err
	}
	if limit, offset := limitAndOffsetFromContext(ctx); limit != 0 && (offset > 0 || len(results) == int(limit)) {
		countTotal(ctx, query)
	} else {
		setTotal(ctx, uint64(len(results)))
	}
	membershipRights := make(map[*ttnpb.OrganizationOrUserIdentifiers]*ttnpb.Rights, len(results))
	for _, result := range results {
		chain := result.GetMembershipChain()
		var ids *ttnpb.OrganizationOrUserIdentifiers
		if chain.OrganizationIdentifiers != nil {
			ids = chain.OrganizationIdentifiers.GetOrganizationOrUserIdentifiers()
		} else {
			ids = chain.UserIdentifiers.GetOrganizationOrUserIdentifiers()
		}
		membershipRights[ids] = chain.RightsOnEntity
	}
	return membershipRights, nil
}

var errMembershipNotFound = errors.DefineNotFound(
	"membership_not_found",
	"account `{account_id}` is not a member of `{entity_type}` `{entity_id}`",
)

func (s *membershipStore) GetMember(ctx context.Context, id *ttnpb.OrganizationOrUserIdentifiers, entityID *ttnpb.EntityIdentifiers) (*ttnpb.Rights, error) {
	defer trace.StartRegion(ctx, "get membership").End()
	query := s.queryWithDirectMemberships(ctx, entityID.EntityType(), entityID.IDString()).Where(
		fmt.Sprintf(`"direct_accounts"."account_type" = '%s' AND "direct_accounts"."uid" = ?`, id.EntityType()),
		id.IDString(),
	)
	var results []membershipChain
	executedQuery := query.Scan(&results)
	if err := executedQuery.Error; err != nil {
		return nil, err
	}
	if len(results) != 1 {
		return nil, errMembershipNotFound.WithAttributes(
			"account_id", id.IDString(),
			"entity_type", entityID.EntityType(),
			"entity_id", entityID.IDString(),
		)
	}
	return results[0].GetMembershipChain().RightsOnEntity, nil
}

var errAccountType = errors.DefineInvalidArgument(
	"account_type",
	"account of type `{account_type}` can not collaborate on `{entity_type}`",
)

func (s *membershipStore) SetMember(ctx context.Context, id *ttnpb.OrganizationOrUserIdentifiers, entityID *ttnpb.EntityIdentifiers, rights *ttnpb.Rights) error {
	defer trace.StartRegion(ctx, "update membership").End()

	var account Account
	err := s.query(ctx, Account{}).Where(Account{
		UID:         id.IDString(),
		AccountType: id.EntityType(),
	}).Find(&account).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errNotFoundForID(id)
		}
		return err
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return err
	}

	entity, err := s.findEntity(ctx, entityID, "id")
	if err != nil {
		return err
	}
	if _, ok := entity.(*Organization); ok && account.AccountType != "user" {
		return errAccountType.WithAttributes("account_type", account.AccountType, "entity_type", "organization")
	}
	if err := ctx.Err(); err != nil { // Early exit if context canceled
		return err
	}

	var (
		membership Membership
		query      = s.query(ctx, Membership{})
	)
	err = s.query(ctx, Membership{}).Where(&Membership{
		AccountID:  account.PrimaryKey(),
		EntityID:   entity.PrimaryKey(),
		EntityType: entityTypeForID(entityID),
	}).First(&membership).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}
		if len(rights.Rights) == 0 {
			return err
		}
		membership = Membership{
			AccountID:  account.PrimaryKey(),
			EntityID:   entity.PrimaryKey(),
			EntityType: entityTypeForID(entityID),
		}
		membership.SetContext(ctx)
	} else {
		query = query.Select("updated_at", "rights")
	}
	if len(rights.Rights) == 0 {
		return query.Delete(&membership).Error
	}
	membership.Rights = Rights(*rights)
	return query.Save(&membership).Error
}

func (s *membershipStore) DeleteEntityMembers(ctx context.Context, entityID *ttnpb.EntityIdentifiers) error {
	defer trace.StartRegion(ctx, "delete entity memberships").End()
	entity, err := s.findDeletedEntity(ctx, entityID, "id")
	if err != nil {
		return err
	}
	return s.query(ctx, Membership{}).Where(&Membership{
		EntityID:   entity.PrimaryKey(),
		EntityType: entityTypeForID(entityID),
	}).Delete(&Membership{}).Error
}

func (s *membershipStore) DeleteAccountMembers(ctx context.Context, id *ttnpb.OrganizationOrUserIdentifiers) error {
	defer trace.StartRegion(ctx, "delete account memberships").End()
	var account Account
	err := s.query(ctx, Account{}, withSoftDeleted()).Where(Account{
		UID:         id.IDString(),
		AccountType: id.EntityType(),
	}).Find(&account).Error
	if err != nil {
		return err
	}
	return s.query(ctx, Membership{}).Where(&Membership{
		AccountID: account.PrimaryKey(),
	}).Delete(&Membership{}).Error
}
