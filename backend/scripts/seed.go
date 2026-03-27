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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	start := time.Now()

	ctx := context.Background()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

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

	users := seedUsers(ctx, userRepo, queries, r)
	posts := seedPostsBatch(ctx, pool, users, r)

	// replies now returned so we can like them too
	replies := seedRepliesBatch(ctx, pool, users, posts, r)

	// merge posts + replies for likes
	allPostIDs := append(posts, replies...)

	seedLikes(ctx, postRepo, users, allPostIDs, r)

	fmt.Println("✅ Seeding complete!")
	fmt.Println("⏱ Total time:", time.Since(start))
}

func randomTime(r *rand.Rand) time.Time {
	return time.Now().Add(-time.Duration(r.Intn(720)) * time.Hour)
}

func seedUsers(ctx context.Context, repo *postgres.UserRepository, q *db.Queries, r *rand.Rand) []string {
	start := time.Now()

	var ids []string

	adjectives := []string{"cool", "fast", "lazy", "happy", "wild", "silent", "noisy", "smart", "crazy", "chill"}
	nouns := []string{"dev", "coder", "builder", "hacker", "gopher", "engineer", "creator", "thinker", "designer"}

	bios := []string{
		"backend dev", "golang enjoyer", "building cool stuff 🚀",
		"coffee + code", "clean architecture fan", "debugging life",
		"ship fast mindset", "open source lover", "learning everyday", "just vibes ✨",
	}

	hash, _ := util.HashPassword("password123")

	for i := 0; i < 30; i++ {
		id := uuid.New().String()

		username := fmt.Sprintf("%s_%s_%d",
			adjectives[r.Intn(len(adjectives))],
			nouns[r.Intn(len(nouns))],
			r.Intn(100000),
		)

		email := fmt.Sprintf("%s_%s@seed.dev", username, id[:8])

		if err := repo.Create(ctx, id, username, email, hash); err != nil {
			log.Println("user error:", err)
			continue
		}

		_ = q.ActivateUser(ctx, id)

		bio := bios[r.Intn(len(bios))]
		avatar := fmt.Sprintf("https://i.pravatar.cc/150?u=%s", id)

		_ = repo.UpdateProfile(ctx, id, bio, avatar)

		ids = append(ids, id)
	}

	fmt.Println("👤 users:", len(ids), "⏱", time.Since(start))
	return ids
}

var samplePosts = []string{
	"hello world 👋", "building twitter clone 🚀", "golang is clean af",
	"debugging pain 😭", "late night coding", "sqlc is goated",
	"clean architecture >>", "why is this not working", "ship fast", "just deployed 🔥",
}

func seedPostsBatch(ctx context.Context, pool *pgxpool.Pool, users []string, r *rand.Rand) []string {
	start := time.Now()

	var ids []string
	batch := &pgx.Batch{}

	for i := 0; i < 120; i++ {
		id := uuid.New().String()
		userID := users[r.Intn(len(users))]
		content := samplePosts[r.Intn(len(samplePosts))]
		createdAt := randomTime(r)

		batch.Queue(
			`INSERT INTO posts (id, user_id, content, created_at)
			 VALUES ($1, $2, $3, $4)`,
			id, userID, content, createdAt,
		)

		ids = append(ids, id)
	}

	br := pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(ids); i++ {
		_, err := br.Exec()
		if err != nil {
			log.Println("batch post error:", err)
		}
	}

	fmt.Println("📝 posts:", len(ids), "⏱", time.Since(start))
	return ids
}

var sampleReplies = []string{
	"totally agree 💯", "this is interesting", "can you explain more?",
	"nah I disagree", "this helped a lot 🙏", "lol same",
	"this is underrated", "good point!", "what do you mean?", "🔥🔥🔥",
}

// deeper recursion up to 10 levels
func generateRepliesRecursive(users []string, rootID string, parentID string, depth int, r *rand.Rand, out *[]post.Post) {
	if depth > 10 {
		return
	}

	var numReplies int

	if depth == 0 {
		numReplies = r.Intn(11) + 10
	} else if depth <= 2 {
		numReplies = r.Intn(5) + 2
	} else {
		numReplies = r.Intn(2)
	}

	for i := 0; i < numReplies; i++ {
		replyID := uuid.New().String()

		p := post.Post{
			ID:           replyID,
			UserID:       users[r.Intn(len(users))],
			Content:      sampleReplies[r.Intn(len(sampleReplies))],
			ParentPostID: &parentID,
			RootPostID:   &rootID,
		}

		*out = append(*out, p)

		// deeper threads but controlled
		var continueChance int
		if depth == 0 {
			continueChance = 70
		} else if depth <= 2 {
			continueChance = 50
		} else {
			continueChance = 30
		}

		if r.Intn(100) < continueChance {
			generateRepliesRecursive(users, rootID, replyID, depth+1, r, out)
		}
	}
}

// now returns reply IDs
func seedRepliesBatch(ctx context.Context, pool *pgxpool.Pool, users, posts []string, r *rand.Rand) []string {
	start := time.Now()

	var allReplies []post.Post
	var replyIDs []string

	for _, postID := range posts {
		if r.Intn(100) > 30 {
			generateRepliesRecursive(users, postID, postID, 0, r, &allReplies)
		}
	}

	batch := &pgx.Batch{}

	for _, p := range allReplies {
		batch.Queue(
			`INSERT INTO posts (id, user_id, content, parent_post_id, root_post_id)
			 VALUES ($1, $2, $3, $4, $5)`,
			p.ID, p.UserID, p.Content, *p.ParentPostID, *p.RootPostID,
		)

		replyIDs = append(replyIDs, p.ID)
	}

	br := pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(allReplies); i++ {
		_, err := br.Exec()
		if err != nil {
			log.Println("batch reply error:", err)
		}
	}

	fmt.Println("💬 replies:", len(allReplies), "⏱", time.Since(start))
	return replyIDs
}

// now likes apply to posts + replies
func seedLikes(ctx context.Context, repo *postgres.PostRepository, users, postIDs []string, r *rand.Rand) {
	start := time.Now()
	total := 0

	for _, userID := range users {
		n := r.Intn(15) + 5

		for i := 0; i < n; i++ {
			postID := postIDs[r.Intn(len(postIDs))]

			if err := repo.CreateLike(ctx, userID, postID); err == nil {
				total++
			}
		}
	}

	fmt.Println("❤️ likes:", total, "⏱", time.Since(start))
}
