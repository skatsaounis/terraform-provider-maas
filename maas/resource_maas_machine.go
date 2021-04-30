package maas

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionutbalutoiu/gomaasclient/client"
	"github.com/ionutbalutoiu/gomaasclient/entity"
)

func resourceMaasMachine() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMachineCreate,
		ReadContext:   resourceMachineRead,
		UpdateContext: resourceMachineUpdate,
		DeleteContext: resourceMachineDelete,

		Schema: map[string]*schema.Schema{
			"power_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"power_parameters": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"pxe_mac_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "amd64/generic",
			},
			"min_hwe_kernel": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pool": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceMachineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	// Create MAAS machine
	machine, err := client.Machines.Create(getMachineCreateParams(d), getMachinePowerParams(d))
	if err != nil {
		return diag.FromErr(err)
	}

	// Save Id
	d.SetId(machine.SystemID)

	// Wait for machine to be ready
	log.Printf("[DEBUG] Waiting for machine (%s) to become ready\n", machine.SystemID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Commissioning", "Testing"},
		Target:     []string{"Ready"},
		Refresh:    getMachineStatusFunc(client, machine.SystemID),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	// Return updated machine
	return resourceMachineUpdate(ctx, d, m)
}

func resourceMachineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	// Get machine
	machine, err := client.Machine.Get(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Set Terraform state
	if err := d.Set("architecture", machine.Architecture); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("min_hwe_kernel", machine.MinHWEKernel); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hostname", machine.Hostname); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain", machine.Domain.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("zone", machine.Zone.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("pool", machine.Pool.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMachineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	// Update machine
	machine, err := client.Machine.Get(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.Machine.Update(machine.SystemID, getMachineUpdateParams(d, machine), getMachinePowerParams(d))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMachineRead(ctx, d, m)
}

func resourceMachineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	// Delete machine
	err := client.Machine.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getMachinePowerParams(d *schema.ResourceData) map[string]string {
	powerParams := d.Get("power_parameters").(map[string]interface{})
	params := make(map[string]string, len(powerParams))
	for k, v := range powerParams {
		params[fmt.Sprintf("power_parameters_%s", k)] = v.(string)
	}
	return params
}

func getMachineCreateParams(d *schema.ResourceData) *entity.MachineParams {
	params := entity.MachineParams{
		PowerType:     d.Get("power_type").(string),
		PXEMacAddress: d.Get("pxe_mac_address").(string),
		Commission:    true,
	}

	if p, ok := d.GetOk("architecture"); ok {
		params.Architecture = p.(string)
	}
	if p, ok := d.GetOk("min_hwe_kernel"); ok {
		params.MinHWEKernel = p.(string)
	}
	if p, ok := d.GetOk("hostname"); ok {
		params.Hostname = p.(string)
	}
	if p, ok := d.GetOk("domain"); ok {
		params.Domain = p.(string)
	}

	return &params
}

func getMachineUpdateParams(d *schema.ResourceData, machine *entity.Machine) *entity.MachineParams {
	params := entity.MachineParams{
		PowerType:    d.Get("power_type").(string),
		CPUCount:     machine.CPUCount,
		Memory:       machine.Memory,
		SwapSize:     machine.SwapSize,
		Architecture: machine.Architecture,
		MinHWEKernel: machine.MinHWEKernel,
		Description:  machine.Description,
	}

	if p, ok := d.GetOk("architecture"); ok {
		params.Architecture = p.(string)
	}
	if p, ok := d.GetOk("min_hwe_kernel"); ok {
		params.MinHWEKernel = p.(string)
	}
	if p, ok := d.GetOk("hostname"); ok {
		params.Hostname = p.(string)
	}
	if p, ok := d.GetOk("domain"); ok {
		params.Domain = p.(string)
	}
	if p, ok := d.GetOk("zone"); ok {
		params.Zone = p.(string)
	}
	if p, ok := d.GetOk("pool"); ok {
		params.Pool = p.(string)
	}

	return &params
}
