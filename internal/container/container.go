package container

import (
	"context"
	"fmt"

	"github.com/flightctl/flightctl/api/v1alpha1"
	"github.com/flightctl/flightctl/internal/util/validation"
)

type BootcHost struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
	Status     Status   `json:"status"`
}

type Metadata struct {
	Name string `json:"name"`
}

type Spec struct {
	Image ImageSpec `json:"image"`
}

type ImageSpec struct {
	Image     string `json:"image"`
	Transport string `json:"transport"`
}

type Status struct {
	Staged   ImageStatus `json:"staged"`
	Booted   ImageStatus `json:"booted"`
	Rollback ImageStatus `json:"rollback"`
	Type     string      `json:"type"`
}

type ImageStatus struct {
	Image        ImageDetails  `json:"image"`
	CachedUpdate *bool         `json:"cachedUpdate"`
	Incompatible bool          `json:"incompatible"`
	Pinned       bool          `json:"pinned"`
	Ostree       OstreeDetails `json:"ostree"`
}

type ImageDetails struct {
	Image       ImageSpec `json:"image"`
	Version     string    `json:"version"`
	Timestamp   string    `json:"timestamp"`
	ImageDigest string    `json:"imageDigest"`
}

type OstreeDetails struct {
	Checksum     string `json:"checksum"`
	DeploySerial int    `json:"deploySerial"`
}

type BootcClient interface {
	Status(ctx context.Context) (*BootcHost, error)
	Switch(ctx context.Context, image string) error
	Apply(ctx context.Context) error
	UsrOverlay(ctx context.Context) error
}

var (
	ErrParsingImage = fmt.Errorf("unable to parse image reference into a valid bootc target")
)

// IsOsImageReconciled returns true if the booted image equals the target for the spec image.
func IsOsImageReconciled(host *BootcHost, desiredSpec *v1alpha1.RenderedDeviceSpec) (bool, error) {
	if desiredSpec.Os == nil {
		return false, nil
	}

	target, err := ImageToBootcTarget(desiredSpec.Os.Image)
	if err != nil {
		return false, err
	}
	// If the booted image equals the desired target, the OS image is reconciled
	return host.GetBootedImage() == target, nil
}

func (b *BootcHost) GetBootedImage() string {
	return b.Status.Booted.Image.Image.Image
}

func (b *BootcHost) GetBootedImageDigest() string {
	return b.Status.Booted.Image.ImageDigest
}

func (b *BootcHost) GetStagedImage() string {
	return b.Status.Staged.Image.Image.Image
}

func (b *BootcHost) GetRollbackImage() string {
	return b.Status.Rollback.Image.Image.Image
}

// Bootc does not accept images with tags AND digests specified - in the case when we
// get both we will use the image digest.
//
// Related underlying issue: https://github.com/containers/image/issues/1736
func ImageToBootcTarget(image string) (string, error) {
	matches := validation.OciImageReferenceRegexp.FindStringSubmatch(image)
	if len(matches) == 0 {
		return image, ErrParsingImage
	}

	// The OciImageReferenceRegexp has 3 capture groups for the base, tag, and digest
	base := matches[1]
	tag := matches[2]
	digest := matches[3]

	if tag != "" && digest != "" {
		return fmt.Sprintf("%s@%s", base, digest), nil
	}

	return image, nil
}
