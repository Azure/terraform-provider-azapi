package tffwdocs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/magodo/terraform-plugin-framework-docs/internal/metadata"
)

type Generator struct {
	metadata metadata.Metadata
}

// Create a Generator by analyzing the schema of the provider and all the registered resources to the provider.
func NewGenerator(ctx context.Context, p provider.Provider) (*Generator, error) {
	metadata, diags := metadata.GetMetadata(ctx, p)
	if diags.HasError() {
		return nil, diagsToError(diags)
	}

	return &Generator{metadata: metadata}, nil
}

type Example = metadata.Example
type ImportId = metadata.ImportId
type ProviderRenderOption = metadata.ProviderRenderOption
type ResourceRenderOption = metadata.ResourceRenderOption
type DataSourceRenderOption = metadata.DataSourceRenderOption
type EphemeralResourceRenderOption = metadata.EphemeralRenderOption
type ActionRenderOption = metadata.ActionRenderOption
type ListResourceRenderOption = metadata.ListRenderOption
type FunctionRenderOption = metadata.FunctionRenderOption

type RenderOptions struct {
	Provider           *ProviderRenderOption
	Resources          map[string]ResourceRenderOption
	DataSources        map[string]DataSourceRenderOption
	EphemeralResources map[string]EphemeralResourceRenderOption
	ListResources      map[string]ListResourceRenderOption
	Actions            map[string]ActionRenderOption
	Functions          map[string]FunctionRenderOption
}

type LintOptions struct {
	SkipProvider           bool
	SkipResources          map[string]bool
	SkipDataSources        map[string]bool
	SkipEphemeralResources map[string]bool
	SkipListResources      map[string]bool
	SkipActions            map[string]bool
	SkipFunctions          map[string]bool
}

// Lint traverses the provider and provider resource type's schema to identify any field
// without a description specified.
func (gen Generator) Lint(opt *LintOptions) error {
	if opt == nil {
		opt = &LintOptions{
			SkipProvider:           false,
			SkipResources:          map[string]bool{},
			SkipDataSources:        map[string]bool{},
			SkipEphemeralResources: map[string]bool{},
			SkipListResources:      map[string]bool{},
			SkipActions:            map[string]bool{},
			SkipFunctions:          map[string]bool{},
		}
	}
	var errs []error

	if !opt.SkipProvider {
		schema := gen.metadata.Provider.Schema
		if schema.Description == "" {
			errs = append(errs, errors.New("Provider has no description"))
		}
		if err := schema.Fields.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Provider:\n%v", err))
		}
		if err := schema.Nested.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Provider:\n%v", err))
		}
	}

	for name, res := range gen.metadata.Resources {
		if opt.SkipResources[name] {
			continue
		}
		if res.Schema.Description == "" {
			errs = append(errs, fmt.Errorf("Resource %q has no description", name))
		}
		if err := res.Schema.Fields.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Resource %q:\n%v", name, err))
		}
		if err := res.Schema.Nested.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Resource %q:\n%v", name, err))
		}
		if res.Identity != nil {
			if err := res.Identity.Fields.Lint(); err != nil {
				errs = append(errs, fmt.Errorf("Resource %q Identity:\n%v", name, err))
			}
		}
	}
	for name, ds := range gen.metadata.DataSources {
		if opt.SkipDataSources[name] {
			continue
		}
		if ds.Schema.Description == "" {
			errs = append(errs, fmt.Errorf("Data Source %q has no description", name))
		}
		if err := ds.Schema.Fields.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Data Source %q:\n%v", name, err))
		}
		if err := ds.Schema.Nested.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Data Source %q:\n%v", name, err))
		}
	}
	for name, res := range gen.metadata.Ephemerals {
		if opt.SkipEphemeralResources[name] {
			continue
		}
		if res.Schema.Description == "" {
			errs = append(errs, errors.New("Ephemeral Resource %q has no description"))
		}
		if err := res.Schema.Fields.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Ephemeral Resource %q:\n%v", name, err))
		}
		if err := res.Schema.Nested.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Ephemeral Resource %q:\n%v", name, err))
		}
	}
	for name, res := range gen.metadata.Lists {
		if opt.SkipListResources[name] {
			continue
		}
		if res.Schema.Description == "" {
			errs = append(errs, fmt.Errorf("List Resource %q has no description", name))
		}
		if err := res.Schema.Fields.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("List Resource %q:\n%v", name, err))
		}
		if err := res.Schema.Nested.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("List Resource %q:\n%v", name, err))
		}
	}
	for name, act := range gen.metadata.Actions {
		if opt.SkipActions[name] {
			continue
		}
		if act.Schema.Description == "" {
			errs = append(errs, fmt.Errorf("Action %q has no description", name))
		}
		if err := act.Schema.Fields.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Action %q:\n%v", name, err))
		}
		if err := act.Schema.Nested.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Action %q:\n%v", name, err))
		}
	}
	for name, fun := range gen.metadata.Functions {
		if opt.SkipFunctions[name] {
			continue
		}
		if fun.Schema.Description == "" {
			errs = append(errs, fmt.Errorf("Function %q has no description", name))
		}
		if err := fun.Schema.Parameters.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Function %q parameter:\n%v", name, err))
		}
		if err := fun.Schema.Objects.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Function %q parameter nested objects:\n%v", name, err))
		}
		if err := fun.Schema.ReturnObjects.Lint(); err != nil {
			errs = append(errs, fmt.Errorf("Function %q return nested objects:\n%v", name, err))
		}
	}

	return errors.Join(errs...)
}

// RenderProvider renders the provider document.
// If option is specified, it will override the ProviderWithRenderOption.RenderOption() if is implemented.
func (gen Generator) RenderProvider(ctx context.Context, w io.Writer, option *ProviderRenderOption) error {
	rr, err := gen.metadata.NewProviderRender(option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

// RenderResource renders the resource document.
// If option is specified, it will override the ResourceWithRenderOption.RenderOption() if is implemented.
func (gen Generator) RenderResource(ctx context.Context, w io.Writer, resourceType string, option *ResourceRenderOption) error {
	rr, err := gen.metadata.NewResourceRender(resourceType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

// RenderDataSource renders the data source document.
// If option is specified, it will override the DataSourceWithRenderOption.RenderOption() if is implemented.
func (gen Generator) RenderDataSource(ctx context.Context, w io.Writer, dataSourceType string, option *DataSourceRenderOption) error {
	rr, err := gen.metadata.NewDataSourceRender(dataSourceType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

// RenderEphemeralResource renders the ephemeral resource document.
// If option is specified, it will override the EphemeralResourceWithRenderOption.RenderOption() if is implemented.
func (gen Generator) RenderEphemeralResource(ctx context.Context, w io.Writer, ephemeralResourceType string, option *EphemeralResourceRenderOption) error {
	rr, err := gen.metadata.NewEphemeralRender(ephemeralResourceType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

// RenderAction renders the action document.
// If option is specified, it will override the ActionWithRenderOption.RenderOption() if is implemented.
func (gen Generator) RenderAction(ctx context.Context, w io.Writer, actionType string, option *ActionRenderOption) error {
	rr, err := gen.metadata.NewActionRender(actionType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

// RenderListResource renders the list resource document.
// If option is specified, it will override the ListResourceWithRenderOption.RenderOption() if is implemented.
func (gen Generator) RenderListResource(ctx context.Context, w io.Writer, listResourceType string, option *ListResourceRenderOption) error {
	rr, err := gen.metadata.NewListRender(listResourceType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

// RenderFunction renders the function document.
// If option is specified, it will override the FunctionWithRenderOption.RenderOption() if is implemented.
func (gen Generator) RenderFunction(ctx context.Context, w io.Writer, functionName string, option *FunctionRenderOption) error {
	rr, err := gen.metadata.NewFunctionRender(functionName, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

// WriteAll generates the documents for the provider and all registered resources, actions, functions, etc.
// and write to the specified document directory, which should be existing.
func (gen Generator) WriteAll(ctx context.Context, docDir string, opts *RenderOptions) error {
	if opts == nil {
		opts = &RenderOptions{
			Provider:           nil,
			Resources:          map[string]ResourceRenderOption{},
			DataSources:        map[string]DataSourceRenderOption{},
			EphemeralResources: map[string]EphemeralResourceRenderOption{},
			ListResources:      map[string]ListResourceRenderOption{},
			Actions:            map[string]ActionRenderOption{},
			Functions:          map[string]FunctionRenderOption{},
		}
	}

	// Validate render options
	for name := range opts.Resources {
		if _, ok := gen.metadata.Resources[name]; !ok {
			return fmt.Errorf("invalid render option: unknown resource type: %v", name)
		}
	}
	for name := range opts.DataSources {
		if _, ok := gen.metadata.DataSources[name]; !ok {
			return fmt.Errorf("invalid render option: unknown data source type: %v", name)
		}
	}
	for name := range opts.EphemeralResources {
		if _, ok := gen.metadata.Ephemerals[name]; !ok {
			return fmt.Errorf("invalid render option: unknown ephemeral resource type: %v", name)
		}
	}
	for name := range opts.ListResources {
		if _, ok := gen.metadata.Lists[name]; !ok {
			return fmt.Errorf("invalid render option: unknown list resource type: %v", name)
		}
	}
	for name := range opts.Actions {
		if _, ok := gen.metadata.Actions[name]; !ok {
			return fmt.Errorf("invalid render option: unknown action type: %v", name)
		}
	}
	for name := range opts.Functions {
		if _, ok := gen.metadata.Functions[name]; !ok {
			return fmt.Errorf("invalid render option: unknown function type: %v", name)
		}
	}

	// Provider
	{
		var buf bytes.Buffer
		if err := gen.RenderProvider(ctx, &buf, opts.Provider); err != nil {
			return fmt.Errorf("render provider: %v", err)
		}
		if err := os.WriteFile(filepath.Join(docDir, "index.md"), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write provider file: %v", err)
		}
	}

	providerPrefix := gen.metadata.ProviderName + "_"

	// Resources
	for name := range gen.metadata.Resources {
		var buf bytes.Buffer
		var opt *ResourceRenderOption
		if optt, ok := opts.Resources[name]; ok {
			opt = &optt
		}
		if err := gen.RenderResource(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render resource %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "resources")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", strings.TrimPrefix(name, providerPrefix))), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write resource file for %q: %v", name, err)
		}
	}

	// DataSource
	for name := range gen.metadata.DataSources {
		var buf bytes.Buffer
		var opt *DataSourceRenderOption
		if optt, ok := opts.DataSources[name]; ok {
			opt = &optt
		}
		if err := gen.RenderDataSource(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render data source %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "data-sources")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", strings.TrimPrefix(name, providerPrefix))), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write data source file for %q: %v", name, err)
		}
	}

	// Ephemeral Resource
	for name := range gen.metadata.Ephemerals {
		var buf bytes.Buffer
		var opt *EphemeralResourceRenderOption
		if optt, ok := opts.EphemeralResources[name]; ok {
			opt = &optt
		}
		if err := gen.RenderEphemeralResource(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render ephemeral resource %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "ephemeral-resources")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", strings.TrimPrefix(name, providerPrefix))), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write ephemeral resource file for %q: %v", name, err)
		}
	}

	// List Resource
	for name := range gen.metadata.Lists {
		var buf bytes.Buffer
		var opt *ListResourceRenderOption
		if optt, ok := opts.ListResources[name]; ok {
			opt = &optt
		}
		if err := gen.RenderListResource(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render list resource %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "list-resources")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", strings.TrimPrefix(name, providerPrefix))), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write list resource file for %q: %v", name, err)
		}
	}

	// Action
	for name := range gen.metadata.Actions {
		var buf bytes.Buffer
		var opt *ActionRenderOption
		if optt, ok := opts.Actions[name]; ok {
			opt = &optt
		}
		if err := gen.RenderAction(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render action %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "actions")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", strings.TrimPrefix(name, providerPrefix))), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write action file for %q: %v", name, err)
		}
	}

	// Function
	for name := range gen.metadata.Functions {
		var buf bytes.Buffer
		var opt *FunctionRenderOption
		if optt, ok := opts.Functions[name]; ok {
			opt = &optt
		}
		if err := gen.RenderFunction(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render function %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "functions")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", strings.TrimPrefix(name, providerPrefix))), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write function file for %q: %v", name, err)
		}
	}

	return nil
}
