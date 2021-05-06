package workflow

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/antonmedv/expr"
	uuid "github.com/satori/go.uuid"
)

var workflows []*WorkFlow = make([]*WorkFlow, 0)
var workflowInstances []*WorkFlowInstance = make([]*WorkFlowInstance, 0)

func GetWorkFlowByName(name string) (*WorkFlow, error) {
	for _, workflow := range workflows {
		if workflow.Name == name {
			return workflow, nil
		}
	}
	return nil, fmt.Errorf("WorkFlow %s is not found.", name)
}

func GetWorkFlowInstance(id string) (*WorkFlowInstance, error) {
	for _, inst := range workflowInstances {
		if inst.Id == id {
			return inst, nil
		}
	}
	return nil, fmt.Errorf("WorkFlowInstance %s is not found.", id)
}

// workflow
func (wf *WorkFlow) Instance() (inst *WorkFlowInstance) {
	inst = &WorkFlowInstance{
		Id:            uuid.NewV4().String(),
		WorkFlow:      wf,
		TaskInstances: make([]*TaskInstance, 0),
		Status:        FLOW_INACTIVE,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Mu:            sync.Mutex{},
		Env:           make(map[string]interface{}),
	}

	// 系统函数
	for k, v := range GetFunctionMap() {
		inst.Env[k] = v
	}

	for k, v := range wf.Input {
		inst.Env[k] = v
	}

	inst.Env["currentWorkFlow"] = inst

	workflowInstances = append(workflowInstances, inst)
	return
}

func (inst *WorkFlowInstance) Start() (err error) {
	if inst.Status != FLOW_INACTIVE {
		err = fmt.Errorf("workflow %s instance id %s had been started.", inst.WorkFlow.Name, inst.Id)
		return
	}
	task := inst.WorkFlow.Tasks[0]
	taskInst, err := inst.InstanceTask(&task)

	return taskInst.Excute()
}

func (inst *WorkFlowInstance) InstanceTask(task *Task) (taskInst *TaskInstance, err error) {
	taskInst = &TaskInstance{
		Id:               uuid.NewV4().String(),
		Task:             task,
		WorkFlowInstance: inst,
		Status:           TASK_INACTIVE,
		TaskResult:       make(map[string]interface{}),
		InputTransitions: make([]*TransitionInstance, 0),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	for idx, sourceTask := range inst.WorkFlow.Tasks {
		for _, tran := range sourceTask.Transitions {
			if tran.Target == task.Name {
				tranInst, _ := tran.Instance(&inst.WorkFlow.Tasks[idx], task)
				taskInst.InputTransitions = append(taskInst.InputTransitions, tranInst)
			}
		}
	}

	inst.Mu.Lock()
	inst.TaskInstances = append(inst.TaskInstances, taskInst)
	inst.Mu.Unlock()

	return
}

func (inst *WorkFlowInstance) GetTaskInstancesByName(name string) (taskInstances []*TaskInstance, err error) {
	taskInstances = make([]*TaskInstance, 0)
	for _, t := range inst.TaskInstances {
		if t.Task.Name == name {
			taskInstances = append(taskInstances, t)
		}
	}

	err = fmt.Errorf("Task %s instance not found in the workflow %s, id %s", name, inst.WorkFlow.Name, inst.Id)
	return
}

func (inst *WorkFlowInstance) GetTaskByName(name string) (task *Task, err error) {
	for idx, t := range inst.WorkFlow.Tasks {
		if t.Name == name {
			task = &inst.WorkFlow.Tasks[idx]
			return
		}
	}

	err = fmt.Errorf("Task %s not found in the workflow %s, id %s", name, inst.WorkFlow.Name, inst.Id)
	return
}

// task
func (inst *TaskInstance) InstanceAction() (actionInst *ActionInstance) {
	actionInst = &ActionInstance{
		Id:               uuid.NewV4().String(),
		TaskInstance:     inst,
		WorkFlowInstance: inst.WorkFlowInstance,
		Action:           inst.Task.Action,
		Status:           ACTION_INACTIVE,
		ActionResult:     make(map[string]interface{}),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	inst.ActionInstance = actionInst
	return
}

func (inst *TaskInstance) GetTransition(source *Task) (tranInst *TransitionInstance, err error) {
	for _, tran := range inst.InputTransitions {
		if tran.Source.Name == source.Name {
			tranInst = tran
			return
		}
	}

	err = fmt.Errorf("Task %s not found transition from task %s.", inst.Task.Name, source.Name)
	return
}

func (inst *TaskInstance) Excute() (err error) {
	inst.Status = TASK_ACTIVE
	actionInst := inst.InstanceAction()

	err = actionInst.Run()
	return
}

func (inst *TaskInstance) ExcuteNextTask() (err error) {
	for _, tran := range inst.Task.Transitions {
		taskInstances, err := inst.WorkFlowInstance.GetTaskInstancesByName(tran.Target)
		var taskInst *TaskInstance
		var dropFlag = false
		for _, t := range taskInstances {
			if t.Status == TASK_INACTIVE {
				taskInst = t
				break
			}
			if t.Status == TASK_ACTIVE {
				// 已存在已激活的任务
				// 这条转移丢弃
				dropFlag = true
			}
			if t.Status == TASK_DONE {
				continue
			}
		}
		if dropFlag {
			continue
		}

		if taskInst == nil {
			task, err := inst.WorkFlowInstance.GetTaskByName(tran.Target)
			if err != nil {
				return err
			}
			taskInst, err = inst.WorkFlowInstance.InstanceTask(task)
		}

		tranInst, err := taskInst.GetTransition(inst.Task)
		if err != nil {
			return err
		}

		result, err := tranInst.Condition.Eval(inst)
		if err != nil {
			return err
		}
		tranInst.Status = TRAN_ACTIVE

		tranInst.Results = append(tranInst.Results, result)

		// 决策执行下一任务
		switch taskInst.Task.Join {
		case JOIN_ALL:
			var flag bool = true
			for _, t := range taskInst.InputTransitions {
				if t.Status != TRAN_ACTIVE {
					flag = false
					break
				}
				if !t.Results[0] {
					flag = false
					break
				}
			}
			if flag {
				taskInst.Excute()
				for _, t := range taskInst.InputTransitions {
					t.Results = t.Results[1:]
				}
			}
		case JOIN_ANY:
			var flag bool = false
			for _, t := range taskInst.InputTransitions {
				if t.Status == TRAN_ACTIVE {
					result := t.Results[0]
					if result {
						flag = true
						break
					}
				}
			}
			if flag {
				taskInst.Excute()
				for _, t := range taskInst.InputTransitions {
					if t.Status == TRAN_ACTIVE {
						result := t.Results[0]
						if result {
							t.Results = t.Results[1:]
						}
					}
				}
			}
		case JOIN_ONE:
			var flag bool = false
			var num int = 0
			for _, t := range taskInst.InputTransitions {
				if t.Status == TRAN_ACTIVE {
					result := t.Results[0]
					if result {
						flag = true
						num++
					}
				}
			}
			if flag && num == 1 {
				taskInst.Excute()
				for _, t := range taskInst.InputTransitions {
					if t.Status == TRAN_ACTIVE {
						result := t.Results[0]
						if result {
							t.Results = t.Results[1:]
						}
					}
				}
			}
		default:
			var flag bool = false
			var num int = 0
			for _, t := range taskInst.InputTransitions {
				if t.Status == TRAN_ACTIVE {
					result := t.Results[0]
					if result {
						flag = true
						num++
					}
				}
			}

			var joinNum int
			joinNum, err = strconv.Atoi(taskInst.Task.Join)
			if err != nil {
				return err
			}
			if flag && num == joinNum {
				taskInst.Excute()
				for _, t := range taskInst.InputTransitions {
					if t.Status == TRAN_ACTIVE {
						result := t.Results[0]
						if result {
							t.Results = t.Results[1:]
						}
					}
				}
			}
		}
	}
	inst.Status = TASK_DONE
	return
}

// transition
func (tran *Transition) Instance(sourceTask *Task, targetTask *Task) (tranInst *TransitionInstance, err error) {
	tranInst = &TransitionInstance{
		Id:        uuid.NewV4().String(),
		Source:    sourceTask,
		Target:    targetTask,
		Condition: tran.Condition,
		Results:   make([]bool, 0),
		Status:    TRAN_INACTIVE,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return
}

// action
func (inst *ActionInstance) Run() (err error) {
	env := inst.WorkFlowInstance.Env
	for k, v := range inst.GetActionMap() {
		env[k] = v
	}

	_, err = expr.Eval(string(inst.Action), env)

	return
}

func (inst *ActionInstance) SetResult(result map[string]interface{}) (err error) {
	inst.ActionResult = result
	inst.TaskInstance.TaskResult = result
	inst.Status = ACTION_DONE
	inst.TaskInstance.Status = TASK_DONE

	err = inst.TaskInstance.ExcuteNextTask()
	return
}

// condition
func (cond Condition) Eval(taskInst *TaskInstance) (result bool, err error) {
	env := taskInst.WorkFlowInstance.Env
	env["currentTask"] = taskInst

	if taskInst.TaskResult != nil {
		for k, v := range taskInst.TaskResult {
			env[k] = v
		}
	}

	out, err := expr.Eval(string(cond), env)

	result = out.(bool)
	if err != nil {
		return
	}
	return
}
