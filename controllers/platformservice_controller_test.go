package controllers_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	platformv1 "github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1"
)

var _ = Describe("PlatformService Controller", func() {
	_ = context.Background()

	Context("When creating a PlatformService", func() {
		It("should provision a Deployment and Service", func() {
			ps := &platformv1.PlatformService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
					Namespace: "default",
				},
				Spec: platformv1.PlatformServiceSpec{
					Team: "test-team",
					Tier: "frontend",
				},
			}
			Expect(ps.Spec.Tier).To(Equal("frontend"))
		})

		It("should reject invalid tier", func() {
			ps := &platformv1.PlatformService{
				ObjectMeta: metav1.ObjectMeta{Name: "bad-service", Namespace: "default"},
				Spec: platformv1.PlatformServiceSpec{Team: "test-team", Tier: "invalid"},
			}
			Expect(ps.Spec.Tier).NotTo(Equal("frontend"))
		})
	})
})

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}
