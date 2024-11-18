package link

import (
	"net/http"
	"server/configs"
	"server/pkg/event"
	"server/pkg/logger"
	"server/pkg/middleware"
	"server/pkg/request"
	"server/pkg/response"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	*configs.Config
	LinkRepository *LinkRepository
	EventBus       *event.EventBys
}

type LinkHandler struct {
	*configs.Config
	LinkRepository *LinkRepository
	EventBus       *event.EventBys
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		Config:         deps.Config,
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}

	router.Handle("/link/{hash}", middleware.IsAuthed(handler.Get(), deps.Config))
	router.Handle("POST /link", middleware.IsAuthed(handler.Create(), deps.Config))
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.IsAuthed(handler.Delete(), deps.Config))
	router.Handle("/link", middleware.IsAuthed(handler.GetAll(), deps.Config))
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("LinkHandler: Create")

		_, ok := req.Context().Value(middleware.ContextEmailKey).(string)

		if !ok {
			http.Error(w, ErrNoAuthorized, http.StatusUnauthorized)
			return
		}

		body, err := request.HandleBody[LinkCreateRequest](&w, req)
		if err != nil {
			response.Json(w, err.Error(), 402)
			return
		}

		link := NewLink(body.Url)
		for {
			existedLink, _ := handler.LinkRepository.Get(link.Hash)
			if existedLink == nil {
				break
			}
			link.GenerateHash()
		}

		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, createdLink, 201)
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("LinkHandler: Update")

		body, err := request.HandleBody[LinkUpdateRequest](&w, req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.LinkRepository.checkIsExist(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, link, 201)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("LinkHandler: Delete")

		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.LinkRepository.checkIsExist(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.LinkRepository.Delete(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, "success", 200)
	}
}

func (handler *LinkHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("LinkHandler: Get")

		hash := req.PathValue("hash")
		link, err := handler.LinkRepository.Get(hash)

		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})
		http.Redirect(w, req, link.Url, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		limit := req.URL.Query().Get("limit")
		offset := req.URL.Query().Get("offset")

		logger.Log("limit: ", limit)

		limitInt, err := strconv.Atoi(limit)

		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		offsetInt, err := strconv.Atoi(offset)

		if err != nil {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}

		count := handler.LinkRepository.Count()

		links := handler.LinkRepository.GetAll(limitInt, offsetInt)

		linksResponse := &LinkGetAllResponse{
			Links: links,
			Count: count,
		}

		response.Json(w, linksResponse, 200)
	}
}
