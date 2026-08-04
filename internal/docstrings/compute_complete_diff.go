package docstrings

const (
	computeCompleteDiffStr = `Whether to compute a complete diff that surfaces writable properties which exist in the remote resource but are missing from %sbody%s. When enabled, out-of-band changes to such properties (for example, a subnet delegation added directly in Azure) are detected as drift instead of being silently ignored. This takes precedence over %signore_missing_property%s. Read-only properties and service-applied default (zero) values are not surfaced, to keep the plan convergent.`
)

// ComputeCompleteDiff returns the docstring for the compute_complete_diff schema attribute.
func ComputeCompleteDiff() string {
	return addBackquotes(computeCompleteDiffStr)
}
