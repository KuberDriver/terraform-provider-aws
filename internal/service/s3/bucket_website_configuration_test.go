package s3_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfs3 "github.com/hashicorp/terraform-provider-aws/internal/service/s3"
)

func TestAccS3BucketWebsiteConfiguration_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_website_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, s3.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckBucketWebsiteConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketWebsiteConfigurationBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "aws_s3_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "index_document.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index_document.0.suffix", "index.html"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccS3BucketWebsiteConfiguration_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_website_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, s3.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckBucketWebsiteConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketWebsiteConfigurationBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfs3.ResourceBucketWebsiteConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccS3BucketWebsiteConfiguration_update(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_website_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, s3.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckBucketWebsiteConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketWebsiteConfigurationBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
				),
			},
			{
				Config: testAccBucketWebsiteConfigurationUpdateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "aws_s3_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "index_document.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "index_document.0.suffix", "index.html"),
					resource.TestCheckResourceAttr(resourceName, "error_document.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "error_document.0.key", "error.html"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccS3BucketWebsiteConfiguration_Redirect(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_website_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, s3.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckBucketWebsiteConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketWebsiteConfigurationConfig_Redirect(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "aws_s3_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "redirect_all_requests_to.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "redirect_all_requests_to.0.host_name", "example.com"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccS3BucketWebsiteConfiguration_RoutingRules_ConditionAndRedirect(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_website_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, s3.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckBucketWebsiteConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketWebsiteConfigurationConfig_RoutingRules_OptionalRedirection(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "routing_rule.*", map[string]string{
						"condition.#":                        "1",
						"condition.0.key_prefix_equals":      "docs/",
						"redirect.#":                         "1",
						"redirect.0.replace_key_prefix_with": "documents/",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketWebsiteConfigurationConfig_RoutingRules_RedirectErrors(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "routing_rule.*", map[string]string{
						"condition.#": "1",
						"condition.0.http_error_code_returned_equals": "404",
						"redirect.#":                         "1",
						"redirect.0.replace_key_prefix_with": "report-404",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketWebsiteConfigurationConfig_RoutingRules_RedirectToPage(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "routing_rule.*", map[string]string{
						"condition.#":                   "1",
						"condition.0.key_prefix_equals": "images/",
						"redirect.#":                    "1",
						"redirect.0.replace_key_with":   "errorpage.html",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccS3BucketWebsiteConfiguration_RoutingRules_MultipleRules(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_website_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, s3.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckBucketWebsiteConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketWebsiteConfigurationConfig_RoutingRules_MultipleRules(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "routing_rule.*", map[string]string{
						"condition.#":                   "1",
						"condition.0.key_prefix_equals": "docs/",
						"redirect.#":                    "1",
						"redirect.0.replace_key_with":   "errorpage.html",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "routing_rule.*", map[string]string{
						"condition.#":                   "1",
						"condition.0.key_prefix_equals": "images/",
						"redirect.#":                    "1",
						"redirect.0.replace_key_with":   "errorpage.html",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketWebsiteConfigurationBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
				),
			},
		},
	})
}

func TestAccS3BucketWebsiteConfiguration_RoutingRules_RedirectOnly(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_website_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, s3.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckBucketWebsiteConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketWebsiteConfigurationConfig_RoutingRules_RedirectOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "routing_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "routing_rule.*", map[string]string{
						"redirect.#":                  "1",
						"redirect.0.protocol":         s3.ProtocolHttps,
						"redirect.0.replace_key_with": "errorpage.html",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBucketWebsiteConfigurationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_s3_bucket_website_configuration" {
			continue
		}

		input := &s3.GetBucketWebsiteInput{
			Bucket: aws.String(rs.Primary.ID),
		}

		output, err := conn.GetBucketWebsite(input)

		if tfawserr.ErrCodeEquals(err, s3.ErrCodeNoSuchBucket, tfs3.ErrCodeNoSuchWebsiteConfiguration) {
			continue
		}

		if err != nil {
			return fmt.Errorf("error getting S3 bucket website configuration (%s): %w", rs.Primary.ID, err)
		}

		if output != nil {
			return fmt.Errorf("S3 bucket website configuration (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckBucketWebsiteConfigurationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn

		input := &s3.GetBucketWebsiteInput{
			Bucket: aws.String(rs.Primary.ID),
		}

		output, err := conn.GetBucketWebsite(input)

		if err != nil {
			return fmt.Errorf("error getting S3 bucket website configuration (%s): %w", rs.Primary.ID, err)
		}

		if output == nil {
			return fmt.Errorf("S3 Bucket website configuration (%s) not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccBucketWebsiteConfigurationBasicConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  acl    = "public-read"

  lifecycle {
    ignore_changes = [
      website
    ]
  }
}

resource "aws_s3_bucket_website_configuration" "test" {
  bucket = aws_s3_bucket.test.id
  index_document {
    suffix = "index.html"
  }
}
`, rName)
}

func testAccBucketWebsiteConfigurationUpdateConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  acl    = "public-read"

  lifecycle {
    ignore_changes = [
      website
    ]
  }
}

resource "aws_s3_bucket_website_configuration" "test" {
  bucket = aws_s3_bucket.test.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}
`, rName)
}

func testAccBucketWebsiteConfigurationConfig_Redirect(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  acl    = "public-read"

  lifecycle {
    ignore_changes = [
      website
    ]
  }
}

resource "aws_s3_bucket_website_configuration" "test" {
  bucket = aws_s3_bucket.test.id
  redirect_all_requests_to {
    host_name = "example.com"
  }
}
`, rName)
}

func testAccBucketWebsiteConfigurationConfig_RoutingRules_OptionalRedirection(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  acl    = "public-read"

  lifecycle {
    ignore_changes = [
      website
    ]
  }
}

resource "aws_s3_bucket_website_configuration" "test" {
  bucket = aws_s3_bucket.test.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }

  routing_rule {
    condition {
      key_prefix_equals = "docs/"
    }
    redirect {
      replace_key_prefix_with = "documents/"
    }
  }
}
`, rName)
}

func testAccBucketWebsiteConfigurationConfig_RoutingRules_RedirectErrors(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  acl    = "public-read"

  lifecycle {
    ignore_changes = [
      website
    ]
  }
}

resource "aws_s3_bucket_website_configuration" "test" {
  bucket = aws_s3_bucket.test.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }

  routing_rule {
    condition {
      http_error_code_returned_equals = "404"
    }
    redirect {
      replace_key_prefix_with = "report-404"
    }
  }
}
`, rName))
}

func testAccBucketWebsiteConfigurationConfig_RoutingRules_RedirectToPage(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  acl    = "public-read"

  lifecycle {
    ignore_changes = [
      website
    ]
  }
}

resource "aws_s3_bucket_website_configuration" "test" {
  bucket = aws_s3_bucket.test.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }

  routing_rule {
    condition {
      key_prefix_equals = "images/"
    }
    redirect {
      replace_key_with = "errorpage.html"
    }
  }
}
`, rName)
}

func testAccBucketWebsiteConfigurationConfig_RoutingRules_RedirectOnly(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  acl    = "public-read"

  lifecycle {
    ignore_changes = [
      website
    ]
  }
}

resource "aws_s3_bucket_website_configuration" "test" {
  bucket = aws_s3_bucket.test.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }

  routing_rule {
    redirect {
      protocol         = "https"
      replace_key_with = "errorpage.html"
    }
  }
}
`, rName)
}

func testAccBucketWebsiteConfigurationConfig_RoutingRules_MultipleRules(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  acl    = "public-read"

  lifecycle {
    ignore_changes = [
      website
    ]
  }
}

resource "aws_s3_bucket_website_configuration" "test" {
  bucket = aws_s3_bucket.test.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }

  routing_rule {
    condition {
      key_prefix_equals = "images/"
    }
    redirect {
      replace_key_with = "errorpage.html"
    }
  }

  routing_rule {
    condition {
      key_prefix_equals = "docs/"
    }
    redirect {
      replace_key_with = "errorpage.html"
    }
  }
}
`, rName)
}
