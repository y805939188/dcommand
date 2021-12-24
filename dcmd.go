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
	Name          string
	Flags         []*Flag // operator 可以有 flags
	Desc          string
	Passed        bool // 是否传了这个参数
	OwnerCommand  *Command
	flagParamsMap map[string][]string
}

type Command struct {
	Name           string
	Operators      []*Operator
	Flags          []*Flag // command 也可以有 flags
	Desc           string
	handlers       []func(string, *DCommand) error
	paramsHandlers []func(string, *DCommand, ...interface{}) error
	flagParamsMap  map[string][]string
}

type DCommand struct {
	Commands            []*Command
	currentCommand      string
	currentOperator     string
	currentCustomParams []interface{}
}

func (fc *DCommand) Flag(name string, other ...string) *DCommand {
	if name == "" {
		fmt.Println("Need a name")
		return fc
	}
	if fc.Commands == nil {
		fmt.Println("Need to call Command before calling Operator")
		return fc
	}
	if len(fc.Commands) == 0 {
		fmt.Println("There is currently no Command")
		return fc
	}
	for _, cmd := range fc.Commands {
		if fc.currentCommand == cmd.Name {
			if fc.currentOperator == "" {
				// 说明是直接给 command 添加 flag
				if cmd.Flags == nil {
					cmd.Flags = []*Flag{
						{
							Long:         name,
							Params:       []string{},
							OwnerCommand: cmd,
						},
					}
					for i, p := range other {
						switch i {
						case 0:
							cmd.Flags[0].Short = p
							continue
						case 1:
							cmd.Flags[0].Desc = p
						}
					}
				} else {
					temp := &Flag{
						Long:         name,
						Params:       []string{},
						OwnerCommand: cmd,
					}
					for i, p := range other {
						switch i {
						case 0:
							temp.Short = p
							continue
						case 1:
							temp.Desc = p
						}
					}
					cmd.Flags = append(cmd.Flags, temp)
				}
			} else {
				for _, operator := range cmd.Operators {
					if fc.currentOperator == operator.Name {
						// 说明是给 operator 添加 flag
						if operator.Flags == nil || len(operator.Flags) == 0 {
							operator.Flags = []*Flag{
								{
									Long:          name,
									Params:        []string{},
									OwnerOperator: operator,
								},
							}
							for i, p := range other {
								switch i {
								case 0:
									operator.Flags[0].Short = p
									continue
								case 1:
									operator.Flags[0].Desc = p
								}
							}
						} else {
							temp := &Flag{
								Long:          name,
								Params:        []string{},
								OwnerOperator: operator,
							}
							for i, p := range other {
								switch i {
								case 0:
									temp.Short = p
									continue
								case 1:
									temp.Desc = p
								}
							}
							operator.Flags = append(operator.Flags, temp)
						}
						break
					}
				}
			}
			break
		}
	}
	return fc
}

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
		if fc.currentCommand == cmd.Name {
			if cmd.Operators == nil || len(cmd.Operators) == 0 {
				cmd.Operators = []*Operator{
					{
						Name:          name,
						Flags:         nil,
						OwnerCommand:  cmd,
						flagParamsMap: make(map[string][]string),
					},
				}
			} else {
				cmd.Operators = append(cmd.Operators, &Operator{
					Name:          name,
					Flags:         nil,
					OwnerCommand:  cmd,
					flagParamsMap: make(map[string][]string),
				})
			}
			break
		}
	}
	return fc
}

func (fc *DCommand) GetFlagIfExistInCommand(flagName string, isLong bool) []*Flag {
	if fc.Commands == nil {
		return nil
	}

	var res []*Flag
	for _, cmd := range fc.Commands {
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

func (fc *DCommand) GetOperatorIfExist(operatorName string) []*Operator {
	if fc.Commands == nil {
		return nil
	}

	var res []*Operator
	for _, cmd := range fc.Commands {
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

func (fc *DCommand) IsFlag(str string) (bool, isLong bool, pureFlag string) {
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

func (fc *DCommand) SetFlagParamsForCommand(flag string, params []string, command *Command) error {

	isFlag, isLone, pureFlag := fc.IsFlag(flag)

	if !isFlag {
		return fmt.Errorf("The str(%s) is not a flag", flag)
	}

	flags := fc.GetFlagIfExistInCommand(pureFlag, isLone)

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

/**
 * cmd --xxx a b c --yyy d e f
 * cmd a b c --xxx d e f -yyy g h i
 */
func (fc *DCommand) Execute(cmd []string) error {

	commandName := cmd[0]
	cmd = cmd[1:]
	command := fc.GetCommandIfExist(commandName)
	if command == nil {
		return fmt.Errorf("command name is nil")
	}

	if len(cmd) != 0 {
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

	return nil
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
	} else {
		fc.Commands = append(fc.Commands, &Command{
			Name:           name,
			Operators:      []*Operator{},
			flagParamsMap:  make(map[string][]string),
			handlers:       []func(string, *DCommand) error{},
			paramsHandlers: []func(string, *DCommand, ...interface{}) error{},
		})
		fc.currentCommand = name
	}

	return fc
}
