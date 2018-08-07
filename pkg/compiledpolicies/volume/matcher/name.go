package matcher

import (
	"fmt"
	"regexp"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/compiledpolicies/utils"
)

func init() {
	compilers = append(compilers, newNameMatcher)
}

func newNameMatcher(policy *v1.Policy) (Matcher, error) {
	name := policy.GetFields().GetVolumePolicy().GetName()
	if name == "" {
		return nil, nil
	}

	nameRegex, err := utils.CompileStringRegex(name)
	if err != nil {
		return nil, err
	}
	matcher := &nameMatcherImpl{nameRegex}
	return matcher.match, nil
}

type nameMatcherImpl struct {
	nameRegex *regexp.Regexp
}

func (p *nameMatcherImpl) match(volume *v1.Volume) []*v1.Alert_Violation {
	var violations []*v1.Alert_Violation
	if p.nameRegex.MatchString(volume.GetName()) {
		v := &v1.Alert_Violation{
			Message: fmt.Sprintf("Volume name matched: %s", p.nameRegex),
		}
		violations = append(violations, v)
	}
	return violations
}
