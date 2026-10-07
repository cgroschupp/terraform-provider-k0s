package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestK0sctlConfigSSH(t *testing.T) {
	for _, withBastion := range []bool{false, true} {
		name := "direct"
		if withBastion {
			name = "bastion"
		}
		t.Run(name, func(t *testing.T) {
			host := ClusterResourceModelHost{
				Role: types.StringValue("controller+worker"),
				SSH: ClusterResourceModelHostSSH{
					Address: types.StringValue("10.0.1.10"),
					Port:    types.Int64Value(22),
					User:    types.StringValue("root"),
					KeyPath: types.StringValue("~/.ssh/node"),
					Key:     types.StringNull(),
				},
				InstallFlags: types.ListNull(types.StringType),
				Environment:  types.MapNull(types.StringType),
			}
			if withBastion {
				host.SSH.Bastion = &ClusterResourceModelHostSSHBastion{
					Address: types.StringValue("bastion.example.com"),
					Port:    types.Int64Value(2222),
					User:    types.StringValue("ubuntu"),
					KeyPath: types.StringValue("~/.ssh/bastion"),
					Key:     types.StringNull(),
				}
			}
			var diags diag.Diagnostics
			config := getK0sctlConfig(context.Background(), &diags, &ClusterResourceModel{
				Name:    types.StringValue("test"),
				Version: types.StringValue("1.27.2+k0s.0"),
				Hosts:   []ClusterResourceModelHost{host},
			})
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			ssh := config.Spec.Hosts[0].SSH
			if ssh == nil || ssh.Address != "10.0.1.10" || ssh.Port != 22 || ssh.User != "root" || ssh.KeyPath == nil || *ssh.KeyPath != "~/.ssh/node" {
				t.Fatalf("unexpected SSH configuration: %+v", ssh)
			}
			if !withBastion {
				if ssh.Bastion != nil {
					t.Fatal("unexpected bastion for direct connection")
				}
				return
			}
			bastion := ssh.Bastion
			if bastion == nil || bastion.Address != "bastion.example.com" || bastion.Port != 2222 || bastion.User != "ubuntu" || bastion.KeyPath == nil || *bastion.KeyPath != "~/.ssh/bastion" {
				t.Fatalf("unexpected bastion configuration: %+v", bastion)
			}
		})
	}
}

func TestAccClusterResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccClusterResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("k0s_cluster.test", "id", "test"),
					resource.TestCheckResourceAttr("k0s_cluster.test", "name", "test"),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      "k0s_cluster.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// This is not normally necessary, but is here because this
			// 	// example code does not have an actual upstream service.
			// 	// Once the Read method is able to refresh information from
			// 	// the upstream service, this can be removed.
			// 	// ImportStateVerifyIgnore: []string{"configurable_attribute"},
			// },
			// // Update and Read testing
			// {
			// 	Config: testAccClusterResourceConfig(),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("k0s_cluster.test", "configurable_attribute", "two"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccClusterResourceConfig() string {
	return `
resource "k0s_cluster" "test" {
  name = "test"
	version = "1.27.2+k0s.0"

	hosts = [
		{
			role = "controller+worker"

			ssh = {
				address  = "127.0.0.1"
				port     = 9022
				user     = "root"
				key_path = "id_rsa_k0s"
			}
		}
	]
}
`
}
