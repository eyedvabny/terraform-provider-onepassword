package onepassword

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceItemDocument() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemDocumentRead,
		Create: resourceItemDocumentCreate,
		Delete: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemDocumentRead(d, meta)
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vault": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceItemDocumentRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultId := d.Get("vault").(string)
	err, v := m.onePassClient.ReadItem(getId(d), vaultId)
	log.Printf("[DEBUG] %v", v)
	if err != nil {
		return err
	}
	if v.Template != Category2Template(DocumentCategory) {
		return errors.New("Item is not from " + string(DocumentCategory))
	}

	d.SetId(v.Uuid)
	d.Set("name", v.Overview.Title)
	d.Set("tags", v.Overview.Tags)
	d.Set("vault", v.Vault)

	if err, content := m.onePassClient.ReadDocument(v.Uuid); err != nil {
		return err
	} else {
		d.Set("content", content)
	}
	return nil
}

func resourceItemDocumentCreate(d *schema.ResourceData, meta interface{}) error {
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(DocumentCategory),
		Overview: Overview{
			Title: d.Get("name").(string),
			Tags:  ParseTags(d),
		},
	}
	m := meta.(*Meta)
	err := m.onePassClient.CreateDocument(item, d.Get("file_path").(string))
	if err != nil {
		return err
	}
	d.SetId(item.Uuid)
	return nil
}
