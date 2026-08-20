// Package schema exposes resource schema capabilities backed by AzAPI's embedded Bicep type definitions.
//
// The embedded definitions are generated from Bicep types and periodically synchronized with
// azure-rest-api-specs. This package intentionally exposes capability results rather than the
// generated schema model so consumers do not depend on implementation details.
package schema
