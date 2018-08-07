package matcher

import (
	"fmt"
	"regexp"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/compiledpolicies/utils"
)

func init() {
	compilers = append(compilers, newArgsMatcher)
}

func newArgsMatcher(policy *v1.Policy) (Matcher, error) {
	args := policy.GetFields().GetArgs()
	if args == "" {
		return nil, nil
	}

	argsRegex, err := utils.CompileStringRegex(args)
	if err != nil {
		return nil, err
	}

	matcher := &argsMatcherImpl{argsRegex}
	return matcher.match, nil
}

type argsMatcherImpl struct {
	argsRegex *regexp.Regexp
}

func (p *argsMatcherImpl) match(config *v1.ContainerConfig) []*v1.Alert_Violation {
	var violations []*v1.Alert_Violation
	if !p.matchArg(config.GetArgs()) {
		v := &v1.Alert_Violation{
			Message: fmt.Sprintf("Args matched configs policy: %s", p.argsRegex),
		}
		violations = append(violations, v)
	}
	return violations
}

func (p *argsMatcherImpl) matchArg(args []string) bool {
	for _, arg := range args {
		if p.argsRegex.MatchString(arg) {
			return true
		}
	}
	return false
}
