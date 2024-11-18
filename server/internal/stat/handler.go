package stat

import (
	"net/http"
	"server/configs"
	"server/pkg/middleware"
	"server/pkg/response"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandlerDeps struct {
	Config         *configs.Config
	StatRepository *StatRepository
}

type StatHandler struct {
	StatRepository *StatRepository
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}

	router.Handle("/stat", middleware.IsAuthed(handler.GetStat(), deps.Config))
}

func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		from, err := time.Parse(time.DateOnly, req.URL.Query().Get("from"))

		if err != nil {
			http.Error(w, ErrFromDate, http.StatusBadRequest)
			return
		}

		to, err := time.Parse(time.DateOnly, req.URL.Query().Get("to"))

		if err != nil {
			http.Error(w, ErrToDate, http.StatusBadRequest)
			return
		}

		by := req.URL.Query().Get("by")

		if by != GroupByMonth && by != GroupByDay {
			http.Error(w, ErrBy, http.StatusBadRequest)
			return
		}

		stats := handler.StatRepository.GetStats(by, from, to)

		response.Json(w, stats, 200)
	}
}
