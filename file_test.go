package ini

import (
	"os"
	"strings"
	"testing"
)

func Test_file_ParseEnvironmentVariables(t *testing.T) {
	testIni := `
[test1]
value1 = {{ .Env.GO_INI_TEST_ONE }}
value2 = "{{ .Env.GO_INI_TEST_TWO }}"
value3[] = {{ .Env.GO_INI_TEST_THREE_ONE }}
value3[] = {{ .Env.GO_INI_TEST_THREE_TWO }}
value4 = {{ .Env.GO_INI_TEST_FOUR }}

[test2]
valueInt = {{ .Env.GO_INI_TEST_INT }}
valueBool = {{ .Env.GO_INI_TEST_BOOL }}
`

	err := os.Setenv("GO_INI_TEST_ONE", "ONE")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST_TWO", "TWO")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST_THREE_ONE", "THREE_ONE")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST_THREE_TWO", "THREE_TWO")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST_FOUR", "{{ .Env.GO_INI_TEST_RUNTIME_ENV_VAR }}")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST_INT", "42")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST_BOOL", "true")
	if err != nil {
		t.Fatal(err)
	}

	file, err := Load(strings.NewReader(testIni))
	if err != nil {
		t.Fatal(err)
	}

	err = file.ParseEnvironmentVariables()
	if err != nil {
		t.Fatal(err)
	}

	checkStr(t, file, "test1", "value1", "ONE")
	checkStr(t, file, "test1", "value2", "TWO")
	checkArr(t, file, "test1", "value3", []string{"THREE_ONE", "THREE_TWO"})

	// Test is validating that we can use a Go template value as a variable to support the use-case where we parse
	// the ini file once when building with additional values added at runtime.
	checkStr(t, file, "test1", "value4", "{{ .Env.GO_INI_TEST_RUNTIME_ENV_VAR }}")

	checkInt(t, file, "test2", "valueInt", 42)
	checkBool(t, file, "test2", "valueBool", true)
}
