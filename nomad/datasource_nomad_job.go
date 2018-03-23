package nomad

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceJob() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJobRead,
		Schema: map[string]*schema.Schema{

			"id": {
				Description: "Job ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			// computed attributes
			"name": {
				Description: "Job Name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"namepsace": {
				Description: "Job Namespace",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "Job Type",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"version": {
				Description: "Job Version",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"region": {
				Description: "Job Region",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"datacenters": &schema.Schema{
				Description: "Job Datacenters",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Description: "Deployment Status",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"status_description": {
				Description: "Deployment Status Description",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"submit_time": {
				Description: "Job Submit Time",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"create_index": {
				Description: "Create Index",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"modify_index": {
				Description: "Modify Index",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"job_modify_index": {
				Description: "Job Modify Index",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"stop": {
				Description: "Job Stopped",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"priority": {
				Description: "Job Priority",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"parent_id": {
				Description: "Job Parent ID",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"task_groups": {
				Description: "Job Task Groups",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"placed_canaries": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_revert": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"promoted": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"desired_canaries": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"desired_total": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"placed_alloc": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"healthy_alloc": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"unhealthy_alloc": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"stable": &schema.Schema{
				Description: "Job Stable",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"all_at_once": &schema.Schema{
				Description: "Job All At Once",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"constraints": {
				Description: "Job Constraints",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ltarget": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rtarget": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"operand": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"update_strategy": {
				Description: "Job Update Policy",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stagger": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_parallel": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"health_check": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"min_healthy_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"healthy_deadline": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_revert": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"canary": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"periodic_config": {
				Description: "Job Periodic Configuration",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"spec_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"prohibit_overlap": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"vault_token": {
				Description: "Vault Token",
				Computed:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceJobRead(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	id := d.Get("id").(string)
	log.Printf("[DEBUG] Getting job status: %q", id)
	job, _, err := client.Jobs().Info(id, nil)
	if err != nil {
		// As of Nomad 0.4.1, the API client returns an error for 404
		// rather than a nil result, so we must check this way.
		if strings.Contains(err.Error(), "404") {
			return err
		}

		return fmt.Errorf("error checking for job: %#v", err)
	}

	d.SetId(*job.ID)
	d.Set("name", job.Name)
	d.Set("type", job.Type)
	d.Set("version", job.Version)
	d.Set("namespace", job.Namespace)
	d.Set("region", job.Region)
	d.Set("datacenters", job.Datacenters)
	d.Set("status", job.Status)
	d.Set("status_description", job.StatusDescription)
	d.Set("submit_time", job.SubmitTime)
	d.Set("create_index", job.CreateIndex)
	d.Set("modify_index", job.ModifyIndex)
	d.Set("job_modify_index", job.JobModifyIndex)
	d.Set("stop", job.Stop)
	d.Set("priority", job.Priority)
	d.Set("parent_id", job.ParentID)
	d.Set("task_groups", job.TaskGroups)
	d.Set("stable", job.Stable)
	d.Set("all_at_once", job.AllAtOnce)
	d.Set("contraints", job.Constraints)
	d.Set("update_strategy", job.Update)
	d.Set("periodic_config", job.Periodic)
	d.Set("vault_token", job.VaultToken)

	return nil
}
