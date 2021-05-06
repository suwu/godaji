package workflow

import (
	"sync"
	"time"
)

// workflow 的 status
// INACTIVE 工作流未激活
// DONE 工作流已完成
const (
	FLOW_INACTIVE = iota
	FLOW_ACTIVE
	FLOW_PAUSE
	FLOW_STOP
	FLOW_DONE
)

const (
	TASK_INACTIVE = iota
	TASK_ACTIVE
	TASK_DONE
)

const (
	TRAN_INACTIVE = iota
	TRAN_ACTIVE
)

const (
	ACTION_INACTIVE = iota
	ACTION_ACTIVE
	ACTION_DONE
)

const (
	JOIN_ALL = "all"
	JOIN_ANY = "any"
	JOIN_ONE = "one"
)

type Transition struct {
	Condition Condition
	Target    string
}

type Task struct {
	Name        string
	Description string
	Action      Action
	Repeat      bool
	Join        string
	Transitions []Transition
}

type Action string
type Condition string

type WorkFlow struct {
	Name        string
	Description string
	Input       Input
	Output      Output
	Tasks       []Task
}

type Input map[string]interface{}
type Output map[string]interface{}

type WorkFlowInstance struct {
	Id            string
	WorkFlow      *WorkFlow
	TaskInstances []*TaskInstance
	Status        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Mu            sync.Mutex
	Env           map[string]interface{}
}

type TaskInstance struct {
	Id               string
	Task             *Task
	WorkFlowInstance *WorkFlowInstance
	ActionInstance   *ActionInstance
	Status           int
	TaskResult       map[string]interface{}
	InputTransitions []*TransitionInstance
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type TransitionInstance struct {
	Id        string
	Source    *Task
	Target    *Task
	Condition Condition
	Results   []bool
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ActionInstance struct {
	Id               string
	TaskInstance     *TaskInstance
	WorkFlowInstance *WorkFlowInstance
	Action           Action
	Status           int
	ActionResult     map[string]interface{}
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
