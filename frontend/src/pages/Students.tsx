import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import {
  deleteStudent,
  getStudents,
} from "../api/studentApi";

import {
  getClassrooms,
  type Classroom,
} from "../api/classroomApi";

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
  const navigate = useNavigate();

  const [students, setStudents] = useState<Student[]>([]);
  const [classrooms, setClassrooms] = useState<Classroom[]>([]);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [selectedClassroomID, setSelectedClassroomID] =
    useState<string>("");

  const [currentPage, setCurrentPage] =
    useState(1);

  const studentsPerPage = 20;

  // ------------------------------------------------
  // LOAD STUDENTS + CLASSROOMS
  // ------------------------------------------------

  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        setError("");

        const [studentsData, classroomsData] =
          await Promise.all([
            getStudents(),
            getClassrooms(),
          ]);

        console.log(
          "Students response:",
          studentsData,
        );

        console.log(
          "Classrooms response:",
          classroomsData,
        );

        if (Array.isArray(studentsData)) {
          setStudents(studentsData);
        } else if (
          studentsData &&
          Array.isArray(studentsData.students)
        ) {
          setStudents(studentsData.students);
        } else {
          setStudents([]);
        }

        if (Array.isArray(classroomsData)) {
          setClassrooms(classroomsData);
        } else {
          setClassrooms([]);
        }
      } catch (err) {
        console.error(
          "Students page error:",
          err,
        );

        setError(
          "Failed to load students",
        );
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

  const handleEditStudent = (
    studentID: string,
  ) => {
    navigate(
      `/students/${studentID}/edit`,
    );
  };

  // ------------------------------------------------
  // DELETE STUDENT
  // ------------------------------------------------

  const handleDeleteStudent = async (
    student: Student,
  ) => {
    const confirmed = window.confirm(
      `Are you sure you want to delete ${student.name}?`,
    );

    if (!confirmed) {
      return;
    }

    try {
      setError("");

      await deleteStudent(student.id);

      setStudents((currentStudents) =>
        currentStudents.filter(
          (item) =>
            item.id !== student.id,
        ),
      );
    } catch (err) {
      console.error(
        "Delete student error:",
        err,
      );

      setError(
        "Failed to delete student",
      );
    }
  };

  // ------------------------------------------------
  // GET CLASSROOM NAME
  // ------------------------------------------------

  const getClassroomName = (
    classroomID: string,
  ) => {
    const classroom = classrooms.find(
      (item) =>
        item.id === classroomID,
    );

    if (!classroom) {
      return "Classroom not found";
    }

    return `${classroom.name} - Section ${classroom.section}`;
  };

  // ------------------------------------------------
  // STUDENT COUNT BY CLASSROOM
  // ------------------------------------------------

  const getStudentCount = (
    classroomID: string,
  ) => {
    return students.filter(
      (student) =>
        student.classroom_id === classroomID,
    ).length;
  };

  // ------------------------------------------------
  // SORT CLASSROOMS
  // ------------------------------------------------

  const sortedClassrooms =
    [...classrooms].sort(
      (a, b) => {
        const classA =
          parseInt(a.name) || 0;

        const classB =
          parseInt(b.name) || 0;

        if (classA !== classB) {
          return classA - classB;
        }

        return a.section.localeCompare(
          b.section,
        );
      },
    );

  // ------------------------------------------------
  // FILTER STUDENTS BY CLASSROOM
  // ------------------------------------------------

  const filteredStudents =
    selectedClassroomID
      ? students.filter(
          (student) =>
            student.classroom_id ===
            selectedClassroomID,
        )
      : students;

  // ------------------------------------------------
  // SELECTED CLASSROOM
  // ------------------------------------------------

  const selectedClassroom =
    classrooms.find(
      (classroom) =>
        classroom.id ===
        selectedClassroomID,
    );

  // ------------------------------------------------
  // PAGINATION
  // ------------------------------------------------

  const totalPages = Math.ceil(
    filteredStudents.length /
      studentsPerPage,
  );

  const startIndex =
    (currentPage - 1) *
    studentsPerPage;

  const endIndex =
    startIndex +
    studentsPerPage;

  const paginatedStudents =
    filteredStudents.slice(
      startIndex,
      endIndex,
    );

  // ------------------------------------------------
  // CLASSROOM CLICK
  // ------------------------------------------------

  const handleClassroomClick = (
    classroomID: string,
  ) => {
    setSelectedClassroomID(
      classroomID,
    );

    setCurrentPage(1);
  };

  // ------------------------------------------------
  // ALL STUDENTS
  // ------------------------------------------------

  const handleAllStudents = () => {
    setSelectedClassroomID("");

    setCurrentPage(1);
  };

  return (
    <div className="min-h-screen bg-gray-100 p-8">

      {/* PAGE HEADER */}

      <div className="flex justify-between items-center">

        <div>
          <h1 className="text-3xl font-bold">
            Students
          </h1>

          <p className="text-gray-600 mt-2">
            Manage students of your school.
          </p>
        </div>

        <button
          type="button"
          onClick={handleAddStudent}
          className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2 rounded-lg"
        >
          + Add Student
        </button>

      </div>

      {/* LOADING */}

      {loading && (
        <p className="mt-8 text-gray-500">
          Loading students...
        </p>
      )}

      {/* ERROR */}

      {error && (
        <p className="mt-8 text-red-600">
          {error}
        </p>
      )}

      {!loading && !error && (
        <>

          {/* STUDENTS BY CLASSROOM */}

          <div className="bg-white rounded-xl shadow mt-8 p-6">

            <div className="flex justify-between items-center mb-5">

              <h2 className="text-xl font-semibold">
                Students by Classroom
              </h2>

              <button
                type="button"
                onClick={handleAllStudents}
                className={`px-4 py-2 rounded-lg ${
                  selectedClassroomID === ""
                    ? "bg-blue-600 text-white"
                    : "bg-gray-200 text-gray-700 hover:bg-gray-300"
                }`}
              >
                All Students
              </button>

            </div>

            {sortedClassrooms.length === 0 ? (
              <p className="text-gray-500">
                No classrooms found.
              </p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">

                {sortedClassrooms.map(
                  (classroom) => {
                    const studentCount =
                      getStudentCount(
                        classroom.id,
                      );

                    const isSelected =
                      selectedClassroomID ===
                      classroom.id;

                    return (
                      <div
                        key={classroom.id}
                        onClick={() =>
                          handleClassroomClick(
                            classroom.id,
                          )
                        }
                        className={`border rounded-lg p-4 cursor-pointer transition ${
                          isSelected
                            ? "border-blue-600 bg-blue-50"
                            : "hover:bg-gray-50"
                        }`}
                      >

                        <h3 className="font-semibold text-lg">
                          {classroom.name}
                        </h3>

                        <p className="text-gray-500 mt-1">
                          Section{" "}
                          {classroom.section}
                        </p>

                        <p className="text-2xl font-bold mt-3">
                          {studentCount} /{" "}
                          {classroom.capacity}
                        </p>

                        <p className="text-sm text-gray-500">
                          Students
                        </p>

                      </div>
                    );
                  },
                )}

              </div>
            )}

          </div>

          {/* STUDENT LIST */}

          <div className="bg-white rounded-xl shadow mt-8">

            <div className="p-6 border-b">

              <h2 className="text-xl font-semibold">
                Student List
              </h2>

              {selectedClassroom ? (
                <p className="text-blue-600 mt-2 font-medium">
                  {selectedClassroom.name}
                  {" - Section "}
                  {selectedClassroom.section}
                </p>
              ) : (
                <p className="text-gray-500 mt-2">
                  All Classrooms
                </p>
              )}

              <p className="text-gray-500 mt-1">
                Total Students:{" "}
                {filteredStudents.length}
              </p>

            </div>

            {/* NO STUDENTS */}

            {filteredStudents.length === 0 ? (
              <div className="p-8 text-center text-gray-500">
                No students found in this classroom.
              </div>
            ) : (
              <>

                {/* TABLE */}

                <div className="overflow-x-auto">

                  <table className="w-full">

                    <thead className="bg-gray-50">
                      <tr>

                        <th className="text-left p-4">
                          Roll No
                        </th>

                        <th className="text-left p-4">
                          Name
                        </th>

                        <th className="text-left p-4">
                          Age
                        </th>

                        <th className="text-left p-4">
                          Gender
                        </th>

                        <th className="text-left p-4">
                          Classroom
                        </th>

                        <th className="text-left p-4">
                          Actions
                        </th>

                      </tr>
                    </thead>

                    <tbody>

                      {paginatedStudents.map(
                        (student) => (
                          <tr
                            key={student.id}
                            className="border-t hover:bg-gray-50"
                          >

                            <td className="p-4">
                              {student.roll_number}
                            </td>

                            <td className="p-4 font-medium">
                              {student.name}
                            </td>

                            <td className="p-4">
                              {student.age}
                            </td>

                            <td className="p-4">
                              {student.gender}
                            </td>

                            <td className="p-4">
                              {getClassroomName(
                                student.classroom_id,
                              )}
                            </td>

                            <td className="p-4">

                              <button
                                type="button"
                                onClick={() =>
                                  handleEditStudent(
                                    student.id,
                                  )
                                }
                                className="text-blue-600 hover:underline mr-4"
                              >
                                Edit
                              </button>

                              <button
                                type="button"
                                onClick={() =>
                                  handleDeleteStudent(
                                    student,
                                  )
                                }
                                className="text-red-600 hover:underline"
                              >
                                Delete
                              </button>

                            </td>

                          </tr>
                        ),
                      )}

                    </tbody>

                  </table>

                </div>

                {/* PAGINATION */}

                <div className="flex justify-between items-center p-5 border-t">

                  <p className="text-gray-600">

                    Showing{" "}
                    {startIndex + 1}
                    {" - "}
                    {Math.min(
                      endIndex,
                      filteredStudents.length,
                    )}
                    {" of "}
                    {filteredStudents.length}
                    {" students"}

                  </p>

                  <div className="flex items-center gap-3">

                    <button
                      type="button"
                      disabled={
                        currentPage === 1
                      }
                      onClick={() =>
                        setCurrentPage(
                          (page) =>
                            page - 1,
                        )
                      }
                      className="px-4 py-2 border rounded-lg hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed"
                    >
                      Previous
                    </button>

                    <span className="text-gray-700">
                      Page{" "}
                      {currentPage}
                      {" of "}
                      {totalPages || 1}
                    </span>

                    <button
                      type="button"
                      disabled={
                        currentPage >=
                        totalPages
                      }
                      onClick={() =>
                        setCurrentPage(
                          (page) =>
                            page + 1,
                        )
                      }
                      className="px-4 py-2 border rounded-lg hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed"
                    >
                      Next
                    </button>

                  </div>

                </div>

              </>
            )}

          </div>

        </>
      )}

    </div>
  );
}

export default Students;