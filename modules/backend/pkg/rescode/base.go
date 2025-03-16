package rescode

import (
	"net/http"
)

var (
	ValidationFailed = New(1000, http.StatusUnprocessableEntity, "validation_failed")
	Failed           = New(1001, http.StatusInternalServerError, "failed")
	NotFound         = New(1002, http.StatusNotFound, "not_found")
	PermissionDenied = New(1003, http.StatusForbidden, "permission_denied")

	SlugRequired = New(1100, http.StatusBadRequest, "slug_required")
	SlugInvalid  = New(1101, http.StatusBadRequest, "slug_invalid")
	IDRequired   = New(1102, http.StatusBadRequest, "id_required")
	IDInvalid    = New(1103, http.StatusBadRequest, "id_invalid")

	SurveyNotFound = New(1200, http.StatusNotFound, "survey_not_found")
)
