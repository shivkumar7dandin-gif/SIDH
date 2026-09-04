import { useEffect, useState, type FormEvent } from "react";

import {
  createClassroom,
  getClassrooms,
  updateClassroom,
  type Classroom,
} from "../api/classroomApi";

function Classrooms() {
  const [classrooms, setClassrooms] = useState<Classroom[]>([]);

  const [name, setName] = useState("");
  const [section, setSection] = useState("");
  const [capacity, setCapacity] = useState("");

  const [editingID, setEditingID] = useState<string | null>(null);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const loadClassrooms = async () => {
    try {
      setLoading(true);
      setError("");

      const data = await getClassrooms();

      if (Array.isArray(data)) {
        setClassrooms(data);
      } else {
        setClassrooms([]);
      }
    } catch (err) {
      console.error("Load classrooms error:", err);
      setError("Failed to load classrooms");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadClassrooms();
  }, []);

  const handleEdit = (classroom: Classroom) => {
    setEditingID(classroom.id);
    setName(classroom.name);
    setSection(classroom.section);
    setCapacity(String(classroom.capacity));
    setError("");
  };

  const handleCancelEdit = () => {
    setEditingID(null);
    setName("");
    setSection("");
    setCapacity("");
    setError("");
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      setError("");

      const classroomData = {
        name,
        section,
        capacity: Number(capacity),
      };

      if (editingID) {
        await updateClassroom(editingID, classroomData);
      } else {
        await createClassroom(classroomData);
      }

      setName("");
      setSection("");
      setCapacity("");
      setEditingID(null);

      await loadClassrooms();
    } catch (err) {
      console.error("Classroom save error:", err);

      if (editingID) {
        setError("Failed to update classroom");
      } else {
        setError("Failed to create classroom");
      }
    }
  };

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Classrooms</h1>

      <div className="bg-white p-6 rounded-lg shadow mb-8">
        <h2 className="text-xl font-semibold mb-4">
          {editingID ? "Edit Classroom" : "Add Classroom"}
        </h2>

        {error && <p className="text-red-600 mb-4">{error}</p>}

        <form
          onSubmit={handleSubmit}
          className="grid grid-cols-1 md:grid-cols-5 gap-4"
        >
          {/* CLASS */}

          <select
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            className="border rounded-lg p-3"
          >
            <option value="">Select Class</option>

            <option value="1st standard">1st standard</option>

            <option value="2nd standard">2nd standard</option>

            <option value="3rd standard">3rd standard</option>

            <option value="4th standard">4th standard</option>

            <option value="5th standard">5th standard</option>

            <option value="6th standard">6th standard</option>

            <option value="7th standard">7th standard</option>

            <option value="8th standard">8th standard</option>

            <option value="9th standard">9th standard</option>

            <option value="10th standard">10th standard</option>

            <option value="11th standard">11th standard</option>

            <option value="12th standard">12th standard</option>
          </select>

          {/* SECTION */}

          <select
            value={section}
            onChange={(e) => setSection(e.target.value)}
            required
            className="border rounded-lg p-3"
          >
            <option value="">Select Section</option>

            <option value="A">Section A</option>

            <option value="B">Section B</option>

            <option value="C">Section C</option>

            <option value="D">Section D</option>
          </select>

          {/* CAPACITY */}

          <input
            type="number"
            placeholder="Capacity"
            value={capacity}
            onChange={(e) => setCapacity(e.target.value)}
            min={1}
            max={60}
            required
            className="border rounded-lg p-3"
          />

          {/* ADD / UPDATE */}

          <button
            type="submit"
            className="bg-blue-600 hover:bg-blue-700 text-white rounded-lg px-5 py-3"
          >
            {editingID ? "Update Classroom" : "+ Add Classroom"}
          </button>

          {/* CANCEL EDIT */}

          {editingID && (
            <button
              type="button"
              onClick={handleCancelEdit}
              className="bg-gray-500 hover:bg-gray-600 text-white rounded-lg px-5 py-3"
            >
              Cancel
            </button>
          )}
        </form>
      </div>

      {/* CLASSROOM TABLE */}

      <div className="bg-white rounded-lg shadow overflow-hidden">
        {loading ? (
          <p className="p-6">Loading classrooms...</p>
        ) : classrooms.length === 0 ? (
          <p className="p-6">No classrooms found.</p>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="text-left p-4">Class</th>

                <th className="text-left p-4">Section</th>

                <th className="text-left p-4">Capacity</th>

                <th className="text-left p-4">Actions</th>
              </tr>
            </thead>

            <tbody>
              {classrooms.map((classroom) => (
                <tr key={classroom.id} className="border-t">
                  <td className="p-4">{classroom.name}</td>

                  <td className="p-4">{classroom.section}</td>

                  <td className="p-4">{classroom.capacity}</td>

                  <td className="p-4">
                    <button
                      type="button"
                      onClick={() => handleEdit(classroom)}
                      className="text-blue-600 hover:underline"
                    >
                      Edit
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

export default Classrooms;
