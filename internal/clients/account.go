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
	subscriptionId *string
	mutex          sync.Mutex
}

func NewResourceManagerAccount(subscriptionId string) ResourceManagerAccount {
	if subscriptionId == "" {
		return ResourceManagerAccount{}
	}
	return ResourceManagerAccount{
		subscriptionId: &subscriptionId,
	}
}

func (account ResourceManagerAccount) GetSubscriptionId() string {
	account.mutex.Lock()
	defer account.mutex.Unlock()
	if account.subscriptionId != nil {
		return *account.subscriptionId
	}

	subscriptionId, err := getDefaultSubscriptionID()
	if err != nil {
		log.Printf("[DEBUG] Error getting default subscription ID: %s", err)
	}

	account.subscriptionId = &subscriptionId

	return *account.subscriptionId
}

// getDefaultSubscriptionID tries to determine the default subscription
func getDefaultSubscriptionID() (string, error) {
	var account struct {
		SubscriptionID string `json:"id"`
	}
	err := jsonUnmarshalAzCmd(&account, "account", "show")
	if err != nil {
		return "", fmt.Errorf("obtaining subscription ID: %s", err)
	}

	return account.SubscriptionID, nil
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
