import axios from "axios";

const API_URL = "http://localhost:8080/api/v1";

export type Teacher = {
  id: string;
  college_id: string;
  name: string;
  age: number;
  gender: string;
  email: string;
  phone: string;
  subject: string;
  username: string;
};

export type CreateTeacherData = {
  name: string;
  age: number;
  gender: string;
  email: string;
  phone: string;
  subject: string;
  username: string;
  password: string;
};

export type UpdateTeacherData = {
  name: string;
  age: number;
  gender: string;
  email: string;
  phone: string;
  subject: string;
};

// GET ALL TEACHERS
export const getTeachers = async (): Promise<Teacher[]> => {
  const token = localStorage.getItem("token");

  const response = await axios.get<Teacher[]>(
    `${API_URL}/teachers`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    },
  );

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
export const getTeacherById = async (
  id: string,
): Promise<Teacher> => {
  const token = localStorage.getItem("token");

  const response = await axios.get<Teacher>(
    `${API_URL}/teachers/${id}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    },
  );

  return response.data;
};

// UPDATE TEACHER
export const updateTeacher = async (
  id: string,
  teacherData: UpdateTeacherData,
): Promise<Teacher> => {
  const token = localStorage.getItem("token");

  const response = await axios.put<Teacher>(
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

  const response = await axios.delete(
    `${API_URL}/teachers/${id}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    },
  );

  return response.data;
};