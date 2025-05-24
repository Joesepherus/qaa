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
	"qaa/controllers/authController"
	"qaa/controllers/questionsController"
	"qaa/controllers/trainingsController"
	"qaa/middlewares/authMiddleware"
	"qaa/services/answersService"
	"qaa/services/questionsService"
	"qaa/services/trainingsService"
	"qaa/services/usersService"
	"qaa/templates"
	"qaa/types/questionsTypes"
	"qaa/types/userTypes"
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

func isLoggedIn(user *userTypes.User, err error, w http.ResponseWriter, r *http.Request) bool {
	if user == nil {
		return false
	}
	return true
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

	email, ok := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := usersService.GetUserByEmail(email)
	loggedIn := isLoggedIn(user, err, w, r)

	if !loggedIn {
		switch r.URL.Path {
		case "/":
			templateLocation = templates.BaseLocation + "/index.html"
			pageTitle = "Trading Alerts"
		case "/health":
			healthHandler(w)
			return
		case "/error":
			templateLocation = templates.BaseLocation + "/error.html"
			pageTitle = "Error"
			message := r.URL.Query().Get("message")
			data["Message"] = message
		default:
			templateLocation = templates.BaseLocation + "/404.html"
			pageTitle = "Page not found"
		}
	} else {
		if ok {
			data["Email"] = email
		}
		switch r.URL.Path {
		case "/":
			templateLocation = templates.BaseLocation + "/index.html"
			pageTitle = "Trading Alerts"
		case "/random":

			question, err := questionsService.GetRandomQuestion(user.ID)
			if err == nil {
				data["Question"] = question
			}

			trainings, err := trainingsService.GetTrainings(user.ID)
			if err == nil {
				data["Trainings"] = trainings
			}

			templateLocation = templates.BaseLocation + "/random.html"
			pageTitle = "Random Question"
		case "/questions":
			questions, err := questionsService.GetQuestions(user.ID)
			if err == nil {
				data["Questions"] = questions
			}

			trainings, err := trainingsService.GetTrainings(user.ID)
			if err == nil {
				data["Trainings"] = trainings
			}

			templateLocation = templates.BaseLocation + "/questions.html"
			pageTitle = "Questions"
		case "/trainings":

			trainings, err := trainingsService.GetTrainings(user.ID)
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

			answer, err := answersService.GetAnswerById(user.ID, answerId)
			if err == nil {
				data["Answer"] = answer
			}

			question, err := questionsService.GetQuestionById(user.ID, answer.QuestionID)
			if err == nil {
				data["Question"] = question
			}

			templateLocation = templates.BaseLocation + "/feedback.html"
			pageTitle = "Feedback"
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

				answer, err := answersService.GetAnswerById(user.ID, answerId)
				if err == nil {
					data["Answer"] = answer
				}

				if answer.Feedback != nil {
					http.Redirect(w, r, "/random", http.StatusSeeOther)
				}

				question, err := questionsService.GetQuestionById(user.ID, answer.QuestionID)
				if err == nil {
					data["Question"] = question
					data["CorrectAnswer"] = template.HTML(question.CorrectAnswer)
				}

				templateLocation = templates.BaseLocation + "/feedback.html"
				pageTitle = "Feedback"

			}else if strings.HasPrefix(r.URL.Path, "/random/") {
				trainings, err := trainingsService.GetTrainings(user.ID)
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

				training, err := trainingsService.GetTrainingById(user.ID, trainingId)
				var question questionsTypes.Question

				if training.ID != 0 {
					question, err = questionsService.GetRandomQuestionWithTraining(user.ID, trainingId)
					data["TrainingId"] = trainingId
				} else {
					question, err = questionsService.GetRandomQuestion(user.ID)
				}

				if err == nil {
					data["Question"] = question
				}

				templateLocation = templates.BaseLocation + "/random.html"
				pageTitle = "Random Question"
			} else if strings.HasPrefix(r.URL.Path, "/questions/") {
				questions, err := questionsService.GetQuestions(user.ID)
				if err == nil {
					data["Questions"] = questions
				}

				trainings, err := trainingsService.GetTrainings(user.ID)
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

				selectedQuestion, err := questionsService.GetQuestionById(user.ID, questionId)
				if err != nil {
					http.Redirect(w, r, "/questions", http.StatusSeeOther)
				}
				data["SelectedQuestion"] = selectedQuestion

				templateLocation = templates.BaseLocation + "/questions.html"
				pageTitle = "Questions"
			} else if strings.HasPrefix(r.URL.Path, "/trainings/") {
				trainings, err := trainingsService.GetTrainings(user.ID)
				if err == nil {
					data["Trainings"] = trainings
				}

				pathParts := strings.Split(r.URL.Path, "/")
				if len(pathParts) < 3 || pathParts[1] != "trainings" {
					http.Error(w, "Invalid URL format", http.StatusBadRequest)
					return
				}

				trainingIdStr := pathParts[2]
				trainingId, err := strconv.Atoi(trainingIdStr)

				selectedTraining, err := trainingsService.GetTrainingById(user.ID, trainingId)
				if err != nil {
					http.Redirect(w, r, "/trainings", http.StatusSeeOther)
				}
				data["SelectedTraining"] = selectedTraining

				templateLocation = templates.BaseLocation + "/trainings.html"
				pageTitle = "Trainings"
			} else {
				templateLocation = templates.BaseLocation + "/404.html"
				pageTitle = "Page not found"
			}

		}
	}
	data["Title"] = pageTitle
	data["Content"] = templateLocation

	templates.RenderTemplate(w, r, templateLocation, data)
}

func RestApi() {
	port := 8092
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

	http.Handle("/api/questions/random", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(questionsController.GetRandomQuestion)))
	http.Handle("/api/questions/save-question", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(questionsController.SaveQuestion)))
	http.Handle("/api/questions/edit-question", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(questionsController.EditQuestion)))
	http.Handle("/api/questions/delete-question", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(questionsController.DeleteQuestion)))
	http.Handle("/api/trainings/save-training", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(trainingsController.SaveTraining)))
	http.Handle("/api/training/edit-training", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(trainingsController.EditTraining)))
	http.Handle("/api/training/delete-training", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(trainingsController.DeleteTraining)))
	http.Handle("/api/answers/save-answer", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(answersController.SaveAnswer)))
	http.Handle("/api/answers/feedback", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(answersController.UpdateFeedbackOnAnswer)))

	// Authentication routes
	http.Handle("/api/sign-up", authMiddleware.TokenCheckMiddleware(http.HandlerFunc(authController.SignUp)))
	http.Handle("/api/login", authMiddleware.TokenCheckMiddleware(http.HandlerFunc(authController.Login)))
	http.Handle("/api/logout", authMiddleware.TokenCheckMiddleware(http.HandlerFunc(authController.Logout)))
	http.Handle("/api/reset-password", authMiddleware.TokenCheckMiddleware(http.HandlerFunc(authController.ResetPassword)))
	http.Handle("/api/set-password", authMiddleware.TokenCheckMiddleware(http.HandlerFunc(authController.SetPassword)))
	http.Handle("/api/delete-account", authMiddleware.TokenCheckMiddleware(http.HandlerFunc(authController.DeleteAccount)))

	http.Handle("/", authMiddleware.TokenCheckMiddleware(http.HandlerFunc(PageHandler)))

	// Serve static files (CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Starting server on :%d...\n", port)
	log.Fatal(server.ListenAndServe())
}
