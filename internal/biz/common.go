// Copyright 2023 Nautes Authors
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

package biz

const (
	// CodeRepo naming prefix
	RepoPrefix = "repo-"
	// Key data default store user
	DefaultUser = "default"
	// DeployKey fingerprint data Key
	Fingerprint = "fingerprint"
	// DeployKey fingerprint data ID
	DeployKeyID = "id"
	// AccessToken fingerprint data ID
	AccessTokenID = "id"
	// Secret Repo stores git related data engine name
	SecretsGitEngine = "git"
	// Secret Repo stores deploykey key
	SecretsDeployKey = "deploykey"
	// Secret Repo stores access token key
	SecretsAccessToken = "accesstoken"
	// Access token key path name
	AccessTokenName = "accesstoken-api"
	// Git read-only and read-write permissions
	ReadOnly  DeployKeyType = "readonly"
	ReadWrite DeployKeyType = "readwrite"
	// Access token scope permissions
	APIPermission AccessTokenPermission = "api"
	// Access token authorization role
	Developer  AccessLevelValue = 30
	Maintainer AccessLevelValue = 40
	Owner      AccessLevelValue = 50
)
