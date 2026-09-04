package main

import (
	"log"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	assessmentHandler "github.com/shivkumar7dandin-gif/students-api/internal/assessment/handler"
	assessmentRepository "github.com/shivkumar7dandin-gif/students-api/internal/assessment/repository"
	assessmentService "github.com/shivkumar7dandin-gif/students-api/internal/assessment/service"

	attendanceHandler "github.com/shivkumar7dandin-gif/students-api/internal/attendance/handler"
	attendanceRepository "github.com/shivkumar7dandin-gif/students-api/internal/attendance/repository"
	attendanceService "github.com/shivkumar7dandin-gif/students-api/internal/attendance/service"

	authHandler "github.com/shivkumar7dandin-gif/students-api/internal/auth/handler"
	authMiddleware "github.com/shivkumar7dandin-gif/students-api/internal/auth/middleware"
	authService "github.com/shivkumar7dandin-gif/students-api/internal/auth/service"

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

	teacherHandler "github.com/shivkumar7dandin-gif/students-api/internal/teacher/handler"
	teacherRepository "github.com/shivkumar7dandin-gif/students-api/internal/teacher/repository"
	teacherService "github.com/shivkumar7dandin-gif/students-api/internal/teacher/service"

	userRepository "github.com/shivkumar7dandin-gif/students-api/internal/user/repository"
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
	// USER
	// =========================

	userRepo := userRepository.NewUserRepository(db)

	// =========================
	// COLLEGE
	// =========================

	collegeRepo := collegeRepository.NewCollegeRepository(db)

	teacherRepo := teacherRepository.NewTeacherRepository(db)

	classroomRepo := classroomRepository.NewClassroomRepository(db)

	studentRepo := studentRepository.NewStudentRepository(db)

	collegeSvc := collegeService.NewCollegeService(
		collegeRepo,
		teacherRepo,
		classroomRepo,
		studentRepo,
	)

	collegeH := collegeHandler.NewCollegeHandler(collegeSvc)
	// =========================
	// AUTH
	// =========================

	jwtSecret := "my-super-secret-jwt-key"

	authSvc := authService.NewAuthService(
		userRepo,
		jwtSecret,
	)

	authH := authHandler.NewAuthHandler(
		authSvc,
	)

	// =========================
	// TEACHER
	// =========================

	//teacherRepo := teacherRepository.NewTeacherRepository(db)

	teacherSvc := teacherService.NewTeacherService(
		teacherRepo,
		userRepo,
	)

	teacherH := teacherHandler.NewTeacherHandler(
		teacherSvc,
	)

	// =========================
	// CLASSROOM
	// =========================

	//classroomRepo := classroomRepository.NewClassroomRepository(db)

	classroomSvc := classroomService.NewClassroomService(
		classroomRepo,
	)

	classroomH := classroomHandler.NewClassroomHandler(
		classroomSvc,
	)

	// =========================
	// STUDENT
	// =========================

	//studentRepo := studentRepository.NewStudentRepository(db)

	studentSvc := studentService.NewStudentService(
		studentRepo,
		classroomRepo,
		userRepo,
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

	studentH := studentHandler.NewStudentHandler(
		studentSvc,
		attendanceSvc,
		assessmentSvc,
	)

	// =========================
	// ROUTER
	// =========================

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},

		ExposeHeaders: []string{
			"Content-Length",
		},

		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}))

	// =========================
	// HEALTH CHECK
	// =========================

	router.GET("/health", func(c *gin.Context) {
		c.JSON(
			200,
			gin.H{
				"status": "OK",
			},
		)
	})

	// =========================
	// API VERSION
	// =========================

	api := router.Group("/api/v1")

	// =========================
	// PUBLIC AUTH ROUTES
	// =========================

	auth := api.Group("/auth")
	{
		auth.POST(
			"/login",
			authH.Login,
		)
	}

	// =========================
	// PUBLIC COLLEGE ROUTES
	// =========================

	// Public registration for now.
	// Later this should be protected by super-admin.
	collegesPublic := api.Group("/colleges")
	{
		collegesPublic.POST(
			"",
			collegeH.Create,
		)
	}

	// =========================
	// JWT PROTECTED ROUTES
	// =========================

	protected := api.Group("")

	protected.Use(
		authMiddleware.AuthMiddleware(
			jwtSecret,
		),
	)

	// =========================
	// ROLE GROUPS
	// =========================

	// Only college admin
	adminOnly := protected.Group("")

	adminOnly.Use(
		authMiddleware.RequireRole(
			"college_admin",
		),
	)

	// College admin + teacher
	teacherAndAdmin := protected.Group("")

	teacherAndAdmin.Use(
		authMiddleware.RequireRole(
			"college_admin",
			"teacher",
		),
	)

	// College admin + teacher + student
	allUsers := protected.Group("")

	allUsers.Use(
		authMiddleware.RequireRole(
			"college_admin",
			"teacher",
			"student",
		),
	)

	// =========================
	// COLLEGE ROUTES
	// =========================

	adminOnly.GET(
		"/colleges/me/dashboard",
		collegeH.GetDashboard,
	)

	collegesRead := allUsers.Group("/colleges")
	{
		collegesRead.GET(
			"",
			collegeH.GetAll,
		)

		collegesRead.GET(
			"/:id",
			collegeH.GetByID,
		)
	}

	// =========================
	// TEACHER ROUTES
	// =========================

	// All authenticated roles can read teacher data.
	teachersRead := allUsers.Group("/teachers")
	{
		teachersRead.GET(
			"",
			teacherH.GetAll,
		)

		teachersRead.GET(
			"/:id",
			teacherH.GetByID,
		)
	}

	// Only college admin can manage teachers.
	teachersAdmin := adminOnly.Group("/teachers")
	{
		teachersAdmin.POST(
			"",
			teacherH.Create,
		)

		teachersAdmin.PUT(
			"/:id",
			teacherH.Update,
		)

		teachersAdmin.DELETE(
			"/:id",
			teacherH.Delete,
		)
	}

	teacherOnly := protected.Group("")
	teacherOnly.Use(
		authMiddleware.RequireRole("teacher"),
	)

	teacherOnly.GET(
		"/teachers/me",
		teacherH.GetMe,
	)
	// =========================
	// CLASSROOM ROUTES
	// =========================

	classroomsRead := allUsers.Group("/classrooms")
	{
		classroomsRead.GET(
			"",
			classroomH.GetAll,
		)

		classroomsRead.GET(
			"/:id",
			classroomH.GetByID,
		)
	}

	classroomsAdmin := adminOnly.Group("/classrooms")
	{
		classroomsAdmin.POST(
			"",
			classroomH.Create,
		)

		classroomsAdmin.PUT(
			"/:id",
			classroomH.Update,
		)

		classroomsAdmin.DELETE(
			"/:id",
			classroomH.Delete,
		)
	}

	// =========================
	// STUDENT ROUTES
	// =========================

	studentsRead := allUsers.Group("/students")
	{
		studentsRead.GET(
			"",
			studentH.GetAll,
		)

		studentsRead.GET(
			"/:id",
			studentH.GetByID,
		)
	}

	studentsAdmin := adminOnly.Group("/students")
	{
		studentsAdmin.POST(
			"",
			studentH.Create,
		)

		studentsAdmin.PUT(
			"/:id",
			studentH.Update,
		)

		studentsAdmin.DELETE(
			"/:id",
			studentH.Delete,
		)
	}
	studentOnly := protected.Group("")
	studentOnly.Use(
		authMiddleware.RequireRole("student"),
	)

	studentOnly.GET(
		"/students/me",
		studentH.GetMe,
	)

	// =========================
	// ATTENDANCE ROUTES
	// =========================

	attendanceRead := allUsers.Group("/attendance")
	{
		attendanceRead.GET(
			"",
			attendanceH.GetAll,
		)

		attendanceRead.GET(
			"/student/:studentId",
			attendanceH.GetByStudent,
		)

		attendanceRead.GET(
			"/student/:studentId/summary",
			attendanceH.GetSummary,
		)
	}

	attendanceWrite := teacherAndAdmin.Group("/attendance")
	{
		attendanceWrite.POST(
			"",
			attendanceH.Create,
		)
	}

	// =========================
	// ASSESSMENT ROUTES
	// =========================

	assessmentRead := allUsers.Group("/assessments")
	{
		assessmentRead.GET(
			"",
			assessmentH.GetAll,
		)

		assessmentRead.GET(
			"/student/:studentId",
			assessmentH.GetByStudent,
		)
	}

	assessmentWrite := teacherAndAdmin.Group("/assessments")
	{
		assessmentWrite.POST(
			"",
			assessmentH.Create,
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
