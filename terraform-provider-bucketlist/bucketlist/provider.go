package bucketlist

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
      "bucketlist_user": resourceUser(),
    },
		DataSourcesMap: map[string]*schema.Resource{
			"bucketlist_users": dataSourceUsers(),
		},
	}
}
