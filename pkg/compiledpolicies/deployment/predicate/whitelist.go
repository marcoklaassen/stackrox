package predicate

import (
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/scopecomp"
)

func init() {
	compilers = append(compilers, newWhitelistPredicate)
}

func newWhitelistPredicate(policy *v1.Policy) (Predicate, error) {
	var predicate Predicate
	for _, whitelist := range policy.GetWhitelists() {
		if whitelist.GetDeployment() != nil {
			wrap := &whitelistWrapper{whitelist.GetDeployment()}
			predicate = predicate.And(wrap.shouldProcess)
		}
	}
	return predicate, nil
}

type whitelistWrapper struct {
	whitelist *v1.Whitelist_Deployment
}

func (w *whitelistWrapper) shouldProcess(deployment *v1.Deployment) bool {
	return !MatchesWhitelist(w.whitelist, deployment)
}

// MatchesWhitelist returns true if the given deployment matches the given whitelist.
func MatchesWhitelist(whitelist *v1.Whitelist_Deployment, deployment *v1.Deployment) bool {
	if whitelist == nil {
		return false
	}
	if whitelist.GetScope() != nil && !scopecomp.WithinScope(whitelist.GetScope(), deployment) {
		return false
	}
	if whitelist.GetName() != "" && whitelist.GetName() != deployment.GetName() {
		return false
	}
	return true
}
