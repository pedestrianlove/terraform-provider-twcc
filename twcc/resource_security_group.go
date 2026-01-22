package twcc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type SecurityGroupCreateBody struct {
	Name        string `json:"name"`
	Project     string `json:"project"`
	Description string `json:"description,omitempty"`
}

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Delete: resourceSecurityGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"platform": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"security_group_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ethertype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_ip_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_range_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*PConfig)
	name := d.Get("name").(string)
	platform := d.Get("platform").(string)
	project := d.Get("project").(string)
	description := d.Get("description").(string)

	resourcePath := fmt.Sprintf("api/v3/%s/security_groups/", platform)

	body := SecurityGroupCreateBody{
		Name:        name,
		Project:     project,
		Description: description,
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

	if err != nil {
		return fmt.Errorf("Error creating twcc_security_group %s on %s: %v", name, platform, err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(response), &data)

	if err != nil {
		return err
	}

	securityGroupID := data["id"].(string)
	d.SetId(securityGroupID)

	log.Printf("[DEBUG] Created twcc_security_group: %s", securityGroupID)

	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*PConfig)
	securityGroupID := d.Id()
	platform := d.Get("platform").(string)
	projectID := d.Get("project").(string)

	resourcePath := fmt.Sprintf(
		"api/v3/%s/security_groups/?project=%s&sg=%s",
		platform,
		projectID,
		securityGroupID,
	)

	response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

	if err != nil {
		return fmt.Errorf("Unable to retrieve security group %s on %s: %v", securityGroupID, platform, err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(response), &data)

	if err != nil {
		return fmt.Errorf("Unable to retrieve security group json data: %v", err)
	}

	log.Printf("[DEBUG] Retrieved twcc_security_group %s", securityGroupID)
	d.Set("name", data["name"])
	if createTime, ok := data["create_time"].(string); ok {
		d.Set("create_time", createTime)
	}
	if sgType, ok := data["type"].(string); ok {
		d.Set("type", sgType)
	}
	if user, ok := data["user"].(map[string]interface{}); ok {
		d.Set("user", user)
	}
	if security_group_rules, ok := data["security_group_rules"].([]interface{}); ok {
		rulesInfo := flattenSecurityGroupRulesInfo(security_group_rules)
		d.Set("security_group_rules", rulesInfo)
	}

	return nil
}

func resourceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*PConfig)
	securityGroupID := d.Id()
	platform := d.Get("platform").(string)
	projectID := d.Get("project").(string)

	resourcePath := fmt.Sprintf(
		"api/v3/%s/security_groups/%s/?project=%s",
		platform,
		securityGroupID,
		projectID,
	)

	_, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

	if err != nil {
		return fmt.Errorf("Unable to delete security group %s on %s: %v", securityGroupID, platform, err)
	}

	d.SetId("")
	return nil
}
