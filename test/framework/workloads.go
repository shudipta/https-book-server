package framework

import (
	"strings"

	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

const (
	Deployment = 1
	ReplicationController = 2
	ReplicaSet = 3
	DaemonSet = 4
	Job = 5
	CronJob = 6
	StatefulSet = 7
)

func int32Ptr(i int32) *int32 { return &i }

func newObjectMeta(name, namespace string, labels map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels: labels,
	}
}

func newSelector(labels map[string]string) *metav1.LabelSelector {
	return &metav1.LabelSelector{MatchLabels: labels}
}

func newPodTemplateSpec(
	secret string, labels map[string]string,
	containers []core.Container) core.PodTemplateSpec {
	return core.PodTemplateSpec{
		ObjectMeta: newObjectMeta("", "", labels),
		Spec: core.PodSpec{
			Containers: containers,
			ImagePullSecrets: []core.LocalObjectReference{
				{
					Name: secret,
				},
			},
		},
	}
}

func newDeployment(
	name, namespace string,
	labels map[string]string,
	containers []core.Container, secret string) *appsv1.Deployment {

		return &appsv1.Deployment{
		ObjectMeta: newObjectMeta(name, namespace, labels),
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: newSelector(labels),
			Template: newPodTemplateSpec(secret, labels, containers),
		},
	}
}

func newReplicationController(
	name, namespace string,
	labels map[string]string,
	containers []core.Container, secret string) *core.ReplicationController {

	podTemplateSpec := newPodTemplateSpec(secret, labels, containers)
	return &core.ReplicationController{
		ObjectMeta: newObjectMeta(name, namespace, labels),
		Spec: core.ReplicationControllerSpec{
			Replicas: int32Ptr(1),
			Selector: labels,
			Template: &podTemplateSpec,
		},
	}
}

func newReplicaSet(
	name, namespace string,
	labels map[string]string,
	containers []core.Container, secret string) *extensions.ReplicaSet {

	return &extensions.ReplicaSet{
		ObjectMeta: newObjectMeta(name, namespace, labels),
		Spec: extensions.ReplicaSetSpec{
			Replicas: int32Ptr(1),
			Selector: newSelector(labels),
			Template: newPodTemplateSpec(secret, labels, containers),
		},
	}
}

func newDaemonSet(
	name, namespace string,
	labels map[string]string,
	containers []core.Container, secret string) *extensions.DaemonSet {

	return &extensions.DaemonSet{
		ObjectMeta: newObjectMeta(name, namespace, labels),
		Spec: extensions.DaemonSetSpec{
			Selector: newSelector(labels),
			Template: newPodTemplateSpec(secret, labels, containers),
		},
	}
}

func newJob(
	name, namespace string,
	labels map[string]string,
	containers []core.Container, secret string) *batchv1.Job {

	podTemplateSpec := newPodTemplateSpec(secret, labels, containers)
	podTemplateSpec.Spec.RestartPolicy = "OnFailure"
	return &batchv1.Job{
		ObjectMeta: newObjectMeta(name, namespace, labels),
		Spec: batchv1.JobSpec{
			Selector: newSelector(labels),
			Template: podTemplateSpec,
		},
	}
}

func newCronJob(
	name, namespace string,
	labels map[string]string,
	containers []core.Container, secret string) *batchv1beta1.CronJob {

	podTemplateSpec := newPodTemplateSpec(secret, labels, containers)
	podTemplateSpec.Spec.RestartPolicy = "OnFailure"
	return &batchv1beta1.CronJob{
		ObjectMeta: newObjectMeta(name, namespace, labels),
		Spec: batchv1beta1.CronJobSpec{
			Schedule: "*/1 * * * *",
			JobTemplate: batchv1beta1.JobTemplateSpec{
				ObjectMeta: newObjectMeta("", "", labels),
				Spec: batchv1.JobSpec{
					Selector: newSelector(labels),
					Template: podTemplateSpec,
				},
			},
		},
	}
}

func newStatefulSet(
	name, namespace string,
	labels map[string]string,
	containers []core.Container, secret string) *appsv1.StatefulSet {

	return &appsv1.StatefulSet{
		ObjectMeta: newObjectMeta(name, namespace, labels),
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    int32Ptr(1),
			Selector: newSelector(labels),
			Template: newPodTemplateSpec(secret, labels, containers),
		},
	}
}

func (f *Invocation) NewService(
	name, namespace string,
	labels map[string]string) *core.Service {

	return &core.Service{
		ObjectMeta: newObjectMeta(name, namespace, labels),
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{
					Protocol: core.ProtocolTCP,
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
				},
			},
			Type:     core.ServiceTypeNodePort,
			Selector: labels,
		},
	}
}

func (f *Invocation) NewSecret(name, namespace, data string, labels map[string]string) *core.Secret {
	return &core.Secret{
		ObjectMeta: newObjectMeta(name, namespace, labels),
		StringData: map[string]string{
			".dockerconfigjson": data,
		},
		Type: "kubernetes.io/dockerconfigjson",
	}
}

func (f *Invocation) NewWorkload(
	name, namespace string,
	labels map[string]string,
	containers []core.Container, secret string, workloadCode int) runtime.Object {

	switch workloadCode {
	case Deployment:
		return newDeployment(name, namespace, labels, containers, secret)
	case ReplicationController:
		return newReplicationController(name, namespace, labels, containers, secret)
	case ReplicaSet:
		return newReplicaSet(name, namespace, labels, containers, secret)
	case DaemonSet:
		return newDaemonSet(name, namespace, labels, containers, secret)
	case Job:
		return newJob(name, namespace, labels, containers, secret)
	case CronJob:
		return newCronJob(name, namespace, labels, containers, secret)
	case StatefulSet:
		return newStatefulSet(name, namespace, labels, containers, secret)
	default:
		return nil
	}
}

func create(root *Framework, obj runtime.Object) error {
	var err error
	switch t := obj.(type) {
	case *core.Pod:
		_, err = root.KubeClient.CoreV1().Pods(t.Namespace).Create(t)
		// ReplicationController
	case *core.ReplicationController:
		_, err = root.KubeClient.CoreV1().ReplicationControllers(t.Namespace).Create(t)
		// Deployment
	case *extensions.Deployment:
		_, err = root.KubeClient.ExtensionsV1beta1().Deployments(t.Namespace).Create(t)
	case *appsv1beta1.Deployment:
		_, err = root.KubeClient.AppsV1beta1().Deployments(t.Namespace).Create(t)
	case *appsv1beta2.Deployment:
		_, err = root.KubeClient.AppsV1beta2().Deployments(t.Namespace).Create(t)
	case *appsv1.Deployment:
		_, err = root.KubeClient.AppsV1().Deployments(t.Namespace).Create(t)
		// DaemonSet
	case *extensions.DaemonSet:
		_, err = root.KubeClient.ExtensionsV1beta1().DaemonSets(t.Namespace).Create(t)
	case *appsv1beta2.DaemonSet:
		_, err = root.KubeClient.AppsV1beta2().DaemonSets(t.Namespace).Create(t)
	case *appsv1.DaemonSet:
		_, err = root.KubeClient.AppsV1().DaemonSets(t.Namespace).Create(t)
		// ReplicaSet
	case *extensions.ReplicaSet:
		_, err = root.KubeClient.ExtensionsV1beta1().ReplicaSets(t.Namespace).Create(t)
	case *appsv1beta2.ReplicaSet:
		_, err = root.KubeClient.AppsV1beta2().ReplicaSets(t.Namespace).Create(t)
	case *appsv1.ReplicaSet:
		_, err = root.KubeClient.AppsV1().ReplicaSets(t.Namespace).Create(t)
		// StatefulSet
	case *appsv1beta1.StatefulSet:
		_, err = root.KubeClient.AppsV1beta1().StatefulSets(t.Namespace).Create(t)
	case *appsv1beta2.StatefulSet:
		_, err = root.KubeClient.AppsV1beta2().StatefulSets(t.Namespace).Create(t)
	case *appsv1.StatefulSet:
		_, err = root.KubeClient.AppsV1().StatefulSets(t.Namespace).Create(t)
		// Job
	case *batchv1.Job:
		_, err = root.KubeClient.BatchV1().Jobs(t.Namespace).Create(t)
		// CronJob
	case *batchv1beta1.CronJob:
		_, err = root.KubeClient.BatchV1beta1().CronJobs(t.Namespace).Create(t)
	default:
		err = fmt.Errorf("the object is not a pod or does not have a pod template")
	}
	
	return err
}

func update(root *Framework, obj runtime.Object) error {
	var err error
	switch t := obj.(type) {
	case *core.Pod:
		_, err = root.KubeClient.CoreV1().Pods(t.Namespace).Update(t)
		// ReplicationController
	case *core.ReplicationController:
		_, err = root.KubeClient.CoreV1().ReplicationControllers(t.Namespace).Update(t)
		// Deployment
	case *extensions.Deployment:
		_, err = root.KubeClient.ExtensionsV1beta1().Deployments(t.Namespace).Update(t)
	case *appsv1beta1.Deployment:
		_, err = root.KubeClient.AppsV1beta1().Deployments(t.Namespace).Update(t)
	case *appsv1beta2.Deployment:
		_, err = root.KubeClient.AppsV1beta2().Deployments(t.Namespace).Update(t)
	case *appsv1.Deployment:
		_, err = root.KubeClient.AppsV1().Deployments(t.Namespace).Update(t)
		// DaemonSet
	case *extensions.DaemonSet:
		_, err = root.KubeClient.ExtensionsV1beta1().DaemonSets(t.Namespace).Update(t)
	case *appsv1beta2.DaemonSet:
		_, err = root.KubeClient.AppsV1beta2().DaemonSets(t.Namespace).Update(t)
	case *appsv1.DaemonSet:
		_, err = root.KubeClient.AppsV1().DaemonSets(t.Namespace).Update(t)
		// ReplicaSet
	case *extensions.ReplicaSet:
		_, err = root.KubeClient.ExtensionsV1beta1().ReplicaSets(t.Namespace).Update(t)
	case *appsv1beta2.ReplicaSet:
		_, err = root.KubeClient.AppsV1beta2().ReplicaSets(t.Namespace).Update(t)
	case *appsv1.ReplicaSet:
		_, err = root.KubeClient.AppsV1().ReplicaSets(t.Namespace).Update(t)
		// StatefulSet
	case *appsv1beta1.StatefulSet:
		_, err = root.KubeClient.AppsV1beta1().StatefulSets(t.Namespace).Update(t)
	case *appsv1beta2.StatefulSet:
		_, err = root.KubeClient.AppsV1beta2().StatefulSets(t.Namespace).Update(t)
	case *appsv1.StatefulSet:
		_, err = root.KubeClient.AppsV1().StatefulSets(t.Namespace).Update(t)
		// Job
	case *batchv1.Job:
		_, err = root.KubeClient.BatchV1().Jobs(t.Namespace).Update(t)
		// CronJob
	case *batchv1beta1.CronJob:
		_, err = root.KubeClient.BatchV1beta1().CronJobs(t.Namespace).Update(t)
	default:
		err = fmt.Errorf("the object is not a pod or does not have a pod template")
	}

	return err
}

func (f *Invocation) EventuallyCreateWithVulnerableImage(root *Framework, obj runtime.Object) GomegaAsyncAssertion {
	return Eventually(
		func() bool {
			err := create(root, obj)
			Expect(err).To(HaveOccurred())

			return strings.Contains(err.Error(), "contains vulnerabilities")
		},
		time.Minute,
		time.Millisecond * 5,
	)
}

func (f *Invocation) EventuallyUpdateWithVulnerableImage(root *Framework, obj runtime.Object) GomegaAsyncAssertion {
	return Eventually(
		func() bool {
			err := update(root, obj)
			Expect(err).To(HaveOccurred())

			return strings.Contains(err.Error(), "contains vulnerabilities")
		},
		time.Minute,
		time.Millisecond * 5,
	)
}

func (f *Invocation) EventuallyCreateWithNonVulnerableImage(root *Framework, obj runtime.Object) GomegaAsyncAssertion {
	return Eventually(
		func() error {
			return create(root, obj)
		},
		time.Minute,
		time.Millisecond * 5,
	)
}

func (f *Invocation) deleteAllDeployments() {
	deployments, err := f.KubeClient.AppsV1beta1().Deployments(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, deploy := range deployments.Items {
		err := f.KubeClient.AppsV1beta1().Deployments(deploy.Namespace).Delete(deploy.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) deleteAllReplicationControllers() {
	replicationcontrollers, err := f.KubeClient.CoreV1().ReplicationControllers(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, rc := range replicationcontrollers.Items {
		err := f.KubeClient.CoreV1().ReplicationControllers(rc.Namespace).Delete(rc.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) deleteAllReplicaSets() {
	replicasets, err := f.KubeClient.ExtensionsV1beta1().ReplicaSets(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, rs := range replicasets.Items {
		err := f.KubeClient.ExtensionsV1beta1().ReplicaSets(rs.Namespace).Delete(rs.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) deleteAllDaemonSet() {
	daemonsets, err := f.KubeClient.ExtensionsV1beta1().DaemonSets(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, ds := range daemonsets.Items {
		err := f.KubeClient.ExtensionsV1beta1().DaemonSets(ds.Namespace).Delete(ds.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) deleteAllJobs() {
	jobs, err := f.KubeClient.BatchV1().Jobs(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, job := range jobs.Items {
		err := f.KubeClient.BatchV1().Jobs(job.Namespace).Delete(job.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) deleteAllCronJobs() {
	cronJobs, err := f.KubeClient.BatchV1beta1().CronJobs(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, cronJob := range cronJobs.Items {
		err := f.KubeClient.BatchV1beta1().CronJobs(cronJob.Namespace).Delete(cronJob.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) deleteAllStatefulSets() {
	statefulsets, err := f.KubeClient.AppsV1beta1().StatefulSets(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, sts := range statefulsets.Items {
		err := f.KubeClient.AppsV1beta1().StatefulSets(sts.Namespace).Delete(sts.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) DeleteWorkloads(workloadCode int) {
	switch workloadCode {
	case Deployment:
		f.deleteAllDeployments()
	case ReplicationController:
		f.deleteAllReplicationControllers()
	case ReplicaSet:
		f.deleteAllReplicaSets()
	case DaemonSet:
		f.deleteAllDaemonSet()
	case Job:
		f.deleteAllJobs()
	case CronJob:
		f.deleteAllCronJobs()
	case StatefulSet:
		f.deleteAllStatefulSets()
	}
}

func (f *Invocation) DeleteAllServices() {
	services, err := f.KubeClient.CoreV1().Services(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, svc := range services.Items {
		err := f.KubeClient.CoreV1().Services(svc.Namespace).Delete(svc.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) DeleteAllSecrets() {
	secrets, err := f.KubeClient.CoreV1().Secrets(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, value := range secrets.Items {
		err := f.KubeClient.CoreV1().Secrets(value.Namespace).Delete(value.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}
