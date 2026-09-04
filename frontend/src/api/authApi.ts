import axios from "axios";

const API_URL = "http://localhost:8080/api/v1";

export const loginUser = async (username: string, password: string) => {
  const response = await axios.post(`${API_URL}/auth/login`, {
    username,
    password,
  });

  return response.data;
};
