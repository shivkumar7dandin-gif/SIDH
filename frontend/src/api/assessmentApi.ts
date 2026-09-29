const API_URL = "http://localhost:8080/api/v1";

export interface Assessment {
  id?: string;
  college_id?: string;
  student_id: string;
  classroom_id: string;
  academic_year: string;
  assessment_date: string;
  assessment_no: number;
  subject: string;
  total_marks: number;
  obtained_marks: number;
  percentage?: number;
  result?: string;
}

export async function createAssessment(
  assessment: Assessment
) {
  const token = localStorage.getItem("token");

  const response = await fetch(
    `${API_URL}/assessments`,
    {
      method: "POST",

      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },

      body: JSON.stringify(assessment),
    }
  );

  const data = await response.json();

  if (!response.ok) {
    throw new Error(
      data.error || "Failed to create assessment"
    );
  }

  return data;
}

export async function getStudentAssessments(
  studentId: string,
  academicYear: string
) {
  const token = localStorage.getItem("token");

  const response = await fetch(
    `${API_URL}/assessments/student/${studentId}?academic_year=${encodeURIComponent(
      academicYear
    )}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  const data = await response.json();

  if (!response.ok) {
    throw new Error(
      data.error || "Failed to get assessments"
    );
  }

  return data;
}