package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var req model.CreateTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.Create(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(json)
	case "PUT":
		var req model.UpdateTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.ID == 0 || req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.Update(r.Context(), &req)
		if err != nil {
			var notFound *model.ErrNotFound
			if errors.As(err, &notFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(json)
	case "GET":
		prevID := r.URL.Query().Get("prev_id")
		size := r.URL.Query().Get("size")
		if prevID == "" {
			prevID = "0"
		}
		if size == "" {
			size = "0"
		}

		prevIDInt, err := strconv.Atoi(prevID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req := &model.ReadTODORequest{
			PrevID: int64(prevIDInt),
			Size:   int64(sizeInt),
		}

		res, err := h.Read(r.Context(), req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(json)
	default:
		fmt.Fprint(w, "Method not allowed.\n")
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	res, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.CreateTODOResponse{
		TODO: *res,
	}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}

	res := make([]model.TODO, 0)
	for _, todo := range todos {
		res = append(res, *todo)
	}

	return &model.ReadTODOResponse{
		TODOs: res,
	}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	res, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.UpdateTODOResponse{
		TODO: *res,
	}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
