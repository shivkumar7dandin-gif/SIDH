import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createStudent } from "../api/studentApi";
import { getClassrooms } from "../api/classroomApi";

type Classroom = {
  id: string;
  name: string;
  section: string;
  capacity: number;
};

function AddStudent() {
  const navigate = useNavigate();

  const [classrooms, setClassrooms] = useState<Classroom[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const [form, setForm] = useState({
    name: "",
    age: "",
    roll_number: "",
    gender: "",
    classroom_id: "",
    username: "",
    password: "",

    address: {
      house_no: "",
      street: "",
      village: "",
      city: "",
      state: "",
      pincode: "",
    },
  });

  useEffect(() => {
    const loadClassrooms = async () => {
      try {
        const data = await getClassrooms();

        console.log("Classrooms response:", data);

        if (Array.isArray(data)) {
          setClassrooms(data);
        } else if (Array.isArray(data?.classrooms)) {
          setClassrooms(data.classrooms);
        } else {
          setClassrooms([]);
        }
      } catch (err) {
        console.error("Failed to load classrooms:", err);
        setError("Failed to load classrooms");
      }
    };

    loadClassrooms();
  }, []);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;

    setForm((previous) => ({
      ...previous,
      [name]: value,
    }));
  };

  const handleAddressChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    setForm((previous) => ({
      ...previous,

      address: {
        ...previous.address,
        [name]: value,
      },
    }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setLoading(true);
    setError("");

    try {
      const studentData = {
        name: form.name,
        age: Number(form.age),
        roll_number: Number(form.roll_number),
        gender: form.gender,
        classroom_id: form.classroom_id,
        username: form.username,
        password: form.password,
        address: form.address,
      };

      console.log("Creating student:", studentData);

      await createStudent(studentData);

      navigate("/students");
    } catch (err) {
      console.error("Create student error:", err);

      setError("Failed to create student");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      {/* PAGE TITLE */}
      <h1 className="text-3xl font-bold">Add Student</h1>

      <p className="text-gray-600 mt-2">Register a new student.</p>

      {/* ERROR */}
      {error && (
        <div className="mt-5 bg-red-100 text-red-700 p-3 rounded-lg">
          {error}
        </div>
      )}

      {/* FORM */}
      <form
        onSubmit={handleSubmit}
        className="bg-white rounded-xl shadow mt-8 p-8 max-w-5xl"
      >
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* STUDENT NAME */}
          <div>
            <label className="block font-medium mb-2">Student Name</label>

            <input
              type="text"
              name="name"
              value={form.name}
              onChange={handleChange}
              placeholder="Enter student name"
              required
              className="w-full border rounded-lg p-3"
            />
          </div>

          {/* AGE */}
          <div>
            <label className="block font-medium mb-2">Age</label>

            <input
              type="number"
              name="age"
              value={form.age}
              onChange={handleChange}
              placeholder="Enter age"
              required
              className="w-full border rounded-lg p-3"
            />
          </div>

          {/* ROLL NUMBER */}
          <div>
            <label className="block font-medium mb-2">Roll Number</label>

            <input
              type="number"
              name="roll_number"
              value={form.roll_number}
              onChange={handleChange}
              placeholder="Enter roll number"
              required
              className="w-full border rounded-lg p-3"
            />
          </div>

          {/* GENDER */}
          <div>
            <label className="block font-medium mb-2">Gender</label>

            <select
              name="gender"
              value={form.gender}
              onChange={handleChange}
              required
              className="w-full border rounded-lg p-3"
            >
              <option value="">Select gender</option>

              <option value="Male">Male</option>

              <option value="Female">Female</option>

              <option value="Other">Other</option>
            </select>
          </div>

          {/* CLASSROOM */}
          <div>
            <label className="block font-medium mb-2">Classroom</label>

            <select
              name="classroom_id"
              value={form.classroom_id}
              onChange={handleChange}
              required
              className="w-full border rounded-lg p-3"
            >
              <option value="">Select classroom</option>

              {classrooms.map((classroom) => (
                <option key={classroom.id} value={classroom.name}>
                  {classroom.name}
                  {classroom.section ? ` - Section ${classroom.section}` : ""}
                </option>
              ))}
            </select>
          </div>

          {/* USERNAME */}
          <div>
            <label className="block font-medium mb-2">Username</label>

            <input
              type="text"
              name="username"
              value={form.username}
              onChange={handleChange}
              placeholder="Create student username"
              required
              className="w-full border rounded-lg p-3"
            />
          </div>

          {/* PASSWORD */}
          <div>
            <label className="block font-medium mb-2">Password</label>

            <input
              type="password"
              name="password"
              value={form.password}
              onChange={handleChange}
              placeholder="Create student password"
              required
              className="w-full border rounded-lg p-3"
            />
          </div>
        </div>

        {/* ADDRESS */}
        <h2 className="text-xl font-semibold mt-8 mb-5">Address</h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block font-medium mb-2">House No</label>

            <input
              type="text"
              name="house_no"
              value={form.address.house_no}
              onChange={handleAddressChange}
              className="w-full border rounded-lg p-3"
            />
          </div>

          <div>
            <label className="block font-medium mb-2">Street</label>

            <input
              type="text"
              name="street"
              value={form.address.street}
              onChange={handleAddressChange}
              className="w-full border rounded-lg p-3"
            />
          </div>

          <div>
            <label className="block font-medium mb-2">Village</label>

            <input
              type="text"
              name="village"
              value={form.address.village}
              onChange={handleAddressChange}
              className="w-full border rounded-lg p-3"
            />
          </div>

          <div>
            <label className="block font-medium mb-2">City</label>

            <input
              type="text"
              name="city"
              value={form.address.city}
              onChange={handleAddressChange}
              className="w-full border rounded-lg p-3"
            />
          </div>

          <div>
            <label className="block font-medium mb-2">State</label>

            <input
              type="text"
              name="state"
              value={form.address.state}
              onChange={handleAddressChange}
              className="w-full border rounded-lg p-3"
            />
          </div>

          <div>
            <label className="block font-medium mb-2">Pincode</label>

            <input
              type="text"
              name="pincode"
              value={form.address.pincode}
              onChange={handleAddressChange}
              className="w-full border rounded-lg p-3"
            />
          </div>
        </div>

        {/* BUTTONS */}
        <div className="flex gap-4 mt-8">
          <button
            type="submit"
            disabled={loading}
            className="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-6 py-3 rounded-lg"
          >
            {loading ? "Creating..." : "Create Student"}
          </button>

          <button
            type="button"
            onClick={() => navigate("/students")}
            className="bg-gray-500 hover:bg-gray-600 text-white px-6 py-3 rounded-lg"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}

export default AddStudent;
