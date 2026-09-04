import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import {
  deleteTeacher,
  getTeachers,
} from "../api/teacherApi";

import type { Teacher } from "../api/teacherApi";

function Teachers() {
  const navigate = useNavigate();

  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const loadTeachers = async () => {
    try {
      setLoading(true);
      setError("");

      const data = await getTeachers();

      setTeachers(data);
    } catch (err) {
      console.error("Get teachers error:", err);

      setError("Failed to load teachers");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTeachers();
  }, []);

  const handleDelete = async (id: string) => {
    const confirmed = window.confirm(
      "Are you sure you want to delete this teacher?",
    );

    if (!confirmed) {
      return;
    }

    try {
      await deleteTeacher(id);

      await loadTeachers();
    } catch (err) {
      console.error("Delete teacher error:", err);

      setError("Failed to delete teacher");
    }
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold">
          Teachers
        </h1>

        <button
          onClick={() => navigate("/teachers/add")}
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Add Teacher
        </button>
      </div>

      {error && (
        <div className="mb-4 bg-red-100 text-red-700 p-3 rounded">
          {error}
        </div>
      )}

      {loading ? (
        <p>Loading teachers...</p>
      ) : teachers.length === 0 ? (
        <div className="bg-white shadow rounded-lg p-6">
          <p>No teachers found.</p>
        </div>
      ) : (
        <div className="bg-white shadow rounded-lg overflow-x-auto">
          <table className="w-full border-collapse">
            <thead className="bg-gray-100">
              <tr>
                <th className="border px-4 py-3 text-left">
                  Name
                </th>

                <th className="border px-4 py-3 text-left">
                  Age
                </th>

                <th className="border px-4 py-3 text-left">
                  Gender
                </th>

                <th className="border px-4 py-3 text-left">
                  Subject
                </th>

                <th className="border px-4 py-3 text-left">
                  Email
                </th>

                <th className="border px-4 py-3 text-left">
                  Phone
                </th>

                <th className="border px-4 py-3 text-left">
                  Username
                </th>

                <th className="border px-4 py-3 text-left">
                  Actions
                </th>
              </tr>
            </thead>

            <tbody>
              {teachers.map((teacher) => (
                <tr key={teacher.id}>
                  <td className="border px-4 py-3">
                    {teacher.name}
                  </td>

                  <td className="border px-4 py-3">
                    {teacher.age}
                  </td>

                  <td className="border px-4 py-3">
                    {teacher.gender}
                  </td>

                  <td className="border px-4 py-3">
                    {teacher.subject}
                  </td>

                  <td className="border px-4 py-3">
                    {teacher.email}
                  </td>

                  <td className="border px-4 py-3">
                    {teacher.phone}
                  </td>

                  <td className="border px-4 py-3">
                    {teacher.username}
                  </td>

                  <td className="border px-4 py-3">
                    <div className="flex gap-2">
                      <button
                        onClick={() =>
                          navigate(
                            `/teachers/${teacher.id}/edit`,
                          )
                        }
                        className="bg-yellow-500 text-white px-3 py-1 rounded hover:bg-yellow-600"
                      >
                        Edit
                      </button>

                      <button
                        onClick={() =>
                          handleDelete(teacher.id)
                        }
                        className="bg-red-600 text-white px-3 py-1 rounded hover:bg-red-700"
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default Teachers;