package workflow

import "fmt"

func (inst *ActionInstance) Echo(value string) {
	fmt.Printf("echo %s\n", value)
	inst.SetResult(nil)
}

func (inst *ActionInstance) GetActionMap() map[string]interface{} {
	return map[string]interface{}{
		"echo": inst.Echo,
	}
}
