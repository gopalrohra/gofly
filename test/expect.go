package test

type Expectations struct {
	ShouldHave ShouldHave
}
type ShouldHave struct {
	Key   string
	Value interface{}
}
type ExpectationChecker struct {
	response     map[string]interface{}
	expectations Expectations
}

func (e *ExpectationChecker) shouldHave(s ShouldHave) bool {
	v, ok := e.response[s.Key]
	if !ok {
		return false
	}
	if v == s.Value {
		return true
	}
	return false
}
