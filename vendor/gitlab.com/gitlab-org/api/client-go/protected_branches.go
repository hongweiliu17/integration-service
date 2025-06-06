//
// Copyright 2021, Sander van Harmelen, Michael Lihs
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
//

package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
)

type (
	ProtectedBranchesServiceInterface interface {
		ListProtectedBranches(pid interface{}, opt *ListProtectedBranchesOptions, options ...RequestOptionFunc) ([]*ProtectedBranch, *Response, error)
		GetProtectedBranch(pid interface{}, branch string, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error)
		ProtectRepositoryBranches(pid interface{}, opt *ProtectRepositoryBranchesOptions, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error)
		UnprotectRepositoryBranches(pid interface{}, branch string, options ...RequestOptionFunc) (*Response, error)
		UpdateProtectedBranch(pid interface{}, branch string, opt *UpdateProtectedBranchOptions, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error)
		RequireCodeOwnerApprovals(pid interface{}, branch string, opt *RequireCodeOwnerApprovalsOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// ProtectedBranchesService handles communication with the protected branch
	// related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/protected_branches/
	ProtectedBranchesService struct {
		client *Client
	}
)

var _ ProtectedBranchesServiceInterface = (*ProtectedBranchesService)(nil)

// ProtectedBranch represents a protected branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#list-protected-branches
type ProtectedBranch struct {
	ID                        int                        `json:"id"`
	Name                      string                     `json:"name"`
	PushAccessLevels          []*BranchAccessDescription `json:"push_access_levels"`
	MergeAccessLevels         []*BranchAccessDescription `json:"merge_access_levels"`
	UnprotectAccessLevels     []*BranchAccessDescription `json:"unprotect_access_levels"`
	AllowForcePush            bool                       `json:"allow_force_push"`
	CodeOwnerApprovalRequired bool                       `json:"code_owner_approval_required"`
}

// BranchAccessDescription represents the access description for a protected
// branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#list-protected-branches
type BranchAccessDescription struct {
	ID                     int              `json:"id"`
	AccessLevel            AccessLevelValue `json:"access_level"`
	AccessLevelDescription string           `json:"access_level_description"`
	DeployKeyID            int              `json:"deploy_key_id"`
	UserID                 int              `json:"user_id"`
	GroupID                int              `json:"group_id"`
}

// ListProtectedBranchesOptions represents the available ListProtectedBranches()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#list-protected-branches
type ListProtectedBranchesOptions struct {
	ListOptions
	Search *string `url:"search,omitempty" json:"search,omitempty"`
}

// ListProtectedBranches gets a list of protected branches from a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#list-protected-branches
func (s *ProtectedBranchesService) ListProtectedBranches(pid interface{}, opt *ListProtectedBranchesOptions, options ...RequestOptionFunc) ([]*ProtectedBranch, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/protected_branches", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var p []*ProtectedBranch
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// GetProtectedBranch gets a single protected branch or wildcard protected branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#get-a-single-protected-branch-or-wildcard-protected-branch
func (s *ProtectedBranchesService) GetProtectedBranch(pid interface{}, branch string, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/protected_branches/%s", PathEscape(project), url.PathEscape(branch))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	p := new(ProtectedBranch)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// ProtectRepositoryBranchesOptions represents the available
// ProtectRepositoryBranches() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#protect-repository-branches
type ProtectRepositoryBranchesOptions struct {
	Name                      *string                     `url:"name,omitempty" json:"name,omitempty"`
	PushAccessLevel           *AccessLevelValue           `url:"push_access_level,omitempty" json:"push_access_level,omitempty"`
	MergeAccessLevel          *AccessLevelValue           `url:"merge_access_level,omitempty" json:"merge_access_level,omitempty"`
	UnprotectAccessLevel      *AccessLevelValue           `url:"unprotect_access_level,omitempty" json:"unprotect_access_level,omitempty"`
	AllowForcePush            *bool                       `url:"allow_force_push,omitempty" json:"allow_force_push,omitempty"`
	AllowedToPush             *[]*BranchPermissionOptions `url:"allowed_to_push,omitempty" json:"allowed_to_push,omitempty"`
	AllowedToMerge            *[]*BranchPermissionOptions `url:"allowed_to_merge,omitempty" json:"allowed_to_merge,omitempty"`
	AllowedToUnprotect        *[]*BranchPermissionOptions `url:"allowed_to_unprotect,omitempty" json:"allowed_to_unprotect,omitempty"`
	CodeOwnerApprovalRequired *bool                       `url:"code_owner_approval_required,omitempty" json:"code_owner_approval_required,omitempty"`
}

// BranchPermissionOptions represents a branch permission option.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#protect-repository-branches
type BranchPermissionOptions struct {
	ID          *int              `url:"id,omitempty" json:"id,omitempty"`
	UserID      *int              `url:"user_id,omitempty" json:"user_id,omitempty"`
	GroupID     *int              `url:"group_id,omitempty" json:"group_id,omitempty"`
	DeployKeyID *int              `url:"deploy_key_id,omitempty" json:"deploy_key_id,omitempty"`
	AccessLevel *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
	Destroy     *bool             `url:"_destroy,omitempty" json:"_destroy,omitempty"`
}

// ProtectRepositoryBranches protects a single repository branch or several
// project repository branches using a wildcard protected branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#protect-repository-branches
func (s *ProtectedBranchesService) ProtectRepositoryBranches(pid interface{}, opt *ProtectRepositoryBranchesOptions, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/protected_branches", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	p := new(ProtectedBranch)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// UnprotectRepositoryBranches unprotects the given protected branch or wildcard
// protected branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#unprotect-repository-branches
func (s *ProtectedBranchesService) UnprotectRepositoryBranches(pid interface{}, branch string, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/protected_branches/%s", PathEscape(project), url.PathEscape(branch))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UpdateProtectedBranchOptions represents the available
// UpdateProtectedBranch() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#update-a-protected-branch
type UpdateProtectedBranchOptions struct {
	Name                      *string                     `url:"name,omitempty" json:"name,omitempty"`
	AllowForcePush            *bool                       `url:"allow_force_push,omitempty" json:"allow_force_push,omitempty"`
	CodeOwnerApprovalRequired *bool                       `url:"code_owner_approval_required,omitempty" json:"code_owner_approval_required,omitempty"`
	AllowedToPush             *[]*BranchPermissionOptions `url:"allowed_to_push,omitempty" json:"allowed_to_push,omitempty"`
	AllowedToMerge            *[]*BranchPermissionOptions `url:"allowed_to_merge,omitempty" json:"allowed_to_merge,omitempty"`
	AllowedToUnprotect        *[]*BranchPermissionOptions `url:"allowed_to_unprotect,omitempty" json:"allowed_to_unprotect,omitempty"`
}

// UpdateProtectedBranch updates a protected branch.
//
// Gitlab API docs:
// https://docs.gitlab.com/api/protected_branches/#update-a-protected-branch
func (s *ProtectedBranchesService) UpdateProtectedBranch(pid interface{}, branch string, opt *UpdateProtectedBranchOptions, options ...RequestOptionFunc) (*ProtectedBranch, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/protected_branches/%s", PathEscape(project), url.PathEscape(branch))

	req, err := s.client.NewRequest(http.MethodPatch, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	p := new(ProtectedBranch)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

// RequireCodeOwnerApprovalsOptions represents the available
// RequireCodeOwnerApprovals() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_branches/#update-a-protected-branch
type RequireCodeOwnerApprovalsOptions struct {
	CodeOwnerApprovalRequired *bool `url:"code_owner_approval_required,omitempty" json:"code_owner_approval_required,omitempty"`
}

// RequireCodeOwnerApprovals updates the code owner approval option.
//
// Deprecated: Use UpdateProtectedBranch() instead.
//
// Gitlab API docs:
// https://docs.gitlab.com/api/protected_branches/#update-a-protected-branch
func (s *ProtectedBranchesService) RequireCodeOwnerApprovals(pid interface{}, branch string, opt *RequireCodeOwnerApprovalsOptions, options ...RequestOptionFunc) (*Response, error) {
	updateOptions := &UpdateProtectedBranchOptions{
		CodeOwnerApprovalRequired: opt.CodeOwnerApprovalRequired,
	}
	_, req, err := s.UpdateProtectedBranch(pid, branch, updateOptions, options...)
	return req, err
}
