package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum/models"

	_ "github.com/mattn/go-sqlite3"
)

type postSQL struct {
	db *sql.DB
}

// NewPostsSQL create new database struct
func NewPostsSQL(db *sql.DB) *postSQL {
	return &postSQL{db: db}
}

// CreatePost
// INSERT INTO posts (user_id, date, title, content) values(  1, "2023-05-01 13:35:04.556898354+06:00" , "post about JS", "JavaScript is a scripting or programming language that allows you to implement complex features on web pages — every time a web page does more than just sit there and display static information for you to look at — displaying timely content updates, interactive maps, animated 2D,3D graphics, scrolling video jukeboxes, etc. — you can bet that JavaScript is probably involved. It is the third layer of the layer cake of standard web technologies, two of which (HTML and CSS) we have covered in much more detail in other parts of the Learning Area.");
func (r *postSQL) CreatePost(post models.Post, currentUserId int) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (user_id, date, title, content,image_name) values (?,?,?,?,?) RETURNING id`, postsTable)
	row := r.db.QueryRow(query, post.AuthorID, post.Date, post.Title, post.Content, post.ImageName)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	err := InsertPostActivity(r, currentUserId, id, 0, post.AuthorID, "create")
	if err != nil {
		return id, err
	}
	return id, nil
}

// AddCategoriesToPost
// INSERT INTO posts_categories (post_id, category_id) values (1, 2);
func (r *postSQL) AddCategoryToPost(postId, catId int) error {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (post_id, category_id) values (?, ?) RETURNING id`, categoriesToPostsTable)
	row := r.db.QueryRow(query, postId, catId)
	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil
}

// SELECT p.*, group_concat(c.name, ", "), (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = true) as Likes, (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = false) as Dislikes, (SELECT pl.id FROM  posts_likes pl WHERE pl.post_id = p.id and type = true and pl.user_id = 1) as like_id FROM posts p LEFT JOIN posts_categories pc ON p.id = pc.post_id LEFT JOIN categories c ON c.id = pc.category_id group by p.id;

// GetAllPosts
// SELECT p.*, group_concat(c.name, ", "), (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = true) as Likes, (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = false) as Dislikes, FROM posts p LEFT JOIN posts_categories pc ON p.id = pc.post_id LEFT JOIN categories c ON c.id = pc.category_id group by p.id;
func (r *postSQL) GetAllPosts(currentUserId int) ([]models.Post, error) {
	var posts []models.Post
	query := fmt.Sprintf(`
	
	SELECT p.*, u.username as user_name, group_concat(c.name, ", ") as categories, 
	(SELECT Count(*) FROM %s  pl WHERE pl.post_id = p.id and type = true) as likes,
	(SELECT Count(*) FROM %s  pl WHERE pl.post_id = p.id and type = false) as dislikes,
	(SELECT pl.id FROM %s pl WHERE pl.post_id = p.id and type = true and pl.user_id = ?) as my_like_id 
	FROM %s  p 
	LEFT JOIN %s pc ON p.id = pc.post_id 
	LEFT JOIN %s c ON c.id = pc.category_id 
	LEFT JOIN %s u ON u.id = p.user_id  
	group by p.id;`, postsLikesTable, postsLikesTable, postsLikesTable, postsTable, categoriesToPostsTable, categoriesTable, usersTable)

	rows, err := r.db.Query(query, currentUserId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no posts found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("can't get posts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var myLikeId sql.NullInt32
		var post models.Post
		if err = rows.Scan(
			&post.Id,
			&post.AuthorID,
			&post.Date,
			&post.Title,
			&post.Content,
			&post.ImageName,
			&post.AuthorName,
			&post.Categories,
			&post.Likes,
			&post.Dislikes,
			&myLikeId,
		); err != nil {
			return nil, fmt.Errorf("can't scan posts: %w", err)
		}
		post.MyLikeId = int(myLikeId.Int32)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get posts: %w", err)
	}
	return posts, nil
}

// GetPostById
// SELECT p.*, group_concat(c.name, ", "), (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = true) as Likes, (SELECT Count(*) FROM posts_likes pl WHERE pl.post_id = p.id and type = false) as Dislikes FROM posts p LEFT JOIN posts_categories pc ON p.id = pc.post_id LEFT JOIN categories c ON c.id = pc.category_id WHERE p.id = 7;
func (r *postSQL) GetPostById(userId int) (models.Post, error) {
	var post models.Post
	query := fmt.Sprintf(`
	SELECT p.*, u.username as user_name, group_concat(c.name, ", ") as categories, 
	(SELECT Count(*) FROM %s  pl WHERE pl.post_id = p.id and type = true) as likes,
	(SELECT Count(*) FROM %s  pl WHERE pl.post_id = p.id and type = false) as dislikes 
	FROM %s p 
	LEFT JOIN %s pc ON p.id = pc.post_id 
	LEFT JOIN %s c ON c.id = pc.category_id 
	LEFT JOIN %s u ON u.id = p.user_id  
	WHERE p.id = ?`, postsLikesTable, postsLikesTable, postsTable, categoriesToPostsTable, categoriesTable, usersTable)

	err := r.db.QueryRow(query, userId).Scan(
		&post.Id,
		&post.AuthorID,
		&post.Date,
		&post.Title,
		&post.Content,
		&post.ImageName,
		&post.AuthorName,
		&post.Categories,
		&post.Likes,
		&post.Dislikes)
	if err == sql.ErrNoRows {
		return post, err
	}
	return post, nil
}

func (r *postSQL) SetPost(post models.Post, currentUserId int) (models.Post, error) {
	query := fmt.Sprintf(`UPDATE %s SET title=?, content=?, image_name=? where id=?`, postsTable)
	if _, err := r.db.Exec(query, post.Title, post.Content, post.ImageName, post.Id); err != nil {
		return post, fmt.Errorf("can't set post: %w", err)
	}
	err := InsertPostActivity(r, currentUserId, post.Id, 0, post.AuthorID, "edit")
	if err != nil {
		return post, err
	}
	return post, nil
}
func (r *postSQL) DeletePostById(postId int) error {
	query := fmt.Sprintf(`DELETE FROM %s where id=?`, postsTable)
	if _, err := r.db.Exec(query, postId); err != nil {
		return fmt.Errorf("can't set post: %w", err)
	}
	return nil
}

func (r *postSQL) DeleteCategoriesToPost(postId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE post_id=?`, categoriesToPostsTable)
	if _, err := r.db.Exec(query, postId); err != nil {
		return fmt.Errorf("can't set post: %w", err)
	}
	return nil
}

func InsertPostActivity(r *postSQL, userId, postId, commentId, authorId int, action string) error {
	now := time.Now()
	query := fmt.Sprintf(`INSERT INTO %s (user_id, post_id, comment_id, author_id, action, is_shown, created_at) values (?,?,?,?,?,?,?) RETURNING id`, activityTable)
	if _, err := r.db.Exec(query, userId, postId, commentId, authorId, action, false, now); err != nil {
		return fmt.Errorf("can't set activity: %w", err)
	}
	return nil
}
