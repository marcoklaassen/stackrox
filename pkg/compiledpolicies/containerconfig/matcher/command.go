package matcher

import (
	"fmt"
	"regexp"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/compiledpolicies/utils"
)

func init() {
	compilers = append(compilers, newCommandMatcher)
}

func newCommandMatcher(policy *v1.Policy) (Matcher, error) {
	commands := policy.GetFields().GetCommand()
	if commands == "" {
		return nil, nil
	}

	commandsRegex, err := utils.CompileStringRegex(commands)
	if err != nil {
		return nil, err
	}
	matcher := &commandMatcherImpl{commandsRegex}
	return matcher.match, nil
}

type commandMatcherImpl struct {
	commandsRegex *regexp.Regexp
}

func (p *commandMatcherImpl) match(config *v1.ContainerConfig) []*v1.Alert_Violation {
	var violations []*v1.Alert_Violation
	if !p.matchArg(config.GetCommand()) {
		v := &v1.Alert_Violation{
			Message: fmt.Sprintf("Commands matched configs policy: %s", p.commandsRegex),
		}
		violations = append(violations, v)
	}
	return violations
}

func (p *commandMatcherImpl) matchArg(commands []string) bool {
	for _, arg := range commands {
		if p.commandsRegex.MatchString(arg) {
			return true
		}
	}
	return false
}
