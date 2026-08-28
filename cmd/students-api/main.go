package main

import (
	"log"

	"github.com/gin-gonic/gin"

	assessmentHandler "github.com/shivkumar7dandin-gif/students-api/internal/assessment/handler"
	assessmentRepository "github.com/shivkumar7dandin-gif/students-api/internal/assessment/repository"
	assessmentService "github.com/shivkumar7dandin-gif/students-api/internal/assessment/service"

	attendanceHandler "github.com/shivkumar7dandin-gif/students-api/internal/attendance/handler"
	attendanceRepository "github.com/shivkumar7dandin-gif/students-api/internal/attendance/repository"
	attendanceService "github.com/shivkumar7dandin-gif/students-api/internal/attendance/service"

	classroomHandler "github.com/shivkumar7dandin-gif/students-api/internal/classroom/handler"
	classroomRepository "github.com/shivkumar7dandin-gif/students-api/internal/classroom/repository"
	classroomService "github.com/shivkumar7dandin-gif/students-api/internal/classroom/service"

	"github.com/shivkumar7dandin-gif/students-api/internal/config"
	"github.com/shivkumar7dandin-gif/students-api/internal/database"

	studentHandler "github.com/shivkumar7dandin-gif/students-api/internal/student/handler"
	studentRepository "github.com/shivkumar7dandin-gif/students-api/internal/student/repository"
	studentService "github.com/shivkumar7dandin-gif/students-api/internal/student/service"
)

func main() {

	cfg := config.MustLoad()

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

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	api := router.Group("/api/v1")

	// CLASSROOM ROUTES
	classrooms := api.Group("/classrooms")
	{
		classrooms.POST("", classroomH.Create)
		classrooms.GET("", classroomH.GetAll)
		classrooms.GET("/:id", classroomH.GetByID)
		classrooms.PUT("/:id", classroomH.Update)
		classrooms.DELETE("/:id", classroomH.Delete)
	}

	// STUDENT ROUTES
	students := api.Group("/students")
	{
		students.POST("", studentH.Create)
		students.GET("", studentH.GetAll)
		students.GET("/:id", studentH.GetByID)
		students.PUT("/:id", studentH.Update)
		students.DELETE("/:id", studentH.Delete)
	}

	// ATTENDANCE ROUTES
	attendance := api.Group("/attendance")
	{
		attendance.POST("", attendanceH.Create)
		attendance.GET("", attendanceH.GetAll)
		attendance.GET(
			"/student/:studentId",
			attendanceH.GetByStudent,
		)

		attendance.GET(
			"/student/:studentId/summary",
			attendanceH.GetSummary,
		)

	}

	// ASSESSMENT ROUTES
	assessments := api.Group("/assessments")
	{
		assessments.POST("", assessmentH.Create)
		assessments.GET("", assessmentH.GetAll)
		assessments.GET(
			"/student/:studentId",
			assessmentH.GetByStudent,
		)
	}

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
