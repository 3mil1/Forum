package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/pkg/logger"
	"log"
	"strings"
)

type PostRepo struct {
	storage *Storage
}

var (
	postColumns   = "user_id, content, subject, parent_id, image"
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
	query := fmt.Sprintf(`INSERT INTO posts (%s) VALUES ($1, $2, $3, $4, $5) returning id, created_at`, postColumns)
	row := tx.QueryRow(query, p.UserId, p.Content, p.Subject, id, p.ImagePath)
	err = row.Scan(&p.Id, &p.CreatedAt)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.ErrorLogger.Println(err)
			return nil, appError.DataBaseError(err)
		}
		logger.ErrorLogger.Println(err)
		return nil, appError.SystemError(err)
	}
	if len(p.Categories) != 0 {
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
	}
	if err := tx.Commit(); err != nil {
		logger.ErrorLogger.Println(err)
		return nil, appError.DataBaseError(err)
	}
	return p, nil
}

func (pr *PostRepo) ShowAll() ([]models.PostAndMarks, error) {
	query := fmt.Sprintf(`SELECT p.id,
       p.user_id,
       u.login,
       p.content,
       p.subject,
       p.created_at,
       COALESCE(p.parent_id, 0)      as parent_id,
       coalesce(dislike, 0)          as dislike,
       coalesce(like, 0)             as dislike,
       group_concat(distinct c.name) as category_name
FROM posts p
         LEFT JOIN (
    Select post_id,
           sum(case when not mark then 1 else 0 end) AS dislike,
           sum(case when mark then 1 else 0 end)     AS like
    FROM likes_dislikes
    group by post_id
) as ld ON p.id = ld.post_id
         INNER JOIN posts_categories pc on p.id = pc.post_id
         INNER JOIN categories c on c.id = pc.category_id
         INNER JOIN users u on u.id = p.user_id
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
		err := rows.Scan(&p.Id, &p.Post.UserId, &p.UserLogin, &p.Content, &p.Subject, &p.CreatedAt, &p.ParentId, &p.Dislikes, &p.Likes, &p.Categories)
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
	query := `SELECT p.id,
       p.user_id,
       u.login,
       p.content,
       p.subject,
       p.created_at,
       p.image,
       COALESCE(p.parent_id, 0)      as parent_id,
       coalesce(dislike, 0)          as dislike,
       coalesce(like, 0)             as dislike,
       group_concat(distinct c.name) as category_name
FROM posts p
         LEFT JOIN (
    Select post_id,
           sum(case when not mark then 1 else 0 end) AS dislike,
           sum(case when mark then 1 else 0 end)     AS like
    FROM likes_dislikes
    group by post_id
) as ld ON p.id = ld.post_id
         INNER JOIN posts_categories pc on p.id = pc.post_id
         INNER JOIN categories c on c.id = pc.category_id
         INNER JOIN users u on u.id = p.user_id
WHERE p.id =$1`
	row := pr.storage.db.QueryRow(query, id)
	var post models.PostAndMarks
	err := row.Scan(&post.Id, &post.UserId, &post.UserLogin, &post.Content, &post.Subject, &post.CreatedAt, &post.ImagePath, &post.ParentId, &post.Dislikes, &post.Likes, &post.Categories)
	if err != nil {
		return nil, appError.NotFoundError(err, "cannot find post")
	}
	return &post, nil
}

func (pr *PostRepo) FindByUserId(id string) ([]models.PostAndMarks, error) {
	query := fmt.Sprintf(`SELECT p.id,
       p.user_id,
       u.login,
       p.content,
       p.subject,
       p.created_at,
       COALESCE(p.parent_id, 0)      as parent_id,
       coalesce(dislike, 0)          as dislike,
       coalesce(like, 0)             as dislike,
       group_concat(distinct c.name) as category_name
FROM posts p
         LEFT JOIN (
    Select post_id,
           sum(case when not mark then 1 else 0 end) AS dislike,
           sum(case when mark then 1 else 0 end)     AS like
    FROM likes_dislikes
    group by post_id
) as ld ON p.id = ld.post_id
         INNER JOIN posts_categories pc on p.id = pc.post_id
         INNER JOIN categories c on c.id = pc.category_id
         INNER JOIN users u on u.id = p.user_id
WHERE u.id =$1
group by p.id
ORDER BY p.created_at desc`)
	rows, err := pr.storage.db.Query(query, id)
	if err != nil {
		return nil, appError.SystemError(err)
	}
	defer rows.Close()
	var posts []models.PostAndMarks
	for rows.Next() {
		var post models.PostAndMarks
		err := rows.Scan(&post.Id, &post.UserId, &post.UserLogin, &post.Content, &post.Subject, &post.CreatedAt, &post.ParentId, &post.Dislikes, &post.Likes, &post.Categories)
		if err != nil {
			logger.ErrorLogger.Println(err)
			continue
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (pr PostRepo) GetMark(m *models.Mark) (*bool, error) {
	query := fmt.Sprintf("SELECT mark FROM likes_dislikes WHERE post_id=$1 and user_id=$2")
	row := pr.storage.db.QueryRow(query, m.PostId, m.UserId)
	var mark *bool
	if err := row.Scan(&mark); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mark, nil
}

func (pr *PostRepo) AddMark(m *models.Mark) (*models.Mark, error) {
	query := fmt.Sprintf("INSERT INTO likes_dislikes (%s) VALUES ($1, $2, $3)", markerColumns)
	if _, err := pr.storage.db.Exec(query, m.PostId, m.UserId, m.Mark); err != nil {
		return nil, err
	}
	return m, nil
}

func (pr *PostRepo) UpdateMark(m *models.Mark) (*models.Mark, error) {
	query := fmt.Sprintf("UPDATE likes_dislikes SET mark=$1 WHERE post_id=$2 and user_id=$3")
	if _, err := pr.storage.db.Exec(query, m.Mark, m.PostId, m.UserId); err != nil {
		return nil, err
	}
	return m, nil
}

func (pr PostRepo) DeleteMark(m *models.Mark) (*models.Mark, error) {
	query := fmt.Sprintf("DELETE FROM likes_dislikes WHERE post_id=$1 and user_id=$2")
	if _, err := pr.storage.db.Exec(query, m.PostId, m.UserId); err != nil {
		return nil, err
	}
	return nil, nil
}

func (pr *PostRepo) FindAllCommentsToPost(postID int) ([]models.PostAndMarks, error) {
	query := fmt.Sprintf(`with recursive cte (id, user_id, parent_id, content, created_at) as (
    select id, user_id, parent_id, content, created_at
    from posts
    where parent_id =$1
    union all
    select p.id,
           p.user_id,
           p.parent_id,
           p.content,
           p.created_at
    from posts p
             inner join cte on p.parent_id = cte.id
)
select cte.id,
       cte.user_id,
       u.login,
       cte.content,
       cte.created_at,
       cte.parent_id,
       coalesce(sum(case when not ld.mark then 1 else 0 end), 0) AS dislike,
       coalesce(sum(case when ld.mark then 1 else 0 end), 0)     AS like
from cte
         LEFT JOIN users u on cte.user_id = u.id
         LEFT JOIN likes_dislikes ld on cte.id = ld.post_id
group by cte.id`)
	rows, err := pr.storage.db.Query(query, postID)
	if err != nil {
		return nil, appError.SystemError(err)
	}
	defer rows.Close()

	var comments []models.PostAndMarks
	for rows.Next() {
		var p models.PostAndMarks
		if err := rows.Scan(&p.Id, &p.UserId, &p.UserLogin, &p.Content, &p.CreatedAt, &p.ParentId, &p.Dislikes, &p.Likes); err != nil {
			log.Println(err)
			continue
		}
		comments = append(comments, p)
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

func (pr *PostRepo) FindByCategory(cat int) ([]models.PostAndMarks, error) {
	query := `SELECT p.id,
       p.user_id,
       u.login,
       p.content,
       p.subject,
       p.created_at,
       COALESCE(p.parent_id, 0)      as parent_id,
       coalesce(dislike, 0)          as dislike,
       coalesce(like, 0)             as dislike,
       group_concat(distinct c.name) as category_name
FROM posts p
         LEFT JOIN (
    Select post_id,
           sum(case when not mark then 1 else 0 end) AS dislike,
           sum(case when mark then 1 else 0 end)     AS like
    FROM likes_dislikes
    group by post_id
) as ld ON p.id = ld.post_id
         INNER JOIN posts_categories pc on p.id = pc.post_id
         INNER JOIN categories c on c.id = pc.category_id
         INNER JOIN users u on u.id = p.user_id
		WHERE p.parent_id is null and p.id IN (SELECT post_id FROM posts_categories WHERE category_id=$1)
		group by p.id`
	rows, err := pr.storage.db.Query(query, cat)
	if err != nil {
		return nil, appError.SystemError(err)
	}
	defer rows.Close()
	var posts []models.PostAndMarks
	for rows.Next() {
		var post models.PostAndMarks
		err := rows.Scan(&post.Id, &post.UserId, &post.UserLogin, &post.Content, &post.Subject, &post.CreatedAt, &post.ParentId, &post.Dislikes, &post.Likes, &post.Categories)
		if err != nil {
			logger.ErrorLogger.Println(err)
			continue
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (pr *PostRepo) FindAllLiked(userID string) ([]models.PostAndMarks, error) {
	query := `SELECT p.id,
       p.user_id,
       u.login,
       p.content,
       p.subject,
       p.created_at,
       COALESCE(p.parent_id, 0)      as parent_id,
       coalesce(dislike, 0)          as dislike,
       coalesce(like, 0)             as dislike,
       group_concat(distinct c.name) as category_name
FROM posts p
         LEFT JOIN (
    Select post_id,
           sum(case when not mark then 1 else 0 end) AS dislike,
           sum(case when mark then 1 else 0 end)     AS like
    FROM likes_dislikes
    group by post_id
) as ld ON p.id = ld.post_id
         INNER JOIN posts_categories pc on p.id = pc.post_id
         INNER JOIN categories c on c.id = pc.category_id
         INNER JOIN users u on u.id = p.user_id
WHERE p.parent_id is null and p.id IN (SELECT post_id FROM likes_dislikes WHERE mark=1 and user_id=$1)
group by p.id
ORDER BY p.created_at desc`

	rows, err := pr.storage.db.Query(query, userID)
	if err != nil {
		return nil, appError.SystemError(err)
	}
	defer rows.Close()
	var posts []models.PostAndMarks
	for rows.Next() {
		var post models.PostAndMarks
		err := rows.Scan(&post.Id, &post.UserId, &post.UserLogin, &post.Content, &post.Subject, &post.CreatedAt, &post.ParentId, &post.Dislikes, &post.Likes, &post.Categories)
		if err != nil {
			logger.ErrorLogger.Println(err)
			continue
		}
		posts = append(posts, post)
	}
	return posts, nil
}
