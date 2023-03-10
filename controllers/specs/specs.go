package specs

import (
	"math/rand"

	configmaps "example.com/pg/apis/configmaps/v1alpha1"
	deployments "example.com/pg/apis/deployments/v1alpha1"
	secrets "example.com/pg/apis/secrets/v1alpha1"
	services "example.com/pg/apis/services/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const numOfReplicas = 1
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$&*"

// BuildDeployment creates a kubernetes deployment specification
func BuildDeployment(ns string, deploy *deployments.Deployment) *appsv1.Deployment {
	name := "pg-deploy"
	replicas := int32(numOfReplicas)
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: buildMetadata(name, ns, deploy.APIVersion, deploy.Kind, deploy.UID),
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: buildLabels(name),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: buildMetadata(name, ns, deploy.APIVersion, deploy.Kind, deploy.UID),
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "pg",
							Image: "postgres:" + deploy.Spec.PgVersion,
							Ports: []v1.ContainerPort{
								{
									Protocol:      v1.ProtocolTCP,
									ContainerPort: deploy.Spec.ContainerPort,
								},
							},
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceMemory: resource.MustParse(deploy.Spec.MemoryLimit),
									v1.ResourceCPU:    resource.MustParse(deploy.Spec.CpuLimit),
								},
								Requests: v1.ResourceList{
									v1.ResourceMemory: resource.MustParse(deploy.Spec.MemoryRequest),
									v1.ResourceCPU:    resource.MustParse(deploy.Spec.CpuRequest),
								},
							},
							Env: []v1.EnvVar{
								getEnvVarSecretSource("CERT", "cert", "key.crt"),
							},
						},
					},
					RestartPolicy: v1.RestartPolicyAlways,
				},
			},
		},
	}
}

// BuildConfigMap will build a kubernetes config map for postgres
func BuildConfigMap(ns string, cm *configmaps.ConfigMap) *v1.ConfigMap {
	name := "pg-cm"
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "batch/v1/beta1",
		},
		ObjectMeta: buildMetadata(name, ns, cm.APIVersion, cm.Kind, cm.UID),
		Data: map[string]string{
			"pg_hba.conf":     "###",
			"postgresql.conf": "data_directory = /var/lib/postgresql/data/data-directory",
		},
	}
}

// BuildService in kubernetes with pgDeploy port
func BuildService(ns string, service *services.Service) *v1.Service {
	name := "pg-service"
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: buildMetadata(name, ns, service.APIVersion, service.Kind, service.UID),
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeNodePort,
			Ports:    buildServicePorts(service),
			Selector: buildLabels(name),
		},
	}
}

func BuildSecret(ns string, secret *secrets.Secret) *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: buildMetadata("pg-secret", ns, secret.APIVersion, secret.Kind, secret.UID),
		StringData: createSecret(secret),
	}
}

//
// ------------------------------------------------------------------------------------------------------- helpers -----------------------------------------------------------------------------
//

func buildServicePorts(serivce *services.Service) []v1.ServicePort {
	var ports []v1.ServicePort

	for _, s := range serivce.Spec.Ports {
		sp := v1.ServicePort{
			Name:     s.Name,
			Protocol: v1.Protocol(s.Protocol),
			Port:     s.Port,
		}
		ports = append(ports, sp)
	}

	return ports
}

func getEnvVarSecretSource(envName, name, secret string) v1.EnvVar {
	return v1.EnvVar{
		Name: envName,
		ValueFrom: &v1.EnvVarSource{
			SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: name,
				},
				Key: secret,
			},
		},
	}
}

func getEnvVarConfigMapSource(configName, fileName string) v1.EnvVar {
	return v1.EnvVar{
		Name: configName,
		ValueFrom: &v1.EnvVarSource{
			ConfigMapKeyRef: &v1.ConfigMapKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: configName,
				},
				Key: fileName,
			},
		},
	}
}

func buildMetadata(name, namespace, apiVersion, kind string, uid types.UID) metav1.ObjectMeta {
	controlled := true
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    buildLabels(name),
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: apiVersion,
				Kind:       kind,
				Name:       name,
				UID:        uid,
				Controller: &controlled,
			},
		},
	}
}

func buildLabels(name string) map[string]string {
	m := make(map[string]string)
	m["app"] = "pg"
	m["app.kubernetes.io/name"] = name
	m["app.kubernetes.io/component"] = name
	return m
}

func createSecret(secret *secrets.Secret) map[string]string {
	m := make(map[string]string)
	m["postgres_password"] = randStringBytes(12)
	m["postgres_db"] = secret.Spec.Database
	m["postgres_user"] = secret.Spec.User
	return m
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
