import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import { deleteStudent, getStudents } from "../api/studentApi";

import { getClassrooms, type Classroom } from "../api/classroomApi";

type Student = {
  id: string;
  college_id: string;
  name: string;
  age: number;
  roll_number: number;
  gender: string;
  classroom_id: string;
};

function Students() {
  const [students, setStudents] = useState<Student[]>([]);
  const [classrooms, setClassrooms] = useState<Classroom[]>([]);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const navigate = useNavigate();

  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        setError("");

        const [studentsData, classroomsData] = await Promise.all([
          getStudents(),
          getClassrooms(),
        ]);

        console.log("Students response:", studentsData);

        console.log("Classrooms response:", classroomsData);

        // Students
        if (Array.isArray(studentsData)) {
          setStudents(studentsData);
        } else if (studentsData && Array.isArray(studentsData.students)) {
          setStudents(studentsData.students);
        } else {
          setStudents([]);
        }

        // Classrooms
        if (Array.isArray(classroomsData)) {
          setClassrooms(classroomsData);
        } else {
          setClassrooms([]);
        }
      } catch (err) {
        console.error("Students page error:", err);

        setError("Failed to load students");
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  // ------------------------------------------------
  // ADD STUDENT
  // ------------------------------------------------

  const handleAddStudent = () => {
    navigate("/students/add");
  };

  // ------------------------------------------------
  // EDIT STUDENT
  // ------------------------------------------------

  const handleEditStudent = (studentID: string) => {
    navigate(`/students/${studentID}/edit`);
  };

  // ------------------------------------------------
  // DELETE STUDENT
  // ------------------------------------------------

  const handleDeleteStudent = async (student: Student) => {
    const confirmed = window.confirm(
      `Are you sure you want to delete ${student.name}?`,
    );

    if (!confirmed) {
      return;
    }

    try {
      setError("");

      await deleteStudent(student.id);

      // Remove deleted student from table
      setStudents((currentStudents) =>
        currentStudents.filter((item) => item.id !== student.id),
      );
    } catch (err) {
      console.error("Delete student error:", err);

      setError("Failed to delete student");
    }
  };

  // ------------------------------------------------
  // GET CLASSROOM NAME
  // ------------------------------------------------

  const getClassroomName = (classroomID: string) => {
    const classroom = classrooms.find((item) => item.id === classroomID);

    if (!classroom) {
      return "Classroom not found";
    }

    return `${classroom.name} - Section ${classroom.section}`;
  };

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      {/* PAGE HEADER */}

      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Students</h1>

          <p className="text-gray-600 mt-2">Manage students of your school.</p>
        </div>

        <button
          onClick={handleAddStudent}
          className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2 rounded-lg"
        >
          + Add Student
        </button>
      </div>

      {/* LOADING */}

      {loading && <p className="mt-8 text-gray-500">Loading students...</p>}

      {/* ERROR */}

      {error && <p className="mt-8 text-red-600">{error}</p>}

      {/* STUDENT LIST */}

      {!loading && !error && (
        <div className="bg-white rounded-xl shadow mt-8">
          <div className="p-6 border-b">
            <h2 className="text-xl font-semibold">Student List</h2>

            <p className="text-gray-500 mt-2">
              Total Students: {students.length}
            </p>
          </div>

          {/* NO STUDENTS */}

          {students.length === 0 ? (
            <div className="p-8 text-center text-gray-500">
              No students found.
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="text-left p-4">Roll No</th>

                    <th className="text-left p-4">Name</th>

                    <th className="text-left p-4">Age</th>

                    <th className="text-left p-4">Gender</th>

                    <th className="text-left p-4">Classroom</th>

                    <th className="text-left p-4">Actions</th>
                  </tr>
                </thead>

                <tbody>
                  {students.map((student) => (
                    <tr key={student.id} className="border-t hover:bg-gray-50">
                      <td className="p-4">{student.roll_number}</td>

                      <td className="p-4 font-medium">{student.name}</td>

                      <td className="p-4">{student.age}</td>

                      <td className="p-4">{student.gender}</td>

                      <td className="p-4">
                        {getClassroomName(student.classroom_id)}
                      </td>

                      <td className="p-4">
                        {/* EDIT */}

                        <button
                          type="button"
                          onClick={() => handleEditStudent(student.id)}
                          className="text-blue-600 hover:underline mr-4"
                        >
                          Edit
                        </button>

                        {/* DELETE */}

                        <button
                          type="button"
                          onClick={() => handleDeleteStudent(student)}
                          className="text-red-600 hover:underline"
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default Students;
