package iamgo

import (
	"context"
	"testing"

	"cloud.google.com/go/iam/apiv1/iampb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
)

type provider func() (*schema.Provider, error)

func testIAMGoProvider(client *mockIamService) provider {
	return func() (*schema.Provider, error) {
		return &schema.Provider{
			Schema: map[string]*schema.Schema{},
			ResourcesMap: map[string]*schema.Resource{
				"iam-go_member": resourceMember(),
			},
			ConfigureContextFunc: testProviderConfigure(client),
		}, nil
	}
}

func testProviderConfigure(client *mockIamService) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics
		return newPolicyUpdate(client), diags
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func newMockClient() *mockIamService {
	return &mockIamService{
		make(map[string]*iampb.Policy),
	}
}

var _ iampb.IAMPolicyClient = &mockIamService{}

type mockIamService struct {
	policies map[string]*iampb.Policy
}

func (m mockIamService) SetIamPolicy(
	_ context.Context,
	req *iampb.SetIamPolicyRequest,
	_ ...grpc.CallOption,
) (*iampb.Policy, error) {
	m.policies[req.GetResource()] = req.Policy
	return req.Policy, nil
}

func (m mockIamService) GetIamPolicy(
	_ context.Context,
	req *iampb.GetIamPolicyRequest,
	_ ...grpc.CallOption,
) (*iampb.Policy, error) {
	if val, ok := m.policies[req.GetResource()]; ok {
		return val, nil
	}
	return &iampb.Policy{}, nil
}

func (m mockIamService) TestIamPermissions(
	_ context.Context,
	_ *iampb.TestIamPermissionsRequest,
	_ ...grpc.CallOption,
) (*iampb.TestIamPermissionsResponse, error) {
	panic("implement me")
}
