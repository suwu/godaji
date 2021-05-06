package workflow

func Instance(name string) (inst *WorkFlowInstance, err error) {
	workflow, err := GetWorkFlowByName(name)
	if err != nil {
		return
	}
	inst = workflow.Instance()
	return
}

func Start(id string) (err error) {
	inst, err := GetWorkFlowInstance(id)
	if err != nil {
		return
	}
	inst.Start()
	return
}

// func Pause(id string) error {
// 	inst, err := GetWorkFlowInstance(id)
// 	if err != nil {
// 		return err
// 	}
// 	inst.Pause()
// 	return nil
// }

// func Resume(id string) error {
// 	inst, err := GetWorkFlowInstance(id)
// 	if err != nil {
// 		return err
// 	}
// 	inst.Resume()
// 	return nil
// }

// func Stop(id string) error {
// 	inst, err := GetWorkFlowInstance(id)
// 	if err != nil {
// 		return err
// 	}
// 	inst.Stop()
// 	return nil
// }
