import axios from "axios";

const API_URL = "http://localhost:8080/api/v1";

export type Timetable = {
  id: string;
  college_id: string;
  teacher_id: string;
  classroom_id: string;
  subject: string;
  day_of_week: string;
  start_time: string;
  end_time: string;
  academic_year: string;
};

export type CreateTimetableRequest = {
  teacher_id: string;
  classroom_id: string;
  subject: string;
  day_of_week: string;
  start_time: string;
  end_time: string;
  academic_year: string;
};

// ========================================
// CREATE TIMETABLE
// ========================================

export const createTimetable = async (
  timetableData: CreateTimetableRequest,
): Promise<Timetable> => {
  const token = localStorage.getItem("token");

  const response = await axios.post<{
    message: string;
    timetable: Timetable;
  }>(`${API_URL}/timetables`, timetableData, {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  return response.data.timetable;
};

// ========================================
// GET ALL TIMETABLES
// ========================================

export const getTimetables = async (
  academicYear: string,
): Promise<Timetable[]> => {
  const token = localStorage.getItem("token");

  const response = await axios.get<{
    timetables: Timetable[];
  }>(`${API_URL}/timetables`, {
    params: {
      academic_year: academicYear,
    },
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data.timetables ?? [];
};

// ========================================
// GET TIMETABLE BY DAY
// ========================================

export const getTimetableByDay = async (
  day: string,
  academicYear: string,
): Promise<Timetable[]> => {
  const token = localStorage.getItem("token");

  const response = await axios.get<{
    timetables: Timetable[];
  }>(`${API_URL}/timetables/day/${day}`, {
    params: {
      academic_year: academicYear,
    },
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data.timetables ?? [];
};

// ========================================
// GET TEACHER TIMETABLE BY DAY
// ========================================

export const getTeacherTimetableByDay = async (
  teacherId: string,
  day: string,
  academicYear: string,
): Promise<Timetable[]> => {
  const token = localStorage.getItem("token");

  const response = await axios.get<{
    timetables: Timetable[];
  }>(`${API_URL}/timetables/teacher/${teacherId}`, {
    params: {
      day,
      academic_year: academicYear,
    },
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data.timetables ?? [];
};

// ========================================
// UPDATE TIMETABLE
// ========================================

export const updateTimetable = async (
  timetableId: string,
  timetableData: CreateTimetableRequest,
): Promise<Timetable> => {
  const token = localStorage.getItem("token");

  const response = await axios.put<{
    message: string;
    timetable: Timetable;
  }>(`${API_URL}/timetables/${timetableId}`, timetableData, {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  return response.data.timetable;
};

// ========================================
// DELETE TIMETABLE
// ========================================

export const deleteTimetable = async (timetableId: string): Promise<void> => {
  const token = localStorage.getItem("token");

  await axios.delete(`${API_URL}/timetables/${timetableId}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
};
