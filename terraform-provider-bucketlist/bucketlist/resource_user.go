package bucketlist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"surname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type User struct {
	FirstName string `json:"first_name"`
	Surname   string `json:"surname"`
	UserEmail string `json:"email"`
	Password  string `json:"password"`
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	firstName := d.Get("first_name").(string)
	surname := d.Get("surname").(string)
	emailAddress := d.Get("email").(string)
	password := d.Get("password").(string)

	oi := User{
		FirstName: firstName,
		Surname:   surname,
		UserEmail: emailAddress,
		Password:  password,
	}
	client := &http.Client{Timeout: 10 * time.Second}

	o := createUser(client, oi)
	if o != nil {
		return o
	}
	d.SetId(oi.UserEmail)

	resourceUserRead(ctx, d, m)
	return diags
}

func createUser(client *http.Client, user User) diag.Diagnostics {
	var newURL url.URL
	var baseURL url.URL
	var diags diag.Diagnostics
	newURL.Host = "192.168.64.2:31071"
	newURL.Path = "signup"
	newURL.Scheme = "http"

	newURL.RawQuery = fmt.Sprintf("first_name=%s&surname=%s&email=%s&password=%s", user.FirstName, user.Surname, user.UserEmail, user.Password)
	u := baseURL.ResolveReference(&newURL)
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Go-http-client/1.1")

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	firstName := d.Get("first_name")

	var newURL url.URL
	var BaseURL url.URL

	newURL.Path = fmt.Sprintf("search/%s", firstName)
	newURL.Host = "192.168.64.2:31071"
	newURL.Scheme = "http"

	u := BaseURL.ResolveReference(&newURL)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Go-http-client/1.1")

	var searchedUsers []User
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	ois := make([]interface{}, len(searchedUsers))

	for i, eachuser := range searchedUsers {
		eu := make(map[string]interface{})

		eu["first_name"] = eachuser.FirstName
		eu["surname"] = eachuser.Surname
		eu["email"] = eachuser.UserEmail
		eu["password"] = eachuser.Password

		ois[i] = eu
	}

	for _, i := range ois {
		if err := d.Set("first_name", i.(map[string]interface{})["first_name"]); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("surname", i.(map[string]interface{})["surname"]); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("email", i.(map[string]interface{})["email"]); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("password", i.(map[string]interface{})["password"]); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var updatedUser *User
	var newURL url.URL
	var BaseURL url.URL

	firstName := d.Get("first_name").(string)
	newURL.Path = fmt.Sprintf("users/%s", firstName)
	newURL.Host = "192.168.64.2:31071"
	newURL.Scheme = "http"

	if d.HasChange("surname") {
		newURL.RawQuery = fmt.Sprintf("surname=%s", d.Get("surname").(string))
		updatedUser = &User{
			Surname: d.Get("surname").(string),
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}
	if d.HasChange("email") {
		newURL.RawQuery = fmt.Sprintf("email=%s", d.Get("email").(string))
		updatedUser = &User{
			UserEmail: d.Get("email").(string),
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}
	if d.HasChange("password") {
		newURL.RawQuery = fmt.Sprintf("password=%s", d.Get("password").(string))
		updatedUser = &User{
			Password: d.Get("password").(string),
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	if updatedUser != nil {
		u := BaseURL.ResolveReference(&newURL)
		req, err := http.NewRequest("PUT", u.String(), nil)
		if err != nil {
			return diag.FromErr(err)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "Go-http-client/1.1")

		client := &http.Client{Timeout: 10 * time.Second}

		resp, err := client.Do(req)
		if err != nil {
			return diag.FromErr(err)
		}
		defer resp.Body.Close()
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var newURL url.URL
	var BaseURL url.URL

	firstName := d.Get("first_name")

	newURL.Path = fmt.Sprintf("users/%s", firstName)
	newURL.Host = "192.168.64.2:31071"
	newURL.Scheme = "http"

	u := BaseURL.ResolveReference(&newURL)
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Go-http-client/1.1")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
