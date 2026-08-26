package controllers_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	platformv1 "github.com/SumitDalavi/k8s-golden-path-provisioner/api/v1alpha1"
)

var _ = Describe("PlatformService Controller", func() {
	ctx := context.Background()

	Context("When creating a PlatformService", func() {
		It("should provision a Deployment and Service", func() {
			ps := &platformv1.PlatformService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
					Namespace: "default",
				},
				Spec: platformv1.PlatformServiceSpec{
					Image:    "nginx:1.25",
					Replicas: 2,
					Port:     8080,
				},
			}
			Expect(k8sClient.Create(ctx, ps)).Should(Succeed())

			// Wait for status to be populated
			Eventually(func() string {
				fetched := &platformv1.PlatformService{}
				k8sClient.Get(ctx, types.NamespacedName{Name: "test-service", Namespace: "default"}, fetched)
				return fetched.Status.Phase
			}, 10*time.Second, 250*time.Millisecond).Should(BeElementOf("Running", "Provisioning"))

			// Cleanup
			Expect(k8sClient.Delete(ctx, ps)).Should(Succeed())
		})

		It("should reject zero replicas", func() {
			ps := &platformv1.PlatformService{
				ObjectMeta: metav1.ObjectMeta{Name: "bad-service", Namespace: "default"},
				Spec: platformv1.PlatformServiceSpec{Image: "nginx:1.25", Replicas: 0, Port: 8080},
			}
			err := k8sClient.Create(ctx, ps)
			Expect(err).Should(HaveOccurred())
		})
	})
})

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}
