package functions_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/functions"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_Snake2CamelFunction(t *testing.T) {
	testCases := []struct {
		name          string
		request       function.RunRequest
		expected      attr.Value
		expectedError *function.FuncError
	}{
		{
			name: "simple-snake-case",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.ObjectValueMust(
						map[string]attr.Type{
							"snake_case": types.StringType,
						},
						map[string]attr.Value{
							"snake_case": types.StringValue("value1"),
						},
					)),
				}),
			},
			expected: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"snakeCase": types.StringType,
				},
				map[string]attr.Value{
					"snakeCase": types.StringValue("value1"),
				},
			)),
		},
		{
			name: "multiple-fields",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.ObjectValueMust(
						map[string]attr.Type{
							"first_name":   types.StringType,
							"last_name":    types.StringType,
							"phone_number": types.StringType,
						},
						map[string]attr.Value{
							"first_name":   types.StringValue("John"),
							"last_name":    types.StringValue("Doe"),
							"phone_number": types.StringValue("123-456-7890"),
						},
					))}),
			},
			expected: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"firstName":   types.StringType,
					"lastName":    types.StringType,
					"phoneNumber": types.StringType,
				},
				map[string]attr.Value{
					"firstName":   types.StringValue("John"),
					"lastName":    types.StringValue("Doe"),
					"phoneNumber": types.StringValue("123-456-7890"),
				},
			)),
		},
		{
			name: "nested-object",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.ObjectValueMust(
						map[string]attr.Type{
							"outer_field": types.StringType,
							"nested_object": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"inner_field": types.StringType,
								},
							},
						},
						map[string]attr.Value{
							"outer_field": types.StringValue("outerValue"),
							"nested_object": types.ObjectValueMust(
								map[string]attr.Type{
									"inner_field": types.StringType,
								},
								map[string]attr.Value{
									"inner_field": types.StringValue("innerValue"),
								},
							),
						},
					))}),
			},
			expected: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"outerField": types.StringType,
					"nestedObject": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"innerField": types.StringType,
						},
					},
				},
				map[string]attr.Value{
					"outerField": types.StringValue("outerValue"),
					"nestedObject": types.ObjectValueMust(
						map[string]attr.Type{
							"innerField": types.StringType,
						},
						map[string]attr.Value{
							"innerField": types.StringValue("innerValue"),
						},
					),
				},
			)),
		},
		{
			name: "already-camelcase",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.ObjectValueMust(
						map[string]attr.Type{
							"alreadycamel": types.StringType,
						},
						map[string]attr.Value{
							"alreadycamel": types.StringValue("value"),
						},
					))}),
			},
			expected: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"alreadycamel": types.StringType,
				},
				map[string]attr.Value{
					"alreadycamel": types.StringValue("value"),
				},
			)),
		},
		{
			name: "multiple-underscores",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.ObjectValueMust(
						map[string]attr.Type{
							"field_with_many_words": types.StringType,
						},
						map[string]attr.Value{
							"field_with_many_words": types.StringValue("value"),
						},
					))}),
			},
			expected: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"fieldWithManyWords": types.StringType,
				},
				map[string]attr.Value{
					"fieldWithManyWords": types.StringValue("value"),
				},
			)),
		},
		{
			name: "different-value-types",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.ObjectValueMust(
						map[string]attr.Type{
							"string_field": types.StringType,
							"number_field": types.NumberType,
							"bool_field":   types.BoolType,
						},
						map[string]attr.Value{
							"string_field": types.StringValue("text"),
							"number_field": types.NumberValue(big.NewFloat(42)),
							"bool_field":   types.BoolValue(true),
						},
					))}),
			},
			expected: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"stringField": types.StringType,
					"numberField": types.NumberType,
					"boolField":   types.BoolType,
				},
				map[string]attr.Value{
					"stringField": types.StringValue("text"),
					"numberField": types.NumberValue(big.NewFloat(42)),
					"boolField":   types.BoolValue(true),
				},
			)),
		},
		{
			name: "map-input-simple",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"string_field": types.StringValue("text"),
							"number_field": types.StringValue("42"),
							"bool_field":   types.StringValue("true"),
						},
					))}),
			},
			expected: types.DynamicValue(types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"stringField": types.StringValue("text"),
					"numberField": types.StringValue("42"),
					"boolField":   types.StringValue("true"),
				},
			)),
		},
		{
			name: "map-of-objects-with-nested-objects",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.MapValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"user_name": types.StringType,
								"user_details": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"email_address": types.StringType,
										"phone_number":  types.StringType,
									},
								},
							},
						},
						map[string]attr.Value{
							"first_user": types.ObjectValueMust(
								map[string]attr.Type{
									"user_name": types.StringType,
									"user_details": types.ObjectType{
										AttrTypes: map[string]attr.Type{
											"email_address": types.StringType,
											"phone_number":  types.StringType,
										},
									},
								},
								map[string]attr.Value{
									"user_name": types.StringValue("john_doe"),
									"user_details": types.ObjectValueMust(
										map[string]attr.Type{
											"email_address": types.StringType,
											"phone_number":  types.StringType,
										},
										map[string]attr.Value{
											"email_address": types.StringValue("john@example.com"),
											"phone_number":  types.StringValue("123-456-7890"),
										},
									),
								},
							),
							"second_user": types.ObjectValueMust(
								map[string]attr.Type{
									"user_name": types.StringType,
									"user_details": types.ObjectType{
										AttrTypes: map[string]attr.Type{
											"email_address": types.StringType,
											"phone_number":  types.StringType,
										},
									},
								},
								map[string]attr.Value{
									"user_name": types.StringValue("jane_smith"),
									"user_details": types.ObjectValueMust(
										map[string]attr.Type{
											"email_address": types.StringType,
											"phone_number":  types.StringType,
										},
										map[string]attr.Value{
											"email_address": types.StringValue("jane@example.com"),
											"phone_number":  types.StringValue("098-765-4321"),
										},
									),
								},
							),
						},
					))}),
			},
			expected: types.DynamicValue(types.MapValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"userName": types.StringType,
						"userDetails": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"emailAddress": types.StringType,
								"phoneNumber":  types.StringType,
							},
						},
					},
				},
				map[string]attr.Value{
					"firstUser": types.ObjectValueMust(
						map[string]attr.Type{
							"userName": types.StringType,
							"userDetails": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"emailAddress": types.StringType,
									"phoneNumber":  types.StringType,
								},
							},
						},
						map[string]attr.Value{
							"userName": types.StringValue("john_doe"),
							"userDetails": types.ObjectValueMust(
								map[string]attr.Type{
									"emailAddress": types.StringType,
									"phoneNumber":  types.StringType,
								},
								map[string]attr.Value{
									"emailAddress": types.StringValue("john@example.com"),
									"phoneNumber":  types.StringValue("123-456-7890"),
								},
							),
						},
					),
					"secondUser": types.ObjectValueMust(
						map[string]attr.Type{
							"userName": types.StringType,
							"userDetails": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"emailAddress": types.StringType,
									"phoneNumber":  types.StringType,
								},
							},
						},
						map[string]attr.Value{
							"userName": types.StringValue("jane_smith"),
							"userDetails": types.ObjectValueMust(
								map[string]attr.Type{
									"emailAddress": types.StringType,
									"phoneNumber":  types.StringType,
								},
								map[string]attr.Value{
									"emailAddress": types.StringValue("jane@example.com"),
									"phoneNumber":  types.StringValue("098-765-4321"),
								},
							),
						},
					),
				},
			)),
		},
		{
			name: "dynamic-unknown",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicUnknown(),
				}),
			},
			expected: types.DynamicUnknown(),
		},
		{
			name: "dynamic-with-unknown-object",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.ObjectUnknown(
						map[string]attr.Type{
							"snake_case_field": types.StringType,
							"another_field":    types.NumberType,
						},
					)),
				}),
			},
			expected: types.DynamicValue(types.ObjectUnknown(
				map[string]attr.Type{
					"snakeCaseField": types.StringType,
					"anotherField":   types.NumberType,
				},
			)),
		},
		{
			name: "partially-unknown-object",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.ObjectValueMust(
						map[string]attr.Type{
							"known_field":   types.StringType,
							"unknown_field": types.StringType,
							"nested_object": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"inner_known":   types.StringType,
									"inner_unknown": types.NumberType,
								},
							},
						},
						map[string]attr.Value{
							"known_field":   types.StringValue("known_value"),
							"unknown_field": types.StringUnknown(),
							"nested_object": types.ObjectValueMust(
								map[string]attr.Type{
									"inner_known":   types.StringType,
									"inner_unknown": types.NumberType,
								},
								map[string]attr.Value{
									"inner_known":   types.StringValue("nested_known"),
									"inner_unknown": types.NumberUnknown(),
								},
							),
						},
					)),
				}),
			},
			expected: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"knownField":   types.StringType,
					"unknownField": types.StringType,
					"nestedObject": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"innerKnown":   types.StringType,
							"innerUnknown": types.NumberType,
						},
					},
				},
				map[string]attr.Value{
					"knownField":   types.StringValue("known_value"),
					"unknownField": types.StringUnknown(),
					"nestedObject": types.ObjectValueMust(
						map[string]attr.Type{
							"innerKnown":   types.StringType,
							"innerUnknown": types.NumberType,
						},
						map[string]attr.Value{
							"innerKnown":   types.StringValue("nested_known"),
							"innerUnknown": types.NumberUnknown(),
						},
					),
				},
			)),
		},
		{
			name: "deeply-nested-structure",
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.DynamicValue(types.ObjectValueMust(
						map[string]attr.Type{
							"company_name": types.StringType,
							"employee_map": types.MapType{
								ElemType: types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"full_name":    types.StringType,
										"job_position": types.StringType,
										"contact_info": types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"email_address": types.StringType,
												"home_address": types.ObjectType{
													AttrTypes: map[string]attr.Type{
														"street_name": types.StringType,
														"city_name":   types.StringType,
														"zip_code":    types.StringType,
													},
												},
											},
										},
									},
								},
							},
						},
						map[string]attr.Value{
							"company_name": types.StringValue("Tech Corp"),
							"employee_map": types.MapValueMust(
								types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"full_name":    types.StringType,
										"job_position": types.StringType,
										"contact_info": types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"email_address": types.StringType,
												"home_address": types.ObjectType{
													AttrTypes: map[string]attr.Type{
														"street_name": types.StringType,
														"city_name":   types.StringType,
														"zip_code":    types.StringType,
													},
												},
											},
										},
									},
								},
								map[string]attr.Value{
									"emp_001": types.ObjectValueMust(
										map[string]attr.Type{
											"full_name":    types.StringType,
											"job_position": types.StringType,
											"contact_info": types.ObjectType{
												AttrTypes: map[string]attr.Type{
													"email_address": types.StringType,
													"home_address": types.ObjectType{
														AttrTypes: map[string]attr.Type{
															"street_name": types.StringType,
															"city_name":   types.StringType,
															"zip_code":    types.StringType,
														},
													},
												},
											},
										},
										map[string]attr.Value{
											"full_name":    types.StringValue("Alice Johnson"),
											"job_position": types.StringValue("Senior Engineer"),
											"contact_info": types.ObjectValueMust(
												map[string]attr.Type{
													"email_address": types.StringType,
													"home_address": types.ObjectType{
														AttrTypes: map[string]attr.Type{
															"street_name": types.StringType,
															"city_name":   types.StringType,
															"zip_code":    types.StringType,
														},
													},
												},
												map[string]attr.Value{
													"email_address": types.StringValue("alice@example.com"),
													"home_address": types.ObjectValueMust(
														map[string]attr.Type{
															"street_name": types.StringType,
															"city_name":   types.StringType,
															"zip_code":    types.StringType,
														},
														map[string]attr.Value{
															"street_name": types.StringValue("123 Main St"),
															"city_name":   types.StringValue("Seattle"),
															"zip_code":    types.StringValue("98101"),
														},
													),
												},
											),
										},
									),
								},
							),
						},
					))}),
			},
			expected: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"companyName": types.StringType,
					"employeeMap": types.MapType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"fullName":    types.StringType,
								"jobPosition": types.StringType,
								"contactInfo": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"emailAddress": types.StringType,
										"homeAddress": types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"streetName": types.StringType,
												"cityName":   types.StringType,
												"zipCode":    types.StringType,
											},
										},
									},
								},
							},
						},
					},
				},
				map[string]attr.Value{
					"companyName": types.StringValue("Tech Corp"),
					"employeeMap": types.MapValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"fullName":    types.StringType,
								"jobPosition": types.StringType,
								"contactInfo": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"emailAddress": types.StringType,
										"homeAddress": types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"streetName": types.StringType,
												"cityName":   types.StringType,
												"zipCode":    types.StringType,
											},
										},
									},
								},
							},
						},
						map[string]attr.Value{
							"emp001": types.ObjectValueMust(
								map[string]attr.Type{
									"fullName":    types.StringType,
									"jobPosition": types.StringType,
									"contactInfo": types.ObjectType{
										AttrTypes: map[string]attr.Type{
											"emailAddress": types.StringType,
											"homeAddress": types.ObjectType{
												AttrTypes: map[string]attr.Type{
													"streetName": types.StringType,
													"cityName":   types.StringType,
													"zipCode":    types.StringType,
												},
											},
										},
									},
								},
								map[string]attr.Value{
									"fullName":    types.StringValue("Alice Johnson"),
									"jobPosition": types.StringValue("Senior Engineer"),
									"contactInfo": types.ObjectValueMust(
										map[string]attr.Type{
											"emailAddress": types.StringType,
											"homeAddress": types.ObjectType{
												AttrTypes: map[string]attr.Type{
													"streetName": types.StringType,
													"cityName":   types.StringType,
													"zipCode":    types.StringType,
												},
											},
										},
										map[string]attr.Value{
											"emailAddress": types.StringValue("alice@example.com"),
											"homeAddress": types.ObjectValueMust(
												map[string]attr.Type{
													"streetName": types.StringType,
													"cityName":   types.StringType,
													"zipCode":    types.StringType,
												},
												map[string]attr.Value{
													"streetName": types.StringValue("123 Main St"),
													"cityName":   types.StringValue("Seattle"),
													"zipCode":    types.StringValue("98101"),
												},
											),
										},
									),
								},
							),
						},
					),
				},
			)),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := function.RunResponse{
				Result: function.NewResultData(types.DynamicUnknown()),
			}

			snake2CamelFunction := functions.Snake2CamelFunction{}
			snake2CamelFunction.Run(context.Background(), testCase.request, &got)

			if testCase.expectedError != nil {
				if got.Error == nil {
					t.Fatal("expected error, got none")
				}
				if diff := cmp.Diff(got.Error, testCase.expectedError); diff != "" {
					t.Errorf("unexpected error difference: %s", diff)
				}
				return
			}

			if got.Error != nil {
				t.Fatalf("unexpected error: %s", got.Error)
			}

			result := got.Result.Value()

			if !testCase.expected.Equal(result) {
				t.Errorf("expected %v, got %v", testCase.expected, result)
			}
		})
	}
}
