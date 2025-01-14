package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	rp "runtime/pprof"
	"runtime/trace"
	"sync"

	"github.com/ArthurGopher/AuthorizationJWT/auth"
	"github.com/ArthurGopher/AuthorizationJWT/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/ArthurGopher/AuthorizationJWT/docs"
)

var (
	cpuProfileFile *os.File
	traceFile      *os.File
	profileMutex   sync.Mutex
	isProfiling    bool
)

func main() {
	// cоздаем JWT-токен с секретным ключом
	auth.TokenAuth = jwtauth.New("HS256", []byte("secret-key"), nil)
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	r.Post("/api/register", auth.UsersRegisterHandler)
	r.Post("/api/login", auth.LoginHandler)

	// Группа маршрутов с профилированием
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth)) // Проверка токена
		r.Use(auth.AuthMiddleware)              // Middleware для авторизации

		r.Post("/api/address/search", handlers.SearchHandler)
		r.Post("/api/address/geocode", handlers.GeocodeHandler)

		// добавляем маршруты для pprof
		r.Route("/mycustompath/pprof", func(r chi.Router) {
			r.HandleFunc("/", pprof.Index)
			r.HandleFunc("/cmdline", pprof.Cmdline)
			r.HandleFunc("/profile", pprof.Profile)
			r.HandleFunc("/symbol", pprof.Symbol)
			r.HandleFunc("/trace", pprof.Trace)
			r.HandleFunc("/allocs", pprof.Handler("allocs").ServeHTTP)
			r.HandleFunc("/block", pprof.Handler("block").ServeHTTP)
			r.HandleFunc("/goroutine", pprof.Handler("goroutine").ServeHTTP)
			r.HandleFunc("/heap", pprof.Handler("heap").ServeHTTP)
			r.HandleFunc("/mutex", pprof.Handler("mutex").ServeHTTP)
			r.HandleFunc("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		})

		// API для запуска профилирования CPU и трассировки
		r.Post("/start-profiling", startProfilingHandler)
		r.Post("/stop-profiling", stopProfilingHandler)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

func startProfilingHandler(w http.ResponseWriter, r *http.Request) {
	profileMutex.Lock()
	defer profileMutex.Unlock()

	if isProfiling {
		http.Error(w, "Профилирование уже запущено", http.StatusBadRequest)
		return
	}

	var err error
	cpuProfileFile, err = os.Create("cpu_profile.prof")
	if err != nil {
		http.Error(w, "Не удалось создать файл профиля CPU", http.StatusInternalServerError)
		return
	}

	if err := rp.StartCPUProfile(cpuProfileFile); err != nil {
		http.Error(w, "Не удалось запустить профилирование CPU", http.StatusInternalServerError)
		return
	}

	traceFile, err = os.Create("trace.out")
	if err != nil {
		http.Error(w, "Не удалось создать файл трассировки", http.StatusInternalServerError)
		return
	}

	if err := trace.Start(traceFile); err != nil {
		http.Error(w, "Не удалось запустить трассировку", http.StatusInternalServerError)
		return
	}

	isProfiling = true
	w.Write([]byte("Профилирование начато"))
}

func stopProfilingHandler(w http.ResponseWriter, r *http.Request) {
	profileMutex.Lock()
	defer profileMutex.Unlock()

	if !isProfiling {
		http.Error(w, "Профилирование не запущено", http.StatusBadRequest)
		return
	}

	rp.StopCPUProfile()
	if cpuProfileFile != nil {
		cpuProfileFile.Close()
	}

	trace.Stop()
	if traceFile != nil {
		traceFile.Close()
	}

	isProfiling = false
	w.Write([]byte("Профилирование остановлено"))
}
