package bucketlist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"users": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
            "first_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"surname": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
            "password": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users", "http://192.168.64.2:31071"), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	users := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		return diag.FromErr(err)
	}

  flattenedUsers := flattenUsersData(users)

	if err := d.Set("users", flattenedUsers); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenUsersData(orderItems []map[string]interface{}) []interface{} {
  if orderItems != nil {
    ois := make([]interface{}, len(orderItems))

    for i, orderItem := range orderItems {
      oi := make(map[string]interface{})

      oi["first_name"] = orderItem["first_name"]
      oi["surname"]    = orderItem["surname"]
      oi["email"]      = orderItem["email"]
      oi["password"]   = orderItem["password"]

      ois[i] = oi
    }

    return ois
  }

  return make([]interface{}, 0)
}