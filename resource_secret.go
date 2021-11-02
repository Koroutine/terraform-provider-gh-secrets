package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/nacl/box"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGithubActionsSecretCreateOrUpdate,
		ReadContext:   resourceGithubActionsSecretRead,
		UpdateContext: resourceGithubActionsSecretCreateOrUpdate,
		DeleteContext: resourceGithubActionsSecretDelete,

		Schema: map[string]*schema.Schema{
			"repo": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSecretNameFunc,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGithubActionsSecretCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Meta).client

	repoPath := d.Get("repo").(string)
	secretName := d.Get("name").(string)
	value := d.Get("value").(string)
	owner, repo, err := getDetails(repoPath)
	if err != nil {
		return diag.FromErr(err)
	}

	var encryptedValue string

	keyId, publicKey, err := getPublicKeyDetails(owner, repo, m)
	if err != nil {
		return diag.FromErr(err)
	}

	encryptedBytes, err := encryptPlaintext(value, publicKey)
	if err != nil {
		return diag.FromErr(err)
	}
	encryptedValue = base64.StdEncoding.EncodeToString(encryptedBytes)

	// Create an EncryptedSecret and encrypt the plaintext value into it
	eSecret := &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyId,
		EncryptedValue: encryptedValue,
	}

	_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, eSecret)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildThreePartID(owner, repo, secretName))
	return resourceGithubActionsSecretRead(ctx, d, m)
}

func resourceGithubActionsSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Meta).client

	owner, repoName, secretName, err := parseThreePartID(d.Id(), "owner", "repository", "secret_name")
	if err != nil {
		return diag.FromErr(err)
	}

	secret, _, err := client.Actions.GetRepoSecret(ctx, owner, repoName, secretName)
	if err != nil {
		if ghErr, ok := err.(*github.ErrorResponse); ok {
			if ghErr.Response.StatusCode == http.StatusNotFound {
				log.Printf("[WARN] Removing actions secret %s from state because it no longer exists in GitHub",
					d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	d.Set("value", d.Get("value"))
	d.Set("created_at", secret.CreatedAt.String())

	// This is a drift detection mechanism based on timestamps.
	//
	// If we do not currently store the "updated_at" field, it means we've only
	// just created the resource and the value is most likely what we want it to
	// be.
	//
	// If the resource is changed externally in the meantime then reading back
	// the last update timestamp will return a result different than the
	// timestamp we've persisted in the state. In that case, we can no longer
	// trust that the value (which we don't see) is equal to what we've declared
	// previously.
	//
	// The only solution to enforce consistency between is to mark the resource
	// as deleted (unset the ID) in order to fix potential drift by recreating
	// the resource.
	if updatedAt, ok := d.GetOk("updated_at"); ok && updatedAt != secret.UpdatedAt.String() {
		log.Printf("[WARN] The secret %s has been externally updated in GitHub", d.Id())
		d.SetId("")
	} else if !ok {
		d.Set("updated_at", secret.UpdatedAt.String())
	}

	return nil
}

func resourceGithubActionsSecretDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Meta).client

	owner, repoName, secretName, err := parseThreePartID(d.Id(), "owner", "repository", "secret_name")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting secret: %s", d.Id())
	_, err = client.Actions.DeleteRepoSecret(ctx, owner, repoName, secretName)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getPublicKeyDetails(owner, repository string, meta interface{}) (keyId, pkValue string, err error) {
	client := meta.(*Meta).client
	ctx := context.Background()

	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repository)
	if err != nil {
		return keyId, pkValue, err
	}

	return publicKey.GetKeyID(), publicKey.GetKey(), err
}

func encryptPlaintext(plaintext, publicKeyB64 string) ([]byte, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return nil, err
	}

	var publicKeyBytes32 [32]byte
	copiedLen := copy(publicKeyBytes32[:], publicKeyBytes)
	if copiedLen == 0 {
		return nil, fmt.Errorf("could not convert publicKey to bytes")
	}

	plaintextBytes := []byte(plaintext)
	var encryptedBytes []byte

	cipherText, err := box.SealAnonymous(encryptedBytes, plaintextBytes, &publicKeyBytes32, nil)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

func getDetails(repoPath string) (owner, repoName string, err error) {

	values := strings.Split(repoPath, "/")

	if len(values) != 2 {
		return "", "", fmt.Errorf("bad repo name format: %s. Should be `owner/name`", repoPath)
	}

	return values[0], values[1], err

}
