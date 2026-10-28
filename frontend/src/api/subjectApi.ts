import axios from "axios";

const API_URL = "http://localhost:8080/api/v1";

export type Subject = {
  id: string;
  college_id: string;
  name: string;
  code: string;
  status: string;
  created_at: string;
  updated_at: string;
};

export type CreateSubjectData = {
  name: string;
  code: string;
};

export type DeleteSubjectResponse = {
  message: string;
};

const getAuthHeaders = () => {
  const token = localStorage.getItem("token");

  return {
    Authorization: `Bearer ${token}`,
  };
};

export const getSubjects = async (): Promise<Subject[]> => {
  const response = await axios.get<Subject[]>(`${API_URL}/subjects`, {
    headers: getAuthHeaders(),
  });

  return response.data;
};

export const createSubject = async (
  data: CreateSubjectData,
): Promise<Subject> => {
  const response = await axios.post<Subject>(`${API_URL}/subjects`, data, {
    headers: {
      ...getAuthHeaders(),
      "Content-Type": "application/json",
    },
  });

  return response.data;
};

export const deleteSubject = async (
  id: string,
): Promise<DeleteSubjectResponse> => {
  const response = await axios.delete<DeleteSubjectResponse>(
    `${API_URL}/subjects/${id}`,
    {
      headers: getAuthHeaders(),
    },
  );

  return response.data;
};
