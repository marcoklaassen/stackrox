package google

import (
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1alpha1"
)

func (c *googleScanner) convertComponent(occurrence *containeranalysis.Occurrence) (string, *v1.ImageScanComponent) {
	location := occurrence.GetInstallation().GetLocation()[0]
	version := location.GetVersion()
	component := &v1.ImageScanComponent{
		Name:    occurrence.GetInstallation().GetName(),
		Version: version.GetName() + "-" + version.GetRevision(),
	}
	return location.GetCpeUri(), component
}

func (c *googleScanner) convertVulnerability(occurrence *containeranalysis.Occurrence, note *containeranalysis.Note) (string, string, *v1.Vulnerability) {
	packageIssues := occurrence.GetVulnerabilityDetails().GetPackageIssue()
	if len(packageIssues) == 0 {
		return "", "", nil
	}
	affectedLocation := packageIssues[0].GetAffectedLocation()
	var link string
	if len(note.GetRelatedUrl()) != 0 {
		link = note.GetRelatedUrl()[0].GetUrl()
	}
	summary := c.getSummary(note.GetVulnerabilityType().GetDetails(), affectedLocation.GetCpeUri())
	return affectedLocation.GetCpeUri(),
		affectedLocation.GetPackage(),
		&v1.Vulnerability{
			Cve:     note.GetShortDescription(),
			Link:    link,
			Cvss:    occurrence.GetVulnerabilityDetails().GetCvssScore(),
			Summary: summary,
		}
}

// getSummary searches through the details and returns the summary that matches the cpeURI
func (c googleScanner) getSummary(details []*containeranalysis.VulnerabilityType_Detail, cpeURI string) string {
	for _, detail := range details {
		if detail.CpeUri == cpeURI {
			return detail.Description
		}
	}
	return ""
}
