package dcommand

import (
	"fmt"
	"strings"
)

type Flag struct {
	Long          string
	Short         string
	Params        []string
	Desc          string
	OwnerCommand  *Command
	OwnerOperator *Operator
}

type Operator struct {
	Name  string
	Flags []*Flag // operator 可以有 flags
	// FallbackFlag  *Flag
	Desc          string
	Passed        bool // 是否传了这个参数
	OwnerCommand  *Command
	flagParamsMap map[string][]string
}

type Command struct {
	Name      string
	Operators []*Operator
	Flags     []*Flag // command 也可以有 flags
	// FallbackFlag   *Flag
	Desc           string
	handlers       []func(string, *DCommand) error
	paramsHandlers []func(string, *DCommand, ...interface{}) error
	flagParamsMap  map[string][]string
	// allFlagsMap    map[string]bool
	invalidOperators []*Operator
}

type DCommand struct {
	Commands            []*Command
	originCommand       string
	currentCommand      string
	currentOperator     string
	currentCustomParams []interface{}
}

type FlagInfo struct {
	Name    string
	Short   string
	Desc    string
	Default []string
}

func (fc *DCommand) validFlag(name string) error {
	if name == "" {
		return fmt.Errorf("Need a name")
	}
	if fc.Commands == nil {
		return fmt.Errorf("Need to call Command before calling Operator")
	}
	if len(fc.Commands) == 0 {
		return fmt.Errorf("There is currently no Command")
	}
	return nil
}

// func (fc *DCommand) flagInCommand(name string, command *Command) bool {
// 	val, ok := command.allFlagsMap[name]
// 	if !ok {
// 		return false
// 	}
// 	return val
// }

// func (fc *DCommand) existValidFlag(cmd []string, command *Command) bool {
// 	for _, c := range cmd {
// 		yes := fc.flagInCommand(c, command)
// 		if yes {
// 			return true
// 		}
// 	}
// 	return false
// }

func (fc *DCommand) Flag(flagInfo *FlagInfo) *DCommand {
	name := flagInfo.Name
	err := fc.validFlag(name)
	if err != nil {
		fmt.Println(err.Error())
		return fc
	}
	for _, cmd := range fc.Commands {
		if fc.currentCommand != cmd.Name {
			continue
		}
		// cmd.allFlagsMap[name] = true
		// if flagInfo.Short != "" {
		// 	cmd.allFlagsMap[flagInfo.Short] = true
		// }
		if fc.currentOperator == "" {
			// 说明是直接给 command 添加 flag
			temp := &Flag{
				Long:         name,
				Params:       []string{},
				OwnerCommand: cmd,
			}
			temp.Short = flagInfo.Short
			temp.Desc = flagInfo.Desc
			if cmd.Flags == nil {
				cmd.Flags = []*Flag{temp}
			} else {
				cmd.Flags = append(cmd.Flags, temp)
			}
			if len(flagInfo.Default) > 0 {
				err := fc.SetFlagParamsForCommand("--"+name, flagInfo.Default, cmd)
				if err != nil {
					fmt.Println("set default value to command error: ", err.Error())
				}
			}
		} else {
			for _, operator := range cmd.Operators {
				if fc.currentOperator != operator.Name {
					continue
				}
				// 说明是给 operator 添加 flag
				temp := &Flag{
					Long:          name,
					Params:        []string{},
					OwnerOperator: operator,
				}
				temp.Short = flagInfo.Short
				temp.Desc = flagInfo.Desc
				if operator.Flags == nil || len(operator.Flags) == 0 {
					operator.Flags = []*Flag{temp}
				} else {
					operator.Flags = append(operator.Flags, temp)
				}
				if len(flagInfo.Default) > 0 {
					err := fc.SetFlagParamsForOperator("--"+name, flagInfo.Default, operator)
					if err != nil {
						fmt.Println("set default value to operator error: ", err)
					}
				}
				break
			}
		}
		break
	}
	return fc
}

/**
 * default flag 用来指定是否要把某个 flag 作为默认的 flag
 * 比如类似注册了这样的命令:
 *	cmd.Command("test").
 *		Operator("op").
 *	  DefaultFlag(&FlagInfo{Name: "default"})
 *		Flag(&FlagInfo{Name: "aaa"}).
 *		Flag(&FlagInfo{Name: "bbb"})
 *
 * 然后执行: test op abc def ght 这种一个 flag 都没有的命令的话
 * 就把 abc def ght 这坨东西塞给默认的 --default 这个 flag
 */
// func (fc *DCommand) DefaultFlag(name string, desc string) *DCommand {
// 	err := fc.validFlag(name)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return fc
// 	}
// 	for _, cmd := range fc.Commands {
// 		if fc.currentCommand != cmd.Name {
// 			continue
// 		}
// 		temp := &Flag{
// 			Long:   name,
// 			Params: []string{},
// 			Desc:   desc,
// 		}
// 		if fc.currentOperator == "" {
// 			// 说明是直接给 command 添加 flag
// 			temp.OwnerCommand = cmd
// 			cmd.FallbackFlag = temp
// 		} else {
// 			for _, operator := range cmd.Operators {
// 				if fc.currentOperator != operator.Name {
// 					continue
// 				}
// 				// 说明是给 operator 添加 flag
// 				temp.OwnerOperator = operator
// 				operator.FallbackFlag = temp
// 			}
// 		}
// 		break
// 	}
// 	return fc
// }

func (fc *DCommand) Handler(fn func(string, *DCommand) error) *DCommand {
	if fc.Commands == nil {
		fmt.Println("Need to call Command before calling Operator")
		return fc
	}

	if len(fc.Commands) == 0 {
		fmt.Println("There is currently no Command")
		return fc
	}

	if fc.currentCommand == "" {
		fmt.Println("Need to call Command before calling Operator")
		return fc
	}

	for _, command := range fc.Commands {
		if fc.currentCommand == command.Name {
			if command.handlers == nil {
				command.handlers = [](func(string, *DCommand) error){}
			}
			command.handlers = append(command.handlers, fn)
			break
		}
	}
	return fc
}

func (fc *DCommand) WithParamsHandler(fn func(string, *DCommand, ...interface{}) error) *DCommand {
	if fc.Commands == nil {
		fmt.Println("Need to call Command before calling Operator")
		return fc
	}

	if len(fc.Commands) == 0 {
		fmt.Println("There is currently no Command")
		return fc
	}

	if fc.currentCommand == "" {
		fmt.Println("Need to call Command before calling Operator")
		return fc
	}

	for _, command := range fc.Commands {
		if fc.currentCommand == command.Name {
			if command.handlers == nil {
				command.paramsHandlers = [](func(string, *DCommand, ...interface{}) error){}
			}
			command.paramsHandlers = append(command.paramsHandlers, fn)
			break
		}
	}
	return fc
}

func (fc *DCommand) Operator(name string) *DCommand {
	if fc.Commands == nil {
		fmt.Println("Need to call Command before calling Operator")
		return fc
	}
	if len(fc.Commands) == 0 {
		fmt.Println("There is currently no Command")
		return fc
	}

	fc.currentOperator = name
	for _, cmd := range fc.Commands {
		if fc.currentCommand != cmd.Name {
			continue
		}
		temp := &Operator{
			Name:          name,
			Flags:         nil,
			OwnerCommand:  cmd,
			flagParamsMap: make(map[string][]string),
		}
		if cmd.Operators == nil || len(cmd.Operators) == 0 {
			cmd.Operators = []*Operator{temp}
		} else {
			cmd.Operators = append(cmd.Operators, temp)
		}
		break
	}
	return fc
}

func (fc *DCommand) GetFlagIfExistInCommand(flagName string, isLong bool, commands ...*Command) []*Flag {
	if fc.Commands == nil && len(commands) == 0 {
		return nil
	}

	currentCommands := fc.Commands
	if len(commands) > 0 {
		currentCommands = commands
	}

	var res []*Flag
	for _, cmd := range currentCommands {
		if cmd.Flags == nil {
			continue
		}

		for _, flag := range cmd.Flags {
			if isLong && flag.Long == flagName {
				res = append(res, flag)
			}

			if !isLong && flag.Short == flagName {
				res = append(res, flag)
			}
		}
	}

	return res
}

func (fc *DCommand) GetFlagIfExistInOperatorByCommand(flagName string, isLong bool, command *Command) []*Flag {
	if fc.Commands == nil {
		return nil
	}

	if command == nil {
		return nil
	}

	cmd := fc.GetCommandIfExist(command.Name)

	if cmd == nil {
		return nil
	}

	var res []*Flag
	for _, op := range cmd.Operators {
		if op.Flags == nil {
			continue
		}

		for _, f := range op.Flags {
			if (isLong && f.Long == flagName) || (!isLong && f.Short == flagName) {
				res = append(res, f)
			}
		}
	}

	return res
}

func (fc *DCommand) GetFlagIfExistInOperatorByOperator(flagName string, isLong bool, operator *Operator) *Flag {
	if fc.Commands == nil {
		return nil
	}

	if operator == nil {
		return nil
	}

	op := fc.GetOperatorIfExistByCommand(operator.Name, operator.OwnerCommand)

	if op == nil {
		return nil
	}

	for _, f := range op.Flags {
		if (isLong && f.Long == flagName) || (!isLong && f.Short == flagName) {
			return f
		}
	}

	return nil
}

func (fc *DCommand) GetFlagIfExistInOperator(flagName string, isLong bool) []*Flag {
	if fc.Commands == nil {
		return nil
	}

	var res []*Flag
	for _, cmd := range fc.Commands {
		if cmd.Operators == nil {
			continue
		}

		for _, operator := range cmd.Operators {
			if operator.Flags == nil {
				continue
			}

			for _, flag := range operator.Flags {
				if isLong && flag.Long == flagName {
					res = append(res, flag)
				}

				if !isLong && flag.Short == flagName {
					res = append(res, flag)
				}
			}
		}
	}

	return res
}

func (fc *DCommand) GetCommandIfExist(command string) *Command {
	if fc.Commands == nil {
		return nil
	}

	for _, cmd := range fc.Commands {
		if cmd.Name == command {
			return cmd
		}
	}
	return nil
}

func (fc *DCommand) GetOperatorIfExist(operatorName string, commands ...*Command) []*Operator {
	if fc.Commands == nil {
		return nil
	}

	currentCommands := fc.Commands
	if len(commands) > 0 {
		currentCommands = commands
	}

	var res []*Operator
	for _, cmd := range currentCommands {
		if cmd.Operators == nil {
			continue
		}
		for _, operator := range cmd.Operators {
			if operator.Name == operatorName {
				res = append(res, operator)
			}
		}
	}

	return res
}

func (fc *DCommand) GetOperatorIfExistByCommand(operatorName string, command *Command) *Operator {
	if command == nil {
		return nil
	}

	if fc.Commands == nil {
		return nil
	}
	for _, cmd := range fc.Commands {
		if cmd.Name != command.Name {
			continue
		}

		for _, operator := range cmd.Operators {
			if operator.Name == operatorName {
				return operator
			}
		}
	}

	return nil
}

func IsFlag(str string) (bool, isLong bool, pureFlag string) {
	if strings.HasPrefix(str, "--") {
		isLong = true
		pureFlag = strings.Replace(str, "--", "", 1)
		return true, isLong, pureFlag
	}

	if strings.HasPrefix(str, "-") {
		isLong = false
		pureFlag = strings.Replace(str, "-", "", 1)
		return true, false, pureFlag
	}

	return false, false, ""
}

func (fc *DCommand) IsFlag(str string) (bool, isLong bool, pureFlag string) {
	return IsFlag(str)
}

func (fc *DCommand) SetFlagParamsForCommand(flag string, params []string, command *Command) error {

	isFlag, isLone, pureFlag := fc.IsFlag(flag)
	if !isFlag {
		return fmt.Errorf("The str(%s) is not a flag", flag)
	}

	flags := fc.GetFlagIfExistInCommand(pureFlag, isLone, command)
	for _, f := range flags {
		if (isLone && f.Long == pureFlag) || (!isLone && f.Short == pureFlag) {
			f.OwnerCommand.flagParamsMap[flag] = params
			f.Params = params
			break
		}
	}
	return nil
}

func (fc *DCommand) SetFlagParamsForOperator(flag string, params []string, operator *Operator) error {
	isFlag, isLone, pureFlag := fc.IsFlag(flag)

	if !isFlag {
		return fmt.Errorf("The str(%s) is not a flag", flag)
	}

	f := fc.GetFlagIfExistInOperatorByOperator(pureFlag, isLone, operator)
	if f == nil {
		return nil
	}
	if (isLone && f.Long == pureFlag) || (!isLone && f.Short == pureFlag) {
		f.OwnerOperator.flagParamsMap[flag] = params
		f.Params = params
	}
	return nil
}

func (fc *DCommand) GetOriginCommand() string {
	return fc.originCommand
}

// func (fc *DCommand) separateInvalidStr(cmd []string, command *Command) ([]string, []string) {
// 	// 该方法用来分离没有意义的字符
// 	// 不停遍历数组中的每一项, 直到碰到第一个 flag 或碰到第一个 operator 就停下
// 	// 然后看是否有 fallback flag, 有的话把前边弄出来的那些无效的都作为这个 flag 的 params
// 	res := []string{}
// 	if cmd == nil {
// 		return res, cmd
// 	}

// 	for index, str := range cmd {
// 		isFlag, _, _ := fc.IsFlag(str)
// 		if isFlag {
// 			// 如果是一个 flag 就直接返回
// 			return res, cmd[index:]
// 		}
// 		// 如果是一个有效的 operator 的话也直接返回
// 		op := fc.GetOperatorIfExistByCommand(str, command)
// 		if op != nil {
// 			return res, cmd[index:]
// 		}
// 		res = append(res, str)
// 	}
// 	return res, []string{}
// }

/**
 * cmd --xxx a b c --yyy d e f
 * cmd a b c --xxx d e f -yyy g h i
 */
func (fc *DCommand) Execute(cmd []string) error {

	commandName := cmd[0]
	originCmd := strings.Join(cmd, " ")
	fc.originCommand = originCmd
	cmd = cmd[1:]
	command := fc.GetCommandIfExist(commandName)

	if command == nil {
		return fmt.Errorf("command name is nil")
	}

	if len(cmd) != 0 {

		// TODO: 做个 fallback 功能
		// yes := fc.existValidFlag(cmd, command)
		// if !yes {
		// 	// 如果这串 cmd 里一个 flag 都没有的话, 那就看是否有 default flag
		// 	// 也就是把所有的东西都当做 default flag 的参数给 fallback

		// }

		// 有可能就一个光杆命令没有 operator 和 flag
		isFlag, _, _ := fc.IsFlag(cmd[0])
		if isFlag {
			// 如果第一个字符串就是 flag 的话, 那后续就不能再有任何 operator 了
			// 就类似 cmd --xxx abc --yyy d e f 这种
			// 因为在 flg 后面加 operator 的话, 无法分辨出到底是 operator 还是 flag 的 param
			// 所以如果第一个就是 flag 的话, 可以直接把它作为 command 的 flag
			temporaryFlag := ""
			temporaryParams := []string{}
			for i, str := range cmd {
				isFlag, _, _ := fc.IsFlag(str)
				if isFlag {
					if temporaryFlag != "" {
						err := fc.SetFlagParamsForCommand(temporaryFlag, temporaryParams, command)
						if err != nil {
							return err
						}
					}
					temporaryParams = []string{}
					temporaryFlag = str
					if i == len(cmd)-1 {
						err := fc.SetFlagParamsForCommand(temporaryFlag, temporaryParams, command)
						if err != nil {
							return err
						}
					}
				} else {
					temporaryParams = append(temporaryParams, str)
					if i == len(cmd)-1 {
						err := fc.SetFlagParamsForCommand(temporaryFlag, temporaryParams, command)
						if err != nil {
							return err
						}
					}
				}
			}
		} else {
			// 进到这里说明第一个字符串可能是 operator, 当然也可能是一个无意义的乱七八糟的玩意儿
			// 对于这种乱七八糟的东西就看是否有 fallbackFlag
			temporaryFlag := ""
			temporaryParams := []string{}
			temporaryOperator := ""

			for i, str := range cmd {
				isFlag, _, _ := fc.IsFlag(str)
				if isFlag {
					// 当碰到第一个是 flag 类型的 str 时
					// 把 cmd 从当前这个 flag 的位置开始截取下来
					// 然后后面的所有 flag 以及 params 都算作是当前这个 temporaryOperator 的参数
					_temporaryCmd := cmd[i:]
					_temporaryFlag := ""
					_temporaryParams := []string{}
					for i, _str := range _temporaryCmd {
						isFlag, _, _ := fc.IsFlag(_str)

						if isFlag {
							if _temporaryFlag != "" {
								err := fc.SetFlagParamsForOperator(_temporaryFlag, _temporaryParams, fc.GetOperatorIfExistByCommand(temporaryOperator, command))
								if err != nil {
									return err
								}
							}
							_temporaryParams = []string{}
							_temporaryFlag = _str
							if i == len(_temporaryCmd)-1 {
								err := fc.SetFlagParamsForCommand(_temporaryFlag, _temporaryParams, command)
								if err != nil {
									fmt.Println(err.Error())
									return err
								}
							}
						} else {
							_temporaryParams = append(_temporaryParams, _str)
							if i == len(_temporaryCmd)-1 {
								err := fc.SetFlagParamsForOperator(_temporaryFlag, _temporaryParams, fc.GetOperatorIfExistByCommand(temporaryOperator, command))
								if err != nil {
									fmt.Println(err.Error())
									return err
								}
							}
						}
					}
					break
				} else {
					if temporaryFlag != "" && temporaryOperator != "" {
						err := fc.SetFlagParamsForOperator(temporaryFlag, temporaryParams, fc.GetOperatorIfExistByCommand(temporaryFlag, command))
						if err != nil {
							fmt.Println(err.Error())
							return err
						}
					}

					if temporaryFlag != "" {
						temporaryParams = append(temporaryParams, str)
					} else {
						// 进到这里, 一定是 operator
						op := fc.GetOperatorIfExistByCommand(str, command)
						if op != nil {
							op.Passed = true
							temporaryOperator = str
						}
					}
				}
			}
		}
	}

	if fc.currentCustomParams != nil {
		for _, fn := range command.paramsHandlers {
			err := fn(commandName, fc, fc.currentCustomParams...)
			if err != nil {
				return err
			}
		}
	} else {
		for _, fn := range command.handlers {
			err := fn(commandName, fc)
			if err != nil {
				return err
			}
		}
	}

	fc.originCommand = ""
	return nil
}

func (fc *DCommand) ExecuteStr(str string) error {
	return fc.Execute(strings.Split(str, " "))
}

func (fc *DCommand) ExecuteStrWithParams(str string, params ...interface{}) error {
	return fc.ExecuteWithParams(strings.Split(str, " "), params...)
}

func (fc *DCommand) ExecuteWithParams(cmd []string, params ...interface{}) error {
	fc.currentCustomParams = params
	err := fc.Execute(cmd)
	fc.currentCustomParams = nil
	return err
}

func (fc *DCommand) Command(name string) *DCommand {
	if fc.Commands == nil {
		cmd := []*Command{
			{
				Name:           name,
				Operators:      []*Operator{},
				flagParamsMap:  make(map[string][]string),
				handlers:       []func(string, *DCommand) error{},
				paramsHandlers: []func(string, *DCommand, ...interface{}) error{},
			},
		}
		fc.currentCommand = name
		fc.Commands = cmd
		fc.currentOperator = ""
	} else {
		fc.Commands = append(fc.Commands, &Command{
			Name:           name,
			Operators:      []*Operator{},
			flagParamsMap:  make(map[string][]string),
			handlers:       []func(string, *DCommand) error{},
			paramsHandlers: []func(string, *DCommand, ...interface{}) error{},
		})
		fc.currentCommand = name
		fc.currentOperator = ""
	}

	return fc
}
