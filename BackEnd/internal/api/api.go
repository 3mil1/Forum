package api

import (
	_ "forum/docs"
	"forum/internal/service"
	"forum/internal/storage"
	"forum/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

var (
	prefix = "/api"
)

type API struct {
	config  *Config
	router  *http.ServeMux
	service *service.Service
}

func New(config *Config) *API {
	return &API{
		config: config,
		router: http.NewServeMux(),
	}
}

// Start http server/configure loggers, router, database connection and etc
func (api *API) Start() error {

	logger.InfoLogger.Println("Starting the application at port:", api.config.Port)

	api.configureRouter()
	if err := api.configureStore(); err != nil {
		return err
	}
	return http.ListenAndServe(":"+api.config.Port, CorsMW(api.router))
}

func (api *API) configureRouter() {
	api.router.Handle("/swagger/", httpSwagger.WrapHandler)

	// Auth
	api.router.Handle(prefix+"/auth/register", ErrorHandler(api.postUserRegister))
	api.router.Handle(prefix+"/auth/login", ErrorHandler(api.postToAuth))
	api.router.Handle(prefix+"/auth/me", api.UserIdentity(ErrorHandler(api.authMe)))
	api.router.Handle(prefix+"/auth/logout", api.UserIdentity(ErrorHandler(api.logOut)))

	// User
	api.router.Handle(prefix+"/users", api.UserIdentity(ErrorHandler(api.GetAllUsers)))
	//api.router.Handle(prefix+"/users", ErrorHandler(api.GetAllUsers))

	// Posts
	api.router.Handle(prefix+"/post/add", api.UserIdentity(ErrorHandler(api.addPost)))
	api.router.Handle(prefix+"/posts", ErrorHandler(api.allPosts))
	api.router.Handle(prefix+"/post", ErrorHandler(api.findByID))
	api.router.Handle(prefix+"/post/comments", ErrorHandler(api.commentsByPostID))
	api.router.Handle(prefix+"/post/mark", api.UserIdentity(ErrorHandler(api.addMark)))
	api.router.Handle(prefix+"/categories", ErrorHandler(api.showCategories))
	api.router.Handle(prefix+"/category", ErrorHandler(api.findByCategory))
	api.router.Handle(prefix+"/post/user_posts", api.UserIdentity(ErrorHandler(api.findByUser)))
	api.router.Handle(prefix+"/post/like", api.UserIdentity(ErrorHandler(api.findAllLiked)))
}

//configureStore method
func (api *API) configureStore() error {
	st := storage.New(api.config.Storage)
	if err := st.Open(); err != nil {
		return err
	}
	st.AddTables()
	sr := service.New(st)
	api.service = sr
	return nil
}
