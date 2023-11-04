package handler

import (
	"context"
	"encoding/json"
	"log"
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
	ReadTodos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: ReadTodos}, nil
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
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	DeleteTodo := model.DeleteTODOResponse{Message: "Successfully delete todos"}
	return &DeleteTodo, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		var deleteRequest model.DeleteTODORequest
		err := json.NewDecoder(r.Body).Decode(&deleteRequest)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if len(deleteRequest.IDs) == 0 {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		response, err := h.Delete(ctx, &deleteRequest)

		log.Printf("req: %v", deleteRequest)
		log.Printf("resp: %v", response)
		log.Printf("err: %v", err)
		if err != nil {
			http.Error(w, "Not found ID", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(response)
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}

	case http.MethodGet:
		query := r.URL.Query()
		prevIDStr := query.Get("prev_id")
		sizeStr := query.Get("size")

		var prev_id, size int64
		// prev_id のエラーチェック
		if prevIDStr == "" {
			prev_id = 0
		} else {
			prevID, err := strconv.Atoi(prevIDStr)
			prev_id = int64(prevID)
			if err != nil {
				http.Error(w, "Invalid prev_id", http.StatusBadRequest)
				return
			}
		}

		// size のエラーチェック
		if sizeStr == "" {
			size = 5
		} else {
			size_tmp, err := strconv.Atoi(sizeStr)
			size = int64(size_tmp)
			if err != nil {
				http.Error(w, "Invalid size", http.StatusBadRequest)
				return
			}
		}
		readRequest := model.ReadTODORequest{PrevID: int64(prev_id), Size: int64(size)}
		log.Println(readRequest)
		ctx := r.Context()
		readTodo, err := h.Read(ctx, &readRequest)
		if err != nil {
			http.Error(w, "Failed to read TODO", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(readTodo)
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}

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
			http.Error(w, "Failed to update TODO", http.StatusNotFound)
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
