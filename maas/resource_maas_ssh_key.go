package maas

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maas/gomaasclient/client"
)

func resourceMaasSSHKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource to manage MAAS SSH keys.",
		CreateContext: resourceSSHKeyCreate,
		ReadContext:   resourceSSHKeyRead,
		DeleteContext: resourceSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*client.Client)
				id, err := strconv.Atoi(d.Id())
				if err != nil {
					return nil, err
				}
				sshKey, err := client.SSHKey.Get(id)
				if err != nil {
					return nil, err
				}
				d.SetId(fmt.Sprintf("%v", sshKey.ID))
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A public SSH key.",
			},
		},
	}
}

func resourceSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	sshKey, err := client.SSHKeys.Create(d.Get("key").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%v", sshKey.ID))

	return nil
}

func resourceSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(client.SSHKey.Delete(id))
}

func resourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sshKey, err := client.SSHKey.Get(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%v", sshKey.ID))

	d.Set("key", sshKey.Key)

	return nil
}
