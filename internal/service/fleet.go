package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	api "github.com/flightctl/flightctl/api/v1alpha1"
	"github.com/flightctl/flightctl/internal/api/server"
	"github.com/flightctl/flightctl/internal/auth"
	"github.com/flightctl/flightctl/internal/flterrors"
	"github.com/flightctl/flightctl/internal/service/common"
	"github.com/flightctl/flightctl/internal/store"
	"github.com/flightctl/flightctl/internal/store/selector"
	"github.com/flightctl/flightctl/internal/util"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
)

func FleetFromReader(r io.Reader) (*api.Fleet, error) {
	var fleet api.Fleet
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&fleet)
	return &fleet, err
}

// (POST /api/v1/fleets)
func (h *ServiceHandler) CreateFleet(ctx context.Context, request server.CreateFleetRequestObject) (server.CreateFleetResponseObject, error) {
	allowed, err := auth.GetAuthZ().CheckPermission(ctx, "fleets", "create")
	if err != nil {
		h.log.WithError(err).Error("failed to check authorization permission")
		return server.CreateFleet503JSONResponse{Message: AuthorizationServerUnavailable}, nil
	}
	if !allowed {
		return server.CreateFleet403JSONResponse{Message: Forbidden}, nil
	}
	orgId := store.NullOrgId

	// don't set fields that are managed by the service
	request.Body.Status = nil
	common.NilOutManagedObjectMetaProperties(&request.Body.Metadata)
	if request.Body.Spec.Template.Metadata != nil {
		common.NilOutManagedObjectMetaProperties(request.Body.Spec.Template.Metadata)
	}

	if errs := request.Body.Validate(); len(errs) > 0 {
		return server.CreateFleet400JSONResponse{Message: errors.Join(errs...).Error()}, nil
	}

	result, err := h.store.Fleet().Create(ctx, orgId, request.Body, h.callbackManager.FleetUpdatedCallback)
	switch {
	case err == nil:
		return server.CreateFleet201JSONResponse(*result), nil
	case errors.Is(err, flterrors.ErrResourceIsNil), errors.Is(err, flterrors.ErrIllegalResourceVersionFormat):
		return server.CreateFleet400JSONResponse{Message: err.Error()}, nil
	case errors.Is(err, flterrors.ErrDuplicateName):
		return server.CreateFleet409JSONResponse{Message: err.Error()}, nil
	default:
		return nil, err
	}
}

// (GET /api/v1/fleets)
func (h *ServiceHandler) ListFleets(ctx context.Context, request server.ListFleetsRequestObject) (server.ListFleetsResponseObject, error) {
	allowed, err := auth.GetAuthZ().CheckPermission(ctx, "fleets", "list")
	if err != nil {
		h.log.WithError(err).Error("failed to check authorization permission")
		return server.ListFleets503JSONResponse{Message: AuthorizationServerUnavailable}, nil
	}
	if !allowed {
		return server.ListFleets403JSONResponse{Message: Forbidden}, nil
	}
	orgId := store.NullOrgId

	cont, err := store.ParseContinueString(request.Params.Continue)
	if err != nil {
		return server.ListFleets400JSONResponse{Message: fmt.Sprintf("failed to parse continue parameter: %v", err)}, nil
	}

	var fieldSelector *selector.FieldSelector
	if request.Params.FieldSelector != nil {
		if fieldSelector, err = selector.NewFieldSelector(*request.Params.FieldSelector); err != nil {
			return server.ListFleets400JSONResponse{Message: fmt.Sprintf("failed to parse field selector: %v", err)}, nil
		}
	}

	var labelSelector *selector.LabelSelector
	if request.Params.LabelSelector != nil {
		if labelSelector, err = selector.NewLabelSelector(*request.Params.LabelSelector); err != nil {
			return server.ListFleets400JSONResponse{Message: fmt.Sprintf("failed to parse label selector: %v", err)}, nil
		}
	}

	listParams := store.ListParams{
		Limit:         int(swag.Int32Value(request.Params.Limit)),
		Continue:      cont,
		FieldSelector: fieldSelector,
		LabelSelector: labelSelector,
	}
	if listParams.Limit == 0 {
		listParams.Limit = store.MaxRecordsPerListRequest
	}
	if listParams.Limit > store.MaxRecordsPerListRequest {
		return server.ListFleets400JSONResponse{Message: fmt.Sprintf("limit cannot exceed %d", store.MaxRecordsPerListRequest)}, nil
	}

	result, err := h.store.Fleet().List(ctx, orgId, listParams, store.WithDeviceCount(util.DefaultBoolIfNil(request.Params.AddDevicesCount, false)))
	if err == nil {
		return server.ListFleets200JSONResponse(*result), nil
	}

	var se *selector.SelectorError

	switch {
	case selector.AsSelectorError(err, &se):
		return server.ListFleets400JSONResponse{Message: se.Error()}, nil
	default:
		return nil, err
	}
}

// (DELETE /api/v1/fleets)
func (h *ServiceHandler) DeleteFleets(ctx context.Context, request server.DeleteFleetsRequestObject) (server.DeleteFleetsResponseObject, error) {
	allowed, err := auth.GetAuthZ().CheckPermission(ctx, "fleets", "deletecollection")
	if err != nil {
		h.log.WithError(err).Error("failed to check authorization permission")
		return server.DeleteFleets503JSONResponse{Message: AuthorizationServerUnavailable}, nil
	}
	if !allowed {
		return server.DeleteFleets403JSONResponse{Message: Forbidden}, nil
	}
	orgId := store.NullOrgId

	err = h.store.Fleet().DeleteAll(ctx, orgId, h.callbackManager.AllFleetsDeletedCallback)
	switch err {
	case nil:
		return server.DeleteFleets200JSONResponse{}, nil
	default:
		return nil, err
	}
}

// (GET /api/v1/fleets/{name})
func (h *ServiceHandler) ReadFleet(ctx context.Context, request server.ReadFleetRequestObject) (server.ReadFleetResponseObject, error) {
	allowed, err := auth.GetAuthZ().CheckPermission(ctx, "fleets", "get")
	if err != nil {
		h.log.WithError(err).Error("failed to check authorization permission")
		return server.ReadFleet503JSONResponse{Message: AuthorizationServerUnavailable}, nil
	}
	if !allowed {
		return server.ReadFleet403JSONResponse{Message: Forbidden}, nil
	}
	orgId := store.NullOrgId

	result, err := h.store.Fleet().Get(ctx, orgId, request.Name, store.WithSummary(util.DefaultBoolIfNil(request.Params.AddDevicesSummary, false)))
	switch {
	case err == nil:
		return server.ReadFleet200JSONResponse(*result), nil
	case errors.Is(err, flterrors.ErrResourceNotFound):
		return server.ReadFleet404JSONResponse{}, nil
	default:
		return nil, err
	}
}

// (PUT /api/v1/fleets/{name})
func (h *ServiceHandler) ReplaceFleet(ctx context.Context, request server.ReplaceFleetRequestObject) (server.ReplaceFleetResponseObject, error) {
	allowed, err := auth.GetAuthZ().CheckPermission(ctx, "fleets", "update")
	if err != nil {
		h.log.WithError(err).Error("failed to check authorization permission")
		return server.ReplaceFleet503JSONResponse{Message: AuthorizationServerUnavailable}, nil
	}
	if !allowed {
		return server.ReplaceFleet403JSONResponse{Message: Forbidden}, nil
	}
	orgId := store.NullOrgId

	// don't overwrite fields that are managed by the service
	request.Body.Status = nil
	common.NilOutManagedObjectMetaProperties(&request.Body.Metadata)
	if request.Body.Spec.Template.Metadata != nil {
		common.NilOutManagedObjectMetaProperties(request.Body.Spec.Template.Metadata)
	}

	if errs := request.Body.Validate(); len(errs) > 0 {
		return server.ReplaceFleet400JSONResponse{Message: errors.Join(errs...).Error()}, nil
	}
	if request.Name != *request.Body.Metadata.Name {
		return server.ReplaceFleet400JSONResponse{Message: "resource name specified in metadata does not match name in path"}, nil
	}

	result, created, err := h.store.Fleet().CreateOrUpdate(ctx, orgId, request.Body, nil, true, h.callbackManager.FleetUpdatedCallback)
	switch {
	case err == nil:
		if created {
			return server.ReplaceFleet201JSONResponse(*result), nil
		} else {
			return server.ReplaceFleet200JSONResponse(*result), nil
		}
	case errors.Is(err, flterrors.ErrResourceIsNil):
		return server.ReplaceFleet400JSONResponse{Message: err.Error()}, nil
	case errors.Is(err, flterrors.ErrResourceNameIsNil):
		return server.ReplaceFleet400JSONResponse{Message: err.Error()}, nil
	case errors.Is(err, flterrors.ErrResourceNotFound):
		return server.ReplaceFleet404JSONResponse{}, nil
	case errors.Is(err, flterrors.ErrUpdatingResourceWithOwnerNotAllowed), errors.Is(err, flterrors.ErrNoRowsUpdated), errors.Is(err, flterrors.ErrResourceVersionConflict):
		return server.ReplaceFleet409JSONResponse{Message: err.Error()}, nil
	default:
		return nil, err
	}
}

// (DELETE /api/v1/fleets/{name})
func (h *ServiceHandler) DeleteFleet(ctx context.Context, request server.DeleteFleetRequestObject) (server.DeleteFleetResponseObject, error) {
	allowed, err := auth.GetAuthZ().CheckPermission(ctx, "fleets", "delete")
	if err != nil {
		h.log.WithError(err).Error("failed to check authorization permission")
		return server.DeleteFleet503JSONResponse{Message: AuthorizationServerUnavailable}, nil
	}
	if !allowed {
		return server.DeleteFleet403JSONResponse{Message: Forbidden}, nil
	}
	orgId := store.NullOrgId

	f, err := h.store.Fleet().Get(ctx, orgId, request.Name)
	if err == flterrors.ErrResourceNotFound {
		return server.DeleteFleet404JSONResponse{}, nil
	}
	if f.Metadata.Owner != nil {
		// Can't delete via api
		return server.DeleteFleet403JSONResponse{Message: "unauthorized to delete fleet because it is owned by another resource"}, nil
	}

	err = h.store.Fleet().Delete(ctx, orgId, request.Name, h.callbackManager.FleetUpdatedCallback)
	switch {
	case err == nil:
		return server.DeleteFleet200JSONResponse{}, nil
	case errors.Is(err, flterrors.ErrResourceNotFound):
		return server.DeleteFleet404JSONResponse{}, nil
	default:
		return nil, err
	}
}

// (GET /api/v1/fleets/{name}/status)
func (h *ServiceHandler) ReadFleetStatus(ctx context.Context, request server.ReadFleetStatusRequestObject) (server.ReadFleetStatusResponseObject, error) {
	allowed, err := auth.GetAuthZ().CheckPermission(ctx, "fleets/status", "get")
	if err != nil {
		h.log.WithError(err).Error("failed to check authorization permission")
		return server.ReadFleetStatus503JSONResponse{Message: AuthorizationServerUnavailable}, nil
	}
	if !allowed {
		return server.ReadFleetStatus403JSONResponse{Message: Forbidden}, nil
	}
	orgId := store.NullOrgId

	result, err := h.store.Fleet().Get(ctx, orgId, request.Name)
	switch {
	case err == nil:
		return server.ReadFleetStatus200JSONResponse(*result), nil
	case errors.Is(err, flterrors.ErrResourceNotFound):
		return server.ReadFleetStatus404JSONResponse{}, nil
	default:
		return nil, err
	}
}

// (PUT /api/v1/fleets/{name}/status)
func (h *ServiceHandler) ReplaceFleetStatus(ctx context.Context, request server.ReplaceFleetStatusRequestObject) (server.ReplaceFleetStatusResponseObject, error) {
	allowed, err := auth.GetAuthZ().CheckPermission(ctx, "fleets/status", "update")
	if err != nil {
		h.log.WithError(err).Error("failed to check authorization permission")
		return server.ReplaceFleetStatus503JSONResponse{Message: AuthorizationServerUnavailable}, nil
	}
	if !allowed {
		return server.ReplaceFleetStatus403JSONResponse{Message: Forbidden}, nil
	}
	orgId := store.NullOrgId

	result, err := h.store.Fleet().UpdateStatus(ctx, orgId, request.Body)
	switch {
	case err == nil:
		return server.ReplaceFleetStatus200JSONResponse(*result), nil
	case errors.Is(err, flterrors.ErrResourceNotFound):
		return server.ReplaceFleetStatus404JSONResponse{}, nil
	default:
		return nil, err
	}
}

// (PATCH /api/v1/fleets/{name})
// Only metadata.labels and spec can be patched. If we try to patch other fields, HTTP 400 Bad Request is returned.
func (h *ServiceHandler) PatchFleet(ctx context.Context, request server.PatchFleetRequestObject) (server.PatchFleetResponseObject, error) {
	allowed, err := auth.GetAuthZ().CheckPermission(ctx, "fleets", "patch")
	if err != nil {
		h.log.WithError(err).Error("failed to check authorization permission")
		return server.PatchFleet503JSONResponse{Message: AuthorizationServerUnavailable}, nil
	}
	if !allowed {
		return server.PatchFleet403JSONResponse{Message: Forbidden}, nil
	}
	orgId := store.NullOrgId

	currentObj, err := h.store.Fleet().Get(ctx, orgId, request.Name)
	if err != nil {
		switch {
		case errors.Is(err, flterrors.ErrResourceIsNil), errors.Is(err, flterrors.ErrResourceNameIsNil):
			return server.PatchFleet400JSONResponse{Message: err.Error()}, nil
		case errors.Is(err, flterrors.ErrResourceNotFound):
			return server.PatchFleet404JSONResponse{}, nil
		default:
			return nil, err
		}
	}

	newObj := &api.Fleet{}
	err = ApplyJSONPatch(ctx, currentObj, newObj, *request.Body, "/api/v1/fleets/"+request.Name)
	if err != nil {
		return server.PatchFleet400JSONResponse{Message: err.Error()}, nil
	}

	if errs := newObj.Validate(); len(errs) > 0 {
		return server.PatchFleet400JSONResponse{Message: errors.Join(errs...).Error()}, nil
	}
	if newObj.Metadata.Name == nil || *currentObj.Metadata.Name != *newObj.Metadata.Name {
		return server.PatchFleet400JSONResponse{Message: "metadata.name is immutable"}, nil
	}
	if currentObj.ApiVersion != newObj.ApiVersion {
		return server.PatchFleet400JSONResponse{Message: "apiVersion is immutable"}, nil
	}
	if currentObj.Kind != newObj.Kind {
		return server.PatchFleet400JSONResponse{Message: "kind is immutable"}, nil
	}
	if !reflect.DeepEqual(currentObj.Status, newObj.Status) {
		return server.PatchFleet400JSONResponse{Message: "status is immutable"}, nil
	}

	common.NilOutManagedObjectMetaProperties(&newObj.Metadata)
	newObj.Metadata.ResourceVersion = nil

	var updateCallback func(uuid.UUID, *api.Fleet, *api.Fleet)

	if h.callbackManager != nil {
		updateCallback = h.callbackManager.FleetUpdatedCallback
	}
	result, err := h.store.Fleet().Update(ctx, orgId, newObj, nil, true, updateCallback)

	switch {
	case err == nil:
		return server.PatchFleet200JSONResponse(*result), nil
	case errors.Is(err, flterrors.ErrResourceIsNil), errors.Is(err, flterrors.ErrResourceNameIsNil):
		return server.PatchFleet400JSONResponse{Message: err.Error()}, nil
	case errors.Is(err, flterrors.ErrResourceNotFound):
		return server.PatchFleet404JSONResponse{}, nil
	case errors.Is(err, flterrors.ErrNoRowsUpdated), errors.Is(err, flterrors.ErrResourceVersionConflict):
		return server.PatchFleet409JSONResponse{}, nil
	default:
		return nil, err
	}
}

// (PATCH /api/v1/fleets/{name}/status)
func (h *ServiceHandler) PatchFleetStatus(ctx context.Context, request server.PatchFleetStatusRequestObject) (server.PatchFleetStatusResponseObject, error) {
	return nil, fmt.Errorf("not yet implemented")
}
