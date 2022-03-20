package todo

import (
	"context"
	"net/http"
	"time"
)

// Структура для запуска http-сервера. Является абстракцией над структурой server пакета http
type Server struct {
	httpServer *http.Server
}

// Метод запуска сервера
func (s *Server) Run(port string, handler http.Handler) error {

	// Инкапсуляция нужных для запуска значений
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // = 1 Mb
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// Метод, который запускает бесконечный цикл
	// и слушает все входящие запросы для их обработки
	return s.httpServer.ListenAndServe()
}

// Метод, который используется при выходе из приложения
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
