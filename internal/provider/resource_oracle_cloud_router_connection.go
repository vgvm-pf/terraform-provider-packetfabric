package provider

import (
	"context"
	"time"

	"github.com/PacketFabric/terraform-provider-packetfabric/internal/packetfabric"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOracleCloudRouteConn() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOracleCloudRouteConnCreate,
		ReadContext:   resourceOracleCloudRouteConnRead,
		UpdateContext: resourceOracleCloudRouteConnUpdate,
		DeleteContext: resourceOracleCloudRouteConnDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"circuit_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Circuit ID of the target cloud router. This starts with \"PF-L3-CUST-\".",
			},
			"maybe_nat": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Set this to true if you intend to use NAT on this connection. Defaults: false",
			},
			"maybe_dnat": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Set this to true if you intend to use DNAT on this connection. Defaults: false",
			},
			"vc_ocid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "OCID of the FastConnect virtual circuit that you created from the Oracle side.",
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The region in which you created the FastConnect virtual circuit.",
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "A brief description of this connection.",
			},
			"account_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				DefaultFunc:  schema.EnvDefaultFunc("PF_ACCOUNT_ID", nil),
				ValidateFunc: validation.IsUUID,
				Description: "The UUID for the billing account that should be billed. " +
					"Can also be set with the PF_ACCOUNT_ID environment variable.",
			},
			"pop": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The POP in which you want to provision the connection.",
			},
			"zone": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The desired availability zone of the connection.\n\n\tExample: \"A\"",
			},
			"published_quote_line_uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "UUID of the published quote line with which this connection should be associated.",
			},
			"po_number": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 32),
				Description:  "Purchase order number or identifier of a service.",
			},
			"labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Label value linked to an object.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"etl": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Early Termination Liability (ETL) fees apply when terminating a service before its term ends. ETL is prorated to the remaining contract days.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: CloudRouterImportStatePassthroughContext,
		},
	}
}

func resourceOracleCloudRouteConnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*packetfabric.PFClient)
	c.Ctx = ctx

	var diags diag.Diagnostics
	if cid, ok := d.GetOk("circuit_id"); ok {
		oracleRouter := extractOracleCloudRouterConn(d)
		resp, err := c.CreateOracleCloudRouerConnection(oracleRouter, cid.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if err := checkCloudRouterConnectionStatus(c, cid.(string), resp.CloudCircuitID); err != nil {
			return diag.FromErr(err)
		}
		if resp != nil {
			d.SetId(resp.CloudCircuitID)

			if labels, ok := d.GetOk("labels"); ok {
				diagnostics, created := createLabels(c, d.Id(), labels)
				if !created {
					return diagnostics
				}
			}
		}
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Circuit ID not present",
			Detail:   "Please provide a valid Circuit ID.",
		})
	}
	return diags
}

func resourceOracleCloudRouteConnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*packetfabric.PFClient)
	c.Ctx = ctx
	var diags diag.Diagnostics
	if cid, ok := d.GetOk("circuit_id"); ok {
		cloudConnCID := d.Get("id")
		resp, err := c.ReadCloudRouterConnection(cid.(string), cloudConnCID.(string))
		if err != nil {
			diags = diag.FromErr(err)
			return diags
		}

		_ = d.Set("account_uuid", resp.AccountUUID)
		_ = d.Set("circuit_id", resp.CloudRouterCircuitID)
		_ = d.Set("vc_ocid", resp.CloudSettings.VcOcid)
		_ = d.Set("region", resp.CloudSettings.OracleRegion)
		_ = d.Set("description", resp.Description)
		_ = d.Set("pop", resp.CloudProvider.Pop)
		_ = d.Set("zone", resp.Zone)
		_ = d.Set("po_number", resp.PONumber)
		// unsetFields: published_quote_line_uuid
	}

	if _, ok := d.GetOk("labels"); ok {
		labels, err2 := getLabels(c, d.Id())
		if err2 != nil {
			return diag.FromErr(err2)
		}
		_ = d.Set("labels", labels)
	}

	etl, err3 := c.GetEarlyTerminationLiability(d.Id())
	if err3 != nil {
		return diag.FromErr(err3)
	}
	if etl > 0 {
		_ = d.Set("etl", etl)
	}
	return diags
}

func resourceOracleCloudRouteConnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceCloudRouterConnUpdate(ctx, d, m)
}

func resourceOracleCloudRouteConnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceCloudRouterConnDelete(ctx, d, m)
}

func extractOracleCloudRouterConn(d *schema.ResourceData) packetfabric.OracleCloudRouterConn {
	oracleRouter := packetfabric.OracleCloudRouterConn{}
	if maybeNat, ok := d.GetOk("maybe_nat"); ok {
		oracleRouter.MaybeNat = maybeNat.(bool)
	}
	if maybeDNat, ok := d.GetOk("maybe_dnat"); ok {
		oracleRouter.MaybeDNat = maybeDNat.(bool)
	}
	if vcOcid, ok := d.GetOk("vc_ocid"); ok {
		oracleRouter.VcOcid = vcOcid.(string)
	}
	if region, ok := d.GetOk("region"); ok {
		oracleRouter.Region = region.(string)
	}
	if desc, ok := d.GetOk("description"); ok {
		oracleRouter.Description = desc.(string)
	}
	if accountUUID, ok := d.GetOk("account_uuid"); ok {
		oracleRouter.AccountUUID = accountUUID.(string)
	}
	if pop, ok := d.GetOk("pop"); ok {
		oracleRouter.Pop = pop.(string)
	}
	if zone, ok := d.GetOk("zone"); ok {
		oracleRouter.Zone = zone.(string)
	}
	if publishedQuote, ok := d.GetOk("published_quote_line_uuid"); ok {
		oracleRouter.PublishedQuoteLineUUID = publishedQuote.(string)
	}
	if poNumber, ok := d.GetOk("po_number"); ok {
		oracleRouter.PONumber = poNumber.(string)
	}
	return oracleRouter
}
