package common

type Condition struct {
	ConditionType int32
	Params        []int32
}

type Conditions []Condition

func NewCondition(conditionType int32, params ...int32) *Condition {
	return &Condition{
		ConditionType: conditionType,
		Params:        params,
	}
}

func NewConditions() *Conditions {
	return &Conditions{}
}

func (c *Conditions) AddCondition(condition *Condition) {
	*c = append(*c, *condition)
}

func (c *Conditions) Empty() bool {
	return len(*c) <= 0
}

type CompoundConditions struct {
	And *Conditions
	Or  []*Conditions
}

func NewCompoundConditions(And *Conditions, Or ...*Conditions) *CompoundConditions {
	return &CompoundConditions{
		And: And,
		Or:  Or,
	}
}
