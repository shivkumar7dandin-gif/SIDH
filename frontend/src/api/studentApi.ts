import axios from "axios";

const API_URL = "http://localhost:8080/api/v1";

export type Student = {
  id: string;
  college_id: string;
  name: string;
  age: number;
  roll_number: number;
  gender: string;
  classroom_id: string;

  address: {
    house_no: string;
    street: string;
    village: string;
    city: string;
    state: string;
    pincode: string;
  };
};

export type CreateStudentData = {
  name: string;
  age: number;
  roll_number: number;
  gender: string;
  classroom_id: string;
  username: string;
  password: string;

  address: {
    house_no: string;
    street: string;
    village: string;
    city: string;
    state: string;
    pincode: string;
  };
};

export type UpdateStudentData = {
  name: string;
  age: number;
  roll_number: number;
  gender: string;
  classroom_id: string;

  address: {
    house_no: string;
    street: string;
    village: string;
    city: string;
    state: string;
    pincode: string;
  };
};

// GET ALL STUDENTS
export const getStudents = async () => {
  const token = localStorage.getItem("token");

  const response = await axios.get(`${API_URL}/students`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data;
};

// GET STUDENT BY ID
export const getStudentByID = async (id: string) => {
  const token = localStorage.getItem("token");

  const response = await axios.get(`${API_URL}/students/${id}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data;
};

// CREATE STUDENT
export const createStudent = async (studentData: CreateStudentData) => {
  const token = localStorage.getItem("token");

  const response = await axios.post(`${API_URL}/students`, studentData, {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  return response.data;
};

// UPDATE STUDENT
export const updateStudent = async (
  id: string,
  studentData: UpdateStudentData,
) => {
  const token = localStorage.getItem("token");

  const response = await axios.put(`${API_URL}/students/${id}`, studentData, {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  return response.data;
};

// DELETE STUDENT
export const deleteStudent = async (id: string) => {
  const token = localStorage.getItem("token");

  const response = await axios.delete(`${API_URL}/students/${id}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data;
};
