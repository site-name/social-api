package localcachelayer

import (
	"bytes"
	"context"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
)

type LocalCacheRoleStore struct {
	store.RoleStore
	rootStore *LocalCacheStore
}

func (s *LocalCacheRoleStore) handleClusterInvalidateRole(msg *model_helper.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.roleCache.Purge()
	} else {
		s.rootStore.roleCache.Remove(string(msg.Data))
	}
}

func (s *LocalCacheRoleStore) handleClusterInvalidateRolePermissions(msg *model_helper.ClusterMessage) {
	if bytes.Equal(msg.Data, clearCacheMessageData) {
		s.rootStore.rolePermissionsCache.Purge()
	} else {
		s.rootStore.rolePermissionsCache.Remove(string(msg.Data))
	}
}

func (s LocalCacheRoleStore) Upsert(role model.Role) (*model.Role, error) {
	if role.Name != "" {
		defer s.rootStore.doInvalidateCacheCluster(s.rootStore.roleCache, role.Name)
		defer s.rootStore.doClearCacheCluster(s.rootStore.rolePermissionsCache)
	}
	return s.RoleStore.Upsert(role)
}

func (s LocalCacheRoleStore) GetByName(ctx context.Context, name string) (*model.Role, error) {
	var role *model.Role
	if err := s.rootStore.doStandardReadCache(s.rootStore.roleCache, name, &role); err == nil {
		return role, nil
	}

	role, err := s.RoleStore.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	s.rootStore.doStandardAddToCache(s.rootStore.roleCache, name, role)
	return role, nil
}

func (s LocalCacheRoleStore) GetByNames(names []string) (model.RoleSlice, error) {
	var foundRoles []*model.Role
	var rolesToQuery []string

	for _, roleName := range names {
		var role *model.Role
		if err := s.rootStore.doStandardReadCache(s.rootStore.roleCache, roleName, &role); err == nil {
			foundRoles = append(foundRoles, role)
		} else {
			rolesToQuery = append(rolesToQuery, roleName)
		}
	}

	roles, err := s.RoleStore.GetByNames(rolesToQuery)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		s.rootStore.doStandardAddToCache(s.rootStore.roleCache, role.Name, role)
	}

	return append(foundRoles, roles...), nil
}

func (s LocalCacheRoleStore) Delete(roleId string) (*model.Role, error) {
	role, err := s.RoleStore.Delete(roleId)

	if err == nil {
		s.rootStore.doInvalidateCacheCluster(s.rootStore.roleCache, role.Name)
		defer s.rootStore.doClearCacheCluster(s.rootStore.rolePermissionsCache)
	}
	return role, err
}

func (s LocalCacheRoleStore) PermanentDeleteAll() error {
	defer s.rootStore.roleCache.Purge()
	defer s.rootStore.doClearCacheCluster(s.rootStore.roleCache)
	defer s.rootStore.doClearCacheCluster(s.rootStore.rolePermissionsCache)

	return s.RoleStore.PermanentDeleteAll()
}
