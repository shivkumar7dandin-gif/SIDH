import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { getStudentByID, updateStudent } from "../api/studentApi";

import { getClassrooms, type Classroom } from "../api/classroomApi";

function EditStudent() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [classrooms, setClassrooms] = useState<Classroom[]>([]);

  const [name, setName] = useState("");
  const [age, setAge] = useState(0);
  const [rollNumber, setRollNumber] = useState(0);

  const [gender, setGender] = useState("");

  const [classroomID, setClassroomID] = useState("");

  const [houseNo, setHouseNo] = useState("");

  const [street, setStreet] = useState("");

  const [village, setVillage] = useState("");

  const [city, setCity] = useState("");

  const [state, setState] = useState("");

  const [pincode, setPincode] = useState("");

  const [loading, setLoading] = useState(true);

  const [saving, setSaving] = useState(false);

  const [error, setError] = useState("");

  // ------------------------------------------------
  // LOAD STUDENT + CLASSROOMS
  // ------------------------------------------------

  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        setError("");

        if (!id) {
          setError("Invalid student id");
          return;
        }

        const [studentData, classroomData] = await Promise.all([
          getStudentByID(id),
          getClassrooms(),
        ]);

        console.log("Student:", studentData);

        console.log("Classrooms:", classroomData);

        setName(studentData.name || "");
        setAge(studentData.age || 0);

        setRollNumber(studentData.roll_number || 0);

        setGender(studentData.gender || "");

        setClassroomID(studentData.classroom_id || "");

        setHouseNo(studentData.address?.house_no || "");

        setStreet(studentData.address?.street || "");

        setVillage(studentData.address?.village || "");

        setCity(studentData.address?.city || "");

        setState(studentData.address?.state || "");

        setPincode(studentData.address?.pincode || "");

        if (Array.isArray(classroomData)) {
          setClassrooms(classroomData);
        } else {
          setClassrooms([]);
        }
      } catch (err) {
        console.error("Load student error:", err);

        setError("Failed to load student");
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [id]);

  // ------------------------------------------------
  // UPDATE STUDENT
  // ------------------------------------------------

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!id) {
      return;
    }

    try {
      setSaving(true);
      setError("");

      await updateStudent(id, {
        name,
        age,
        roll_number: rollNumber,
        gender,
        classroom_id: classroomID,

        address: {
          house_no: houseNo,
          street,
          village,
          city,
          state,
          pincode,
        },
      });

      navigate("/students");
    } catch (err) {
      console.error("Update student error:", err);

      setError("Failed to update student");
    } finally {
      setSaving(false);
    }
  };

  // ------------------------------------------------
  // LOADING
  // ------------------------------------------------

  if (loading) {
    return <div className="p-8">Loading student...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-6">
          <h1 className="text-3xl font-bold">Edit Student</h1>

          <p className="text-gray-600 mt-2">Update student information.</p>
        </div>

        {error && (
          <div className="bg-red-100 text-red-700 p-4 rounded-lg mb-6">
            {error}
          </div>
        )}

        <form
          onSubmit={handleSubmit}
          className="bg-white rounded-xl shadow p-6"
        >
          {/* BASIC DETAILS */}

          <h2 className="text-xl font-semibold mb-4">Student Details</h2>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            {/* NAME */}

            <div>
              <label className="block mb-2 font-medium">Student Name</label>

              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full border rounded-lg px-4 py-2"
                required
              />
            </div>

            {/* AGE */}

            <div>
              <label className="block mb-2 font-medium">Age</label>

              <input
                type="number"
                value={age}
                onChange={(e) => setAge(Number(e.target.value))}
                className="w-full border rounded-lg px-4 py-2"
                required
              />
            </div>

            {/* ROLL NUMBER */}

            <div>
              <label className="block mb-2 font-medium">Roll Number</label>

              <input
                type="number"
                value={rollNumber}
                onChange={(e) => setRollNumber(Number(e.target.value))}
                className="w-full border rounded-lg px-4 py-2"
                required
              />
            </div>

            {/* GENDER */}

            <div>
              <label className="block mb-2 font-medium">Gender</label>

              <select
                value={gender}
                onChange={(e) => setGender(e.target.value)}
                className="w-full border rounded-lg px-4 py-2"
                required
              >
                <option value="">Select Gender</option>

                <option value="Male">Male</option>

                <option value="Female">Female</option>

                <option value="Other">Other</option>
              </select>
            </div>

            {/* CLASSROOM */}

            <div className="md:col-span-2">
              <label className="block mb-2 font-medium">Classroom</label>

              <select
                value={classroomID}
                onChange={(e) => setClassroomID(e.target.value)}
                className="w-full border rounded-lg px-4 py-2"
                required
              >
                <option value="">Select Classroom</option>

                {classrooms.map((classroom) => (
                  <option key={classroom.id} value={classroom.id}>
                    {classroom.name}
                    {" - Section "}
                    {classroom.section}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* ADDRESS */}

          <h2 className="text-xl font-semibold mt-8 mb-4">Address</h2>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div>
              <label className="block mb-2 font-medium">House No</label>

              <input
                type="text"
                value={houseNo}
                onChange={(e) => setHouseNo(e.target.value)}
                className="w-full border rounded-lg px-4 py-2"
              />
            </div>

            <div>
              <label className="block mb-2 font-medium">Street</label>

              <input
                type="text"
                value={street}
                onChange={(e) => setStreet(e.target.value)}
                className="w-full border rounded-lg px-4 py-2"
              />
            </div>

            <div>
              <label className="block mb-2 font-medium">Village</label>

              <input
                type="text"
                value={village}
                onChange={(e) => setVillage(e.target.value)}
                className="w-full border rounded-lg px-4 py-2"
              />
            </div>

            <div>
              <label className="block mb-2 font-medium">City</label>

              <input
                type="text"
                value={city}
                onChange={(e) => setCity(e.target.value)}
                className="w-full border rounded-lg px-4 py-2"
              />
            </div>

            <div>
              <label className="block mb-2 font-medium">State</label>

              <input
                type="text"
                value={state}
                onChange={(e) => setState(e.target.value)}
                className="w-full border rounded-lg px-4 py-2"
              />
            </div>

            <div>
              <label className="block mb-2 font-medium">Pincode</label>

              <input
                type="text"
                value={pincode}
                onChange={(e) => setPincode(e.target.value)}
                className="w-full border rounded-lg px-4 py-2"
              />
            </div>
          </div>

          {/* BUTTONS */}

          <div className="flex gap-4 mt-8">
            <button
              type="submit"
              disabled={saving}
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg disabled:opacity-50"
            >
              {saving ? "Updating..." : "Update Student"}
            </button>

            <button
              type="button"
              onClick={() => navigate("/students")}
              className="bg-gray-300 hover:bg-gray-400 px-6 py-2 rounded-lg"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default EditStudent;
