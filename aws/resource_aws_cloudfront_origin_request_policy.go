package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAwsCloudFrontOriginRequestPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsCloudFrontOriginRequestPolicyCreate,
		Read:   resourceAwsCloudFrontOriginRequestPolicyRead,
		Update: resourceAwsCloudFrontOriginRequestPolicyUpdate,
		Delete: resourceAwsCloudFrontOriginRequestPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cookies_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_behavior": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"none", "whitelist", "all"}, false),
						},
						"cookies": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"items": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"headers_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_behavior": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"none", "whitelist", "allViewer", "allViewerAndWhitelistCloudFront"}, false),
						},
						"headers": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"items": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"query_strings_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query_string_behavior": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"none", "whitelist", "allExcept", "all"}, false),
						},
						"query_strings": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"items": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAwsCloudFrontOriginRequestPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cloudfrontconn

	request := &cloudfront.CreateOriginRequestPolicyInput{
		OriginRequestPolicyConfig: expandCloudFrontOriginRequestPolicyConfig(d),
	}

	resp, err := conn.CreateOriginRequestPolicy(request)

	if err != nil {
		return err
	}

	d.SetId(aws.StringValue(resp.OriginRequestPolicy.Id))

	return resourceAwsCloudFrontOriginRequestPolicyRead(d, meta)
}

func resourceAwsCloudFrontOriginRequestPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cloudfrontconn
	request := &cloudfront.GetOriginRequestPolicyInput{
		Id: aws.String(d.Id()),
	}

	resp, err := conn.GetOriginRequestPolicy(request)
	if err != nil {
		return err
	}
	d.Set("etag", aws.StringValue(resp.ETag))

	flattenCloudFrontOriginRequestPolicy(d, resp.OriginRequestPolicy.OriginRequestPolicyConfig)

	return nil
}

func resourceAwsCloudFrontOriginRequestPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cloudfrontconn

	request := &cloudfront.UpdateOriginRequestPolicyInput{
		OriginRequestPolicyConfig: expandCloudFrontOriginRequestPolicyConfig(d),
		Id:                        aws.String(d.Id()),
		IfMatch:                   aws.String(d.Get("etag").(string)),
	}

	_, err := conn.UpdateOriginRequestPolicy(request)
	if err != nil {
		return err
	}

	return resourceAwsCloudFrontOriginRequestPolicyRead(d, meta)
}

func resourceAwsCloudFrontOriginRequestPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cloudfrontconn

	request := &cloudfront.DeleteOriginRequestPolicyInput{
		Id:      aws.String(d.Id()),
		IfMatch: aws.String(d.Get("etag").(string)),
	}

	_, err := conn.DeleteOriginRequestPolicy(request)
	if err != nil {
		if isAWSErr(err, cloudfront.ErrCodeNoSuchOriginRequestPolicy, "") {
			return nil
		}
		return err
	}

	return nil
}