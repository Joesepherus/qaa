package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"os"
	"qaa/controllers/answersController"
	"qaa/controllers/questionsController"
	"qaa/controllers/trainingsController"
	"qaa/services/answersService"
	"qaa/services/questionsService"
	"qaa/services/trainingsService"
	"qaa/templates"
	"qaa/types/questionsTypes"
	"strconv"
)

// Health check handler
func healthHandler(w http.ResponseWriter) {
	// Basic health check response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// This handler will only be called if the token is valid
	fmt.Fprintf(w, "Welcome to the protected area!")
}

func PageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET") // Specifies the HTTP methods allowed.
	w.Header().Set("X-Frame-Options", "DENY")                   // Prevents clickjacking
	w.Header().Set("X-Content-Type-Options", "nosniff")         // Prevents MIME sniffing
	w.Header().Set("X-XSS-Protection", "1; mode=block")         // Protects against XSS attacks
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	var templateLocation, pageTitle string

	data := map[string]interface{}{
		"Email": "",
		// Add other default data here if needed
	}

	switch r.URL.Path {
	case "/":
		templateLocation = templates.BaseLocation + "/index.html"
		pageTitle = "Trading Alerts"
	case "/random":
		question, err := questionsService.GetRandomQuestion()
		if err == nil {
			data["Question"] = question
		}

		trainings, err := trainingsService.GetTrainings()
		if err == nil {
			data["Trainings"] = trainings
		}

		templateLocation = templates.BaseLocation + "/random.html"
		pageTitle = "Random Question"
	case "/questions":
		questions, err := questionsService.GetQuestions()
		if err == nil {
			data["Questions"] = questions
		}

		trainings, err := trainingsService.GetTrainings()
		if err == nil {
			data["Trainings"] = trainings
		}

		templateLocation = templates.BaseLocation + "/questions.html"
		pageTitle = "Questions"
	case "/trainings":
		trainings, err := trainingsService.GetTrainings()
		if err == nil {
			data["Trainings"] = trainings
		}

		templateLocation = templates.BaseLocation + "/trainings.html"
		pageTitle = "Trainings"
	case "/training-saved":
		templateLocation = templates.BaseLocation + "/training-saved.html"
		pageTitle = "Training Saved"
	case "/feedback":
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 || pathParts[1] != "feedback" {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		answerIdStr := pathParts[2]
		answerId, err := strconv.Atoi(answerIdStr)

		answer, err := answersService.GetAnswerById(answerId)
		if err == nil {
			data["Answer"] = answer
		}

		question, err := questionsService.GetQuestionById(answer.QuestionID)
		if err == nil {
			data["Question"] = question
		}

		templateLocation = templates.BaseLocation + "/feedback.html"
		pageTitle = "Feedback"
	case "/health":
		healthHandler(w)
		return
	case "/error":
		templateLocation = templates.BaseLocation + "/error.html"
		pageTitle = "Error - Trading Alerts"
		message := r.URL.Query().Get("message")
		data["Message"] = message
	default:
		// Handle /feedback/{number}
		if strings.HasPrefix(r.URL.Path, "/feedback/") {

			pathParts := strings.Split(r.URL.Path, "/")
			if len(pathParts) < 3 || pathParts[1] != "feedback" {
				http.Error(w, "Invalid URL format", http.StatusBadRequest)
				return
			}

			answerIdStr := pathParts[2]
			answerId, err := strconv.Atoi(answerIdStr)

			answer, err := answersService.GetAnswerById(answerId)
			if err == nil {
				data["Answer"] = answer
			}

			if answer.Feedback != nil {
				http.Redirect(w, r, "/random", http.StatusSeeOther)
			}

			question, err := questionsService.GetQuestionById(answer.QuestionID)
			if err == nil {
				data["Question"] = question
				data["CorrectAnswer"] = template.HTML(question.CorrectAnswer)
			}

			templateLocation = templates.BaseLocation + "/feedback.html"
			pageTitle = "Feedback"

		}
		if strings.HasPrefix(r.URL.Path, "/random/") {
			trainings, err := trainingsService.GetTrainings()
			if err == nil {
				data["Trainings"] = trainings
			}

			pathParts := strings.Split(r.URL.Path, "/")
			if len(pathParts) < 3 || pathParts[1] != "random" {
				http.Error(w, "Invalid URL format", http.StatusBadRequest)
				return
			}

			trainingIdStr := pathParts[2]
			trainingId, err := strconv.Atoi(trainingIdStr)

			training, err := trainingsService.GetTrainingById(trainingId)
			var question questionsTypes.Question

			if training.ID != 0 {
				question, err = questionsService.GetRandomQuestionWithTraining(trainingId)
				data["TrainingId"] = trainingId
			} else {
				question, err = questionsService.GetRandomQuestion()
			}

			if err == nil {
				data["Question"] = question
			}

			templateLocation = templates.BaseLocation + "/random.html"
			pageTitle = "Random Question"
		}
		if strings.HasPrefix(r.URL.Path, "/questions/") {
			questions, err := questionsService.GetQuestions()
			if err == nil {
				data["Questions"] = questions
			}

			trainings, err := trainingsService.GetTrainings()
			if err == nil {
				data["Trainings"] = trainings
			}

	        pathParts := strings.Split(r.URL.Path, "/")
			if len(pathParts) < 3 || pathParts[1] != "questions" {
				http.Error(w, "Invalid URL format", http.StatusBadRequest)
				return
			}

			questionIdStr := pathParts[2]
			questionId, err := strconv.Atoi(questionIdStr)

			selectedQuestion, err := questionsService.GetQuestionById(questionId)
            println("selectedQuestion", selectedQuestion.TrainingID)
			if err == nil {
				data["SelectedQuestion"] = selectedQuestion
			}

			templateLocation = templates.BaseLocation + "/questions.html"
			pageTitle = "Questions"
		} else {
			templateLocation = templates.BaseLocation + "/404.html"
			pageTitle = "Page not found"
		}
	}

	data["Title"] = pageTitle
	data["Content"] = templateLocation

	templates.RenderTemplate(w, r, templateLocation, data)
}

func RestApi() {
	port := 8090
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	// Define the server with timeouts
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(port), // Listen on the specified port
		Handler:      nil,
		ReadTimeout:  5 * time.Second,  // Max time to read the request
		WriteTimeout: 10 * time.Second, // Max time to write the response
		IdleTimeout:  15 * time.Second, // Max time for idle connections
	}

	http.HandleFunc("/api/questions/random", questionsController.GetRandomQuestion)
	http.HandleFunc("/api/answers/save-answer", answersController.SaveAnswer)
	http.HandleFunc("/api/answers/feedback", answersController.UpdateFeedbackOnAnswer)
	http.HandleFunc("/api/questions/save-question", questionsController.SaveQuestion)
	http.HandleFunc("/api/questions/edit-question", questionsController.EditQuestion)
	http.HandleFunc("/api/questions/delete-question", questionsController.DeleteQuestion)
	http.HandleFunc("/api/trainings/save-training", trainingsController.SaveTraining)
	http.Handle("/", http.HandlerFunc(PageHandler))

	// Serve static files (CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Starting server on :%d...\n", port)
	log.Fatal(server.ListenAndServe())
}
