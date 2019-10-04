package customerimagebuild

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/acquia/fn-drupal-operator/pkg/customercontainer"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	fnworkflowsv1alpha1 "github.com/acquia/fn-workflows-operator/pkg/apis/fnworkflows/v1alpha1"
)

const (
	buildNamespace = "build"
)

var log = logf.Log.WithName("controller_customerimagebuild")

// Add creates a new CustomerImageBuild Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCustomerImageBuild{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("customerimagebuild-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CustomerImageBuild
	err = c.Watch(&source.Kind{Type: &fnworkflowsv1alpha1.CustomerImageBuild{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resources and requeue the owner CustomerImageBuild
	// err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &fnworkflowsv1alpha1.CustomerImageBuild{},
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

// blank assignment to verify that ReconcileCustomerImageBuild implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCustomerImageBuild{}

// ReconcileCustomerImageBuild reconciles a CustomerImageBuild object
type ReconcileCustomerImageBuild struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CustomerImageBuild object and makes changes based on the state read
// and what is in the CustomerImageBuild.Spec
//
// Note: The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCustomerImageBuild) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CustomerImageBuild")

	// Fetch the CustomerImageBuild instance
	cib := &fnworkflowsv1alpha1.CustomerImageBuild{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cib)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request was queued.
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Create a new Job object
	job := newBuildJob(cib)

	// Set CustomerImageBuild instance as the owner and controller
	if err := controllerutil.SetControllerReference(cib, job, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Create the Job resource (ignoring error if it already exists)
	err = r.client.Create(context.TODO(), job)
	if err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func newBuildJob(cib *fnworkflowsv1alpha1.CustomerImageBuild) *batchv1.Job {
	backoffLimit := cib.Spec.Retries
	emptyDir := corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}
	gitSSHSecret := corev1.SecretVolumeSource{SecretName: "git-ssh"}
	kanikoConfig := corev1.ConfigMapVolumeSource{
		LocalObjectReference: corev1.LocalObjectReference{Name: "kaniko-config"},
	}

	// Hash build parameters to generate a "unique" but consistent name
	sha := sha1.New()
	sha.Write([]byte(cib.Spec.RepoURL))
	sha.Write([]byte(cib.Spec.GitRef))
	sha.Write([]byte(cib.Spec.ImageTag))
	buildID := fmt.Sprintf("%x", sha.Sum(nil))[:16]

	// Create the Job object
	gitURL := gitURLFromGitRepo(cib.Spec.RepoURL)
	imageURL := fmt.Sprintf("%v:%v", cib.Spec.ImageRepo, cib.Spec.ImageTag)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "build-customer-image-" + buildID,
			Namespace: buildNamespace, // TODO: build namespacing? build in environment NS?
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					InitContainers: []corev1.Container{
						imageBuildInitContainer(gitURL, cib.Spec.GitRef),
					},
					Containers: []corev1.Container{
						imageBuildContainer(imageURL),
					},
					Volumes: []corev1.Volume{
						{Name: "workspace", VolumeSource: emptyDir},
						{Name: "tmp", VolumeSource: emptyDir},
						{Name: "git-ssh", VolumeSource: corev1.VolumeSource{Secret: &gitSSHSecret}},
						{Name: "kaniko-config", VolumeSource: corev1.VolumeSource{ConfigMap: &kanikoConfig}},
					},
				},
			},
		},
	}
}

func imageBuildInitContainer(gitURL, ref string) corev1.Container {
	wwwUser := int64(82) // www-data on Alpine

	var refEnvVar corev1.EnvVar
	parts := strings.Split(ref, "/")
	switch parts[1] {
	case "heads":
		refEnvVar = corev1.EnvVar{Name: "GIT_BRANCH", Value: parts[2]}
	case "tags":
		refEnvVar = corev1.EnvVar{Name: "GIT_TAG", Value: parts[2]}
	}

	return corev1.Container{
		Name:    "init",
		Image:   customercontainer.ECRRepoRoot + "build-helper:v0.1.0",
		Command: []string{"sh", "/workspace-files/prep-for-build.sh"},
		Env: []corev1.EnvVar{
			{Name: "ECR_REPO_NAME", Value: ecrRepoNameFromGitURL(gitURL)},
			{Name: "GIT_REPO", Value: gitURL},
			refEnvVar,
			{Name: "GIT_SSH_KEY_FILE", Value: "/etc/git-secret/id_rsa"},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: "workspace", MountPath: "/workspace"},
			{Name: "git-ssh", MountPath: "/etc/git-secret", ReadOnly: true},
		},
		SecurityContext: &corev1.SecurityContext{
			RunAsUser: &wwwUser,
		},
	}
}

func imageBuildContainer(imageURL string) corev1.Container {
	return corev1.Container{
		Name:  "build",
		Image: "gcr.io/kaniko-project/executor:v0.10.0",
		Args: []string{
			"--context=/workspace",
			"--dockerfile=/workspace/build-drupal.Dockerfile",
			"--destination=" + imageURL,
			"--build-arg=PHP_VERSION=7.3", // FIXME
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: "workspace", MountPath: "/workspace"},
			{Name: "tmp", MountPath: "/shared/tmp"},
			{Name: "kaniko-config", MountPath: "/kaniko/.docker/", ReadOnly: true},
		},
	}
}

func gitURLFromGitRepo(gitRepo string) string {
	gitURL := fmt.Sprintf("git@%v.git", gitRepo)
	return strings.Replace(gitURL, "/", ":", 1)
}

func gitRepoFromGitURL(gitURL string) string {
	atIndex := strings.IndexRune(gitURL, '@')
	gitSuffixIndex := strings.LastIndex(gitURL, ".git")

	return strings.Replace(gitURL[atIndex+1 : gitSuffixIndex], ":", "/", 1)
}

func ecrRepoNameFromGitURL(gitURL string) string {
	return customercontainer.CustomerECRRepoPrefix + gitRepoFromGitURL(gitURL)
}
