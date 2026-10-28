import axios from "axios";

const API_URL = "http://localhost:8080/api/v1";

export type TeacherAddress = {
  house_no: string;
  street: string;
  village: string;
  city: string;
  state: string;
  pincode: string;
};

export type Teacher = {
  id: string;
  college_id: string;

  employee_id: string;
  name: string;
  age: number;
  gender: string;

  email: string;
  phone: string;

  subject: string;
  qualification: string;
  experience_years: number;
  joining_date: string;

  address: TeacherAddress;

  username: string;
  status: string;
};

export type CreateTeacherData = {
  employee_id: string;
  name: string;
  age: number;
  gender: string;

  email: string;
  phone: string;

  subject: string;
  qualification: string;
  experience_years: number;
  joining_date: string;

  address: TeacherAddress;

  username: string;
  password: string;
};

export type UpdateTeacherData = {
  employee_id: string;
  name: string;
  age: number;
  gender: string;

  email: string;
  phone: string;

  subject: string;
  qualification: string;
  experience_years: number;
  joining_date: string;

  address: TeacherAddress;
};

export type UpdateTeacherResponse = {
  message: string;
};

export type TeacherClassAssignment = {
  id: string;
  college_id: string;
  teacher_id: string;
  classroom_id: string;
  subject: string;
  created_at: string;
  updated_at: string;
};

export type CreateTeacherAssignmentData = {
  teacher_id: string;
  classroom_id: string;
  subject: string;
};

export type DeleteTeacherAssignmentResponse = {
  message: string;
};

// GET ALL TEACHERS
export const getTeachers = async (): Promise<Teacher[]> => {
  const token = localStorage.getItem("token");

  const response = await axios.get<Teacher[]>(`${API_URL}/teachers`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data;
};

// CREATE TEACHER
export const createTeacher = async (
  teacherData: CreateTeacherData,
): Promise<Teacher> => {
  const token = localStorage.getItem("token");

  const response = await axios.post<Teacher>(
    `${API_URL}/teachers`,
    teacherData,
    {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    },
  );

  return response.data;
};

// GET ONE TEACHER
export const getTeacherById = async (id: string): Promise<Teacher> => {
  const token = localStorage.getItem("token");

  const response = await axios.get<Teacher>(`${API_URL}/teachers/${id}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data;
};

// UPDATE TEACHER
export const updateTeacher = async (
  id: string,
  teacherData: UpdateTeacherData,
): Promise<UpdateTeacherResponse> => {
  const token = localStorage.getItem("token");

  const response = await axios.put<UpdateTeacherResponse>(
    `${API_URL}/teachers/${id}`,
    teacherData,
    {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    },
  );

  return response.data;
};

// DELETE TEACHER
export const deleteTeacher = async (id: string) => {
  const token = localStorage.getItem("token");

  const response = await axios.delete(`${API_URL}/teachers/${id}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data;
};

// UPDATE TEACHER STATUS
export const updateTeacherStatus = async (
  id: string,
  status: "active" | "left",
) => {
  const token = localStorage.getItem("token");

  const response = await axios.patch(
    `${API_URL}/teachers/${id}/status`,
    {
      status,
    },
    {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    },
  );

  return response.data;
};

// =====================================================
// TEACHER CLASS ASSIGNMENTS
// =====================================================

// GET ASSIGNMENTS FOR ONE TEACHER
export const getTeacherAssignments = async (
  teacherId: string,
): Promise<TeacherClassAssignment[]> => {
  const token = localStorage.getItem("token");

  const response = await axios.get<TeacherClassAssignment[]>(
    `${API_URL}/teachers/${teacherId}/assignments`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    },
  );

  return response.data;
};

// CREATE TEACHER CLASS ASSIGNMENT
export const createTeacherAssignment = async (
  assignmentData: CreateTeacherAssignmentData,
): Promise<TeacherClassAssignment> => {
  const token = localStorage.getItem("token");

  const response = await axios.post<TeacherClassAssignment>(
    `${API_URL}/teacher-assignments`,
    assignmentData,
    {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    },
  );

  return response.data;
};

// DELETE TEACHER CLASS ASSIGNMENT
export const deleteTeacherAssignment = async (
  assignmentId: string,
): Promise<DeleteTeacherAssignmentResponse> => {
  const token = localStorage.getItem("token");

  const response = await axios.delete<DeleteTeacherAssignmentResponse>(
    `${API_URL}/teacher-assignments/${assignmentId}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    },
  );

  return response.data;
};
