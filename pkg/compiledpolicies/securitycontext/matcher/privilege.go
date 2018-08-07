package matcher

import (
	"fmt"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
)

func init() {
	compilers = append(compilers, newPrivilegeMatcher)
}

func newPrivilegeMatcher(policy *v1.Policy) (Matcher, error) {
	fields := policy.GetFields()
	if fields.GetSetPrivileged() == nil {
		return nil, nil
	}

	privileged := fields.GetPrivileged()
	matcher := &privilegeMatcherImpl{privileged: &privileged}
	return matcher.match, nil
}

type privilegeMatcherImpl struct {
	privileged *bool
}

func (p *privilegeMatcherImpl) match(security *v1.SecurityContext) []*v1.Alert_Violation {
	if security == nil || security.GetPrivileged() != *p.privileged {
		return nil
	}

	var violations []*v1.Alert_Violation
	violations = append(violations, &v1.Alert_Violation{
		Message: fmt.Sprintf("Container privileged set to %t matched configured policy", security.GetPrivileged()),
	})
	return violations
}
