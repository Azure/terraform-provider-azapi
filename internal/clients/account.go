package clients

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type ObjectIDProvider func(ctx context.Context) (string, error)

type ResourceManagerAccount struct {
	tenantId         *string
	subscriptionId   *string
	objectId         *string
	mutex            *sync.Mutex
	objectIDProvider ObjectIDProvider
}

func NewResourceManagerAccount(tenantId, subscriptionId string, provider ObjectIDProvider) ResourceManagerAccount {
	out := ResourceManagerAccount{
		mutex: &sync.Mutex{},
	}
	if tenantId != "" {
		out.tenantId = &tenantId
	}
	if subscriptionId != "" {
		out.subscriptionId = &subscriptionId
	}
	// We lazy load object ID because it's not always needed and could cause a performance hit
	out.objectIDProvider = provider
	return out
}

func (account *ResourceManagerAccount) GetTenantId() string {
	account.mutex.Lock()
	defer account.mutex.Unlock()

	if account.tenantId != nil {
		return *account.tenantId
	}

	err := account.loadDefaultsFromAzCmd()
	if err != nil {
		log.Printf("[DEBUG] Error getting default tenant ID: %s", err)
	}

	if account.tenantId == nil {
		log.Printf("[DEBUG] No default tenant ID found")
		return ""
	}

	return *account.tenantId
}

func (account *ResourceManagerAccount) GetSubscriptionId() string {
	account.mutex.Lock()
	defer account.mutex.Unlock()

	if account.subscriptionId != nil {
		return *account.subscriptionId
	}

	err := account.loadDefaultsFromAzCmd()
	if err != nil {
		log.Printf("[DEBUG] Error getting default subscription ID: %s", err)
	}

	if account.subscriptionId == nil {
		log.Printf("[DEBUG] No subscription ID found")
		return ""
	}

	return *account.subscriptionId
}

func (account *ResourceManagerAccount) GetObjectId(ctx context.Context) string {
	account.mutex.Lock()
	defer account.mutex.Unlock()

	if account.objectId != nil {
		return *account.objectId
	}

	if account.objectIDProvider != nil {
		objectId, err := account.objectIDProvider(ctx)
		if err != nil {
			log.Printf("[DEBUG] Error getting object ID: %s", err)
		}
		if objectId != "" {
			account.objectId = &objectId
			return *account.objectId
		}
	}

	err := account.loadSignedInUserFromAzCmd()
	if err != nil {
		log.Printf("[DEBUG] Error getting user object ID from az cli: %s", err)
	}
	if account.objectId == nil {
		log.Printf("[DEBUG] No object ID found")
		return ""
	}

	return *account.objectId
}

func (account *ResourceManagerAccount) loadSignedInUserFromAzCmd() error {
	var userModel struct {
		ObjectId string `json:"id"`
	}
	err := jsonUnmarshalAzCmd(&userModel, "ad", "signed-in-user", "show")
	if err != nil {
		return fmt.Errorf("obtaining defaults from az cmd: %s", err)
	}

	account.objectId = &userModel.ObjectId
	return nil
}

func (account *ResourceManagerAccount) loadDefaultsFromAzCmd() error {
	var accountModel struct {
		SubscriptionID string `json:"id"`
		TenantId       string `json:"tenantId"`
	}
	err := jsonUnmarshalAzCmd(&accountModel, "account", "show")
	if err != nil {
		return fmt.Errorf("obtaining defaults from az cmd: %s", err)
	}

	account.tenantId = &accountModel.TenantId
	account.subscriptionId = &accountModel.SubscriptionID
	return nil
}

// jsonUnmarshalAzCmd executes an Azure CLI command and unmarshalls the JSON output.
func jsonUnmarshalAzCmd(i interface{}, arg ...string) error {
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	arg = append(arg, "-o=json")
	cmd := exec.Command("az", arg...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		err := fmt.Errorf("launching Azure CLI: %+v", err)
		if stdErrStr := stderr.String(); stdErrStr != "" {
			err = fmt.Errorf("%s: %s", err, strings.TrimSpace(stdErrStr))
		}
		return err
	}

	if err := cmd.Wait(); err != nil {
		err := fmt.Errorf("running Azure CLI: %+v", err)
		if stdErrStr := stderr.String(); stdErrStr != "" {
			err = fmt.Errorf("%s: %s", err, strings.TrimSpace(stdErrStr))
		}
		return err
	}

	if err := json.Unmarshal(stdout.Bytes(), &i); err != nil {
		return fmt.Errorf("unmarshaling the output of Azure CLI: %v", err)
	}

	return nil
}

func parseTokenClaims(token string) (*tokenClaims, error) {
	// Parse the token to get the claims
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("parseTokenClaims: token does not have 3 parts")
	}
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("parseTokenClaims: error decoding token: %s", err)
	}
	var claims tokenClaims
	err = json.Unmarshal(decoded, &claims)
	if err != nil {
		return nil, fmt.Errorf("parseTokenClaims: error unmarshalling claims: %w", err)
	}
	return &claims, nil
}

type tokenClaims struct {
	Audience          string   `json:"aud"`
	Expires           int64    `json:"exp"`
	IssuedAt          int64    `json:"iat"`
	Issuer            string   `json:"iss"`
	IdentityProvider  string   `json:"idp"`
	ObjectId          string   `json:"oid"`
	Roles             []string `json:"roles"`
	Scopes            string   `json:"scp"`
	Subject           string   `json:"sub"`
	TenantRegionScope string   `json:"tenant_region_scope"`
	TenantId          string   `json:"tid"`
	Version           string   `json:"ver"`

	AppDisplayName string `json:"app_displayname,omitempty"`
	AppId          string `json:"appid,omitempty"`
	IdType         string `json:"idtyp,omitempty"`
}

func ParsedTokenClaimsObjectIDProvider(cred azcore.TokenCredential, cloudCfg cloud.Configuration) ObjectIDProvider {
	return func(ctx context.Context) (string, error) {
		tok, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
			EnableCAE: true,
			Scopes:    []string{cloudCfg.Services[cloud.ResourceManager].Audience + "/.default"}})
		if err != nil {
			return "", fmt.Errorf("getting requesting token from credentials: %w", err)
		}
		if tok.Token == "" {
			return "", errors.New("token is empty")
		}
		cl, err := parseTokenClaims(tok.Token)
		if err != nil {
			return "", fmt.Errorf("getting object id from token: %w", err)
		}
		if cl == nil || cl.ObjectId == "" {
			return "", errors.New("object id is empty")
		}
		return cl.ObjectId, nil
	}
}
