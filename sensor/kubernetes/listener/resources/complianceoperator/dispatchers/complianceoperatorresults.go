package dispatchers

import (
	"github.com/ComplianceAsCode/compliance-operator/pkg/apis/compliance/v1alpha1"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/centralsensor"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/sensor/common/centralcaps"
	"github.com/stackrox/rox/sensor/kubernetes/eventpipeline/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ocpComplianceLabelsKey = "compliance.openshift.io/"
)

var (
	scanNameKey  = ocpComplianceLabelsKey + "scan-name"
	suiteNameKey = ocpComplianceLabelsKey + "suite"
)

// ResultDispatcher handles compliance check result objects
type ResultDispatcher struct{}

// NewResultDispatcher creates and returns a new compliance check result dispatcher.
func NewResultDispatcher() *ResultDispatcher {
	return &ResultDispatcher{}
}

func statusToProtoStatus(status v1alpha1.ComplianceCheckStatus) storage.ComplianceOperatorCheckResult_CheckStatus {
	switch status {
	case v1alpha1.CheckResultPass:
		return storage.ComplianceOperatorCheckResult_PASS
	case v1alpha1.CheckResultFail:
		return storage.ComplianceOperatorCheckResult_FAIL
	case v1alpha1.CheckResultInfo:
		return storage.ComplianceOperatorCheckResult_INFO
	case v1alpha1.CheckResultManual:
		return storage.ComplianceOperatorCheckResult_MANUAL
	case v1alpha1.CheckResultError:
		return storage.ComplianceOperatorCheckResult_ERROR
	case v1alpha1.CheckResultNotApplicable:
		return storage.ComplianceOperatorCheckResult_NOT_APPLICABLE
	case v1alpha1.CheckResultInconsistent:
		return storage.ComplianceOperatorCheckResult_INCONSISTENT
	default:
		return storage.ComplianceOperatorCheckResult_UNSET
	}
}

func statusToV2Status(status v1alpha1.ComplianceCheckStatus) central.ComplianceOperatorCheckResultV2_CheckStatus {
	switch status {
	case v1alpha1.CheckResultPass:
		return central.ComplianceOperatorCheckResultV2_PASS
	case v1alpha1.CheckResultFail:
		return central.ComplianceOperatorCheckResultV2_FAIL
	case v1alpha1.CheckResultInfo:
		return central.ComplianceOperatorCheckResultV2_INFO
	case v1alpha1.CheckResultManual:
		return central.ComplianceOperatorCheckResultV2_MANUAL
	case v1alpha1.CheckResultError:
		return central.ComplianceOperatorCheckResultV2_ERROR
	case v1alpha1.CheckResultNotApplicable:
		return central.ComplianceOperatorCheckResultV2_NOT_APPLICABLE
	case v1alpha1.CheckResultInconsistent:
		return central.ComplianceOperatorCheckResultV2_INCONSISTENT
	default:
		return central.ComplianceOperatorCheckResultV2_UNSET
	}
}

func getScanName(labels map[string]string) string {
	if value, ok := labels[scanNameKey]; ok {
		return value
	}

	return ""
}

func getSuiteName(labels map[string]string) string {
	if value, ok := labels[suiteNameKey]; ok {
		return value
	}

	return ""
}

// ProcessEvent processes a compliance operator check result
func (c *ResultDispatcher) ProcessEvent(obj, _ interface{}, action central.ResourceAction) *component.ResourceEvent {
	var complianceCheckResult v1alpha1.ComplianceCheckResult

	unstructuredObject, ok := obj.(*unstructured.Unstructured)
	if !ok {
		log.Errorf("Not of type 'unstructured': %T", obj)
		return nil
	}

	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObject.Object, &complianceCheckResult); err != nil {
		log.Errorf("error converting unstructured to compliance check result: %v", err)
		return nil
	}

	id := string(complianceCheckResult.UID)
	if centralcaps.Has(centralsensor.ComplianceV2Integrations) {
		events := []*central.SensorEvent{
			{
				Id:     id,
				Action: action,
				Resource: &central.SensorEvent_ComplianceOperatorResultV2{
					ComplianceOperatorResultV2: &central.ComplianceOperatorCheckResultV2{
						Id:           id,
						CheckId:      complianceCheckResult.ID,
						CheckName:    complianceCheckResult.GetName(),
						Status:       statusToV2Status(complianceCheckResult.Status),
						Severity:     severityToV2Severity(complianceCheckResult.Severity),
						Description:  complianceCheckResult.Description,
						Instructions: complianceCheckResult.Instructions,
						Labels:       complianceCheckResult.GetLabels(),
						Annotations:  complianceCheckResult.GetAnnotations(),
						CreatedTime:  protoconv.ConvertTimeToTimestamp(complianceCheckResult.GetCreationTimestamp().Time),
						ScanName:     getScanName(complianceCheckResult.GetLabels()),
						SuiteName:    getSuiteName(complianceCheckResult.GetLabels()),
						Rationale:    complianceCheckResult.Rationale,
						ValuesUsed:   complianceCheckResult.ValuesUsed,
						Warnings:     complianceCheckResult.Warnings,
					},
				},
			},
		}
		return component.NewEvent(events...)
	}

	events := []*central.SensorEvent{
		{
			Id:     id,
			Action: action,
			Resource: &central.SensorEvent_ComplianceOperatorResult{
				ComplianceOperatorResult: &storage.ComplianceOperatorCheckResult{
					Id:           id,
					CheckId:      complianceCheckResult.ID,
					CheckName:    complianceCheckResult.Name,
					Status:       statusToProtoStatus(complianceCheckResult.Status),
					Description:  complianceCheckResult.Description,
					Instructions: complianceCheckResult.Instructions,
					Labels:       complianceCheckResult.Labels,
					Annotations:  complianceCheckResult.Annotations,
				},
			},
		},
	}
	return component.NewEvent(events...)
}
