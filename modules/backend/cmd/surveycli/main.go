package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/mstrYoda/maasanketi.co/domain/survey"
	log "github.com/mstrYoda/maasanketi.co/pkg/logger"
	"github.com/mstrYoda/maasanketi.co/repository"
)

func main() {
	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle shutdown signals
	go func() {
		sig := <-sigChan
		log.Logger().Info(fmt.Sprintf("Received signal %s, shutting down gracefully...", sig))
		cancel()
		// Give ongoing operations a chance to complete before exiting
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	// Define command line flags
	var (
		jsonFile    = flag.String("file", "", "Path to JSON file containing survey data")
		useDefault  = flag.Bool("default", false, "Use default survey template")
		title       = flag.String("title", "", "Survey title (used with --default)")
		description = flag.String("description", "", "Survey description (used with --default)")
		minTime     = flag.Int("min-time", 3, "Minimum completion time in minutes (used with --default)")
		slug        = flag.String("slug", "", "Custom slug for the survey (optional)")
	)

	flag.Parse()

	// Validate flags
	if !*useDefault && *jsonFile == "" {
		fmt.Println("Error: Either --file or --default flag must be provided")
		flag.Usage()
		os.Exit(1)
	}

	if *useDefault && (*title == "" || *description == "") {
		fmt.Println("Error: When using --default, both --title and --description must be provided")
		flag.Usage()
		os.Exit(1)
	}

	// Check if DB_CONN_STR environment variable is set
	if os.Getenv("DB_CONN_STR") == "" {
		fmt.Println("Error: DB_CONN_STR environment variable is not set")
		fmt.Println("Example: export DB_CONN_STR=\"postgres://username:password@localhost:5432/database_name\"")
		os.Exit(1)
	}

	// Initialize repository
	repo, err := repository.New()
	if err != nil {
		log.Logger().Fatal(fmt.Sprintf("failed to create repository: %s", err.Error()))
	}
	defer repo.Close()

	// Create survey
	var s *survey.Survey

	if *useDefault {
		// Use default survey template
		s = DefaultSurvey()
		s.Title = *title
		s.Description = *description
		s.MinCompletionTimeMin = *minTime
	} else {
		// Load survey from JSON file
		jsonData, err := os.ReadFile(*jsonFile)
		if err != nil {
			log.Logger().Fatal(fmt.Sprintf("failed to read JSON file: %s", err.Error()))
		}

		s = &survey.Survey{}
		if err := json.Unmarshal(jsonData, s); err != nil {
			log.Logger().Fatal(fmt.Sprintf("failed to parse JSON: %s", err.Error()))
		}

		// Validate required fields
		if s.Title == "" || s.Description == "" {
			log.Logger().Fatal("survey title and description are required")
		}

		// Set default minimum completion time if not provided
		if s.MinCompletionTimeMin <= 0 {
			s.MinCompletionTimeMin = 3
		}

		// Validate questions
		if len(s.Questions) == 0 {
			log.Logger().Fatal("survey must have at least one question")
		}

		// Validate question types
		for i, q := range s.Questions {
			switch q.Type {
			case survey.QuestionTypeText, survey.QuestionTypeNumber, survey.QuestionTypeSelect, survey.QuestionTypeMultiSelect, survey.QuestionTypeBoolean:
				// Valid type
			default:
				log.Logger().Fatal(fmt.Sprintf("invalid question type for question %d: %s", i+1, q.Type))
			}

			// Validate options for select and multi-select questions
			if (q.Type == survey.QuestionTypeSelect || q.Type == survey.QuestionTypeMultiSelect) && len(q.Options) == 0 {
				log.Logger().Fatal(fmt.Sprintf("question %d (%s) requires options", i+1, q.ID))
			}
		}
	}

	// Generate UUID and set timestamps
	s.ID = uuid.New().String()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	// Set or generate slug
	if *slug != "" {
		s.Slug = *slug
	} else {
		s.GenerateSlug()
	}

	// Create a context with 30-second timeout for database operations
	dbCtx, dbCancel := context.WithTimeout(ctx, 30*time.Second)
	defer dbCancel()

	log.Logger().Info("Saving survey to database...")

	// Save survey to database with timeout context
	if err := repo.Survey.Create(dbCtx, s); err != nil {
		if dbCtx.Err() == context.DeadlineExceeded {
			log.Logger().Fatal("database operation timed out after 30 seconds")
		}
		log.Logger().Fatal(fmt.Sprintf("failed to create survey: %s", err.Error()))
	}

	fmt.Printf("Survey created successfully with ID: %s\n", s.ID)
}
