package iamgo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("IAM_GO_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"iam-go_member": resourceMember(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func setupConnection(ctx context.Context, address string, token string) (*grpc.ClientConn, error) {
	connection, err := Connect(ctx, address, token)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	client, err := setupConnection(ctx, d.Get("address").(string), d.Get("token").(string))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return newPolicyUpdate(iam.NewIAMPolicyClient(client)), diags
}
