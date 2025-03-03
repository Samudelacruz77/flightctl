package v1alpha1

import (
	"strings"
	"time"

	"github.com/flightctl/flightctl/internal/util"
)

// IsDisconnected() is true if the device never updated status or its last status update is older than disconnectTimeout.
func (d *Device) IsDisconnected(disconnectTimeout time.Duration) bool {
	return d == nil || (d.Status != nil && d.Status.LastSeen.Add(disconnectTimeout).Before(time.Now()))
}

// IsManaged() if the device has its owner field set.
func (d *Device) IsManaged() bool {
	return d != nil && d.Metadata.Owner != nil && len(util.DefaultIfNil(d.Metadata.Owner, "")) > 0
}

// IsManagedBy() is true if the device is managed by the given fleet.
func (d *Device) IsManagedBy(f *Fleet) bool {
	if f == nil || !d.IsManaged() {
		return false
	}
	return util.FromPtr(f.Metadata.Name) == strings.TrimPrefix(util.FromPtr(d.Metadata.Owner), "Fleet/")
}

// IsUpdating() is true if the device's agent reports that it is updating.
func (d *Device) IsUpdating() bool {
	return d != nil && d.Status != nil && IsStatusConditionTrue(d.Status.Conditions, DeviceUpdating)
}

// IsRebooting() is true if the device's agent has the updating condition set with state Rebooting.
func (d *Device) IsRebooting() bool {
	if d == nil || d.Status == nil {
		return false
	}
	updatingCondition := FindStatusCondition(d.Status.Conditions, DeviceUpdating)
	if updatingCondition == nil {
		return false
	}
	return updatingCondition.Status == ConditionStatusTrue && updatingCondition.Reason == string(UpdateStateRebooting)
}

// IsUpdatedToDeviceSpec() is true if the device's current rendered version matches its spec's rendered version.
func (d *Device) IsUpdatedToDeviceSpec() bool {
	if d == nil || d.Metadata.Annotations == nil {
		// devices without a rendered version cannot be out-of-date
		return true
	}
	if d.Status == nil {
		// devices without status cannot be up to date
		return false
	}
	renderedVersionString, ok := (*d.Metadata.Annotations)[DeviceAnnotationRenderedVersion]
	if !ok {
		// devices without a rendered version cannot be out-of-date
		return true
	}
	return d.Status.Config.RenderedVersion == renderedVersionString
}

// IsUpdatedToFleetSpec() is true if the IsUpdatedToDeviceSpec() and
// device spec's current rendered version matches its fleet's rendered version.
func (d *Device) IsUpdatedToFleetSpec(f *Fleet) bool {
	if !d.IsManagedBy(f) {
		// a device cannot be up to date relative to a fleet it is not managed by
		return false
	}
	if d.Metadata.Annotations == nil || f.Metadata.Annotations == nil {
		return false
	}
	fleetTemplateVersion, ok := (*f.Metadata.Annotations)[FleetAnnotationTemplateVersion]
	if !ok {
		return false
	}
	deviceTemplateVersion, ok := (*d.Metadata.Annotations)[DeviceAnnotationRenderedTemplateVersion]
	if !ok {
		return false
	}
	return d.IsUpdatedToDeviceSpec() && deviceTemplateVersion == fleetTemplateVersion
}

// IsDecommissioning() is true if the device has added a DeviceDecommissioning ConditionType to its Conditions.
func (d *Device) IsDecommissioning() bool {
	if d.Status == nil || d.Status.Conditions == nil {
		return false
	}

	decommissioningCondition := FindStatusCondition(d.Status.Conditions, DeviceDecommissioning)
	if decommissioningCondition == nil {
		return false
	}
	return decommissioningCondition.Status == ConditionStatusTrue
}
