// =================================================================================
// This file is auto-generated by the GoFrame CLI tool. You may modify it as needed.
// =================================================================================

package dao

import (
	"bamboo-main/internal/dao/internal"
)

// linkGroupDao is the data access object for the table xf_link_group.
// You can define custom methods on it to extend its functionality as needed.
type linkGroupDao struct {
	*internal.LinkGroupDao
}

var (
	// LinkGroup is a globally accessible object for table xf_link_group operations.
	LinkGroup = linkGroupDao{internal.NewLinkGroupDao()}
)

// Add your custom methods and functionality below.
