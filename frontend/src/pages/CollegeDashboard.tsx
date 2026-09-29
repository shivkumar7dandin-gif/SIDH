import { useEffect, useRef, useState, type MouseEvent } from "react";

import { getCollegeDashboard } from "../api/collegeApi";
import { getStudents } from "../api/studentApi";
import { getTeachers, type Teacher } from "../api/teacherApi";
import { getClassrooms, type Classroom } from "../api/classroomApi";

type College = {
  name: string;
  code: string;
  city: string;
  state: string;
};

type Student = {
  id: string;
  college_id: string;
  name: string;
  age: number;
  roll_number: number;
  gender: string;
  classroom_id: string;
};

type TodayAttendance = {
  total_students: number;
  present: number;
  absent: number;
  marked: number;
  not_marked: number;
  percentage: number;

  present_students: Student[];
  absent_students: Student[];
  not_marked_students: Student[];
};

type TeacherStatusSummary = {
  total: number;
  active: number;
  left: number;

  active_teachers: Teacher[];
  left_teachers: Teacher[];
};

type TeacherTodayClass = {
  teacher_id: string;
  teacher_name: string;
  subject: string;
  total_classes: number;
};

type ClassroomAttendance = {
  classroom_id: string;
  classroom_name: string;
  section: string;
  total_students: number;
  present: number;
  absent: number;
  marked: number;
  not_marked: number;
  percentage: number;
  status: string;
};

type DashboardData = {
  college: College | null;

  total_students: number;
  total_teachers: number;
  active_teachers: number;
  left_teachers: number;
  total_classrooms: number;

  today_attendance: TodayAttendance;

  teacher_status: TeacherStatusSummary;

  teacher_today_classes: TeacherTodayClass[];
  classroom_attendance: ClassroomAttendance[];
};

type SelectedView = "students" | "teachers" | "classrooms" | null;

type AttendanceFilter = "all" | "present" | "absent" | "not_marked";

type TeacherFilter = "all" | "active" | "left";

function CollegeDashboard() {
  const [dashboard, setDashboard] = useState<DashboardData | null>(null);

  const [students, setStudents] = useState<Student[]>([]);

  const [teachers, setTeachers] = useState<Teacher[]>([]);

  const [classrooms, setClassrooms] = useState<Classroom[]>([]);

  const [selectedView, setSelectedView] = useState<SelectedView>(null);

  const [attendanceFilter, setAttendanceFilter] =
    useState<AttendanceFilter>("all");

  const [teacherFilter, setTeacherFilter] = useState<TeacherFilter>("all");

  const [currentPage, setCurrentPage] = useState(1);

  const [loadingDetails, setLoadingDetails] = useState(false);

  const [error, setError] = useState("");

  const detailsRef = useRef<HTMLDivElement | null>(null);

  const itemsPerPage = 20;

  // ------------------------------------------------
  // SCROLL TO THE CURRENTLY SELECTED DETAILS SECTION
  // ------------------------------------------------

  useEffect(() => {
    if (!selectedView || loadingDetails) {
      return;
    }

    const timer = window.setTimeout(() => {
      detailsRef.current?.scrollIntoView({
        behavior: "smooth",
        block: "start",
      });
    }, 100);

    return () => window.clearTimeout(timer);
  }, [selectedView, attendanceFilter, teacherFilter, loadingDetails]);

  // ------------------------------------------------
  // LOAD DASHBOARD
  // ------------------------------------------------

  useEffect(() => {
    const loadDashboard = async () => {
      try {
        setError("");

        const data = await getCollegeDashboard();

        console.log("Dashboard response:", data);

        setDashboard(data);
      } catch (err) {
        console.error("Dashboard error:", err);

        setError("Failed to load dashboard");
      }
    };

    loadDashboard();
  }, []);

  // ------------------------------------------------
  // LOAD CLASSROOMS FOR STUDENT DETAILS
  // ------------------------------------------------

  const loadClassroomsForStudents = async () => {
    const classroomData = await getClassrooms();

    if (Array.isArray(classroomData)) {
      setClassrooms(classroomData);
    } else {
      setClassrooms([]);
    }
  };

  // ------------------------------------------------
  // CLICK TOTAL STUDENTS
  // ------------------------------------------------

  const handleStudentsClick = async () => {
    try {
      setLoadingDetails(true);
      setError("");
      setCurrentPage(1);

      // Switch only to student details.
      setSelectedView("students");
      setAttendanceFilter("all");

      // Clear previous teacher selection/details.
      setTeacherFilter("all");
      setTeachers([]);

      const data = await getStudents();

      if (Array.isArray(data)) {
        setStudents(data);
      } else if (data && Array.isArray(data.students)) {
        setStudents(data.students);
      } else {
        setStudents([]);
      }

      await loadClassroomsForStudents();
    } catch (err) {
      console.error("Load students error:", err);

      setError("Failed to load students");
    } finally {
      setLoadingDetails(false);
    }
  };

  // ------------------------------------------------
  // CLICK PRESENT
  // ------------------------------------------------

  const handlePresentClick = async (event: MouseEvent<HTMLDivElement>) => {
    event.stopPropagation();

    try {
      setLoadingDetails(true);
      setError("");
      setCurrentPage(1);

      setSelectedView("students");
      setAttendanceFilter("present");

      // Remove any previous teacher filter/details.
      setTeacherFilter("all");
      setTeachers([]);

      const presentStudents =
        dashboard?.today_attendance?.present_students ?? [];

      setStudents(presentStudents);

      await loadClassroomsForStudents();
    } catch (err) {
      console.error("Load present students error:", err);

      setError("Failed to load present students");
    } finally {
      setLoadingDetails(false);
    }
  };

  // ------------------------------------------------
  // CLICK ABSENT
  // ------------------------------------------------

  const handleAbsentClick = async (event: MouseEvent<HTMLDivElement>) => {
    event.stopPropagation();

    try {
      setLoadingDetails(true);
      setError("");
      setCurrentPage(1);

      setSelectedView("students");
      setAttendanceFilter("absent");

      // Remove any previous teacher filter/details.
      setTeacherFilter("all");
      setTeachers([]);

      const absentStudents = dashboard?.today_attendance?.absent_students ?? [];

      setStudents(absentStudents);

      await loadClassroomsForStudents();
    } catch (err) {
      console.error("Load absent students error:", err);

      setError("Failed to load absent students");
    } finally {
      setLoadingDetails(false);
    }
  };

  // ------------------------------------------------
  // CLICK NOT MARKED
  // ------------------------------------------------

  const handleNotMarkedClick = async (event: MouseEvent<HTMLDivElement>) => {
    event.stopPropagation();

    try {
      setLoadingDetails(true);
      setError("");
      setCurrentPage(1);

      setSelectedView("students");
      setAttendanceFilter("not_marked");

      // Remove any previous teacher filter/details.
      setTeacherFilter("all");
      setTeachers([]);

      const notMarkedStudents =
        dashboard?.today_attendance?.not_marked_students ?? [];

      setStudents(notMarkedStudents);

      await loadClassroomsForStudents();
    } catch (err) {
      console.error("Load not marked students error:", err);

      setError("Failed to load not marked students");
    } finally {
      setLoadingDetails(false);
    }
  };

  // ------------------------------------------------
  // CLICK TOTAL TEACHERS
  // ------------------------------------------------

  const handleTeachersClick = async () => {
    try {
      setLoadingDetails(true);
      setError("");
      setCurrentPage(1);

      // Switch only to teacher details.
      setSelectedView("teachers");
      setTeacherFilter("all");

      // Clear previous student selection/details.
      setAttendanceFilter("all");
      setStudents([]);

      const data = await getTeachers();

      if (Array.isArray(data)) {
        setTeachers(data);
      } else {
        setTeachers([]);
      }
    } catch (err) {
      console.error("Load teachers error:", err);

      setError("Failed to load teachers");
    } finally {
      setLoadingDetails(false);
    }
  };

  // ------------------------------------------------
  // CLICK ACTIVE TEACHERS
  // ------------------------------------------------

  const handleActiveTeachersClick = (event: MouseEvent<HTMLDivElement>) => {
    event.stopPropagation();

    setError("");
    setCurrentPage(1);

    setSelectedView("teachers");
    setTeacherFilter("active");

    // Remove any previous student filter/details.
    setAttendanceFilter("all");
    setStudents([]);

    const activeTeachers = dashboard?.teacher_status?.active_teachers ?? [];

    setTeachers(activeTeachers);
  };

  // ------------------------------------------------
  // CLICK LEFT TEACHERS
  // ------------------------------------------------

  const handleLeftTeachersClick = (event: MouseEvent<HTMLDivElement>) => {
    event.stopPropagation();

    setError("");
    setCurrentPage(1);

    setSelectedView("teachers");
    setTeacherFilter("left");

    // Remove any previous student filter/details.
    setAttendanceFilter("all");
    setStudents([]);

    const leftTeachers = dashboard?.teacher_status?.left_teachers ?? [];

    setTeachers(leftTeachers);
  };

  // ------------------------------------------------
  // CLICK CLASSROOMS
  // ------------------------------------------------

  const handleClassroomsClick = async () => {
    try {
      setLoadingDetails(true);
      setError("");
      setCurrentPage(1);

      // Switch only to classroom details.
      setSelectedView("classrooms");

      // Clear student and teacher selections/details.
      setAttendanceFilter("all");
      setTeacherFilter("all");
      setStudents([]);
      setTeachers([]);

      const data = await getClassrooms();

      if (Array.isArray(data)) {
        setClassrooms(data);
      } else {
        setClassrooms([]);
      }
    } catch (err) {
      console.error("Load classrooms error:", err);

      setError("Failed to load classrooms");
    } finally {
      setLoadingDetails(false);
    }
  };

  // ------------------------------------------------
  // CLASSROOM NAME
  // ------------------------------------------------

  const getClassroomName = (classroomID: string) => {
    const classroom = classrooms.find((item) => item.id === classroomID);

    if (!classroom) {
      return "Classroom not found";
    }

    return `${classroom.name} - Section ${classroom.section}`;
  };

  // ------------------------------------------------
  // SORT CLASSROOMS
  // ------------------------------------------------

  const sortedClassrooms = [...classrooms].sort((a, b) => {
    const classA = parseInt(a.name) || 0;

    const classB = parseInt(b.name) || 0;

    if (classA !== classB) {
      return classA - classB;
    }

    return a.section.localeCompare(b.section);
  });

  // ------------------------------------------------
  // PAGINATION
  // ------------------------------------------------

  let totalItems = 0;

  if (selectedView === "students") {
    totalItems = students.length;
  }

  if (selectedView === "teachers") {
    totalItems = teachers.length;
  }

  if (selectedView === "classrooms") {
    totalItems = sortedClassrooms.length;
  }

  const totalPages = Math.ceil(totalItems / itemsPerPage);

  const startIndex = (currentPage - 1) * itemsPerPage;

  const endIndex = startIndex + itemsPerPage;

  const paginatedStudents = students.slice(startIndex, endIndex);

  const paginatedTeachers = teachers.slice(startIndex, endIndex);

  const paginatedClassrooms = sortedClassrooms.slice(startIndex, endIndex);

  // ------------------------------------------------
  // STUDENT DETAILS TITLE
  // ------------------------------------------------

  const getStudentTitle = () => {
    switch (attendanceFilter) {
      case "present":
        return "Present Students";

      case "absent":
        return "Absent Students";

      case "not_marked":
        return "Attendance Not Marked Students";

      default:
        return "Student Details";
    }
  };

  // ------------------------------------------------
  // TEACHER DETAILS TITLE
  // ------------------------------------------------

  const getTeacherTitle = () => {
    switch (teacherFilter) {
      case "active":
        return "Active Teachers";

      case "left":
        return "Left Teachers";

      default:
        return "Teacher Details";
    }
  };

  // ------------------------------------------------
  // TOTAL CLASSES TODAY
  // ------------------------------------------------

  const totalClassesToday = (dashboard?.teacher_today_classes ?? []).reduce(
    (total, teacher) => total + teacher.total_classes,
    0,
  );

  const classroomAttendance = dashboard?.classroom_attendance ?? [];

  const completedClassrooms = classroomAttendance.filter(
    (classroom) => classroom.status === "Completed",
  ).length;

  const partialClassrooms = classroomAttendance.filter(
    (classroom) => classroom.status === "Partial",
  ).length;

  const notMarkedClassrooms = classroomAttendance.filter(
    (classroom) => classroom.status === "Not Marked",
  ).length;

  const noStudentClassrooms = classroomAttendance.filter(
    (classroom) => classroom.status === "No Students",
  ).length;

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold">Admin Dashboard</h1>

      <p className="mt-2 text-gray-600">
        Welcome to your school management dashboard.
      </p>

      {dashboard?.college && (
        <div className="mt-4">
          <h2 className="text-xl font-semibold">{dashboard.college.name}</h2>

          <p className="text-gray-500">
            {dashboard.college.city}, {dashboard.college.state}
          </p>

          <p className="text-sm text-gray-400">
            College Code: {dashboard.college.code}
          </p>
        </div>
      )}

      {error && <p className="mt-4 text-red-600">{error}</p>}

      {/* ========================================= */}
      {/* DASHBOARD CARDS */}
      {/* ========================================= */}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
        {/* ========================================= */}
        {/* STUDENTS */}
        {/* ========================================= */}

        <div
          onClick={handleStudentsClick}
          className={`bg-white p-6 rounded-xl shadow cursor-pointer transition hover:shadow-lg ${
            selectedView === "students"
              ? "border-2 border-blue-600 bg-blue-50"
              : ""
          }`}
        >
          <h3 className="text-gray-500 font-medium">Total Students</h3>

          <p className="text-4xl font-bold mt-3">
            {dashboard?.total_students ?? 0}
          </p>

          {/* TODAY ATTENDANCE */}

          <div className="mt-5 pt-4 border-t">
            <p className="text-sm font-semibold text-gray-600 mb-3">
              Today's Attendance
            </p>

            <div className="grid grid-cols-3 gap-2 text-center">
              {/* PRESENT */}

              <div
                onClick={handlePresentClick}
                className={`rounded-lg p-2 cursor-pointer transition hover:shadow ${
                  attendanceFilter === "present"
                    ? "bg-green-200 ring-2 ring-green-500"
                    : "bg-green-50 hover:bg-green-100"
                }`}
              >
                <p className="text-xl font-bold text-green-600">
                  {dashboard?.today_attendance?.present ?? 0}
                </p>

                <p className="text-xs text-gray-500">Present</p>
              </div>

              {/* ABSENT */}

              <div
                onClick={handleAbsentClick}
                className={`rounded-lg p-2 cursor-pointer transition hover:shadow ${
                  attendanceFilter === "absent"
                    ? "bg-red-200 ring-2 ring-red-500"
                    : "bg-red-50 hover:bg-red-100"
                }`}
              >
                <p className="text-xl font-bold text-red-600">
                  {dashboard?.today_attendance?.absent ?? 0}
                </p>

                <p className="text-xs text-gray-500">Absent</p>
              </div>

              {/* NOT MARKED */}

              <div
                onClick={handleNotMarkedClick}
                className={`rounded-lg p-2 cursor-pointer transition hover:shadow ${
                  attendanceFilter === "not_marked"
                    ? "bg-yellow-200 ring-2 ring-yellow-500"
                    : "bg-yellow-50 hover:bg-yellow-100"
                }`}
              >
                <p className="text-xl font-bold text-yellow-600">
                  {dashboard?.today_attendance?.not_marked ?? 0}
                </p>

                <p className="text-xs text-gray-500">Not Marked</p>
              </div>
            </div>

            {/* ATTENDANCE PERCENTAGE */}

            <div className="mt-4">
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Attendance</span>

                <span className="font-semibold">
                  {(dashboard?.today_attendance?.percentage ?? 0).toFixed(1)}%
                </span>
              </div>

              <div className="w-full bg-gray-200 rounded-full h-2 mt-2">
                <div
                  className="bg-green-500 h-2 rounded-full"
                  style={{
                    width: `${Math.min(
                      dashboard?.today_attendance?.percentage ?? 0,
                      100,
                    )}%`,
                  }}
                />
              </div>
            </div>
          </div>

          <p className="text-blue-600 text-sm mt-4">
            Click total card to view all students
          </p>
        </div>

        {/* ========================================= */}
        {/* TEACHERS */}
        {/* ========================================= */}

        <div
          onClick={handleTeachersClick}
          className={`bg-white p-6 rounded-xl shadow cursor-pointer transition hover:shadow-lg ${
            selectedView === "teachers"
              ? "border-2 border-blue-600 bg-blue-50"
              : ""
          }`}
        >
          <h3 className="text-gray-500 font-medium">Total Teachers</h3>

          <p className="text-4xl font-bold mt-3">
            {dashboard?.total_teachers ?? 0}
          </p>

          {/* TEACHER STATUS */}

          <div className="mt-5 pt-4 border-t">
            <p className="text-sm font-semibold text-gray-600 mb-3">
              Teacher Status
            </p>

            <div className="grid grid-cols-2 gap-3 text-center">
              {/* ACTIVE */}

              <div
                onClick={handleActiveTeachersClick}
                className={`rounded-lg p-3 cursor-pointer transition hover:shadow ${
                  teacherFilter === "active"
                    ? "bg-green-200 ring-2 ring-green-500"
                    : "bg-green-50 hover:bg-green-100"
                }`}
              >
                <p className="text-2xl font-bold text-green-600">
                  {dashboard?.teacher_status?.active ?? 0}
                </p>

                <p className="text-xs text-gray-500">Active</p>
              </div>

              {/* LEFT */}

              <div
                onClick={handleLeftTeachersClick}
                className={`rounded-lg p-3 cursor-pointer transition hover:shadow ${
                  teacherFilter === "left"
                    ? "bg-red-200 ring-2 ring-red-500"
                    : "bg-red-50 hover:bg-red-100"
                }`}
              >
                <p className="text-2xl font-bold text-red-600">
                  {dashboard?.teacher_status?.left ?? 0}
                </p>

                <p className="text-xs text-gray-500">Left</p>
              </div>
            </div>
          </div>

          <p className="text-blue-600 text-sm mt-4">
            Click total card to view all teachers
          </p>
        </div>

        {/* ========================================= */}
        {/* CLASSROOMS */}
        {/* ========================================= */}

        <div
          onClick={handleClassroomsClick}
          className={`bg-white p-6 rounded-xl shadow cursor-pointer transition hover:shadow-lg ${
            selectedView === "classrooms"
              ? "border-2 border-blue-600 bg-blue-50"
              : ""
          }`}
        >
          <h3 className="text-gray-500 font-medium">Total Classrooms</h3>

          <p className="text-4xl font-bold mt-3">
            {dashboard?.total_classrooms ?? 0}
          </p>

          <div className="mt-5 pt-4 border-t">
            <p className="text-sm font-semibold text-gray-600 mb-3">
              Today's Class Attendance
            </p>

            <div className="grid grid-cols-2 gap-2 text-center">
              <div className="rounded-lg p-2 bg-green-50">
                <p className="text-xl font-bold text-green-600">
                  {completedClassrooms}
                </p>
                <p className="text-xs text-gray-500">Completed</p>
              </div>

              <div className="rounded-lg p-2 bg-yellow-50">
                <p className="text-xl font-bold text-yellow-600">
                  {partialClassrooms}
                </p>
                <p className="text-xs text-gray-500">Partial</p>
              </div>

              <div className="rounded-lg p-2 bg-orange-50">
                <p className="text-xl font-bold text-orange-600">
                  {notMarkedClassrooms}
                </p>
                <p className="text-xs text-gray-500">Not Marked</p>
              </div>

              <div className="rounded-lg p-2 bg-gray-100">
                <p className="text-xl font-bold text-gray-600">
                  {noStudentClassrooms}
                </p>
                <p className="text-xs text-gray-500">No Students</p>
              </div>
            </div>
          </div>

          <p className="text-blue-600 text-sm mt-4">
            Click to view all classroom attendance
          </p>
        </div>

        {/* IMPORTANT: closes DASHBOARD CARDS grid */}
      </div>

      {/* ========================================= */}
      {/* TODAY'S CLASSROOM ATTENDANCE */}
      {/* ========================================= */}

      <div className="bg-white mt-8 rounded-xl shadow overflow-hidden">
        <div className="p-6 border-b">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 className="text-xl font-semibold">
                Today's Classroom Attendance
              </h2>

              <p className="text-gray-500 mt-1">
                Present, absent and not marked students for each classroom
                today.
              </p>
            </div>

            <div className="bg-blue-50 px-5 py-3 rounded-lg min-w-32">
              <p className="text-sm text-gray-500 text-center">
                Total Classrooms
              </p>

              <p className="text-2xl font-bold text-blue-600 text-center">
                {dashboard?.classroom_attendance?.length ?? 0}
              </p>
            </div>
          </div>
        </div>

        {(dashboard?.classroom_attendance ?? []).length === 0 ? (
          <div className="p-6 text-gray-500">
            No classroom attendance data found.
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-100">
                <tr>
                  <th className="text-left p-4">Classroom</th>

                  <th className="text-center p-4">Students</th>

                  <th className="text-center p-4">Present</th>

                  <th className="text-center p-4">Absent</th>

                  <th className="text-center p-4">Not Marked</th>

                  <th className="text-center p-4">Attendance</th>

                  <th className="text-center p-4">Status</th>
                </tr>
              </thead>

              <tbody>
                {(dashboard?.classroom_attendance ?? []).map((classroom) => (
                  <tr
                    key={classroom.classroom_id}
                    className="border-t hover:bg-gray-50"
                  >
                    <td className="p-4 font-medium">
                      {classroom.classroom_name}
                      {classroom.section
                        ? ` - Section ${classroom.section}`
                        : ""}
                    </td>

                    <td className="p-4 text-center">
                      {classroom.total_students}
                    </td>

                    <td className="p-4 text-center">
                      <span className="inline-flex items-center justify-center min-w-10 px-3 py-1 rounded-full bg-green-100 text-green-700 font-semibold">
                        {classroom.present}
                      </span>
                    </td>

                    <td className="p-4 text-center">
                      <span className="inline-flex items-center justify-center min-w-10 px-3 py-1 rounded-full bg-red-100 text-red-700 font-semibold">
                        {classroom.absent}
                      </span>
                    </td>

                    <td className="p-4 text-center">
                      <span className="inline-flex items-center justify-center min-w-10 px-3 py-1 rounded-full bg-yellow-100 text-yellow-700 font-semibold">
                        {classroom.not_marked}
                      </span>
                    </td>

                    <td className="p-4 text-center font-semibold">
                      {classroom.marked > 0
                        ? `${classroom.percentage.toFixed(1)}%`
                        : "-"}
                    </td>

                    <td className="p-4 text-center">
                      {classroom.status === "Completed" ? (
                        <span className="px-3 py-1 text-sm rounded-full bg-green-100 text-green-700">
                          Completed
                        </span>
                      ) : classroom.status === "Partial" ? (
                        <span className="px-3 py-1 text-sm rounded-full bg-yellow-100 text-yellow-700">
                          Partial
                        </span>
                      ) : classroom.status === "No Students" ? (
                        <span className="px-3 py-1 text-sm rounded-full bg-gray-100 text-gray-600">
                          No Students
                        </span>
                      ) : (
                        <span className="px-3 py-1 text-sm rounded-full bg-gray-100 text-gray-600">
                          Not Marked
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* ========================================= */}
      {/* TODAY'S TEACHER CLASSES */}
      {/* ========================================= */}

      <div className="bg-white mt-8 rounded-xl shadow overflow-hidden">
        <div className="p-6 border-b">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 className="text-xl font-semibold">Today's Teacher Classes</h2>

              <p className="text-gray-500 mt-1">
                Number of classes scheduled for each active teacher today.
              </p>
            </div>

            <div className="bg-blue-50 px-5 py-3 rounded-lg min-w-32">
              <p className="text-sm text-gray-500 text-center">Total Classes</p>

              <p className="text-2xl font-bold text-blue-600 text-center">
                {totalClassesToday}
              </p>
            </div>
          </div>
        </div>

        {(dashboard?.teacher_today_classes ?? []).length === 0 ? (
          <div className="p-6 text-gray-500">No active teachers found.</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-100">
                <tr>
                  <th className="text-left p-4">Teacher Name</th>

                  <th className="text-left p-4">Subject</th>

                  <th className="text-center p-4">Classes Today</th>

                  <th className="text-center p-4">Status</th>
                </tr>
              </thead>

              <tbody>
                {(dashboard?.teacher_today_classes ?? []).map((teacher) => (
                  <tr
                    key={teacher.teacher_id}
                    className="border-t hover:bg-gray-50"
                  >
                    <td className="p-4 font-medium">{teacher.teacher_name}</td>

                    <td className="p-4">{teacher.subject || "-"}</td>

                    <td className="p-4 text-center">
                      <span className="inline-flex items-center justify-center min-w-10 px-3 py-1 rounded-full bg-blue-100 text-blue-700 font-semibold">
                        {teacher.total_classes}
                      </span>
                    </td>

                    <td className="p-4 text-center">
                      {teacher.total_classes > 0 ? (
                        <span className="px-3 py-1 text-sm rounded-full bg-green-100 text-green-700">
                          Scheduled
                        </span>
                      ) : (
                        <span className="px-3 py-1 text-sm rounded-full bg-gray-100 text-gray-600">
                          No Classes
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* ========================================= */}
      {/* LOADING */}
      {/* ========================================= */}

      {loadingDetails && (
        <div className="bg-white mt-8 p-6 rounded-xl shadow">
          Loading details...
        </div>
      )}

      <div ref={detailsRef} className="scroll-mt-6">
        {/* ========================================= */}
        {/* STUDENT DETAILS */}
        {/* ========================================= */}

        {!loadingDetails && selectedView === "students" && (
          <div className="bg-white mt-8 rounded-xl shadow overflow-hidden">
            <div className="p-6 border-b">
              <h2 className="text-xl font-semibold">{getStudentTitle()}</h2>

              <p className="text-gray-500 mt-2">
                Total Students: {students.length}
              </p>
            </div>

            {students.length === 0 ? (
              <div className="p-6 text-gray-500">
                {attendanceFilter === "present"
                  ? "No present students found for today."
                  : attendanceFilter === "absent"
                    ? "No absent students found for today."
                    : attendanceFilter === "not_marked"
                      ? "All students attendance is marked."
                      : "No students found."}
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-gray-100">
                    <tr>
                      <th className="text-left p-4">Roll No</th>

                      <th className="text-left p-4">Name</th>

                      <th className="text-left p-4">Age</th>

                      <th className="text-left p-4">Gender</th>

                      <th className="text-left p-4">Classroom</th>
                    </tr>
                  </thead>

                  <tbody>
                    {paginatedStudents.map((student) => (
                      <tr key={student.id} className="border-t">
                        <td className="p-4">{student.roll_number}</td>

                        <td className="p-4">{student.name}</td>

                        <td className="p-4">{student.age}</td>

                        <td className="p-4">{student.gender}</td>

                        <td className="p-4">
                          {getClassroomName(student.classroom_id)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {/* ========================================= */}
        {/* TEACHER DETAILS */}
        {/* ========================================= */}

        {!loadingDetails && selectedView === "teachers" && (
          <div className="bg-white mt-8 rounded-xl shadow overflow-hidden">
            <div className="p-6 border-b">
              <h2 className="text-xl font-semibold">{getTeacherTitle()}</h2>

              <p className="text-gray-500 mt-2">
                Total Teachers: {teachers.length}
              </p>
            </div>

            {teachers.length === 0 ? (
              <div className="p-6 text-gray-500">
                {teacherFilter === "active"
                  ? "No active teachers found."
                  : teacherFilter === "left"
                    ? "No teachers with left status found."
                    : "No teachers found."}
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-gray-100">
                    <tr>
                      <th className="text-left p-4">Name</th>

                      <th className="text-left p-4">Subject</th>

                      <th className="text-left p-4">Email</th>

                      <th className="text-left p-4">Phone</th>

                      <th className="text-left p-4">Status</th>
                    </tr>
                  </thead>

                  <tbody>
                    {paginatedTeachers.map((teacher) => (
                      <tr key={teacher.id} className="border-t">
                        <td className="p-4">{teacher.name}</td>

                        <td className="p-4">{teacher.subject}</td>

                        <td className="p-4">{teacher.email}</td>

                        <td className="p-4">{teacher.phone}</td>

                        <td className="p-4">
                          {teacher.status === "left" ? (
                            <span className="px-3 py-1 text-sm rounded-full bg-red-100 text-red-700">
                              Left
                            </span>
                          ) : (
                            <span className="px-3 py-1 text-sm rounded-full bg-green-100 text-green-700">
                              Active
                            </span>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {/* ========================================= */}
        {/* CLASSROOM DETAILS */}
        {/* ========================================= */}

        {!loadingDetails && selectedView === "classrooms" && (
          <div className="bg-white mt-8 rounded-xl shadow overflow-hidden">
            <div className="p-6 border-b">
              <h2 className="text-xl font-semibold">Classroom Details</h2>

              <p className="text-gray-500 mt-2">
                Total Classrooms: {sortedClassrooms.length}
              </p>
            </div>

            {sortedClassrooms.length === 0 ? (
              <div className="p-6 text-gray-500">No classrooms found.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-gray-100">
                    <tr>
                      <th className="text-left p-4">Class</th>

                      <th className="text-left p-4">Section</th>

                      <th className="text-left p-4">Capacity</th>
                    </tr>
                  </thead>

                  <tbody>
                    {paginatedClassrooms.map((classroom) => (
                      <tr key={classroom.id} className="border-t">
                        <td className="p-4">{classroom.name}</td>

                        <td className="p-4">{classroom.section}</td>

                        <td className="p-4">{classroom.capacity}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>

      {/* ========================================= */}
      {/* PAGINATION */}
      {/* ========================================= */}

      {!loadingDetails && selectedView && totalItems > 0 && (
        <div className="bg-white flex justify-between items-center p-5 border-t rounded-b-xl shadow">
          <p className="text-gray-600">
            Showing {startIndex + 1}
            {" - "}
            {Math.min(endIndex, totalItems)}
            {" of "}
            {totalItems}
          </p>

          <div className="flex items-center gap-3">
            <button
              type="button"
              disabled={currentPage === 1}
              onClick={() => setCurrentPage((current) => current - 1)}
              className="px-4 py-2 border rounded-lg disabled:opacity-40"
            >
              Previous
            </button>

            <span>
              Page {currentPage}
              {" of "}
              {totalPages || 1}
            </span>

            <button
              type="button"
              disabled={currentPage >= totalPages}
              onClick={() => setCurrentPage((current) => current + 1)}
              className="px-4 py-2 border rounded-lg disabled:opacity-40"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </main>
  );
}

export default CollegeDashboard;
