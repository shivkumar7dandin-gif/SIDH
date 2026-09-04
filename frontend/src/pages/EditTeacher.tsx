import axios from "axios";
import { useEffect, useState } from "react";
import {
  useNavigate,
  useParams,
} from "react-router-dom";

import {
  getTeacherById,
  updateTeacher,
} from "../api/teacherApi";

import type { UpdateTeacherData } from "../api/teacherApi";

function EditTeacher() {
  const navigate = useNavigate();
  const { id } = useParams();

  const [form, setForm] =
    useState<UpdateTeacherData>({
      name: "",
      age: 0,
      gender: "",
      email: "",
      phone: "",
      subject: "",
    });

  const [loading, setLoading] =
    useState(true);

  const [saving, setSaving] =
    useState(false);

  const [error, setError] =
    useState("");

  const [success, setSuccess] =
    useState("");

  useEffect(() => {
    const loadTeacher = async () => {
      if (!id) {
        setError("Teacher ID not found");
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        setError("");

        const teacher =
          await getTeacherById(id);

        setForm({
          name: teacher.name,
          age: teacher.age,
          gender: teacher.gender,
          email: teacher.email,
          phone: teacher.phone,
          subject: teacher.subject,
        });
      } catch (err) {
        console.error(
          "Get teacher error:",
          err,
        );

        setError(
          "Failed to load teacher",
        );
      } finally {
        setLoading(false);
      }
    };

    loadTeacher();
  }, [id]);

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement |
        HTMLSelectElement
    >,
  ) => {
    const { name, value } =
      e.target;

    setForm((prev) => ({
      ...prev,
      [name]:
        name === "age"
          ? Number(value)
          : value,
    }));
  };

  const handleSubmit = async (
    e: React.FormEvent,
  ) => {
    e.preventDefault();

    if (!id) {
      return;
    }

    try {
      setSaving(true);
      setError("");
      setSuccess("");

      await updateTeacher(
        id,
        form,
      );

      setSuccess(
        "Teacher updated successfully",
      );

      setTimeout(() => {
        navigate("/teachers");
      }, 2000);
    } catch (err) {
      console.error(
        "Update teacher error:",
        err,
      );

      if (
        axios.isAxiosError(err)
      ) {
        setError(
          err.response?.data?.error ||
            "Failed to update teacher",
        );
      } else {
        setError(
          "Failed to update teacher",
        );
      }
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        Loading teacher...
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="max-w-3xl mx-auto bg-white shadow rounded-lg p-6">
        <h1 className="text-2xl font-semibold mb-6">
          Edit Teacher
        </h1>

        {error && (
          <div className="mb-4 bg-red-100 text-red-700 p-3 rounded">
            {error}
          </div>
        )}

        {success && (
          <div className="mb-4 bg-green-100 text-green-800 p-3 rounded">
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

          <div className="flex gap-3 pt-3">
            <button
              type="submit"
              disabled={saving}
              className="bg-blue-600 text-white px-5 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
            >
              {saving
                ? "Updating..."
                : "Update Teacher"}
            </button>

            <button
              type="button"
              onClick={() =>
                navigate("/teachers")
              }
              className="bg-gray-200 px-5 py-2 rounded hover:bg-gray-300"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default EditTeacher;