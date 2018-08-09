package alicloud

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"access_key": {
				Type:        framework.TypeString,
				Description: "Access key with appropriate permissions.",
			},
			"secret_key": {
				Type:        framework.TypeString,
				Description: "Secret key with appropriate permissions.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.operationConfigCreate,
			logical.UpdateOperation: b.operationConfigCreate,
			logical.ReadOperation:   b.operationConfigRead,
			logical.DeleteOperation: b.operationConfigDelete,
		},
		HelpSynopsis:    pathConfigRootHelpSyn,
		HelpDescription: pathConfigRootHelpDesc,
	}
}

func (b *backend) operationConfigCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Access keys and secrets are generated in pairs. You would never need
	// to update one or the other alone, always both together.
	accessKey := ""
	if accessKeyIfc, ok := data.GetOk("access_key"); ok {
		accessKey = accessKeyIfc.(string)
	} else {
		return nil, errors.New("access_key is required")
	}
	secretKey := ""
	if secretKeyIfc, ok := data.GetOk("secret_key"); ok {
		secretKey = secretKeyIfc.(string)
	} else {
		return nil, errors.New("secret_key is required")
	}
	entry, err := logical.StorageEntryJSON("config", credConfig{
		AccessKey: accessKey,
		SecretKey: secretKey,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) operationConfigRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	creds, err := readCredentials(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if creds == nil {
		return nil, nil
	}

	// NOTE:
	// "secret_key" is intentionally not returned by this endpoint,
	// as we lean away from returning sensitive information unless it's absolutely necessary.
	return &logical.Response{
		Data: map[string]interface{}{
			"access_key": creds.AccessKey,
		},
	}, nil
}

func (b *backend) operationConfigDelete(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, "config"); err != nil {
		return nil, err
	}
	return nil, nil
}

func readCredentials(ctx context.Context, s logical.Storage) (*credConfig, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	creds := &credConfig{}
	if err := entry.DecodeJSON(creds); err != nil {
		return nil, err
	}
	return creds, nil
}

type credConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

const pathConfigRootHelpSyn = `
Configure the access key and secret to use for RAM and STS calls.
`

const pathConfigRootHelpDesc = `
Before doing anything, the AliCloud backend needs credentials that are able
to manage RAM users, policies, and access keys, and that can call STS AssumeRole. 
This endpoint is used to configure those credentials.
`
