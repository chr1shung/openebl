package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"github.com/openebl/openebl/pkg/bu_server/middleware"
	"github.com/openebl/openebl/pkg/bu_server/model"
	"github.com/openebl/openebl/pkg/bu_server/trade_document"
	"github.com/sirupsen/logrus"
)

func (a *API) createFileBasedEBL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appID, _ := ctx.Value(middleware.APPLICATION_ID).(string)
	buID, _ := ctx.Value(middleware.BUSINESS_UNIT_ID).(string)

	var req trade_document.IssueFileBasedEBLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Application = appID
	req.Issuer = buID

	ts := time.Now().Unix()
	result, err := a.fileEBLCtrl.Create(ctx, ts, req)
	if err != nil {
		http.Error(w, err.Error(), model.ErrorToHttpStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.Warnf("createFileBasedEBL failed to encode/write response: %v", err)
	}
}

func (a *API) updateFileBasedEBL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appID, _ := ctx.Value(middleware.APPLICATION_ID).(string)
	buID, _ := ctx.Value(middleware.BUSINESS_UNIT_ID).(string)
	docID := mux.Vars(r)["id"]

	var req trade_document.UpdateFileBasedEBLDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Application = appID
	req.Issuer = buID
	req.ID = docID

	ts := time.Now().Unix()
	result, err := a.fileEBLCtrl.UpdateDraft(ctx, ts, req)
	if err != nil {
		http.Error(w, err.Error(), model.ErrorToHttpStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.Warnf("updateFileBasedEBL failed to encode/write response: %v", err)
	}
}

func (a *API) listFileBasedEBL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appID, _ := ctx.Value(middleware.APPLICATION_ID).(string)
	buID, _ := ctx.Value(middleware.BUSINESS_UNIT_ID).(string)

	var req trade_document.ListFileBasedEBLRequest
	req.Lister = buID
	req.Application = appID
	req.Status = r.URL.Query().Get("status")
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		offset, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil || offset < 0 {
			http.Error(w, "offset is invalid", http.StatusBadRequest)
			return
		}
		req.Offset = int(offset)
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || limit < 1 {
			http.Error(w, "limit is invalid", http.StatusBadRequest)
			return
		}
		req.Limit = int(limit)
	}

	result, err := a.fileEBLCtrl.List(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), model.ErrorToHttpStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.Warnf("listFileBasedEBL failed to encode/write response: %v", err)
	}
}

func (a *API) transferEBL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appID, _ := ctx.Value(middleware.APPLICATION_ID).(string)
	buID, _ := ctx.Value(middleware.BUSINESS_UNIT_ID).(string)

	var req trade_document.TransferEBLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Application = appID
	req.TransferBy = buID
	req.ID = mux.Vars(r)["id"]

	ts := time.Now().Unix()
	result, err := a.fileEBLCtrl.Transfer(ctx, ts, req)
	if err != nil {
		http.Error(w, err.Error(), model.ErrorToHttpStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.Warnf("transferEBL failed to encode/write response: %v", err)
	}
}

func (a *API) amendmentRequestEBL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	appID, _ := ctx.Value(middleware.APPLICATION_ID).(string)
	buID, _ := ctx.Value(middleware.BUSINESS_UNIT_ID).(string)

	var req trade_document.AmendmentRequestEBLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Application = appID
	req.RequestBy = buID
	req.ID = mux.Vars(r)["id"]

	ts := time.Now().Unix()
	result, err := a.fileEBLCtrl.AmendmentRequest(ctx, ts, req)
	if err != nil {
		http.Error(w, err.Error(), model.ErrorToHttpStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.Warnf("amendmentRequestEBL failed to encode/write response: %v", err)
	}
}