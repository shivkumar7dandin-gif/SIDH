const API_URL = "http://localhost:8080/api/v1";

export interface Attendance {
  id?: string;
  college_id?: string;
  student_id: string;
  classroom_id: string;
  academic_year: string;
  date: string;
  attendance: "Present" | "Absent";
}

// -----------------------------------------
// BULK ATTENDANCE REQUEST
// -----------------------------------------

export interface BulkAttendanceRequest {
  classroom_id: string;
  academic_year: string;
  date: string;

  attendances: {
    student_id: string;
    attendance: "Present" | "Absent";
  }[];
}

// -----------------------------------------
// CLASSROOM MONTHLY STUDENT SUMMARY
// -----------------------------------------

export interface ClassroomStudentAttendanceSummary {
  student_id: string;
  roll_number: number;
  name: string;
  present: number;
  absent: number;
  attendance_marked_days: number;
  attendance_percentage: number;
}

// -----------------------------------------
// CLASSROOM MONTHLY ATTENDANCE SUMMARY
// -----------------------------------------

export interface ClassroomMonthlyAttendanceSummary {
  classroom_id: string;
  academic_year: string;
  year: number;
  month: number;
  working_days: number;
  attendance_marked_days: number;
  total_students: number;
  present: number;
  absent: number;
  attendance_percentage: number;
  students: ClassroomStudentAttendanceSummary[];
}

// -----------------------------------------
// CREATE / UPDATE SINGLE ATTENDANCE
// -----------------------------------------

export async function createAttendance(attendance: Attendance) {
  const token = localStorage.getItem("token");

  const response = await fetch(`${API_URL}/attendance`, {
    method: "POST",

    headers: {
      "Content-Type": "application/json",

      Authorization: `Bearer ${token}`,
    },

    body: JSON.stringify(attendance),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Failed to save attendance");
  }

  return data;
}

// -----------------------------------------
// CREATE / UPDATE BULK ATTENDANCE
// -----------------------------------------

export async function createBulkAttendance(request: BulkAttendanceRequest) {
  const token = localStorage.getItem("token");

  const response = await fetch(`${API_URL}/attendance/bulk`, {
    method: "POST",

    headers: {
      "Content-Type": "application/json",

      Authorization: `Bearer ${token}`,
    },

    body: JSON.stringify(request),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Failed to save attendance");
  }

  return data;
}

// -----------------------------------------
// GET ALL ATTENDANCE
// -----------------------------------------

export async function getAttendance(
  academicYear: string,
): Promise<Attendance[]> {
  const token = localStorage.getItem("token");

  const response = await fetch(
    `${API_URL}/attendance?academic_year=${encodeURIComponent(academicYear)}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    },
  );

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Failed to get attendance");
  }

  return data;
}

// -----------------------------------------
// GET ATTENDANCE BY STUDENT
// -----------------------------------------

export async function getStudentAttendance(
  studentId: string,
  academicYear: string,
): Promise<Attendance[]> {
  const token = localStorage.getItem("token");

  const response = await fetch(
    `${API_URL}/attendance/student/${studentId}?academic_year=${encodeURIComponent(
      academicYear,
    )}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    },
  );

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Failed to get attendance");
  }

  return data;
}

// -----------------------------------------
// GET CLASSROOM MONTHLY ATTENDANCE HISTORY
// -----------------------------------------

export async function getClassroomMonthlyAttendance(
  classroomId: string,
  academicYear: string,
  month: string,
): Promise<ClassroomMonthlyAttendanceSummary> {
  const token = localStorage.getItem("token");

  const response = await fetch(
    `${API_URL}/attendance/classroom/${classroomId}/monthly-summary?academic_year=${encodeURIComponent(
      academicYear,
    )}&month=${encodeURIComponent(month)}`,
    {
      method: "GET",

      headers: {
        Authorization: `Bearer ${token}`,
      },
    },
  );

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Failed to get classroom monthly attendance");
  }

  return data;
}
