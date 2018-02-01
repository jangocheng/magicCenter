package def

import "muidea.com/magicCenter/application/common"

// ID 模块ID
const ID = common.AuthorityModuleID

// Name 模块名称
const Name = "Magic Authority"

// Description 模块描述信息
const Description = "Magic 权限管理模块"

// URL 模块Url
const URL = "/authority"

// GetACLByID 查询指定的ACL
const GetACLByID = "/acl/:id"

// QueryACLByModule 查询指定的ACL
const QueryACLByModule = "/acl/"

// PostACL 新增ACL
const PostACL = "/acl/"

// DeleteACL 删除ACL
const DeleteACL = "/acl/:id"

// PutACL 更新ACL（启用、禁用）
const PutACL = "/acl/"

// GetACLAuthGroup 查询指定acl的权限组
const GetACLAuthGroup = "/acl/authgroup/:id"

// PutACLAuthGroup 更新指定acl的权限组
const PutACLAuthGroup = "/acl/authgroup/:id"

// GetModuleACL 查询指定Module的ACL
const GetModuleACL = "/module/acl/:id"

// GetModuleUserAuthGroup 查询拥有指定Module的用户
const GetModuleUserAuthGroup = "/module/user/:id"

// PutModuleUserAuthGroup 更新拥有指定Module的用户
const PutModuleUserAuthGroup = "/module/user/:id"

// GetUserModuleAuthGroup 查询指定用户拥有的Module
const GetUserModuleAuthGroup = "/user/module/:id"

// PutUserModuleAuthGroup 更新指定用户拥有的Module
const PutUserModuleAuthGroup = "/user/module/:id"

// GetUserACL 获取指定用户可访问的ACL
const GetUserACL = "/user/acl/:id"
