package tffwdocs

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func diagToError(d diag.Diagnostic) error {
	if d.Severity() != diag.SeverityError {
		return nil
	}
	return fmt.Errorf("%s: %s", d.Summary(), d.Detail())
}

func diagsToError(diags diag.Diagnostics) error {
	var errs []error

	for _, ediag := range diags.Errors() {
		errs = append(errs, diagToError(ediag))
	}
	return errors.Join(errs...)
}
