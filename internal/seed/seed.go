package seed

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/Reensef/golang-social/internal/store"
)

func GenerateFakeData(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		err := store.Users.Create(ctx, user)
		if err != nil {
			log.Printf("Error seeding users: %v, %v", err, user)
		}
	}

	posts := generatePosts(500, users)
	for _, post := range posts {
		err := store.Posts.Create(ctx, post)
		if err != nil {
			log.Printf("Error seeding posts: %v, %v", err, post)
		}
	}

	comments := generateComments(1000, users, posts)
	for _, comment := range comments {
		err := store.Comments.Create(ctx, comment)
		if err != nil {
			log.Printf("Error seeding comments: %v, %v", err, comment)
		}
	}
}

func generateUsers(count int) []*store.User {
	users := make([]*store.User, count)
	for i := 0; i < count; i++ {
		username := fakeNames[rand.Intn(len(fakeNames))] + fmt.Sprintf("%d", rand.Intn(1000))
		users[i] = &store.User{
			Username: username,
			Email:    username + "@mail.com",
		}
	}
	return users
}

func generatePosts(count int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, count)
	for i := 0; i < count; i++ {
		posts[i] = &store.Post{
			Title:   fakePostsNames[rand.Intn(len(fakePostsNames))],
			Content: fakePostsContent[rand.Intn(len(fakePostsContent))],
			UserID:  users[rand.Intn(len(users))].ID,
		}
	}
	return posts
}

func generateComments(count int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, count)
	for i := 0; i < count; i++ {
		comments[i] = &store.Comment{
			Content: fakePostsComments[rand.Intn(len(fakePostsContent))],
			PostID:  posts[rand.Intn(len(posts))].ID,
			Sender: store.CommentSender{
				ID: users[rand.Intn(len(users))].ID},
		}
	}
	return comments
}
