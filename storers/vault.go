package storers

import (
	"context"
	"strings"

	"github.com/carvers/budget"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	yall "yall.in"
)

type vault struct {
	client *api.Client
	root   string
}

func NewVault(address, root, token string) (vault, error) {
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		return vault{}, errors.Wrap(err, "error creating Vault client")
	}
	client.SetToken(token)
	return vault{
		root:   strings.TrimRight(root, "/"),
		client: client,
	}, nil
}

func stringFromMap(required bool, m map[string]interface{}, key string) (string, error) {
	v, ok := m[key]
	if !ok && required {
		return "", errors.Errorf("%s not in map", key)
	} else if !required && !ok {
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", errors.Errorf("%s not a string", key)
	}
	return s, nil
}

func asdFromMap(m map[string]interface{}) (budget.AccountSensitiveDetails, error) {
	accountID, err := stringFromMap(true, m, "account_id")
	if err != nil {
		return budget.AccountSensitiveDetails{}, err
	}
	bankID, err := stringFromMap(false, m, "bank_id")
	if err != nil {
		return budget.AccountSensitiveDetails{}, err
	}
	userID, err := stringFromMap(true, m, "user_id")
	if err != nil {
		return budget.AccountSensitiveDetails{}, err
	}
	userPass, err := stringFromMap(true, m, "user_pass")
	if err != nil {
		return budget.AccountSensitiveDetails{}, err
	}
	return budget.AccountSensitiveDetails{
		AccountID: accountID,
		BankID:    bankID,
		UserID:    userID,
		UserPass:  userPass,
	}, nil
}

func mapFromASD(asd budget.AccountSensitiveDetails) map[string]interface{} {
	return map[string]interface{}{
		"account_id": asd.AccountID,
		"bank_id":    asd.BankID,
		"user_id":    asd.UserID,
		"user_pass":  asd.UserPass,
	}
}

func (v vault) StoreAccountSensitiveDetails(ctx context.Context, id string, asd budget.AccountSensitiveDetails) error {
	path := v.root + "/" + id
	data := mapFromASD(asd)
	yall.FromContext(ctx).WithField("path", path).WithField("id", id).Debug("writing sensitive account details to Vault")
	_, err := v.client.Logical().Write(path, data)
	return err
}

func (v vault) GetAccountSensitiveDetails(ctx context.Context, id string) (budget.AccountSensitiveDetails, error) {
	path := v.root + "/" + id
	yall.FromContext(ctx).WithField("path", path).WithField("id", id).Debug("reading sensitive account details in Vault")
	resp, err := v.client.Logical().Read(path)
	if err != nil {
		return budget.AccountSensitiveDetails{}, err
	}
	if resp == nil {
		return budget.AccountSensitiveDetails{}, budget.ErrAccountNotFound
	}
	return asdFromMap(resp.Data)
}

func (v vault) DeleteAccountSensitiveDetails(ctx context.Context, id string) error {
	path := v.root + "/" + id
	yall.FromContext(ctx).WithField("path", path).WithField("id", id).Debug("removing sensitive account details from Vault")
	_, err := v.client.Logical().Delete(path)
	return err
}
