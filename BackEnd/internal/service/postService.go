package service

import (
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/internal/storage"
	"strings"
)

type PostService struct {
	storage *storage.Storage
}

func (p *PostService) Create(post *models.Post) (*models.Post, error) {
	trimmedPost := strings.TrimSpace(post.Content)
	if trimmedPost == "" {
		return nil, appError.InvalidArgumentError(nil, "you are trying to create a post without text")
	}
	post.Content = trimmedPost
	if post.Subject == "" && post.ParentId == 0 {
		return nil, appError.InvalidArgumentError(nil, "topic is missing")
	}
	if len(post.Categories) == 0 && post.ParentId == 0 {
		return nil, appError.InvalidArgumentError(nil, "category is missing")
	}
	posts, err := p.storage.Post().Create(post)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *PostService) ShowAll() ([]models.PostAndMarks, error) {
	posts, err := p.storage.Post().ShowAll()
	if err != nil {
		return nil, appError.SystemError(err)
	}
	return posts, nil
}

//FindById Поиск сообщения со всеми комментариями
func (p *PostService) FindById(id int) (*models.PostAndMarks, error) {
	pMC, err := p.storage.Post().FindByID(id)
	if err != nil {
		return nil, err
	}
	//pMC.Categories, err = p.storage.Post().GetCategoriesByPostID(pMC.Id)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//if len(pMC.Categories) == 0 {
	//	logger.InfoLogger.Println("This post has no categories")
	//}
	//if pMC.Comments, err = p.storage.Post().FindAllCommentsToPost(pMC.Id); err != nil {
	//	return nil, err
	//}
	return pMC, nil
}

func (p *PostService) CommentsByPostId(id int) ([]models.PostAndMarks, error) {
	comments, err := p.storage.Post().FindAllCommentsToPost(id)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// FindByUserLogin найти все посты пользователя
func (p *PostService) FindByUserLogin(login string) ([]models.PostAndMarks, error) {
	u, err := p.storage.User().FindByLogin(login)
	if err != nil {
		return nil, err
	}
	fmt.Println(u.ID)
	posts, err := p.storage.Post().FindByUserId(u.ID)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, appError.InvalidArgumentError(nil, "user has no posts")
	}
	return posts, nil
}

// FindByCategoryID найти все сообщения заданной темы
func (p *PostService) FindByCategoryID(cat int) ([]models.PostAndMarks, error) {
	return nil, nil
}

func (p *PostService) AddMark(m *models.Mark) (*models.Mark, error) {
	mark, err := p.storage.Post().AddMark(m)
	if err != nil {
		return nil, appError.SystemError(err)
	}
	return mark, nil
}

func (p *PostService) ShowAllCategories() ([]models.Category, error) {
	return p.storage.Post().ShowAllCategories()
}
