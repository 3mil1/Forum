package storage

import (
	"errors"
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/pkg/logger"
	"github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

type PostRepo struct {
	storage *Storage
}

var (
	postColumns   = "user_id, content, subject, parent_id"
	markerColumns = "post_id, user_id, mark"
)

func (pr *PostRepo) Create(p *models.Post) (*models.Post, error) {
	var id *int
	if p.ParentId != 0 {
		id = &p.ParentId
	}
	tx, err := pr.storage.db.Begin()
	if err != nil {
		return nil, appError.DataBaseError(err)
	}
	query := fmt.Sprintf(`INSERT INTO posts (%s) VALUES ($1, $2, $3, $4) returning id, created_at`, postColumns)
	row := tx.QueryRow(query, p.UserId, p.Content, p.Subject, id)
	err = row.Scan(&p.Id, &p.CreatedAt)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.ErrorLogger.Println(err)
			return nil, appError.DataBaseError(err)
		}
		logger.ErrorLogger.Println(err)
		return nil, appError.SystemError(err)
	}
	var res []string
	for _, v := range p.Categories {
		r := fmt.Sprintf("(%v, %v)", p.Id, v)
		res = append(res, r)
	}
	query = fmt.Sprintf(`INSERT INTO posts_categories (post_id, category_id) VALUES %s`, strings.Join(res, ", "))
	if _, err := tx.Exec(query); err != nil {
		if err := tx.Rollback(); err != nil {
			logger.ErrorLogger.Println(err)
			return nil, appError.DataBaseError(err)
		}
		logger.ErrorLogger.Println(err)
		return nil, appError.SystemError(err)
	}
	if err := tx.Commit(); err != nil {
		logger.ErrorLogger.Println(err)
		return nil, appError.DataBaseError(err)
	}
	return p, nil
}

func (pr *PostRepo) ShowAll() ([]models.PostAndMarks, error) {
	query := fmt.Sprintf(`SELECT
									p.id, 
									p.user_id, 
									p.content, 
									p.subject, 
									p.created_at, 
									COALESCE(p.parent_id, 0),
									coalesce(sum(case when not ld.mark then 1 else 0 end), 0) as dislike, 
									coalesce(sum(case when ld.mark then 1 else 0 end), 0) as like 
								FROM posts p 
								LEFT JOIN likes_dislikes ld on p.id = ld.post_id
								WHERE p.parent_id is null	
								group by p.id`)
	rows, err := pr.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Куда читаем
	var post []models.PostAndMarks
	for rows.Next() {
		p := models.PostAndMarks{}
		err := rows.Scan(&p.Id, &p.Post.UserId, &p.Content, &p.Subject, &p.CreatedAt, &p.ParentId, &p.Dislikes, &p.Likes)
		if err != nil {
			// if database cannot read row
			log.Println(err)
			continue
		}
		post = append(post, p)
	}
	return post, nil
}

func (pr *PostRepo) FindByID(id int) (*models.PostAndMarks, error) {
	query := fmt.Sprintf(`SELECT 
									p.id,
									p.user_id,
									p.content,
									p.subject,
									p.created_at,
									coalesce(p.parent_id, 0),
									coalesce(sum(case when not ld.mark then 1 else 0 end), 0) AS dislike,
									coalesce(sum(case when ld.mark then 1 else 0 end), 0)     AS like
								FROM posts p
								LEFT JOIN likes_dislikes ld ON p.id = ld.post_id
								WHERE id=$1`)
	row := pr.storage.db.QueryRow(query, id)
	var post models.PostAndMarks
	err := row.Scan(&post.Id, &post.UserId, &post.Content, &post.Subject, &post.CreatedAt, &post.ParentId, &post.Dislikes, &post.Likes)
	if err != nil {
		return nil, appError.NotFoundError(err, "cannot find post")
	}
	return &post, nil
}

func (pr *PostRepo) FindByUserId(id string) ([]models.PostAndMarks, error) {
	query := fmt.Sprintf(`SELECT 
									p.id,
									p.user_id,
									p.content,
									p.subject,
									p.created_at,
									coalesce(p.parent_id, 0),
									coalesce(sum(case when not ld.mark then 1 else 0 end), 0) AS dislike,
									coalesce(sum(case when ld.mark then 1 else 0 end), 0)     AS like
								FROM posts p
								LEFT JOIN likes_dislikes ld ON p.id = ld.post_id
								WHERE p.user_id=$1 and p.parent_id is null
								group by p.id`)
	rows, err := pr.storage.db.Query(query, id)
	if err != nil {
		return nil, appError.SystemError(err)
	}
	defer rows.Close()
	var posts []models.PostAndMarks
	for rows.Next() {
		var post models.PostAndMarks
		err := rows.Scan(&post.Id, &post.UserId, &post.Content, &post.Subject, &post.CreatedAt, &post.ParentId, &post.Dislikes, &post.Likes)
		if err != nil {
			logger.ErrorLogger.Println(err)
			continue
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (pr *PostRepo) AddMark(m *models.Mark) (*models.Mark, error) {
	query := fmt.Sprintf("INSERT INTO likes_dislikes (%s) VALUES ($1, $2, $3)", markerColumns)
	if _, err := pr.storage.db.Exec(query, m.PostId, m.UserId, m.Mark); err != nil {
		var sErr sqlite3.Error
		if errors.As(err, &sErr) {

			// Delete value if exist
			if err.Error() == "UNIQUE constraint failed: likes_dislikes.post_id, likes_dislikes.user_id" {
				queryUpdate := fmt.Sprintf("DELETE FROM likes_dislikes WHERE post_id=$1 and user_id=$2")
				if _, err := pr.storage.db.Exec(queryUpdate, m.PostId, m.UserId); err != nil {
					return nil, err
				}
				return nil, nil
			}
		}
		return nil, err
	}
	return m, nil
}

func (pr *PostRepo) FindAllCommentsToPost(postID int) ([]models.Post, error) {
	query := fmt.Sprintf("SELECT %s FROM posts WHERE parent_id=$1", postColumns+", created_at, id")
	rows, err := pr.storage.db.Query(query, postID)
	if err != nil {
		return nil, appError.SystemError(err)
	}
	defer rows.Close()

	var comments []models.Post
	for rows.Next() {
		var c models.Post
		if err := rows.Scan(&c.UserId, &c.Content, &c.Subject, &c.ParentId, &c.CreatedAt, &c.Id); err != nil {
			log.Println(err)
			continue
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (pr *PostRepo) ShowAllCategories() ([]models.Category, error) {
	query := fmt.Sprintf("SELECT id, name FROM categories")
	rows, err := pr.storage.db.Query(query)
	if err != nil {
		return nil, appError.DataBaseError(err)
	}
	var categories []models.Category // nil

	for rows.Next() {
		var s models.Category
		if err := rows.Scan(&s.Id, &s.Name); err != nil {
			log.Println(err)
			continue
		}
		categories = append(categories, s)
	}
	if len(categories) == 0 {
		return nil, appError.NotFoundError(nil, "no categories were found")
	}
	return categories, nil
}

func (pr *PostRepo) GetCategoriesByPostID(id int) ([]models.Category, error) {
	query := fmt.Sprintf(`SELECT 
									c.id, 
									c.name 
								FROM categories c 
								LEFT JOIN posts_categories pc ON c.id = pc.category_id 
								WHERE post_id=$1`)
	rows, err := pr.storage.db.Query(query, id)
	if err != nil {
		return nil, appError.SystemError(err)
	}
	defer rows.Close()
	var cat []models.Category
	for rows.Next() {
		var r models.Category
		if err := rows.Scan(&r.Id, &r.Name); err != nil {
			logger.ErrorLogger.Println(err)
			continue
		}
		cat = append(cat, r)
	}
	return cat, nil
}
