package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/VNGritsai/terraform-provider-cloudlab/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &instanceResource{}
	_ resource.ResourceWithConfigure = &instanceResource{}
)

type instanceResource struct {
	client *client.Client
}

type instanceResourceModel struct {
	Id     types.Int64  `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Cpu    types.Int64  `tfsdk:"cpu"`
	Memory types.Int64  `tfsdk:"memory"`
}

func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

func (r *instanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

func (r *instanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Управляет инстансом в cloudlab.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Уникальный идентификатор инстанса, присваивается сервером.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Имя инстанса.",
			},
			"cpu": schema.Int64Attribute{
				Required:    true,
				Description: "Количество ядер CPU.",
			},
			"memory": schema.Int64Attribute{
				Required:    true,
				Description: "Объём памяти в MB.",
			},
		},
	}
}

func (r *instanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Неожиданный тип ProviderData",
			fmt.Sprintf("Ожидался *client.Client, получен %T", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *instanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan instanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateInstance(ctx, client.Instance{
		Name:   plan.Name.ValueString(),
		Cpu:    int(plan.Cpu.ValueInt64()),
		Memory: int(plan.Memory.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Ошибка создания инстанса", err.Error())
		return
	}

	plan.Id = types.Int64Value(int64(created.Id))
	plan.Name = types.StringValue(created.Name)
	plan.Cpu = types.Int64Value(int64(created.Cpu))
	plan.Memory = types.Int64Value(int64(created.Memory))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *instanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state instanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance, err := r.client.GetInstance(ctx, int(state.Id.ValueInt64()))
	if errors.Is(err, client.ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Ошибка чтения инстанса", err.Error())
		return
	}

	state.Id = types.Int64Value(int64(instance.Id))
	state.Name = types.StringValue(instance.Name)
	state.Cpu = types.Int64Value(int64(instance.Cpu))
	state.Memory = types.Int64Value(int64(instance.Memory))

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *instanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan instanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.UpdateInstance(ctx, client.Instance{
		Id:     int(plan.Id.ValueInt64()),
		Name:   plan.Name.ValueString(),
		Cpu:    int(plan.Cpu.ValueInt64()),
		Memory: int(plan.Memory.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Ошибка обновления инстанса", err.Error())
		return
	}

	plan.Id = types.Int64Value(int64(updated.Id))
	plan.Name = types.StringValue(updated.Name)
	plan.Cpu = types.Int64Value(int64(updated.Cpu))
	plan.Memory = types.Int64Value(int64(updated.Memory))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *instanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state instanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteInstance(ctx, int(state.Id.ValueInt64()))
	if errors.Is(err, client.ErrNotFound) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Ошибка удаления инстанса", err.Error())
		return
	}
}
