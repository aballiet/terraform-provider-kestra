package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

func resourceFlow() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Kestra Flow.",

		CreateContext: resourceFlowCreate,
		ReadContext:   resourceFlowRead,
		UpdateContext: resourceFlowUpdate,
		DeleteContext: resourceFlowDelete,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Description: "The flow namespace.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"flow_id": {
				Description: "The flow id.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"revision": {
				Description: "The flow revision.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"keep_original_source": {
				Description: "Use the content as source code, keeping comment and indentation.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"content": {
				Description:      "The flow full content in yaml string.",
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: isYamlEqualsFlow,
				StateFunc:        stateFn,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFlowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)
	var diags diag.Diagnostics

	yamlSourceCode := d.Get("keep_original_source").(bool)

	if yamlSourceCode == true {
		r, reqErr := c.yamlRequest("POST", "/api/v1/flows", stringToPointer(d.Get("content").(string)))
		if reqErr != nil {
			return diag.FromErr(reqErr.Err)
		}

		errs := flowSourceApiToSchema(r.(map[string]interface{}), d)
		if errs != nil {
			return errs
		}

		return diags
	} else {
		body, err := flowSchemaToApi(d)
		if err != nil {
			return diag.FromErr(err)
		}

		r, reqErr := c.request("POST", "/api/v1/flows", body)
		if reqErr != nil {
			return diag.FromErr(reqErr.Err)
		}

		errs := flowApiToSchema(r.(map[string]interface{}), d)
		if errs != nil {
			return errs
		}

		return diags
	}
}

func resourceFlowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)
	var diags diag.Diagnostics

	namespaceId, flowId := flowConvertId(d.Id())
	yamlSourceCode := d.Get("keep_original_source").(bool)

	if yamlSourceCode == true {
		r, reqErr := c.yamlRequest("GET", fmt.Sprintf("/api/v1/flows/%s/%s?source=true", namespaceId, flowId), nil)
		if reqErr != nil {
			if reqErr.StatusCode == http.StatusNotFound {
				d.SetId("")
				return diags
			}

			return diag.FromErr(reqErr.Err)
		}

		errs := flowSourceApiToSchema(r.(map[string]interface{}), d)
		if errs != nil {
			return errs
		}

		return diags
	} else {
		r, reqErr := c.request("GET", fmt.Sprintf("/api/v1/flows/%s/%s", namespaceId, flowId), nil)
		if reqErr != nil {
			if reqErr.StatusCode == http.StatusNotFound {
				d.SetId("")
				return diags
			}

			return diag.FromErr(reqErr.Err)
		}

		errs := flowApiToSchema(r.(map[string]interface{}), d)
		if errs != nil {
			return errs
		}

		return diags
	}
}

func resourceFlowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)
	var diags diag.Diagnostics

	if d.HasChanges("content") {
		yamlSourceCode := d.Get("keep_original_source").(bool)
		if yamlSourceCode == true {
			r, reqErr := c.yamlRequest(
				"PUT",
				fmt.Sprintf(
					"/api/v1/flows/%s/%s",
					d.Get("namespace").(string),
					d.Get("flow_id").(string),
				),
				stringToPointer(d.Get("content").(string)),
			)
			if reqErr != nil {
				return diag.FromErr(reqErr.Err)
			}

			errs := flowSourceApiToSchema(r.(map[string]interface{}), d)
			if errs != nil {
				return errs
			}

			return diags
		} else {
			body, err := flowSchemaToApi(d)
			if err != nil {
				return diag.FromErr(err)
			}

			namespaceId, flowId := flowConvertId(d.Id())

			r, reqErr := c.request("PUT", fmt.Sprintf("/api/v1/flows/%s/%s", namespaceId, flowId), body)
			if reqErr != nil {
				return diag.FromErr(reqErr.Err)
			}

			errs := flowApiToSchema(r.(map[string]interface{}), d)
			if errs != nil {
				return errs
			}

			return diags
		}
	} else {
		return resourceFlowRead(ctx, d, meta)
	}
}

func resourceFlowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*Client)
	var diags diag.Diagnostics

	namespaceId, flowId := flowConvertId(d.Id())

	_, reqErr := c.request("DELETE", fmt.Sprintf("/api/v1/flows/%s/%s", namespaceId, flowId), nil)
	if reqErr != nil {
		return diag.FromErr(reqErr.Err)
	}

	d.SetId("")

	return diags
}
