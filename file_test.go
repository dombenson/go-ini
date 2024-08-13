package ini

import (
	"os"
	"strings"
	"testing"
)

func Test_EnvironmentVariableOverrides_Enabled(t *testing.T) {
	testIni := `
global1 = wrong_global

[test1]
value1 = wrong1
value2 = wrong2
value3[] = wrong3_1
value3[] = wrong3_2
value3[] = wrong3_3
value4[] = wrong4_1
value4[] = wrong4_2
value4[] = wrong4_3

[test2]
value_int = -1
value_bool = false
`

	err := os.Setenv("GO_INI_GLOBAL1", "correct_global")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST1_VALUE1", "correct1")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST1_VALUE2", "correct2")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST1_VALUE3_1", "correct3_1")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST1_VALUE3_2", "correct3_2")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST1_VALUE4", "[]")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST2_VALUE_INT", "42")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST2_VALUE_BOOL", "true")
	if err != nil {
		t.Fatal(err)
	}

	file, err := Load(strings.NewReader(testIni))
	if err != nil {
		t.Fatal(err)
	}

	file.EnableEnvironmentVariableOverrides("GO_INI")

	checkStr(t, file, "", "global1", "correct_global")
	checkStr(t, file, "test1", "value1", "correct1")
	checkStr(t, file, "test1", "value2", "correct2")
	checkArr(t, file, "test1", "value3", []string{"correct3_1", "correct3_2"})
	checkArr(t, file, "test1", "value4", []string{})

	checkInt(t, file, "test2", "value_int", 42)
	checkBool(t, file, "test2", "value_bool", true)
}

func Test_EnvironmentVariableOverrides_Disabled(t *testing.T) {
	testIni := `
global1 = correct_global

[test1]
value1 = correct1
value2 = correct2
value3[] = correct3_1
value3[] = correct3_2
value3[] = correct3_3
value4[] = correct4_1
value4[] = correct4_2
value4[] = correct4_3

[test2]
value_int = 42
value_bool = true
`

	err := os.Setenv("GO_INI_GLOBAL1", "wrong_global")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST1_VALUE1", "wrong1")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST1_VALUE2", "wrong2")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST1_VALUE3_1", "wrong3_1")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST1_VALUE3_2", "wrong3_2")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST1_VALUE4", "[]")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST2_VALUE_INT", "-1")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("GO_INI_TEST2_VALUE_BOOL", "false")
	if err != nil {
		t.Fatal(err)
	}

	file, err := Load(strings.NewReader(testIni))
	if err != nil {
		t.Fatal(err)
	}

	checkStr(t, file, "", "global1", "correct_global")
	checkStr(t, file, "test1", "value1", "correct1")
	checkStr(t, file, "test1", "value2", "correct2")
	checkArr(t, file, "test1", "value3", []string{"correct3_1", "correct3_2", "correct3_3"})
	checkArr(t, file, "test1", "value4", []string{"correct4_1", "correct4_2", "correct4_3"})

	checkInt(t, file, "test2", "value_int", 42)
	checkBool(t, file, "test2", "value_bool", true)
}
