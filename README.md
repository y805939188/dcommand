# 这是一个命令行解析工具

```golang
cmd := &DCmd{}

	// testCmd := "test1 --xxx a b c --iii 7 8 9 --yyy d e f g --aaa --zzz h i g k --bbb"
	cmd.Command("test1").
		Flag("xxx").
		Flag("yyy").
		Flag("zzz").
		Handler(func(command string, fc *DCmd) {
			// fmt.Println("进来了这里")
			test := assert.New(t)
			_cmd := fc.GetCommandIfExist(command)
			test.Equal(command, "test1")
			test.Equal(_cmd.Name, "test1")
			test.Equal(len(_cmd.Operators), 0)
			test.Equal(len(_cmd.flagParamsMap), 3)
			test.Equal(len(_cmd.Flags), 3)
			fmt.Println("这里的 test 1 的 map 是: ", _cmd.flagParamsMap)
		})

	testCmd := "test1 --xxx a b c --yyy d e f g --zzz h i g k"
	cmd.Execute(strings.Split(testCmd, " "))

	// testCmd := "test2 op1 --xxx a b c --yyy d e f"
	cmd.Command("test2").
		Operator("op1").
		Operator("op2").
		Flag("xxx").
		Flag("zzz").
		Flag("yyy").
		Handler(func(command string, fc *DCmd) {
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
		})

	testCmd = "test2 op1 --xxx a b c d --yyy d e f 7"
	cmd.Execute(strings.Split(testCmd, " "))

	cmd.Command("test3").
		Operator("op0").
		Operator("op1").
		Operator("op2").
		Flag("xxx").
		Flag("zzz").
		Flag("yyy").
		Handler(func(command string, fc *DCmd) {
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
		})

	testCmd = "test3 op1 op2 --xxx a b c --yyy d e f g --aaa wuxiao"
	cmd.Execute(strings.Split(testCmd, " "))
```
