package iiii

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"log"
)

func createMongoService(clientset *kubernetes.Clientset) (*v1.Service, error) {
	service := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "mongo",
			Labels: map[string]string{"name": "mongo"},
		},
		Spec: v1.ServiceSpec{
			Type: "NodePort",
			Ports: []v1.ServicePort{
				{
					Port:       27017,
					TargetPort: intstr.IntOrString{Type: 0, IntVal: 27017},
					NodePort:   32017,
				},
			},
			Selector: map[string]string{"role": "mongo"},
		},
	}
	return clientset.CoreV1().Services(defaultNamespace).Create(&service)
}

func createMongoStatefulSet(clientset *kubernetes.Clientset) (*appsv1.StatefulSet, error) {
	replicas := int32(1)
	terminationGracePeriodSeconds := int64(10)
	statefulSet := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: "mongo"},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "mongo",
			Replicas:    &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"role": "mongo"},
				},
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					Volumes: []v1.Volume{
						{
							Name: "mongo-volume",
							VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{Path: "/home/docker/mongo"},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:         "mongo",
							Image:        "library/mongo:latest",
							Command:      []string{"mongod", "--replSet", "rs0", "--bind_ip", "0.0.0.0"},
							Ports:        []v1.ContainerPort{{ContainerPort: 27017}},
							VolumeMounts: []v1.VolumeMount{{Name: "mongo-volume", MountPath: "/data/db"}},
						},
					},
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"role": "mongo"},
			},
		},
	}
	return clientset.AppsV1().StatefulSets(defaultNamespace).Create(&statefulSet)
}

// StartMongo start mongo in kubernetes
func StartMongo() error {
	log.Println("Starting mongodb in local kubernetes...")
	clientset, clientErr := getKubeClientset()
	if clientErr != nil {
		return clientErr
	}
	log.Println("Creating mongodb service...")
	_, serviceErr := createMongoService(clientset)
	if serviceErr != nil {
		return serviceErr
	}
	log.Println("Creating mongodb statefulset...")
	_, statefulSetErr := createMongoStatefulSet(clientset)
	if statefulSetErr != nil {
		return statefulSetErr
	}
	log.Println("MongoDB is launching now!")
	return nil
}

// Rollback rollback mongo in kubernetes
func RollBack() error {
	log.Println("Rolling back mongodb in local kubernetes...")
	clientset, clientErr := getKubeClientset()
	if clientErr != nil {
		return clientErr
	}
	log.Println("Deleting mongo statefulset...")
	statefulErr := clientset.AppsV1().StatefulSets(defaultNamespace).Delete("mongo", &metav1.DeleteOptions{})
	if statefulErr != nil {
		log.Printf("Error while deleing mongo statefulset! %s\n", statefulErr.Error())
	}
	log.Println("Deleting mongo service...")
	serviceErr := clientset.CoreV1().Services(defaultNamespace).Delete("mongo", &metav1.DeleteOptions{})
	if serviceErr != nil {
		log.Printf("Error while deleing mongo service! %s\n", serviceErr.Error())
	}
	log.Println("Rolled back mongodb in local kubernetes successfully!")
	return nil
}
