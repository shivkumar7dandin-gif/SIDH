import axios from "axios";

const API_URL = "http://localhost:8080/api/v1";

export const getCollegeDashboard = async () => {
  const token = localStorage.getItem("token");

  try {
    const response = await axios.get(`${API_URL}/colleges/me/dashboard`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    return response.data;
  } catch (error: any) {
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("role");

      window.location.reload();
    }

    throw error;
  }
};
