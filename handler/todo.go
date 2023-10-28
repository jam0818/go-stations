package handler

import (
	"context"
	"encoding/json"
	"net/http"

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

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	createdTodo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: *createdTodo}, err
}



// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	UpdateTodo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *UpdateTodo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodPost:
			// リクエストボディからJSONデータをデコード
			var createRequest model.CreateTODORequest
			err := json.NewDecoder(r.Body).Decode(&createRequest)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			// subjectが空文字列の場合、400 Bad Requestを返す
			if createRequest.Subject == "" {
				http.Error(w, "Subject cannot be empty", http.StatusBadRequest)
				return
			}

			// TODOを作成し、データベースに保存
			ctx := r.Context()
			createdTodo, err := h.Create(ctx, &createRequest)
			if err != nil {
				http.Error(w, "Failed to create TODO", http.StatusInternalServerError)
				return
			}
			// log.Println(*&createdTodo)

			// レスポンスをJSONエンコードして返す
			w.Header().Set("Content-Type", "application/json")
			encoder := json.NewEncoder(w)
			err = encoder.Encode(createdTodo)
			if err != nil {
				http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
				return
			}
		case http.MethodPut:
			var updateRequest model.UpdateTODORequest
			err := json.NewDecoder(r.Body).Decode(&updateRequest)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			if updateRequest.Subject == "" {
				http.Error(w, "Subject cannot be empty", http.StatusBadRequest)
				return
			} 
			if updateRequest.ID == 0 {
				http.Error(w, "No id have assigned", http.StatusBadRequest)
				return				
			}
			// TODOを作成し、データベースに保存
			ctx := r.Context()
			updateTodo, err := h.Update(ctx, &updateRequest)
			if err != nil {
				http.Error(w, "Failed to create TODO", http.StatusNotFound)
				return
			}
			// レスポンスをJSONエンコードして返す
			w.Header().Set("Content-Type", "application/json")
			encoder := json.NewEncoder(w)
			err = encoder.Encode(updateTodo)
			if err != nil {
				http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
				return
			}
		
		default:
			// HTTPメソッドがPOST, PUTでない場合、400 Bad Requestを返す
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return 
	}
}
