package test_helper

import (
	"context"
	"fmt"
	"time"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ statecheck.StateCheck = expectTimeFormat{}

type expectTimeFormat struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
	layout          string
}

func (e expectTimeFormat) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	var resource *tfjson.StateResource

	if req.State == nil {
		resp.Error = fmt.Errorf("state is nil")
	}

	if req.State.Values == nil {
		resp.Error = fmt.Errorf("state does not contain any state values")
	}

	if req.State.Values.RootModule == nil {
		resp.Error = fmt.Errorf("state does not contain a root module")
	}

	for _, r := range req.State.Values.RootModule.Resources {
		if e.resourceAddress == r.Address {
			resource = r

			break
		}
	}

	if resource == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddress)

		return
	}

	result, err := tfjsonpath.Traverse(resource.AttributeValues, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	if _, err := time.Parse(e.layout, result.(string)); err != nil {
		resp.Error = err
	}
}

func ExpectTimeFormat(resourceAddress string, attributePath tfjsonpath.Path, layout string) statecheck.StateCheck {
	return expectTimeFormat{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		layout:          layout,
	}
}
