import { useEffect, useState } from "react";
import { getCollegeDashboard } from "../api/collegeApi";

type College = {
  name: string;
  code: string;
  city: string;
  state: string;
};

type DashboardData = {
  college: College | null;
  total_students: number;
  total_teachers: number;
  total_classrooms: number;
};

function CollegeDashboard() {
  const [dashboard, setDashboard] = useState<DashboardData | null>(null);

  const [error, setError] = useState("");

  useEffect(() => {
    const loadDashboard = async () => {
      try {
        const data = await getCollegeDashboard();

        console.log("Dashboard response:", data);

        setDashboard(data);
      } catch (err) {
        console.error("Dashboard error:", err);

        setError("Failed to load dashboard");
      }
    };

    loadDashboard();
  }, []);

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold">College Admin Dashboard</h1>

      <p className="mt-2 text-gray-600">
        Welcome to your school management dashboard.
      </p>

      {dashboard?.college && (
        <div className="mt-4">
          <h2 className="text-xl font-semibold">{dashboard.college.name}</h2>

          <p className="text-gray-500">
            {dashboard.college.city}, {dashboard.college.state}
          </p>

          <p className="text-sm text-gray-400">
            College Code: {dashboard.college.code}
          </p>
        </div>
      )}

      {error && <p className="mt-4 text-red-600">{error}</p>}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
        <div className="bg-white p-6 rounded-xl shadow">
          <h3 className="text-gray-500 font-medium">Total Students</h3>

          <p className="text-4xl font-bold mt-3">
            {dashboard?.total_students ?? 0}
          </p>
        </div>

        <div className="bg-white p-6 rounded-xl shadow">
          <h3 className="text-gray-500 font-medium">Total Teachers</h3>

          <p className="text-4xl font-bold mt-3">
            {dashboard?.total_teachers ?? 0}
          </p>
        </div>

        <div className="bg-white p-6 rounded-xl shadow">
          <h3 className="text-gray-500 font-medium">Total Classrooms</h3>

          <p className="text-4xl font-bold mt-3">
            {dashboard?.total_classrooms ?? 0}
          </p>
        </div>
      </div>
    </main>
  );
}

export default CollegeDashboard;
