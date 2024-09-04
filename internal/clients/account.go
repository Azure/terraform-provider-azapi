package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type ResourceManagerAccount struct {
	tenantId       *string
	subscriptionId *string
	mutex          *sync.Mutex
}

func NewResourceManagerAccount(tenantId, subscriptionId string) ResourceManagerAccount {
	out := ResourceManagerAccount{
		mutex: &sync.Mutex{},
	}
	if tenantId != "" {
		out.tenantId = &tenantId
	}
	if subscriptionId != "" {
		out.subscriptionId = &subscriptionId
	}
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
		log.Printf("[DEBUG] No default subscription ID found")
		return ""
	}
	return *account.subscriptionId
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
