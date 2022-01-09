package api

import (
	"encoding/json"
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/pkg/logger"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// @Summary      Add Post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Router       /post/new [post]
func (api *API) addPost(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("Post to Add_POST POST /api/post/new")
	var postFromRequest models.Post
	ct := r.Header.Get("Content-Type")
	if strings.Contains(ct, "multipart/form-data") {
		err := r.ParseMultipartForm(21 << 20)
		if err != nil {
			logger.ErrorLogger.Println(err.Error())
			return err
		}
		fileName, err := api.addImg(r)
		if err != nil {
			return err
		}
		strs := strings.Split(r.FormValue("categories"), ",")
		categories := make([]int, len(strs))
		for i := range categories {
			categories[i], _ = strconv.Atoi(strs[i])
		}
		postFromRequest = models.Post{
			ImagePath:  fileName,
			Categories: categories,
			Content:    r.FormValue("content"),
			Subject:    r.FormValue("subject"),
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&postFromRequest); err != nil {
			logger.ErrorLogger.Println("cannot unmarshal json")
			return appError.InvalidArgumentError(err, "wrong request format")
		}
	}

	// read from context
	user, _ := r.Context().Value("values").(userContext)
	fmt.Println("from context:", user)

	postFromRequest.UserId = user.userID

	post, err := api.service.Post().Create(&postFromRequest)
	if err != nil {
		return err
	}

	logger.InfoLogger.Println("Creating Post")
	return json.NewEncoder(w).Encode(post)
}

func (api *API) addImg(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		// check if a file is sent.
		if err.Error() == "http: no such file" {
			return "", nil
		} else {
			return "", appError.UnsupportedError(err, "Wrong file type. Accepted formats are jpeg, png, gif, svg")
		}
	}
	defer file.Close()

	if handler.Size > 20<<20 {
		return "", appError.InvalidArgumentError(err, "Image too large, you can upload files up to 20 MB")
	}

	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		logger.ErrorLogger.Println(err)
		return "", appError.SystemError(err)
	}
	fileName := fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))

	filetype := mime.TypeByExtension(filepath.Ext(handler.Filename))
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" && filetype != "image/svg+xml" {
		return "", appError.UnsupportedError(err, "Wrong file type. Accepted formats are jpeg, png, gif, svg")
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", appError.SystemError(err)
	}

	dst, err := os.Create(fileName)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", appError.SystemError(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", appError.SystemError(err)
	}
	return fileName, nil
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
	if comments == nil {
		logger.InfoLogger.Println("No comments to this post")
		return nil
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

// findByUser show all user's posts
// @Summary      Find posts by user login
// @Tags         posts
// @Produce      json
// @Router       /post/user_posts?login={login} [get]
func (api *API) findByUser(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("GET Posts by user login /api/post/user_posts")

	user, _ := r.Context().Value("values").(userContext)
	fmt.Println("from context:", user)
	posts, err := api.service.Post().FindByUserLogin(user.login)
	if err != nil {
		return err
	}
	logger.InfoLogger.Println("Posts found")
	return json.NewEncoder(w).Encode(posts)
}

// findByCategory show all posts in the category
// @Summary      Find posts by category
// @Tags         posts
// @Produce      json
// @Router       /post{category_id} [get]
func (api *API) findByCategory(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("GET Posts by category /api/category")
	cat := r.URL.Query().Get("category_id")
	var (
		posts []models.PostAndMarks
		err   error
	)
	if cat == "" {
		return api.allPosts(w, r)
	}

	id, err := strconv.Atoi(cat)
	if err != nil {
		return appError.InvalidArgumentError(err, "cannot get category id")
	}
	posts, err = api.service.Post().FindByCategory(id)
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		logger.InfoLogger.Println("No posts with that category")
	} else {
		logger.InfoLogger.Println("Posts found")
	}
	return json.NewEncoder(w).Encode(posts)
}

// findAllLiked show all user's liked posts
// @Summary      Find posts by user_id
// @Tags         posts
// @Produce      json
// @Router       /post/like [get]
func (api *API) findAllLiked(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("GET all liked posts /api/post/like")

	// read from context
	user, _ := r.Context().Value("values").(userContext)
	fmt.Println("from context:", user)
	posts, err := api.service.Post().FindAllLiked(user.userID)
	if err != nil {
		return err
	}
	if len(posts) == 0 {
		logger.InfoLogger.Println("No liked posts")
	} else {
		logger.InfoLogger.Println("Liked posts were found")
	}
	return json.NewEncoder(w).Encode(posts)
}
