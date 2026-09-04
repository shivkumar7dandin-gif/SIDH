import axios from "axios";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

import { createTeacher } from "../api/teacherApi";
import type { CreateTeacherData } from "../api/teacherApi";

function AddTeacher() {
  const navigate = useNavigate();

  const [form, setForm] = useState<CreateTeacherData>({
    name: "",
    age: 0,
    gender: "",
    email: "",
    phone: "",
    subject: "",
    username: "",
    password: "",
  });

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [saving, setSaving] = useState(false);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;

    setForm((prev) => ({
      ...prev,
      [name]: name === "age" ? Number(value) : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      setSaving(true);
      setError("");
      setSuccess("");

      await createTeacher(form);

      // Show success message
      setSuccess("Teacher created successfully");

      // Wait 2.5 seconds before redirect
      setTimeout(() => {
        navigate("/teachers");
      }, 2500);
    } catch (err) {
      console.error("Create teacher error:", err);

      if (axios.isAxiosError(err)) {
        setError(
          err.response?.data?.error ||
            "Failed to create teacher",
        );
      } else {
        setError("Failed to create teacher");
      }
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="p-6">
      <div className="max-w-3xl mx-auto bg-white shadow rounded-lg p-6">
        <h1 className="text-2xl font-semibold mb-6">
          Add Teacher
        </h1>

        {/* ERROR MESSAGE */}
        {error && (
          <div className="mb-4 border border-red-300 bg-red-100 text-red-700 p-4 rounded">
            {error}
          </div>
        )}

        {/* SUCCESS MESSAGE */}
        {success && (
          <div className="mb-4 border border-green-300 bg-green-100 text-green-800 p-4 rounded font-medium">
            ✅ {success}
          </div>
        )}

        <form
          onSubmit={handleSubmit}
          className="space-y-5"
        >
          <div>
            <label className="block mb-1 font-medium">
              Teacher Name
            </label>

            <input
              type="text"
              name="name"
              value={form.name}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
          </div>

          <div>
            <label className="block mb-1 font-medium">
              Age
            </label>

            <input
              type="number"
              name="age"
              value={form.age || ""}
              onChange={handleChange}
              min="1"
              className="w-full border rounded px-3 py-2"
              required
            />
          </div>

          <div>
            <label className="block mb-1 font-medium">
              Gender
            </label>

            <select
              name="gender"
              value={form.gender}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            >
              <option value="">
                Select Gender
              </option>

              <option value="Male">
                Male
              </option>

              <option value="Female">
                Female
              </option>

              <option value="Other">
                Other
              </option>
            </select>
          </div>

          <div>
            <label className="block mb-1 font-medium">
              Email
            </label>

            <input
              type="email"
              name="email"
              value={form.email}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
          </div>

          <div>
            <label className="block mb-1 font-medium">
              Phone
            </label>

            <input
              type="text"
              name="phone"
              value={form.phone}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
          </div>

          <div>
            <label className="block mb-1 font-medium">
              Subject
            </label>

            <input
              type="text"
              name="subject"
              value={form.subject}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
          </div>

          <div>
            <label className="block mb-1 font-medium">
              Username
            </label>

            <input
              type="text"
              name="username"
              value={form.username}
              onChange={handleChange}
              autoComplete="off"
              className="w-full border rounded px-3 py-2"
              required
            />
          </div>

          <div>
            <label className="block mb-1 font-medium">
              Password
            </label>

            <input
              type="password"
              name="password"
              value={form.password}
              onChange={handleChange}
              autoComplete="new-password"
              minLength={8}
              className="w-full border rounded px-3 py-2"
              required
            />
          </div>

          <div className="flex gap-3 pt-3">
            <button
              type="submit"
              disabled={saving}
              className="bg-blue-600 text-white px-5 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
            >
              {saving
                ? "Creating..."
                : "Create Teacher"}
            </button>

            <button
              type="button"
              onClick={() =>
                navigate("/teachers")
              }
              disabled={saving}
              className="bg-gray-200 px-5 py-2 rounded hover:bg-gray-300 disabled:opacity-50"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default AddTeacher;