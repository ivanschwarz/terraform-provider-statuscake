package statuscake

import (
	"fmt"
	"strconv"

	"log"

	"github.com/DreamItGetIT/statuscake"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceStatusCakeSSLTest() *schema.Resource {
	return &schema.Resource{
		Create: CreateSslTest,
		Update: UpdateSslTest,
		Delete: DeleteSslTest,
		Read:   ReadSslTest,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},

			"contact_groups": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"check_rate": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
			},

			"alert_at": {
				Type:     schema.TypeString,
				Required: true,
			},

			"alert_reminder": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"alert_expiry": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"alert_broken": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func CreateSslTest(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	newTest := &statuscake.PartialSsl{
		Domain:         d.Get("domain").(string),
		ContactGroupsC: d.Get("contact_groups").(string),
		Checkrate:      d.Get("check_rate").(int),
		AlertAt:        d.Get("alert_at").(string),
		AlertReminder:  d.Get("alert_reminder").(bool),
		AlertExpiry:    d.Get("alert_expiry").(bool),
		AlertBroken:    d.Get("alert_broken").(bool),
	}

	log.Printf("[DEBUG] Creating new StatusCake SSL Test: %s", d.Get("domain").(string))
	response, err := client.Ssls().Update(newTest)
	if err != nil {
		log.Printf("%+v\n", response)
		return fmt.Errorf("Error creating StatusCake SSL Test: %s", err.Error())
	}

	d.Set("id", fmt.Sprintf("%s", response.Id))
	d.SetId(fmt.Sprintf("%s", response.Id))
	return ReadSslTest(d, meta)
}

func UpdateSslTest(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)
	params := getSslTestInput(d)

	log.Printf("[DEBUG] StatusCake Test Update for %s", d.Id())
	_, err := client.Ssls().Update(params)
	if err != nil {
		return fmt.Errorf("Error Updating StatusCake SSL Test: %s", err.Error())
	}
	return nil
}

func DeleteSslTest(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)
	log.Printf("[DEBUG] Deleting StatusCake SSL Test: %s", d.Id())
	err := client.Ssls().Delete(d.Id())
	if err != nil {
		return err
	}
	return nil
}

func ReadSslTest(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)
	testResp, err := client.Ssls().Detail(d.Id())
	if err != nil {
		return fmt.Errorf("Error Getting StatusCake SSL Test Details for %s: Error: %s", d.Id(), err)
	}
	d.Set("domain", testResp.Domain)
	d.Set("check_rate", testResp.Checkrate)
	d.Set("paused", testResp.Paused)
	d.Set("alert_at", testResp.AlertAt)
	d.Set("alert_reminder", testResp.AlertReminder)
	d.Set("alert_expiry", testResp.AlertExpiry)
	d.Set("alert_broken", testResp.AlertBroken)
	d.Set("alert_mixed", testResp.AlertMixed)
	d.Set("contact_groups", testResp.ContactGroups)
	d.Set("cert_score", testResp.CertScore)
	d.Set("cert_status", testResp.CertStatus)
	d.Set("cipher", testResp.Cipher)
	d.Set("cipher_score", testResp.CipherScore)
	d.Set("valid_from_utc", testResp.ValidFromUtc)
	d.Set("valid_until_utc", testResp.ValidUntilUtc)
	d.Set("mixed_content", testResp.MixedContent)
	d.Set("flags", testResp.Flags)
	d.Set("last_reminder", testResp.LastReminder)
	d.Set("last_updated_utc", testResp.LastUpdatedUtc)

	return nil
}

func getSslTestInput(d *schema.ResourceData) *statuscake.PartialSsl {
	testId, parseErr := strconv.Atoi(d.Id())
	if parseErr != nil {
		log.Printf("[DEBUG] Error Parsing StatusCake SSL Test Id: %s", d.Id())
	}
	test := &statuscake.PartialSsl{
		Id: testId,
	}
	if v, ok := d.GetOk("domain"); ok {
		test.Domain = v.(string)
	}
	if v, ok := d.GetOk("check_rate"); ok {
		test.Checkrate = v.(int)
	}
	if v, ok := d.GetOk("contact_groups"); ok {
		test.ContactGroupsC = v.(string)
	}
	if v, ok := d.GetOk("alert_at"); ok {
		test.AlertAt = v.(string)
	}
	if v, ok := d.GetOk("alert_reminder"); ok {
		test.AlertReminder = v.(bool)
	}
	if v, ok := d.GetOk("alert_expiry"); ok {
		test.AlertExpiry = v.(bool)
	}
	if v, ok := d.GetOk("alert_broken"); ok {
		test.AlertBroken = v.(bool)
	}

	return test
}
