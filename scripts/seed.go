package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/akshzyx/gorum/internal/config"
	db "github.com/akshzyx/gorum/internal/db/sqlc/generated"
	"github.com/akshzyx/gorum/internal/domain/post"
	"github.com/akshzyx/gorum/internal/repository/postgres"
	"github.com/akshzyx/gorum/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	cfg := config.LoadConfig()

	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := db.New(pool)

	userRepo := postgres.NewUserRepository(queries)
	postRepo := postgres.NewPostRepository(queries)

	fmt.Println("🌱 Seeding started...")

	users := seedUsers(ctx, userRepo, queries)
	posts := seedPosts(ctx, postRepo, users)
	seedReplies(ctx, postRepo, users, posts) // 🔥 NEW
	seedLikes(ctx, postRepo, users, posts)

	fmt.Println("✅ Seeding complete!")
}

func seedUsers(ctx context.Context, repo *postgres.UserRepository, q *db.Queries) []string {
	var ids []string

	bios := []string{
		"backend dev",
		"golang enjoyer",
		"building cool stuff 🚀",
		"coffee + code",
		"clean architecture fan",
		"debugging life",
		"ship fast mindset",
		"open source lover",
		"learning everyday",
		"just vibes ✨",
	}

	for i := 0; i < 30; i++ {
		id := uuid.New().String()
		username := fmt.Sprintf("user_%d", i+1)
		email := fmt.Sprintf("user_%d@example.com", i+1)

		hash, _ := util.HashPassword("password123")

		err := repo.Create(ctx, id, username, email, hash)
		if err != nil {
			log.Println("user error:", err)
			continue
		}

		// mark as verified
		if err := q.ActivateUser(ctx, id); err != nil {
			log.Println("activate user error:", err)
		}

		// add bio + avatar
		bio := bios[rand.Intn(len(bios))]
		avatar := fmt.Sprintf("https://i.pravatar.cc/150?u=%s", id)

		err = repo.UpdateProfile(ctx, id, bio, avatar)
		if err != nil {
			log.Println("profile update error:", err)
		}

		ids = append(ids, id)
	}

	fmt.Println("👤 users created:", len(ids))
	return ids
}

var samplePosts = []string{
	"hello world 👋",
	"building twitter clone 🚀",
	"golang is clean af",
	"debugging pain 😭",
	"late night coding",
	"sqlc is goated",
	"clean architecture >>",
	"why is this not working",
	"ship fast",
	"just deployed 🔥",
}

func seedPosts(ctx context.Context, repo *postgres.PostRepository, users []string) []string {
	var ids []string

	for i := 0; i < 120; i++ {
		p := &post.Post{
			ID:      uuid.New().String(),
			UserID:  users[rand.Intn(len(users))],
			Content: samplePosts[rand.Intn(len(samplePosts))],
		}

		err := repo.Create(ctx, p)
		if err != nil {
			log.Println("post error:", err)
			continue
		}

		ids = append(ids, p.ID)
	}

	fmt.Println("📝 posts created:", len(ids))
	return ids
}

var sampleReplies = []string{
	"totally agree 💯",
	"this is interesting",
	"can you explain more?",
	"nah I disagree",
	"this helped a lot 🙏",
	"lol same",
	"this is underrated",
	"good point!",
	"what do you mean?",
	"🔥🔥🔥",
}

func seedReplies(ctx context.Context, repo *postgres.PostRepository, users, posts []string) {
	total := 0

	for _, rootPostID := range posts {
		// 70% of posts get replies
		if rand.Intn(100) > 70 {
			continue
		}

		numReplies := rand.Intn(5) + 1

		for i := 0; i < numReplies; i++ {
			replyID := uuid.New().String()
			userID := users[rand.Intn(len(users))]
			content := sampleReplies[rand.Intn(len(sampleReplies))]

			rootID := rootPostID
			parentID := rootPostID

			p := &post.Post{
				ID:           replyID,
				UserID:       userID,
				Content:      content,
				ParentPostID: &parentID,
				RootPostID:   &rootID,
			}

			err := repo.CreateReply(ctx, p)
			if err != nil {
				log.Println("reply error:", err)
				continue
			}

			total++

			// nested replies
			if rand.Intn(100) < 40 {
				nestedCount := rand.Intn(3) + 1

				for j := 0; j < nestedCount; j++ {
					nestedID := uuid.New().String()
					nestedUser := users[rand.Intn(len(users))]
					nestedContent := sampleReplies[rand.Intn(len(sampleReplies))]

					parent := replyID
					root := rootPostID

					nested := &post.Post{
						ID:           nestedID,
						UserID:       nestedUser,
						Content:      nestedContent,
						ParentPostID: &parent,
						RootPostID:   &root,
					}

					err := repo.CreateReply(ctx, nested)
					if err != nil {
						log.Println("nested reply error:", err)
						continue
					}

					total++
				}
			}
		}
	}

	fmt.Println("💬 replies created:", total)
}

func seedLikes(ctx context.Context, repo *postgres.PostRepository, users, posts []string) {
	total := 0

	for _, userID := range users {
		n := rand.Intn(10) + 5

		for i := 0; i < n; i++ {
			postID := posts[rand.Intn(len(posts))]

			err := repo.CreateLike(ctx, userID, postID)
			if err == nil {
				total++
			}
		}
	}

	fmt.Println("❤️ likes created:", total)
}
