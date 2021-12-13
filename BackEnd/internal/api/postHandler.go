package api

import (
	"encoding/json"
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/pkg/logger"
	"net/http"
	"strconv"
)

// @Summary      Add Post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Router       /post/new [post]
func (api *API) addPost(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("Post to Add_POST POST /api/post/new")

	var postFromJson models.Post
	err := json.NewDecoder(r.Body).Decode(&postFromJson)
	if err != nil {
		logger.InfoLogger.Println("Invalid json received from client")
		return appError.NewAppError(err, "Provided json is invalid", http.StatusBadRequest)
	}

	// read from context
	user, _ := r.Context().Value("values").(userContext)
	fmt.Println("from context:", user)

	postFromJson.UserId = user.userID

	post, err := api.service.Post().Create(&postFromJson)
	if err != nil {
		fmt.Println(err)
		return err
	}

	logger.InfoLogger.Println("Creating Post")
	return json.NewEncoder(w).Encode(post)
}

// @Summary      Show All Posts
// @Tags         posts
// @Produce      json
// @Router       /posts [get]
func (api *API) allPosts(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("GET to Show posts GET /api/posts")

	allPosts, err := api.service.Post().ShowAll()
	if err != nil {
		logger.InfoLogger.Println("allPosts handler:", err)
		return err
	}
	logger.InfoLogger.Println("Get All Posts GET /api/posts")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(allPosts)
}

// @Summary      Find post by ID
// @Tags         posts
// @Produce      json
// @Router       /post{id} [get]
func (api *API) findByID(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("GET Find Post by ID /api/post")

	id := r.URL.Query().Get("id")
	pID, err := strconv.Atoi(id)
	if err != nil {
		return appError.InvalidArgumentError(err, "cannot get post id")
	}
	post, err := api.service.Post().FindById(pID)
	if err != nil {
		logger.InfoLogger.Println("findByID handler:", err)
		return appError.NewAppError(err, err.Error(), http.StatusBadRequest)
	}
	logger.InfoLogger.Println("Post found")
	return json.NewEncoder(w).Encode(post)
}

func (api *API) commentsByPostID(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("GET Find Comments by postID /api/post/comments")

	id := r.URL.Query().Get("id")
	pID, err := strconv.Atoi(id)
	if err != nil {
		return appError.InvalidArgumentError(err, "cannot get post id")
	}

	comments, err := api.service.Post().CommentsByPostId(pID)
	if err != nil {
		logger.InfoLogger.Println("findByID handler:", err)
		return appError.NewAppError(err, err.Error(), http.StatusBadRequest)
	}

	logger.InfoLogger.Println("Comments found")
	return json.NewEncoder(w).Encode(comments)
}

// @Summary      Add mark
// @Tags         posts
// @Accept       json
// @Produce      json
// @Router       /post/mark [post]
func (api *API) addMark(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("POST add mark POST /api/post/mark")

	var markFromJson models.Mark
	err := json.NewDecoder(r.Body).Decode(&markFromJson)
	if err != nil {
		logger.InfoLogger.Println("Invalid json received from client")
		return appError.NewAppError(err, "Provided json is invalid", http.StatusBadRequest)
	}

	m, err := api.service.Post().AddMark(&markFromJson)
	if err != nil {
		logger.InfoLogger.Println("addMark handler:", err)
		return err
	}
	if m == nil {
		w.WriteHeader(http.StatusNoContent)
		logger.InfoLogger.Println("Mark deleted")
	} else {
		w.WriteHeader(http.StatusOK)
		logger.InfoLogger.Println("Mark added")
	}

	return nil
}

func (api *API) showCategories(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("GET to Show posts GET /api/categories")

	categories, err := api.service.Post().ShowAllCategories()
	if err != nil {
		logger.InfoLogger.Println("showCategories handler:", err)
		return err
	}
	w.WriteHeader(http.StatusOK)

	logger.InfoLogger.Println("Get All Posts GET /api/categories")
	return json.NewEncoder(w).Encode(categories)

}

// FindByUser show all user's posts
// @Summary      Find posts by user ID
// @Tags         posts
// @Produce      json
// @Router       /post{login} [get]
func (api *API) findByUser(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("GET Posts by userID /api/post/user_posts")
	login := r.URL.Query().Get("login")
	posts, err := api.service.Post().FindByUserLogin(login)
	if err != nil {
		return err
	}
	logger.InfoLogger.Println("Posts found")
	return json.NewEncoder(w).Encode(posts)
}

// FindByCategory show all posts in the category
// @Summary      Find posts by category
// @Tags         posts
// @Produce      json
// @Router       /post{name} [get]
func (api *API) findByCategory(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("GET Posts by category /api/category")
	cat := r.URL.Query().Get("category_id")
	id, err := strconv.Atoi(cat)
	if err != nil {
		return appError.InvalidArgumentError(err, "cannot get category id")
	}
	posts, err := api.service.Post().FindByCategory(id)
	if err != nil {
		return err
	}
	logger.InfoLogger.Println("Posts found")
	return json.NewEncoder(w).Encode(posts)
}
