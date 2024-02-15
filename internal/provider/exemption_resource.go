package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

type exemptionResource struct {
	client *client.Client
}

var (
	_ resource.Resource                = &exemptionResource{}
	_ resource.ResourceWithConfigure   = &exemptionResource{}
	_ resource.ResourceWithImportState = &exemptionResource{}
)

func newExemptionResource() resource.Resource {
	return &exemptionResource{}
}

func (e *exemptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*resourceData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	e.client = data.client
}

func (e *exemptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_exemption", req.ProviderTypeName)
}

func (e *exemptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "A UILD that uniquely identifies the exemption.",
			},
			"metric": schema.StringAttribute{
				Required:    true,
				Description: "The name of the metric to be aggregated.",
			},
			"keep_labels": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     defaultEmptyList{},
				Description: "The array of labels to keep; labels not in this array will be aggregated.",
			},
			"created_at": schema.Int64Attribute{
				Computed:    true,
				Description: "Unix timestamp of when this exemption was created.",
			},
			"updated_at": schema.Int64Attribute{
				Computed:    true,
				Description: "Unix timestamp of when this exemption was last updated.",
			},
		},
	}
}

func (e *exemptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.ExemptionTF
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ex, err := e.client.CreateExemption(plan.ToAPIReq())
	if err != nil {
		resp.Diagnostics.AddError("Unable to create exemption", err.Error())
		return
	}

	state := ex.ToTF()
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (e *exemptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.ExemptionTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ex, err := e.client.ReadExemption(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read exemption", err.Error())
		return
	}

	tf := ex.ToTF()
	resp.Diagnostics.Append(resp.State.Set(ctx, &tf)...)
}

func (e *exemptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.ExemptionTF
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state model.ExemptionTF
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ex := plan.ToAPIReq()
	ex.ID = state.ID.ValueString()

	err := e.client.UpdateExemption(ex)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update exemption", err.Error())
		return
	}

	ex, err = e.client.ReadExemption(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read exemption after updating", err.Error())
		return
	}

	state = ex.ToTF()
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (e *exemptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.ExemptionTF
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := e.client.DeleteExemption(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete exemption", err.Error())
	}
}

func (e *exemptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
