import axios from "axios";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { getTeacherById, updateTeacher } from "../api/teacherApi";

import type { UpdateTeacherData } from "../api/teacherApi";

function EditTeacher() {
  const navigate = useNavigate();
  const { id } = useParams();

  const [form, setForm] = useState<UpdateTeacherData>({
    employee_id: "",
    name: "",
    age: 0,
    gender: "",
    email: "",
    phone: "",
    subject: "",
    qualification: "",
    experience_years: 0,
    joining_date: "",
    address: {
      house_no: "",
      street: "",
      village: "",
      city: "",
      state: "",
      pincode: "",
    },
  });

  const [loading, setLoading] = useState(true);

  const [saving, setSaving] = useState(false);

  const [error, setError] = useState("");

  const [success, setSuccess] = useState("");

  const inputClass =
    "w-full border border-gray-300 rounded-lg px-3 py-2.5 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500";

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

        const teacher = await getTeacherById(id);

        setForm({
          employee_id: teacher.employee_id || "",

          name: teacher.name || "",

          age: teacher.age || 0,

          gender: teacher.gender || "",

          email: teacher.email || "",

          phone: teacher.phone || "",

          subject: teacher.subject || "",

          qualification: teacher.qualification || "",

          experience_years: teacher.experience_years || 0,

          joining_date: teacher.joining_date || "",

          address: {
            house_no: teacher.address?.house_no || "",

            street: teacher.address?.street || "",

            village: teacher.address?.village || "",

            city: teacher.address?.city || "",

            state: teacher.address?.state || "",

            pincode: teacher.address?.pincode || "",
          },
        });
      } catch (err) {
        console.error("Get teacher error:", err);

        setError("Failed to load teacher");
      } finally {
        setLoading(false);
      }
    };

    loadTeacher();
  }, [id]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;

    setForm((prev) => ({
      ...prev,

      [name]:
        name === "age" || name === "experience_years" ? Number(value) : value,
    }));
  };

  const handleAddressChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    setForm((prev) => ({
      ...prev,

      address: {
        ...prev.address,
        [name]: value,
      },
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!id) {
      return;
    }

    try {
      setSaving(true);
      setError("");
      setSuccess("");

      await updateTeacher(id, form);

      setSuccess("Teacher profile updated successfully");

      setTimeout(() => {
        navigate("/teachers");
      }, 1500);
    } catch (err) {
      console.error("Update teacher error:", err);

      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Failed to update teacher");
      } else {
        setError("Failed to update teacher");
      }
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <div className="p-6">Loading teacher profile...</div>;
  }

  return (
    <div className="p-6">
      <div className="max-w-5xl mx-auto">
        {/* HEADER */}
        <div className="mb-6">
          <button
            type="button"
            onClick={() => navigate("/teachers")}
            className="text-blue-600 hover:text-blue-800 mb-3"
          >
            ← Back to Teachers
          </button>

          <h1 className="text-3xl font-semibold text-gray-900">Edit Teacher</h1>

          <p className="text-gray-500 mt-1">
            Update teacher personal and professional information.
          </p>
        </div>

        {error && (
          <div className="mb-5 border border-red-300 bg-red-50 text-red-700 p-4 rounded-lg">
            {error}
          </div>
        )}

        {success && (
          <div className="mb-5 border border-green-300 bg-green-50 text-green-800 p-4 rounded-lg font-medium">
            ✓ {success}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* TEACHER INFORMATION */}
          <section className="bg-white shadow-sm border rounded-xl p-6">
            <div className="mb-5">
              <h2 className="text-lg font-semibold">Teacher Information</h2>

              <p className="text-sm text-gray-500">
                Basic information about the teacher.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div>
                <label className="block mb-1.5 font-medium">
                  Employee ID *
                </label>

                <input
                  type="text"
                  name="employee_id"
                  value={form.employee_id}
                  onChange={handleChange}
                  placeholder="Example: TCH-001"
                  className={inputClass}
                  required
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">
                  Teacher Name *
                </label>

                <input
                  type="text"
                  name="name"
                  value={form.name}
                  onChange={handleChange}
                  className={inputClass}
                  required
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">Age *</label>

                <input
                  type="number"
                  name="age"
                  value={form.age || ""}
                  onChange={handleChange}
                  min="1"
                  className={inputClass}
                  required
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">Gender *</label>

                <select
                  name="gender"
                  value={form.gender}
                  onChange={handleChange}
                  className={inputClass}
                  required
                >
                  <option value="">Select Gender</option>

                  <option value="Male">Male</option>

                  <option value="Female">Female</option>

                  <option value="Other">Other</option>
                </select>
              </div>
            </div>
          </section>

          {/* PROFESSIONAL INFORMATION */}
          <section className="bg-white shadow-sm border rounded-xl p-6">
            <div className="mb-5">
              <h2 className="text-lg font-semibold">
                Professional Information
              </h2>

              <p className="text-sm text-gray-500">
                Employment and teaching details.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div>
                <label className="block mb-1.5 font-medium">Subject *</label>

                <input
                  type="text"
                  name="subject"
                  value={form.subject}
                  onChange={handleChange}
                  className={inputClass}
                  required
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">
                  Qualification *
                </label>

                <input
                  type="text"
                  name="qualification"
                  value={form.qualification}
                  onChange={handleChange}
                  placeholder="Example: M.Sc, B.Ed"
                  className={inputClass}
                  required
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">
                  Experience (Years) *
                </label>

                <input
                  type="number"
                  name="experience_years"
                  value={form.experience_years}
                  onChange={handleChange}
                  min="0"
                  className={inputClass}
                  required
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">
                  Joining Date *
                </label>

                <input
                  type="date"
                  name="joining_date"
                  value={form.joining_date}
                  onChange={handleChange}
                  className={inputClass}
                  required
                />
              </div>
            </div>
          </section>

          {/* CONTACT INFORMATION */}
          <section className="bg-white shadow-sm border rounded-xl p-6">
            <div className="mb-5">
              <h2 className="text-lg font-semibold">Contact Information</h2>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div>
                <label className="block mb-1.5 font-medium">Email *</label>

                <input
                  type="email"
                  name="email"
                  value={form.email}
                  onChange={handleChange}
                  className={inputClass}
                  required
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">Phone *</label>

                <input
                  type="text"
                  name="phone"
                  value={form.phone}
                  onChange={handleChange}
                  className={inputClass}
                  required
                />
              </div>
            </div>
          </section>

          {/* ADDRESS */}
          <section className="bg-white shadow-sm border rounded-xl p-6">
            <div className="mb-5">
              <h2 className="text-lg font-semibold">Address</h2>

              <p className="text-sm text-gray-500">
                City, state and pincode are required.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div>
                <label className="block mb-1.5 font-medium">House No</label>

                <input
                  type="text"
                  name="house_no"
                  value={form.address.house_no}
                  onChange={handleAddressChange}
                  className={inputClass}
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">Street</label>

                <input
                  type="text"
                  name="street"
                  value={form.address.street}
                  onChange={handleAddressChange}
                  className={inputClass}
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">Village</label>

                <input
                  type="text"
                  name="village"
                  value={form.address.village}
                  onChange={handleAddressChange}
                  className={inputClass}
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">City *</label>

                <input
                  type="text"
                  name="city"
                  value={form.address.city}
                  onChange={handleAddressChange}
                  className={inputClass}
                  required
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">State *</label>

                <input
                  type="text"
                  name="state"
                  value={form.address.state}
                  onChange={handleAddressChange}
                  className={inputClass}
                  required
                />
              </div>

              <div>
                <label className="block mb-1.5 font-medium">Pincode *</label>

                <input
                  type="text"
                  name="pincode"
                  value={form.address.pincode}
                  onChange={handleAddressChange}
                  className={inputClass}
                  required
                />
              </div>
            </div>
          </section>

          {/* ACTIONS */}
          <div className="bg-white border shadow-sm rounded-xl p-5 flex justify-end gap-3">
            <button
              type="button"
              onClick={() => navigate("/teachers")}
              disabled={saving}
              className="bg-gray-100 text-gray-700 px-5 py-2.5 rounded-lg hover:bg-gray-200 disabled:opacity-50"
            >
              Cancel
            </button>

            <button
              type="submit"
              disabled={saving}
              className="bg-blue-600 text-white px-6 py-2.5 rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {saving ? "Updating..." : "Update Teacher"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default EditTeacher;
