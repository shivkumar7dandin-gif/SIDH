package main

import (
	"log"

	"github.com/gin-gonic/gin"

	authHandler "github.com/shivkumar7dandin-gif/students-api/internal/auth/handler"
	authMiddleware "github.com/shivkumar7dandin-gif/students-api/internal/auth/middleware"
	authService "github.com/shivkumar7dandin-gif/students-api/internal/auth/service"

	assessmentHandler "github.com/shivkumar7dandin-gif/students-api/internal/assessment/handler"
	assessmentRepository "github.com/shivkumar7dandin-gif/students-api/internal/assessment/repository"
	assessmentService "github.com/shivkumar7dandin-gif/students-api/internal/assessment/service"

	attendanceHandler "github.com/shivkumar7dandin-gif/students-api/internal/attendance/handler"
	attendanceRepository "github.com/shivkumar7dandin-gif/students-api/internal/attendance/repository"
	attendanceService "github.com/shivkumar7dandin-gif/students-api/internal/attendance/service"

	classroomHandler "github.com/shivkumar7dandin-gif/students-api/internal/classroom/handler"
	classroomRepository "github.com/shivkumar7dandin-gif/students-api/internal/classroom/repository"
	classroomService "github.com/shivkumar7dandin-gif/students-api/internal/classroom/service"

	collegeHandler "github.com/shivkumar7dandin-gif/students-api/internal/college/handler"
	collegeRepository "github.com/shivkumar7dandin-gif/students-api/internal/college/repository"
	collegeService "github.com/shivkumar7dandin-gif/students-api/internal/college/service"

	"github.com/shivkumar7dandin-gif/students-api/internal/config"
	"github.com/shivkumar7dandin-gif/students-api/internal/database"

	studentHandler "github.com/shivkumar7dandin-gif/students-api/internal/student/handler"
	studentRepository "github.com/shivkumar7dandin-gif/students-api/internal/student/repository"
	studentService "github.com/shivkumar7dandin-gif/students-api/internal/student/service"
)

func main() {

	// =========================
	// CONFIG
	// =========================

	cfg := config.MustLoad()

	// =========================
	// DATABASE
	// =========================

	mongoClient, err := database.Connect(
		cfg.Storage.MongoURI,
	)
	if err != nil {
		log.Fatal("MongoDB connection failed: ", err)
	}

	log.Println("MongoDB connected successfully")

	db := mongoClient.Database(
		cfg.Storage.Database,
	)

	// =========================
	// COLLEGE
	// =========================

	collegeRepo := collegeRepository.NewCollegeRepository(db)

	collegeSvc := collegeService.NewCollegeService(
		collegeRepo,
	)

	collegeH := collegeHandler.NewCollegeHandler(
		collegeSvc,
	)

	// =========================
	// AUTH
	// =========================

	jwtSecret := "my-super-secret-jwt-key"

	authSvc := authService.NewAuthService(
		collegeRepo,
		jwtSecret,
	)

	authH := authHandler.NewAuthHandler(
		authSvc,
	)

	// =========================
	// CLASSROOM
	// =========================

	classroomRepo := classroomRepository.NewClassroomRepository(db)

	classroomSvc := classroomService.NewClassroomService(
		classroomRepo,
	)

	classroomH := classroomHandler.NewClassroomHandler(
		classroomSvc,
	)

	// =========================
	// STUDENT
	// =========================

	studentRepo := studentRepository.NewStudentRepository(db)

	studentSvc := studentService.NewStudentService(
		studentRepo,
		classroomRepo,
	)

	studentH := studentHandler.NewStudentHandler(
		studentSvc,
	)

	// =========================
	// ATTENDANCE
	// =========================

	attendanceRepo := attendanceRepository.NewAttendanceRepository(db)

	attendanceSvc := attendanceService.NewAttendanceService(
		attendanceRepo,
	)

	attendanceH := attendanceHandler.NewAttendanceHandler(
		attendanceSvc,
	)

	// =========================
	// ASSESSMENT
	// =========================

	assessmentRepo := assessmentRepository.NewAssessmentRepository(db)

	assessmentSvc := assessmentService.NewAssessmentService(
		assessmentRepo,
	)

	assessmentH := assessmentHandler.NewAssessmentHandler(
		assessmentSvc,
	)

	// =========================
	// ROUTER
	// =========================

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	api := router.Group("/api/v1")

	// =========================
	// AUTH ROUTES
	// =========================

	auth := api.Group("/auth")
	{
		auth.POST("/login", authH.Login)
	}

	protected := api.Group("")
	protected.Use(
		authMiddleware.AuthMiddleware(jwtSecret),
	)

	// =========================
	// COLLEGE ROUTES
	// =========================

	colleges := api.Group("/colleges")
	{
		colleges.POST("", collegeH.Create)
	}

	protectedColleges := protected.Group("/colleges")
	{
		protectedColleges.GET("", collegeH.GetAll)
		protectedColleges.GET("/:id", collegeH.GetByID)
	}

	// =========================
	// CLASSROOM ROUTES
	// =========================

	//classrooms := api.Group("/classrooms")
	classrooms := protected.Group("/classrooms")

	{
		classrooms.POST("", classroomH.Create)
		classrooms.GET("", classroomH.GetAll)
		classrooms.GET("/:id", classroomH.GetByID)
		classrooms.PUT("/:id", classroomH.Update)
		classrooms.DELETE("/:id", classroomH.Delete)
	}

	// =========================
	// STUDENT ROUTES
	// =========================

	//students := api.Group("/students")
	students := protected.Group("/students")
	{
		students.POST("", studentH.Create)
		students.GET("", studentH.GetAll)
		students.GET("/:id", studentH.GetByID)
		students.PUT("/:id", studentH.Update)
		students.DELETE("/:id", studentH.Delete)
	}

	// =========================
	// ATTENDANCE ROUTES
	// =========================

	//attendance := api.Group("/attendance")
	attendance := protected.Group("/attendance")
	{
		attendance.POST("", attendanceH.Create)
		attendance.GET("", attendanceH.GetAll)

		attendance.GET(
			"/student/:studentId",
			attendanceH.GetByStudent,
		)
	}

	// =========================
	// ASSESSMENT ROUTES
	// =========================

	//assessments := api.Group("/assessments")
	assessments := protected.Group("/assessments")
	{
		assessments.POST("", assessmentH.Create)
		assessments.GET("", assessmentH.GetAll)

		assessments.GET(
			"/student/:studentId",
			assessmentH.GetByStudent,
		)
	}

	// =========================
	// START SERVER
	// =========================

	log.Println(
		"Students API running on",
		cfg.HTTPServer.Address,
	)

	if err := router.Run(
		cfg.HTTPServer.Address,
	); err != nil {
		log.Fatal(err)
	}
}
