package acceptance_test

import (
	"context"
	"fmt"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/acceptance"
	"testing"
)

func Test_1(t *testing.T) {
	clients, err := acceptance.BuildTestClient()
	if err != nil {
		t.Fatal(err)
	}
	body, resp, err := clients.ResourceClient.Get(context.TODO(), "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctestRG-henglu1116/providers/Microsoft.Sql/servers/acctesthenglu1116", "2017-03-01-preview")
	if err != nil {
		return
	}
	fmt.Println(body, resp)
}
