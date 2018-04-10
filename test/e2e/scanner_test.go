package e2e_test

import (
	"os"

	"github.com/soter/scanner/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("Image Scanner", func() {
	var (
		f *framework.Invocation

		labels          map[string]string
		name, namespace string
		containers1     []core.Container
		containers2     []core.Container
		containers3     []core.Container
		data1, data2            string
		skip1, skip2            bool

		secret1, secret2                    *core.Secret
		service, svc              *core.Service
		err error

		obj runtime.Object
		str1, str2, str3 string

		ctx1 = func (workloadCode int) {
			Context("When some images are vulnerable", func() {
				BeforeEach(func() {
					By("Creating secret-1")
					_, err := root.KubeClient.CoreV1().Secrets(secret1.Namespace).Create(secret1)
					Expect(err).NotTo(HaveOccurred())

					obj = f.NewWorkload(name, namespace, labels, containers1, secret1.Name, workloadCode)
				})

				It("Shouldn't be created", func() {
					if skip1 {
						Skip("environment var \"DOCKER_CFG_1\" not found")
					}

					By("Should contain vulnerabilities")
					f.EventuallyCreateWithVulnerableImage(root, obj).Should(Equal(true))
				})
			})
		}

		ctx2 = func (workloadCode int, update bool) {
			str1 = "When no image is vulnerable"
			if update {
				str1 = "When some images are vulnerable"
			}
			Context(str1, func() {
				BeforeEach(func() {
					if update {
						By("Creating secret-1")
						_, err := root.KubeClient.CoreV1().Secrets(secret1.Namespace).Create(secret1)
						Expect(err).NotTo(HaveOccurred())
					}

					By("Creating secret-2")
					_, err := root.KubeClient.CoreV1().Secrets(secret2.Namespace).Create(secret2)
					Expect(err).NotTo(HaveOccurred())

					obj = f.NewWorkload(name, namespace, labels, containers2, secret2.Name, workloadCode)
				})

				AfterEach(func() {
					f.DeleteWorkloads(workloadCode)
				})

				str2 = "Should be created"
				if update {
					str2 = "Shouldn't be updated"
				}
				It(str2, func() {
					if skip2 {
						Skip("environment var \"DOCKER_CFG_2\" not found")
					}

					str3 = "No vulnerabilities"
					if update {
						str3 = "Creating"
					}
					By(str3)
					f.EventuallyCreateWithNonVulnerableImage(root, obj).ShouldNot(HaveOccurred())

					if update {
						By("Updating")
						obj = f.NewWorkload(name, namespace, labels, containers1, secret1.Name, workloadCode)
						f.EventuallyUpdateWithVulnerableImage(root, obj).Should(Equal(true))
					}
				})
			})
		}
	)

	BeforeEach(func() {
		f = root.Invoke()
		name = f.App()
		namespace = f.Namespace()
		labels = map[string]string{
			"app": f.App(),
		}

		data1 = ""
		skip1 = false
		if val, ok := os.LookupEnv("DOCKER_CFG_1"); !ok {
			skip1 = true
		} else {
			data1 = val
			secret1 = f.NewSecret(name + "-secret-1", namespace, data1, labels)
		}

		data2 = ""
		skip2 = false
		if val, ok := os.LookupEnv("DOCKER_CFG_2"); !ok {
			skip2 = true
		} else {
			data2 = val
			secret2 = f.NewSecret(name + "-secret-2", namespace, data2, labels)
		}

		containers1 = []core.Container{
			{
				Name:  "label-practice",
				Image: "shudipta/labels",
				Ports: []core.ContainerPort{
					{
						ContainerPort: 10000,
					},
				},
			},
		}
		containers2 = []core.Container{
			{
				Name:  "hello",
				Image: "alittleprogramming/hello:test",
				Ports: []core.ContainerPort{
					{
						ContainerPort: 80,
					},
				},
			},
		}
		containers3 = []core.Container{
			{
				Name:  "nginx",
				Image: "nginx",
				Ports: []core.ContainerPort{
					{
						ContainerPort: 80,
					},
				},
			},
		}
	})

	Describe("Scan images in Deployment", func() {
		AfterEach(func() {
			f.DeleteAllSecrets()
		})

		Context("Creating Deployment with some vulnerable images", func() {
			ctx1(framework.Deployment)
		})

		Context("Creating Deployment with non-vulnerable images", func() {
			ctx2(framework.Deployment, false)
		})

		Context("Updating Deployment", func() {
			ctx2(framework.Deployment, true)
		})
	})

	Describe("Scan images in ReplicationController", func() {
		AfterEach(func() {
			f.DeleteAllSecrets()
		})

		Context("Creating ReplicationController with some vulnerable images", func() {
			ctx1(framework.ReplicationController)
		})

		Context("Creating ReplicationController with non-vulnerable images", func() {
			ctx2(framework.ReplicationController, false)
		})

		Context("Updating ReplicationController", func() {
			ctx2(framework.ReplicationController, true)
		})
	})

	Describe("Scan images in ReplicaSet", func() {
		AfterEach(func() {
			f.DeleteAllSecrets()
		})

		Context("Creating ReplicaSet with some vulnerable images", func() {
			ctx1(framework.ReplicaSet)
		})

		Context("Creating ReplicaSet with non-vulnerable images", func() {
			ctx2(framework.ReplicaSet, false)
		})

		Context("Updating ReplicaSet", func() {
			ctx2(framework.ReplicaSet, true)
		})
	})

	Describe("Scan images in DaemonSet", func() {
		AfterEach(func() {
			f.DeleteAllSecrets()
		})

		Context("Creating DaemonSet with some vulnerable images", func() {
			ctx1(framework.DaemonSet)
		})

		Context("Creating DaemonSet with non-vulnerable images", func() {
			ctx2(framework.DaemonSet, false)
		})

		Context("Updating DaemonSet", func() {
			ctx2(framework.DaemonSet, true)
		})
	})

	Describe("Scan images in Job", func() {
		AfterEach(func() {
			f.DeleteAllSecrets()
		})

		Context("Creating Job with some vulnerable images", func() {
			ctx1(framework.Job)
		})

		Context("Creating Job with some vulnerable images", func() {
			ctx2(framework.Job, false)
		})

		//Context("Updating Job", func() {
		//	ctx2(framework.Job, true)
		//})
	})

	//Describe("Scan images in CronJob", func() {
	//	AfterEach(func() {
	//		f.DeleteAllSecrets()
	//	})
	//
	//	Context("Creating CronJob with some vulnerable images", func() {
	//		ctx1(framework.CronJob)
	//	})
	//
	//	Context("Creating CronJob with some vulnerable images", func() {
	//		ctx2(framework.CronJob, false)
	//	})
	//
	//	Context("Updating CronJob", func() {
	//		ctx2(framework.CronJob, true)
	//	})
	//})

	Describe("Scan images in StatefulSet", func() {
		BeforeEach(func() {
			By("Creating service")
			service = f.NewService(name, namespace, labels)
			svc, err = root.KubeClient.CoreV1().Services(service.Namespace).Create(service)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			f.DeleteAllSecrets()
			f.DeleteAllServices()
		})

		Context("Creating StatefulSet with some vulnerable images", func() {
			ctx1(framework.StatefulSet)
		})

		Context("Creating StatefulSet with some vulnerable images", func() {
			ctx2(framework.StatefulSet, false)
		})

		Context("Updating StatefulSet", func() {
			ctx2(framework.StatefulSet, true)
		})
	})
})
