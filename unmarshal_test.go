package additionaljson

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUnmarshal_HappyCase_All(t *testing.T) {
	testStruct := struct {
		TestString1 string            `json:"test_string_1"`
		All         map[string]string `json:"-" aj:"all"`
		AllAJOnly   map[string]string `aj:"all"`
	}{}

	err := DefaultUnmarshaler.Unmarshal([]byte(`{"test_string_1": "11", "test2": "BB", "test3": "CC"}`), &testStruct)
	if err != nil {
		t.Fatalf("Faild to unmarshal json cause: %s", err)
	}

	expectedTestString1 := "11"
	if testStruct.TestString1 != expectedTestString1 {
		t.Errorf("testStruct>TestString1 is not as expected. Given %q. Expected %q", testStruct.TestString1, expectedTestString1)
	}

	expectedAll := map[string]string{
		"test_string_1": "11",
		"test2":         "BB",
		"test3":         "CC",
	}
	if !reflect.DeepEqual(testStruct.All, expectedAll) {
		t.Errorf("testStruct>All is not as expected. \nGiven %#v. \nExpected %#v", testStruct.All, expectedAll)
	}
	if !reflect.DeepEqual(testStruct.AllAJOnly, expectedAll) {
		t.Errorf("testStruct>AllAJOnly is not as expected. \nGiven %#v. \nExpected %#v", testStruct.AllAJOnly, expectedAll)
	}
}

func TestUnmarshal_HappyCase_Other(t *testing.T) {
	testStruct := struct {
		TestString1 string         `json:"test_string_1"`
		All         map[string]int `json:"-" aj:"other"`
		AllAJOnly   map[string]int `aj:"other"`
	}{}

	err := DefaultUnmarshaler.Unmarshal([]byte(`{"test_string_1": "AA", "test2": 22, "test3": 33}`), &testStruct)
	if err != nil {
		t.Fatalf("Faild to unmarshal json cause: %s", err)
	}

	expectedTestString1 := "AA"
	if testStruct.TestString1 != expectedTestString1 {
		t.Errorf("testStruct>TestString1 is not as expected. Given %q. Expected %q", testStruct.TestString1, expectedTestString1)
	}

	expectedAll := map[string]int{
		"test2": 22,
		"test3": 33,
	}
	if !reflect.DeepEqual(testStruct.All, expectedAll) {
		t.Errorf("testStruct>All is not as expected. \nGiven %#v. \nExpected %#v", testStruct.All, expectedAll)
	}
	if !reflect.DeepEqual(testStruct.AllAJOnly, expectedAll) {
		t.Errorf("testStruct>AllAJOnly is not as expected. \nGiven %#v. \nExpected %#v", testStruct.AllAJOnly, expectedAll)
	}
}

func TestUnmarshal_ErrorPath(t *testing.T) {
	tests := []struct {
		name              string
		json              string
		structure         interface{}
		expectedErrorText string
	}{
		{
			name: "all: mismatched json 'other' type",
			json: `{"field": "field-payload", "1": 33.55}`,
			structure: &struct {
				Field string             `json:"field"`
				Other map[string]float32 `aj:"all"`
			}{},
			expectedErrorText: "failed to unmarshal json field with tag \"all\": json: cannot unmarshal string into Go value of type float32",
		}, {
			name: "all: mismatched type",
			json: `{"field": "field-payload", "1": 33.55}`,
			structure: &struct {
				Field string `json:"field"`
				Other int    `aj:"all"`
			}{},
			expectedErrorText: "failed to unmarshal json field with tag \"all\": json: cannot unmarshal object into Go value of type int",
		}, {
			name: "other: mismatched json 'other' type",
			json: `{"field": "field-payload", "2": "other-field"}`,
			structure: &struct {
				Field string             `json:"field"`
				Other map[string]float32 `aj:"all"`
			}{},
			expectedErrorText: "failed to unmarshal json field with tag \"all\": json: cannot unmarshal string into Go value of type float32",
		}, {
			name: "other: mismatched type",
			json: `{"field": "field-payload", "1": 33.55}`,
			structure: &struct {
				Field string `json:"field"`
				Other string `aj:"all"`
			}{},
			expectedErrorText: "failed to unmarshal json field with tag \"all\": json: cannot unmarshal object into Go value of type string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DefaultUnmarshaler.Unmarshal([]byte(tt.json), tt.structure)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErrorText) {
				t.Errorf("expected response error is not as expected. Expected: %q. Given: %q", tt.expectedErrorText, err)
			}
		})
	}
}
