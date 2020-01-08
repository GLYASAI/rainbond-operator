package status

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	rainbondv1alpha1 "github.com/GLYASAI/rainbond-operator/pkg/apis/rainbond/v1alpha1"
	"github.com/GLYASAI/rainbond-operator/pkg/controller/rainbondcluster/pkg/mock"
	"github.com/GLYASAI/rainbond-operator/pkg/controller/rainbondcluster/types"
)

func TestGenerateRainbondClusterPackageExtractedCondition(t *testing.T) {
	tests := []struct {
		name            string
		rainbondcluster *rainbondv1alpha1.RainbondCluster
		ehistory        *types.ExtractionHistory
		err             error
		want            rainbondv1alpha1.RainbondClusterCondition
	}{
		{
			name: "without installation package",
			rainbondcluster: &rainbondv1alpha1.RainbondCluster{
				Spec: rainbondv1alpha1.RainbondClusterSpec{
					InstallMode: rainbondv1alpha1.InstallationModeWithoutPackage,
				},
			},
			want: rainbondv1alpha1.RainbondClusterCondition{
				Type:   rainbondv1alpha1.PackageExtracted,
				Status: rainbondv1alpha1.ConditionTrue,
				Reason: string(rainbondv1alpha1.InstallationModeWithoutPackage),
			},
		},
		{
			name: "error fetching extraction history",
			rainbondcluster: &rainbondv1alpha1.RainbondCluster{
				Spec: rainbondv1alpha1.RainbondClusterSpec{},
			},
			err: errors.New("foobar"),
			want: rainbondv1alpha1.RainbondClusterCondition{
				Type:   rainbondv1alpha1.PackageExtracted,
				Status: rainbondv1alpha1.ConditionFalse,
				Reason: ErrHistoryFetch,
			},
		},
		{
			name: "history status is false",
			rainbondcluster: &rainbondv1alpha1.RainbondCluster{
				Spec: rainbondv1alpha1.RainbondClusterSpec{},
			},
			ehistory: &types.ExtractionHistory{
				Status: types.HistoryStatusFalse,
				Reason: "foobar",
			},
			want: rainbondv1alpha1.RainbondClusterCondition{
				Type:   rainbondv1alpha1.PackageExtracted,
				Status: rainbondv1alpha1.ConditionFalse,
				Reason: "foobar",
			},
		},
		{
			name: "ok",
			rainbondcluster: &rainbondv1alpha1.RainbondCluster{
				Spec: rainbondv1alpha1.RainbondClusterSpec{},
			},
			ehistory: &types.ExtractionHistory{
				Status: types.HistoryStatusTrue,
			},
			want: rainbondv1alpha1.RainbondClusterCondition{
				Type:   rainbondv1alpha1.PackageExtracted,
				Status: rainbondv1alpha1.ConditionTrue,
			},
		},
		{
			name: "already extracted",
			rainbondcluster: &rainbondv1alpha1.RainbondCluster{
				Spec: rainbondv1alpha1.RainbondClusterSpec{},
				Status: &rainbondv1alpha1.RainbondClusterStatus{
					Conditions: []rainbondv1alpha1.RainbondClusterCondition{
						{
							Type:   rainbondv1alpha1.PackageExtracted,
							Status: rainbondv1alpha1.ConditionTrue,
						},
					},
				},
			},
			want: rainbondv1alpha1.RainbondClusterCondition{
				Type:   rainbondv1alpha1.PackageExtracted,
				Status: rainbondv1alpha1.ConditionTrue,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			historyer := mock.NewMockHistoryInterface(ctrl)
			historyer.EXPECT().ExtractionHistory().Return(tc.ehistory, tc.err).AnyTimes()

			got := GenerateRainbondClusterPackageExtractedCondition(tc.rainbondcluster, historyer)
			assert.Equal(t, tc.want.Type, got.Type)
			assert.Equal(t, tc.want.Status, got.Status)
			assert.Equal(t, tc.want.Reason, got.Reason)
		})
	}
}

func TestGenerateRainbondClusterImagesLoadedCondition(t *testing.T) {
	tests := []struct {
		name            string
		rainbondcluster *rainbondv1alpha1.RainbondCluster
		images          []string
		err             error
		want            rainbondv1alpha1.RainbondClusterCondition
	}{
		{
			name: "Without installation package",
			rainbondcluster: &rainbondv1alpha1.RainbondCluster{
				Spec: rainbondv1alpha1.RainbondClusterSpec{
					InstallMode: rainbondv1alpha1.InstallationModeWithoutPackage,
				},
			},
			want: rainbondv1alpha1.RainbondClusterCondition{
				Type:   rainbondv1alpha1.ImagesLoaded,
				Status: rainbondv1alpha1.ConditionTrue,
				Reason: string(rainbondv1alpha1.InstallationModeWithoutPackage),
			},
		},
		{
			name: "Error getting metadata",
			rainbondcluster: &rainbondv1alpha1.RainbondCluster{
				Spec: rainbondv1alpha1.RainbondClusterSpec{},
			},
			err: fmt.Errorf("foobar"),
			want: rainbondv1alpha1.RainbondClusterCondition{
				Type:   rainbondv1alpha1.ImagesLoaded,
				Status: rainbondv1alpha1.ConditionFalse,
				Reason: ErrGetMetadata,
			},
		},
		{
			name: "Already loaded",
			rainbondcluster: &rainbondv1alpha1.RainbondCluster{
				Spec: rainbondv1alpha1.RainbondClusterSpec{},
				Status: &rainbondv1alpha1.RainbondClusterStatus{
					Conditions: []rainbondv1alpha1.RainbondClusterCondition{
						{
							Type:   rainbondv1alpha1.ImagesLoaded,
							Status: rainbondv1alpha1.ConditionTrue,
						},
					},
				},
			},
			want: rainbondv1alpha1.RainbondClusterCondition{
				Type:   rainbondv1alpha1.ImagesLoaded,
				Status: rainbondv1alpha1.ConditionTrue,
			},
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			packager := mock.NewMockPackageInterface(ctrl)
			packager.EXPECT().GetMetadata().Return(tc.images, tc.err).AnyTimes()

			got := GenerateRainbondClusterImagesLoadedCondition(tc.rainbondcluster, packager)
			assert.Equal(t, tc.want.Type, got.Type)
			assert.Equal(t, tc.want.Status, got.Status)
			assert.Equal(t, tc.want.Reason, got.Reason)
		})
	}
}