import axios from "axios";

const API_URL = "http://localhost:8080/api/v1";

export type Classroom = {
  id: string;
  college_id: string;
  name: string;
  section: string;
  capacity: number;
};

export type CreateClassroomData = {
  name: string;
  section: string;
  capacity: number;
};

// GET CLASSROOMS
export const getClassrooms = async () => {
  const token = localStorage.getItem("token");

  const response = await axios.get(`${API_URL}/classrooms`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.data;
};

// CREATE CLASSROOM
export const createClassroom = async (classroomData: CreateClassroomData) => {
  const token = localStorage.getItem("token");

  const response = await axios.post(`${API_URL}/classrooms`, classroomData, {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  return response.data;
};

export const updateClassroom = async (
  id: string,
  classroomData: CreateClassroomData,
) => {
  const token = localStorage.getItem("token");

  const response = await axios.put(
    `${API_URL}/classrooms/${id}`,
    classroomData,
    {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    },
  );

  return response.data;
};
