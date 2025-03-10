package datastore

import (
	"context"
	"strconv"
	"testing"

	mockIIStore "github.com/stackrox/rox/central/imageintegration/store/mocks"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/features"
	"go.uber.org/mock/gomock"
)

// mustUpdateSourcedAutogenFeature will attempt to update the the sourced autogenerated
// feature to the desired value, if unable to update the feature value will skip the
// the test.
func mustUpdateSourcedAutogenFeature(t *testing.T, desired bool) {
	t.Setenv(features.SourcedAutogeneratedIntegrations.EnvVar(), strconv.FormatBool(desired))

	if features.SourcedAutogeneratedIntegrations.Enabled() != desired {
		t.Log("Skipping test, SourcedAutogeneratedIntegrations feature cannot be set to desired value")
		t.SkipNow()
	}
}

func TestDeleteAutogeneratedRegistries(t *testing.T) {
	ctx := context.Background()

	iis := []*storage.ImageIntegration{
		{Id: "NoAutogen-NoSource", Autogenerated: false, Source: nil},
		{Id: "NoAutogen-Source", Autogenerated: false, Source: &storage.ImageIntegration_Source{}},
		{Id: "Autogen-NoSource", Autogenerated: true, Source: nil},
		{Id: "Autogen-Source", Autogenerated: true, Source: &storage.ImageIntegration_Source{}},
	}

	testCases := []struct {
		name                  string
		expectedDeletes       []string
		autogenEnabled        bool
		sourcedAutogenEnabled bool
	}{
		{
			"delete all autogen integrations when autogen and sourced autogen disabled",
			[]string{"Autogen-NoSource", "Autogen-Source"},
			false,
			false,
		},
		{
			"delete all autogen integrations when autogen disabled and sourced autogen enabled",
			[]string{"Autogen-NoSource", "Autogen-Source"},
			false,
			true,
		},
		{
			"delete only sourced autogen integrations when autogen enabled and sourced autogen disabled",
			[]string{"Autogen-Source"},
			true,
			false,
		},
		{
			"delete no integrations when autogen and sourced autogen enabled",
			[]string{},
			true,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mustUpdateSourcedAutogenFeature(t, tc.sourcedAutogenEnabled)
			t.Setenv(env.AutogeneratedRegistriesDisabled.EnvVar(), strconv.FormatBool(!tc.autogenEnabled))

			ctrl := gomock.NewController(t)
			s := mockIIStore.NewMockStore(ctrl)
			for _, id := range tc.expectedDeletes {
				s.EXPECT().Delete(ctx, id)
			}

			deleteAutogeneratedRegistries(ctx, s, iis)
		})
	}
}
