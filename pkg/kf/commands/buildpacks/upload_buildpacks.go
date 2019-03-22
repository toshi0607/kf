package buildpacks

import (
	"os"
	"path/filepath"

	"github.com/GoogleCloudPlatform/kf/pkg/kf/buildpacks"
	"github.com/GoogleCloudPlatform/kf/pkg/kf/commands/config"
	"github.com/GoogleCloudPlatform/kf/pkg/kf/internal/kf"
	"github.com/spf13/cobra"
)

// BuilderCreator creates a new buildback builder.
type BuilderCreator interface {
	// Create creates and publishes a builder image.
	Create(dir, containerRegistry string) (string, error)
}

// BuildTemplateUploader uploads a build template
type BuildTemplateUploader interface {
	// UploadBuildTemplate uploads a buildpack build template with the name
	// "buildpack".
	UploadBuildTemplate(imageName string, opts ...buildpacks.UploadBuildTemplateOption) error
}

// NewUploadBuildpacks creates a UploadBuildpacks command.
func NewUploadBuildpacks(p *config.KfParams, c BuilderCreator, u BuildTemplateUploader) *cobra.Command {
	var (
		containerRegistry string
		path              string
	)
	var uploadBuildpacksCmd = &cobra.Command{
		Use:   "upload-buildpacks",
		Short: "Create and upload a new buildpacks builder. This is used to set the available buildpacks that are used while pushing an app.",
		Args:  cobra.ExactArgs(0),
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if path != "" {
				var err error
				path, err = filepath.Abs(path)
				if err != nil {
					return err
				}
			} else {
				var err error
				path, err = os.Getwd()
				if err != nil {
					return err
				}
			}

			image, err := c.Create(path, containerRegistry)
			if err != nil {
				cmd.SilenceUsage = !kf.ConfigError(err)
				return err
			}

			if err := u.UploadBuildTemplate(image, buildpacks.WithUploadBuildTemplateNamespace(p.Namespace)); err != nil {
				cmd.SilenceUsage = !kf.ConfigError(err)
				return err
			}

			return nil
		},
	}

	uploadBuildpacksCmd.Flags().StringVar(
		&path,
		"path",
		"",
		"The path the source code lives. Defaults to current directory.",
	)

	uploadBuildpacksCmd.Flags().StringVar(
		&containerRegistry,
		"container-registry",
		"",
		"The container registry to push the resulting container.",
	)

	return uploadBuildpacksCmd
}
