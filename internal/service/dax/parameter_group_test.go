package dax_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dax"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccDAXParameterGroup_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_dax_parameter_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, dax.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckParameterGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterGroupExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "parameters.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccParameterGroupConfig_parameters(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterGroupExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "parameters.#", "2"),
				),
			},
		},
	})
}

func testAccCheckParameterGroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DAXConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_dax_parameter_group" {
				continue
			}

			_, err := conn.DescribeParameterGroupsWithContext(ctx, &dax.DescribeParameterGroupsInput{
				ParameterGroupNames: []*string{aws.String(rs.Primary.ID)},
			})
			if err != nil {
				if tfawserr.ErrCodeEquals(err, dax.ErrCodeParameterGroupNotFoundFault) {
					return nil
				}
				return err
			}
		}
		return nil
	}
}

func testAccCheckParameterGroupExists(ctx context.Context, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DAXConn(ctx)

		_, err := conn.DescribeParameterGroupsWithContext(ctx, &dax.DescribeParameterGroupsInput{
			ParameterGroupNames: []*string{aws.String(rs.Primary.ID)},
		})

		return err
	}
}

func testAccParameterGroupConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_dax_parameter_group" "test" {
  name = "%s"
}
`, rName)
}

func testAccParameterGroupConfig_parameters(rName string) string {
	return fmt.Sprintf(`
resource "aws_dax_parameter_group" "test" {
  name = "%s"

  parameters {
    name  = "query-ttl-millis"
    value = "100000"
  }

  parameters {
    name  = "record-ttl-millis"
    value = "100000"
  }
}
`, rName)
}
