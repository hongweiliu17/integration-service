package v1alpha1

import (
	"github.com/redhat-appstudio/integration-service/api/v2alpha1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this Memcached to the Hub version (v2alpha1).
func (src *IntegrationTestScenario) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v2alpha1.IntegrationTestScenario)
	// TODO- hard-coded now we need to convert the Resolver and Params from v1alpha1
	dst.Spec.ResolverRef = v2alpha1.ResolverRef{
		Resolver: "git",
		Params: []v2alpha1.PipelineParameter{
			{
				Name:  "url",
				Value: "http://url",
			},
			{
				Name:  "pipeline",
				Value: src.Spec.Pipeline,
			},
			{
				Name:  "bundle",
				Value: src.Spec.Bundle,
			},
		},
	}
	return nil
}

// ConvertFrom converts from the Hub version (v2alpha1) to this version.
func (dst *IntegrationTestScenario) ConvertFrom(srcRaw conversion.Hub) error {
	return nil
}
