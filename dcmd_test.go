package dcommand

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuckcmd(t *testing.T) {
	cmd := &DCommand{}

	// testCmd := "test1 --xxx a b c --iii 7 8 9 --yyy d e f g --aaa --zzz h i g k --bbb"
	cmd.Command("test1").
		Flag(&FlagInfo{Name: "xxx"}).
		Flag(&FlagInfo{Name: "yyy"}).
		Flag(&FlagInfo{Name: "zzz"}).
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "test1")
			test.Equal(_cmd.Name, "test1")
			test.Equal(len(_cmd.Operators), 0)
			test.Equal(len(_cmd.flagParamsMap), 3)
			test.Equal(len(_cmd.Flags), 3)
			fmt.Println("这里的 test 1 的 map 是: ", _cmd.flagParamsMap)
			return nil
		})

	testCmd := "test1 --xxx a b c --yyy d e f g --zzz h i g k"
	cmd.Execute(strings.Split(testCmd, " "))

	// testCmd := "test2 op1 --xxx a b c --yyy d e f"
	cmd.Command("test2").
		Operator("op1").
		Operator("op2").
		Flag(&FlagInfo{Name: "xxx"}).
		Flag(&FlagInfo{Name: "zzz"}).
		Flag(&FlagInfo{Name: "yyy"}).
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "test2")
			test.Equal(_cmd.Name, "test2")
			test.Equal(len(_cmd.Operators), 2)
			test.Equal(_cmd.Operators[0].Name, "op1")
			test.Equal(len(_cmd.Operators[0].flagParamsMap), 0)
			test.Equal(_cmd.Operators[1].Name, "op2")
			test.Equal(len(_cmd.Operators[1].flagParamsMap), 0)
			fmt.Println("这里的 op1是: ", _cmd.Operators[0].flagParamsMap)
			fmt.Println("这里的 op2是: ", _cmd.Operators[1].flagParamsMap)

			test.Equal(len(_cmd.Operators[0].Flags), 0)
			test.Equal(len(_cmd.Operators[1].Flags), 3)

			test.Equal(_cmd.Operators[1].Flags[0].Long, "xxx")
			test.Equal(_cmd.Operators[1].Flags[1].Long, "zzz")
			test.Equal(_cmd.Operators[1].Flags[2].Long, "yyy")
			// 因为 xxx 和 yyy 是 op2 的 flag
			test.Equal(len(_cmd.Operators[1].Flags[0].Params), 0)
			test.Equal(len(_cmd.Operators[1].Flags[2].Params), 0)
			return nil
		})

	testCmd = "test2 op1 --xxx a b c d --yyy d e f 7"
	cmd.Execute(strings.Split(testCmd, " "))

	// "test3 op1 op2 --xxx a b c --yyy d e f g --aaa wuxiao"
	cmd.Command("test3").
		Operator("op0").
		Operator("op1").
		Operator("op2").
		Flag(&FlagInfo{Name: "xxx"}).
		Flag(&FlagInfo{Name: "zzz"}).
		Flag(&FlagInfo{Name: "yyy"}).
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "test3")
			test.Equal(_cmd.Name, "test3")
			test.Equal(len(_cmd.Operators), 3)
			test.Equal(_cmd.Operators[0].Name, "op0")
			test.Equal(len(_cmd.Operators[0].flagParamsMap), 0)
			test.Equal(_cmd.Operators[1].Name, "op1")
			test.Equal(len(_cmd.Operators[1].flagParamsMap), 0)
			test.Equal(_cmd.Operators[2].Name, "op2")
			test.Equal(len(_cmd.Operators[2].flagParamsMap), 2)

			xxx, ok := _cmd.Operators[2].flagParamsMap["--xxx"]
			test.Equal(ok, true)
			test.Equal(len(xxx), 3)
			fmt.Println("这里的 xxx 是: ", xxx)

			yyy, ok := _cmd.Operators[2].flagParamsMap["--yyy"]
			test.Equal(ok, true)
			test.Equal(len(yyy), 4)

			fmt.Println("这里的 yyy 是: ", yyy)

			test.Equal(_cmd.Operators[2].Flags[0].Long, "xxx")
			test.Equal(_cmd.Operators[2].Flags[1].Long, "zzz")
			test.Equal(_cmd.Operators[2].Flags[2].Long, "yyy")
			test.Equal(len(_cmd.Operators[2].Flags), 3)
			test.Equal(len(_cmd.Operators[2].Flags[0].Params), 3)
			test.Equal(len(_cmd.Operators[2].Flags[2].Params), 4)

			return nil
		})

	testCmd = "test3 op1 op2 --xxx a b c --yyy d e f g --aaa wuxiao"
	cmd.Execute(strings.Split(testCmd, " "))

	testCmd = "test2 op1 --xxx a b c d --yyy d e f 7"
	cmd.Execute(strings.Split(testCmd, " "))

	cmd.Command("test2-1").
		Operator("op0").
		Operator("op1").
		Operator("op2").
		Flag(&FlagInfo{Name: "xxx"}).
		Flag(&FlagInfo{Name: "zzz"}).
		Flag(&FlagInfo{Name: "yyy"}).
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "test2-1")
			test.Equal(_cmd.Name, "test2-1")
			test.Equal(len(_cmd.Operators), 3)
			test.Equal(_cmd.Operators[0].Name, "op0")
			test.Equal(len(_cmd.Operators[0].flagParamsMap), 0)
			test.Equal(_cmd.Operators[1].Name, "op1")
			test.Equal(len(_cmd.Operators[1].flagParamsMap), 0)
			test.Equal(_cmd.Operators[2].Name, "op2")
			test.Equal(len(_cmd.Operators[2].flagParamsMap), 2)

			xxx, ok := _cmd.Operators[2].flagParamsMap["--xxx"]
			test.Equal(ok, true)
			test.Equal(len(xxx), 3)
			fmt.Println("这里的 xxx 是: ", xxx)

			yyy, ok := _cmd.Operators[2].flagParamsMap["--yyy"]
			test.Equal(ok, true)
			test.Equal(len(yyy), 4)
			fmt.Println("这里的 yyy 是: ", yyy)
			return nil
		}).
		WithParamsHandler(func(command string, fc *DCommand, params ...interface{}) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			for i, p := range params {
				_i := p.(int)
				test.Equal(i+1, _i)
				fmt.Println("这里的参数是: ", _i)
			}

			test.Equal(command, "test2-1")
			test.Equal(_cmd.Name, "test2-1")
			test.Equal(len(_cmd.Operators), 3)
			test.Equal(_cmd.Operators[0].Name, "op0")
			test.Equal(len(_cmd.Operators[0].flagParamsMap), 0)
			test.Equal(_cmd.Operators[1].Name, "op1")
			test.Equal(len(_cmd.Operators[1].flagParamsMap), 0)
			test.Equal(_cmd.Operators[2].Name, "op2")
			test.Equal(len(_cmd.Operators[2].flagParamsMap), 2)

			xxx, ok := _cmd.Operators[2].flagParamsMap["--xxx"]
			test.Equal(ok, true)
			test.Equal(len(xxx), 3)
			fmt.Println("这里的 xxx 是: ", xxx)

			yyy, ok := _cmd.Operators[2].flagParamsMap["--yyy"]
			test.Equal(ok, true)
			test.Equal(len(yyy), 4)
			fmt.Println("这里的 yyy 是: ", yyy)
			return nil
		})

	testCmd = "test2-1 op1 op2 --xxx a b c --yyy d e f g --aaa wuxiao"
	cmd.Execute(strings.Split(testCmd, " "))
	cmd.ExecuteWithParams(strings.Split(testCmd, " "), 1, 2, 3)

	cmd.Command("chahua").
		Operator("publish").
		Flag(&FlagInfo{Name: "add", Short: "a"}).
		Flag(&FlagInfo{Name: "upload", Short: "u"}).
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "chahua")
			test.Equal(_cmd.Name, "chahua")
			test.Equal(len(_cmd.Operators), 1)

			test.Equal(_cmd.Operators[0].Name, "publish")
			test.Equal(len(_cmd.Operators[0].Flags), 2)
			test.Equal(_cmd.Operators[0].Flags[0].Long, "add")
			test.Equal(_cmd.Operators[0].Flags[0].Short, "a")
			test.Equal(_cmd.Operators[0].Flags[1].Long, "upload")
			test.Equal(_cmd.Operators[0].Flags[1].Short, "u")
			test.Equal(len(_cmd.Operators[0].flagParamsMap), 2)
			upload, ok := _cmd.Operators[0].flagParamsMap["--upload"]
			test.Equal(ok, true)
			test.Equal(upload[0], "illustration-x3")
			test.Equal(upload[1], "illustration-x4")
			add, ok := _cmd.Operators[0].flagParamsMap["--add"]
			test.Equal(ok, true)
			test.Equal(add[0], "illustration-x1")
			test.Equal(add[1], "illustration-x2")
			return nil
		})

	testCmd = "chahua publish --upload illustration-x3 illustration-x4 --add illustration-x1 illustration-x2"
	cmd.Execute(strings.Split(testCmd, " "))

	cmd.Command("chahua2").
		Operator("publish").
		Flag(&FlagInfo{Name: "add", Short: "a"}).
		Flag(&FlagInfo{Name: "upload", Short: "u"}).
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "chahua2")
			test.Equal(_cmd.Name, "chahua2")
			test.Equal(len(_cmd.Operators), 1)

			test.Equal(_cmd.Operators[0].Name, "publish")
			test.Equal(len(_cmd.Operators[0].Flags), 2)
			test.Equal(_cmd.Operators[0].Flags[0].Long, "add")
			test.Equal(_cmd.Operators[0].Flags[0].Short, "a")
			test.Equal(_cmd.Operators[0].Flags[1].Long, "upload")
			test.Equal(_cmd.Operators[0].Flags[1].Short, "u")
			test.Equal(len(_cmd.Operators[0].flagParamsMap), 2)
			upload, ok := _cmd.Operators[0].flagParamsMap["-u"]
			test.Equal(ok, true)
			test.Equal(upload[0], "illustration-x3")
			test.Equal(upload[1], "illustration-x4")
			add, ok := _cmd.Operators[0].flagParamsMap["-a"]
			test.Equal(ok, true)
			test.Equal(add[0], "illustration-x1")
			test.Equal(add[1], "illustration-x2")
			return nil
		})

	testCmd = "chahua2 publish -u illustration-x3 illustration-x4 -a illustration-x1 illustration-x2"
	cmd.Execute(strings.Split(testCmd, " "))

	cmd.Command("test-return").
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "test-return")
			test.Equal(_cmd.Name, "test-return")

			return fmt.Errorf("1")
		}).
		WithParamsHandler(func(command string, fc *DCommand, params ...interface{}) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "test-return")
			test.Equal(_cmd.Name, "test-return")
			str := params[0].(string)
			test.Equal(str, "2")
			return fmt.Errorf(str)
		})

	testCmd = "test-return"
	err := cmd.Execute(strings.Split(testCmd, " "))
	if err != nil {
		test := assert.New(t)
		fmt.Println("这里的错误是1111 : ", err.Error())
		test.Equal(err.Error(), "1")
	}

	err = cmd.ExecuteWithParams(strings.Split(testCmd, " "), "2")
	if err != nil {
		test := assert.New(t)
		fmt.Println("这里的错误是2222 : ", err.Error())
		test.Equal(err.Error(), "2")
	}

	testStrCmd := "test-origin I am cmd"
	cmd.Command("test-origin").
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			originCmd := fc.GetOriginCommand()
			fmt.Println("这里原始的 cmd 是: ", originCmd)
			test.Equal(originCmd, testStrCmd)
			return nil
		})
	cmd.ExecuteStr(testStrCmd)

	testStrCmd2 := "test-origin-2 I am cmd 2222"
	cmd.Command("test-origin-2").
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			originCmd := fc.GetOriginCommand()
			fmt.Println("这里原始的 cmd 2222 是: ", originCmd)
			test.Equal(originCmd, testStrCmd2)
			return nil
		})
	cmd.ExecuteStr(testStrCmd2)

	// testCmd := "test1 --iii 7 8 9 --aaa --zzz h i g k --bbb"
	cmd.Command("test-default-1").
		Flag(&FlagInfo{Name: "xxx", Default: []string{"a", "b", "c"}}).
		Flag(&FlagInfo{Name: "yyy", Short: "y", Default: []string{"d", "e", "f"}}).
		Flag(&FlagInfo{Name: "zzz"}).
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "test-default-1")
			test.Equal(_cmd.Name, "test-default-1")
			test.Equal(len(_cmd.Operators), 0)
			test.Equal(len(_cmd.flagParamsMap), 3)
			test.Equal(len(_cmd.Flags), 3)
			fmt.Println("这里的 test-default-1 的 map 是: ", _cmd.flagParamsMap)
			return nil
		})

	testCmd = "test-default-1 --iii 7 8 9 --aaa --zzz h i g k --bbb"
	cmd.Execute(strings.Split(testCmd, " "))

	// testCmd := "test-default-2 op -zzz 7 8 9"
	cmd.Command("test-default-2").
		Operator("op").
		Flag(&FlagInfo{Name: "xxx", Default: []string{"a", "b", "c"}}).
		Flag(&FlagInfo{Name: "zzz", Default: []string{"1", "2", "3"}}).
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "test-default-2")
			test.Equal(_cmd.Name, "test-default-2")
			test.Equal(len(_cmd.Operators), 1)
			test.Equal(_cmd.Operators[0].Name, "op")
			test.Equal(len(_cmd.Operators[0].flagParamsMap), 2)
			fmt.Println("这里的 op是: ", _cmd.Operators[0].flagParamsMap)
			test.Equal(len(_cmd.Operators[0].Flags), 2)

			test.Equal(_cmd.Operators[0].Flags[0].Long, "xxx")
			test.Equal(_cmd.Operators[0].Flags[1].Long, "zzz")

			flag := fc.GetFlagIfExistInOperatorByOperator("xxx", true, _cmd.Operators[0])
			test.Equal(flag.Params[0], "a")
			test.Equal(flag.Params[1], "b")
			test.Equal(flag.Params[2], "c")
			flag = fc.GetFlagIfExistInOperatorByOperator("zzz", true, _cmd.Operators[0])
			test.Equal(flag.Params[0], "7")
			test.Equal(flag.Params[1], "8")
			test.Equal(flag.Params[2], "9")
			return nil
		})

	testCmd = "test-default-2 op --zzz 7 8 9"
	cmd.Execute(strings.Split(testCmd, " "))

	temporaryIndex := 0
	cmd.Command("test-flag-passed").
		Operator("op").
		Flag(&FlagInfo{Name: "xxx"}).
		Handler(func(command string, fc *DCommand) error {
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "test-flag-passed")
			if temporaryIndex == 0 {
				test.Equal(_cmd.Operators[0].Flags[0].Passed, true)
			}
			if temporaryIndex == 1 {
				test.Equal(_cmd.Operators[0].Flags[0].Passed, false)
			}
			return nil
		})
	testCmd = "test-flag-passed op --xxx"
	cmd.Execute(strings.Split(testCmd, " "))
	temporaryIndex = 1

	testCmd = "test-flag-passed op"
	cmd.Execute(strings.Split(testCmd, " "))
}
