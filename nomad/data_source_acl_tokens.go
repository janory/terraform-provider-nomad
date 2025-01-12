// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"fmt"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceACLTokens() *schema.Resource {
	return &schema.Resource{
		Read: aclTokensDataSourceRead,

		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"acl_tokens": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accessor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"roles": {
							Description: "The roles that are applied to the token.",
							Computed:    true,
							Type:        schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the ACL role.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the ACL role.",
									},
								},
							},
						},
						"global": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_time": {
							Description: "The point after which a token is considered revoked and eligible for destruction.",
							Computed:    true,
							Type:        schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func aclTokensDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConfig).client

	qOpts := &api.QueryOptions{
		Prefix: d.Get("prefix").(string),
	}
	tokens, _, err := client.ACLTokens().List(qOpts)
	if err != nil {
		return fmt.Errorf("error while getting the list of tokens: %v", err)
	}

	result := make([]map[string]interface{}, len(tokens))
	for i, t := range tokens {

		var expirationTime string
		if t.ExpirationTime != nil {
			expirationTime = t.ExpirationTime.Format(time.RFC3339)
		}

		roles := make([]map[string]interface{}, len(t.Roles))
		for i, roleLink := range t.Roles {
			roles[i] = map[string]interface{}{"id": roleLink.ID, "name": roleLink.Name}
		}

		result[i] = map[string]interface{}{
			"accessor_id":     t.AccessorID,
			"name":            t.Name,
			"type":            t.Type,
			"policies":        t.Policies,
			"roles":           roles,
			"global":          t.Global,
			"create_time":     t.CreateTime.String(),
			"expiration_time": expirationTime,
		}
	}

	d.SetId("nomad-tokens")
	return d.Set("acl_tokens", result)
}
