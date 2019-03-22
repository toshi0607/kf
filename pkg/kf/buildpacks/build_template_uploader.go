package buildpacks

import (
	"errors"

	build "github.com/knative/build/pkg/apis/build/v1alpha1"
	cbuild "github.com/knative/build/pkg/client/clientset/versioned/typed/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//go:generate go run ../internal/tools/option-builder/option-builder.go options.yml

// BuildFactory returns a client for build.
type BuildFactory func() (cbuild.BuildV1alpha1Interface, error)

// BuildTemplateUploader uploads a new buildpack build template. It should be
// created via NewBuildTemplateUploader.
type BuildTemplateUploader struct {
	f BuildFactory
}

// NewBuildTemplateUploader creates a new BuildTemplateUploader.
func NewBuildTemplateUploader(f BuildFactory) *BuildTemplateUploader {
	return &BuildTemplateUploader{
		f: f,
	}
}

// UploadBuildTemplate uploads a buildpack build template with the name
// "buildpack".
func (u *BuildTemplateUploader) UploadBuildTemplate(imageName string, opts ...UploadBuildTemplateOption) error {
	if imageName == "" {
		return errors.New("image name must not be empty")
	}

	cfg := UploadBuildTemplateOptionDefaults().Extend(opts).toConfig()
	c, err := u.f()
	if err != nil {
		return err
	}

	// TODO: It would be nice if we generated this instead.
	if _, err := u.deployer(c, cfg.Namespace)(&build.BuildTemplate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "build.knative.dev/v1alpha1",
			Kind:       "BuildTemplate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "buildpack",
		},
		Spec: build.BuildTemplateSpec{
			Parameters: []build.ParameterSpec{
				{
					Name:        "IMAGE",
					Description: `The image you wish to create. For example, "repo/example", or "example.com/repo/image"`,
				},
				{
					Name:        "RUN_IMAGE",
					Description: `The run image buildpacks will use as the base for IMAGE.`,
					Default:     u.strToPtr("packs/run:v3alpha2"),
				},
				{
					Name:        "BUILDER_IMAGE",
					Description: `The builder image (must include v3 lifecycle and compatible buildpacks).`,
					Default:     u.strToPtr(imageName),
				},
				{
					Name:        "USE_CRED_HELPERS",
					Description: `Use Docker credential helpers for Google's GCR, Amazon's ECR, or Microsoft's ACR.`,
					Default:     u.strToPtr("true"),
				},
				{
					Name:        "CACHE",
					Description: `The name of the persistent app cache volume`,
					Default:     u.strToPtr("empty-dir"),
				},
				{
					Name:        "USER_ID",
					Description: `The user ID of the builder image user`,
					Default:     u.strToPtr("1000"),
				},
				{
					Name:        "GROUP_ID",
					Description: `The group ID of the builder image user`,
					Default:     u.strToPtr("1000"),
				},
			},
			Steps: []corev1.Container{
				{
					Name:    "prepare",
					Image:   "alpine",
					Command: []string{"/bin/sh"},
					Args: []string{
						"-c",
						`chown -R "${USER_ID}:${GROUP_ID}" "/builder/home" &&
						 chown -R "${USER_ID}:${GROUP_ID}" /layers &&
						 chown -R "${USER_ID}:${GROUP_ID}" /workspace`,
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      "${CACHE}",
						MountPath: "/layers",
					}},
					ImagePullPolicy: "Always",
				},
				{
					Name:    "detect",
					Image:   "${BUILDER_IMAGE}",
					Command: []string{"/lifecycle/detector"},
					Args: []string{
						"-app=/workspace",
						"-group=/layers/group.toml",
						"-plan=/layers/plan.toml",
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      "${CACHE}",
						MountPath: "/layers",
					}},
					ImagePullPolicy: "Always",
				},
				{
					Name:    "analyze",
					Image:   "${BUILDER_IMAGE}",
					Command: []string{"/lifecycle/analyzer"},
					Args: []string{
						"-layers=/layers",
						"-helpers=${USE_CRED_HELPERS}",
						"-group=/layers/group.toml",
						"${IMAGE}",
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      "${CACHE}",
						MountPath: "/layers",
					}},
					ImagePullPolicy: "Always",
				},
				{
					Name:    "build",
					Image:   "${BUILDER_IMAGE}",
					Command: []string{"/lifecycle/builder"},
					Args: []string{
						"-layers=/layers",
						"-app=/workspace",
						"-group=/layers/group.toml",
						"-plan=/layers/plan.toml",
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      "${CACHE}",
						MountPath: "/layers",
					}},
					ImagePullPolicy: "Always",
				},
				{
					Name:    "export",
					Image:   "${BUILDER_IMAGE}",
					Command: []string{"/lifecycle/exporter"},
					Args: []string{
						"-layers=/layers",
						"-helpers=${USE_CRED_HELPERS}",
						"-app=/workspace",
						"-image=${RUN_IMAGE}",
						"-group=/layers/group.toml",
						"${IMAGE}",
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      "${CACHE}",
						MountPath: "/layers",
					}},
					ImagePullPolicy: "Always",
				},
			},
			Volumes: []corev1.Volume{{
				Name: "empty-dir",
			}},
		},
	}); err != nil {
		return err
	}

	return nil
}

type deployer func(*build.BuildTemplate) (*build.BuildTemplate, error)

func (u *BuildTemplateUploader) deployer(c cbuild.BuildV1alpha1Interface, namespace string) deployer {
	builds, err := c.BuildTemplates(namespace).List(metav1.ListOptions{
		FieldSelector: "metadata.name=buildpack",
	})

	if err != nil {
		// Simplify workflow and just return a deployer that will fail with the
		// given error.
		return func(t *build.BuildTemplate) (*build.BuildTemplate, error) {
			return nil, err
		}
	}

	if len(builds.Items) == 0 {
		return func(t *build.BuildTemplate) (*build.BuildTemplate, error) {
			return c.BuildTemplates(namespace).Create(t)
		}
	}

	return func(t *build.BuildTemplate) (*build.BuildTemplate, error) {
		t.ResourceVersion = builds.Items[0].ResourceVersion
		return c.BuildTemplates(namespace).Update(t)
	}
}

func (u *BuildTemplateUploader) strToPtr(s string) *string {
	return &s
}
