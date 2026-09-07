import { useEffect, useState } from "react";

import { getCollegeDashboard } from "../api/collegeApi";
import { getStudents } from "../api/studentApi";
import {
  getTeachers,
  type Teacher,
} from "../api/teacherApi";
import {
  getClassrooms,
  type Classroom,
} from "../api/classroomApi";

type College = {
  name: string;
  code: string;
  city: string;
  state: string;
};

type Student = {
  id: string;
  college_id: string;
  name: string;
  age: number;
  roll_number: number;
  gender: string;
  classroom_id: string;
};

type DashboardData = {
  college: College | null;
  total_students: number;
  total_teachers: number;
  total_classrooms: number;
};

type SelectedView =
  | "students"
  | "teachers"
  | "classrooms"
  | null;

function CollegeDashboard() {
  const [dashboard, setDashboard] =
    useState<DashboardData | null>(null);

  const [students, setStudents] =
    useState<Student[]>([]);

  const [teachers, setTeachers] =
    useState<Teacher[]>([]);

  const [classrooms, setClassrooms] =
    useState<Classroom[]>([]);

  const [selectedView, setSelectedView] =
    useState<SelectedView>(null);

  const [currentPage, setCurrentPage] =
    useState(1);

  const [loadingDetails, setLoadingDetails] =
    useState(false);

  const [error, setError] = useState("");

  const itemsPerPage = 20;

  // ------------------------------------------------
  // LOAD DASHBOARD TOTALS
  // ------------------------------------------------

  useEffect(() => {
    const loadDashboard = async () => {
      try {
        setError("");

        const data =
          await getCollegeDashboard();

        console.log(
          "Dashboard response:",
          data,
        );

        setDashboard(data);
      } catch (err) {
        console.error(
          "Dashboard error:",
          err,
        );

        setError(
          "Failed to load dashboard",
        );
      }
    };

    loadDashboard();
  }, []);

  // ------------------------------------------------
  // CLICK STUDENTS
  // ------------------------------------------------

  const handleStudentsClick =
    async () => {
      try {
        setLoadingDetails(true);
        setError("");
        setCurrentPage(1);
        setSelectedView("students");

        const data =
          await getStudents();

        if (Array.isArray(data)) {
          setStudents(data);
        } else if (
          data &&
          Array.isArray(data.students)
        ) {
          setStudents(data.students);
        } else {
          setStudents([]);
        }

        // Need classroom details to show classroom name
        const classroomData =
          await getClassrooms();

        if (
          Array.isArray(
            classroomData,
          )
        ) {
          setClassrooms(
            classroomData,
          );
        }
      } catch (err) {
        console.error(
          "Load students error:",
          err,
        );

        setError(
          "Failed to load students",
        );
      } finally {
        setLoadingDetails(false);
      }
    };

  // ------------------------------------------------
  // CLICK TEACHERS
  // ------------------------------------------------

  const handleTeachersClick =
    async () => {
      try {
        setLoadingDetails(true);
        setError("");
        setCurrentPage(1);
        setSelectedView("teachers");

        const data =
          await getTeachers();

        setTeachers(data);
      } catch (err) {
        console.error(
          "Load teachers error:",
          err,
        );

        setError(
          "Failed to load teachers",
        );
      } finally {
        setLoadingDetails(false);
      }
    };

  // ------------------------------------------------
  // CLICK CLASSROOMS
  // ------------------------------------------------

  const handleClassroomsClick =
    async () => {
      try {
        setLoadingDetails(true);
        setError("");
        setCurrentPage(1);
        setSelectedView(
          "classrooms",
        );

        const data =
          await getClassrooms();

        if (Array.isArray(data)) {
          setClassrooms(data);
        } else {
          setClassrooms([]);
        }
      } catch (err) {
        console.error(
          "Load classrooms error:",
          err,
        );

        setError(
          "Failed to load classrooms",
        );
      } finally {
        setLoadingDetails(false);
      }
    };

  // ------------------------------------------------
  // CLASSROOM NAME
  // ------------------------------------------------

  const getClassroomName = (
    classroomID: string,
  ) => {
    const classroom =
      classrooms.find(
        (item) =>
          item.id === classroomID,
      );

    if (!classroom) {
      return "Classroom not found";
    }

    return `${classroom.name} - Section ${classroom.section}`;
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

        if (
          classA !== classB
        ) {
          return (
            classA - classB
          );
        }

        return a.section.localeCompare(
          b.section,
        );
      },
    );

  // ------------------------------------------------
  // PAGINATION DATA
  // ------------------------------------------------

  let totalItems = 0;

  if (
    selectedView ===
    "students"
  ) {
    totalItems =
      students.length;
  }

  if (
    selectedView ===
    "teachers"
  ) {
    totalItems =
      teachers.length;
  }

  if (
    selectedView ===
    "classrooms"
  ) {
    totalItems =
      sortedClassrooms.length;
  }

  const totalPages =
    Math.ceil(
      totalItems /
        itemsPerPage,
    );

  const startIndex =
    (currentPage - 1) *
    itemsPerPage;

  const endIndex =
    startIndex +
    itemsPerPage;

  const paginatedStudents =
    students.slice(
      startIndex,
      endIndex,
    );

  const paginatedTeachers =
    teachers.slice(
      startIndex,
      endIndex,
    );

  const paginatedClassrooms =
    sortedClassrooms.slice(
      startIndex,
      endIndex,
    );

  return (
    <main className="p-8">

      <h1 className="text-3xl font-bold">
        College Admin Dashboard
      </h1>

      <p className="mt-2 text-gray-600">
        Welcome to your school management dashboard.
      </p>

      {dashboard?.college && (
        <div className="mt-4">

          <h2 className="text-xl font-semibold">
            {dashboard.college.name}
          </h2>

          <p className="text-gray-500">
            {dashboard.college.city},{" "}
            {dashboard.college.state}
          </p>

          <p className="text-sm text-gray-400">
            College Code:{" "}
            {dashboard.college.code}
          </p>

        </div>
      )}

      {error && (
        <p className="mt-4 text-red-600">
          {error}
        </p>
      )}

      {/* ========================================= */}
      {/* DASHBOARD CARDS */}
      {/* ========================================= */}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">

        {/* STUDENTS */}

        <div
          onClick={
            handleStudentsClick
          }
          className={`bg-white p-6 rounded-xl shadow cursor-pointer transition hover:shadow-lg ${
            selectedView ===
            "students"
              ? "border-2 border-blue-600 bg-blue-50"
              : ""
          }`}
        >
          <h3 className="text-gray-500 font-medium">
            Total Students
          </h3>

          <p className="text-4xl font-bold mt-3">
            {dashboard
              ?.total_students ??
              0}
          </p>

          <p className="text-blue-600 text-sm mt-3">
            Click to view details
          </p>
        </div>

        {/* TEACHERS */}

        <div
          onClick={
            handleTeachersClick
          }
          className={`bg-white p-6 rounded-xl shadow cursor-pointer transition hover:shadow-lg ${
            selectedView ===
            "teachers"
              ? "border-2 border-blue-600 bg-blue-50"
              : ""
          }`}
        >
          <h3 className="text-gray-500 font-medium">
            Total Teachers
          </h3>

          <p className="text-4xl font-bold mt-3">
            {dashboard
              ?.total_teachers ??
              0}
          </p>

          <p className="text-blue-600 text-sm mt-3">
            Click to view details
          </p>
        </div>

        {/* CLASSROOMS */}

        <div
          onClick={
            handleClassroomsClick
          }
          className={`bg-white p-6 rounded-xl shadow cursor-pointer transition hover:shadow-lg ${
            selectedView ===
            "classrooms"
              ? "border-2 border-blue-600 bg-blue-50"
              : ""
          }`}
        >
          <h3 className="text-gray-500 font-medium">
            Total Classrooms
          </h3>

          <p className="text-4xl font-bold mt-3">
            {dashboard
              ?.total_classrooms ??
              0}
          </p>

          <p className="text-blue-600 text-sm mt-3">
            Click to view details
          </p>
        </div>

      </div>

      {/* ========================================= */}
      {/* LOADING DETAILS */}
      {/* ========================================= */}

      {loadingDetails && (
        <div className="bg-white mt-8 p-6 rounded-xl shadow">
          Loading details...
        </div>
      )}

      {/* ========================================= */}
      {/* STUDENT DETAILS */}
      {/* ========================================= */}

      {!loadingDetails &&
        selectedView ===
          "students" && (
          <div className="bg-white mt-8 rounded-xl shadow overflow-hidden">

            <div className="p-6 border-b">

              <h2 className="text-xl font-semibold">
                Student Details
              </h2>

              <p className="text-gray-500 mt-2">
                Total Students:{" "}
                {
                  students.length
                }
              </p>

            </div>

            {students.length ===
            0 ? (
              <div className="p-6 text-gray-500">
                No students found.
              </div>
            ) : (
              <div className="overflow-x-auto">

                <table className="w-full">

                  <thead className="bg-gray-100">
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
                    </tr>
                  </thead>

                  <tbody>
                    {paginatedStudents.map(
                      (
                        student,
                      ) => (
                        <tr
                          key={
                            student.id
                          }
                          className="border-t"
                        >
                          <td className="p-4">
                            {
                              student.roll_number
                            }
                          </td>

                          <td className="p-4">
                            {
                              student.name
                            }
                          </td>

                          <td className="p-4">
                            {
                              student.age
                            }
                          </td>

                          <td className="p-4">
                            {
                              student.gender
                            }
                          </td>

                          <td className="p-4">
                            {getClassroomName(
                              student.classroom_id,
                            )}
                          </td>
                        </tr>
                      ),
                    )}
                  </tbody>

                </table>

              </div>
            )}

          </div>
        )}

      {/* ========================================= */}
      {/* TEACHER DETAILS */}
      {/* ========================================= */}

      {!loadingDetails &&
        selectedView ===
          "teachers" && (
          <div className="bg-white mt-8 rounded-xl shadow overflow-hidden">

            <div className="p-6 border-b">
              <h2 className="text-xl font-semibold">
                Teacher Details
              </h2>

              <p className="text-gray-500 mt-2">
                Total Teachers:{" "}
                {
                  teachers.length
                }
              </p>
            </div>

            {teachers.length ===
            0 ? (
              <div className="p-6 text-gray-500">
                No teachers found.
              </div>
            ) : (
              <div className="overflow-x-auto">

                <table className="w-full">

                  <thead className="bg-gray-100">
                    <tr>
                      <th className="text-left p-4">
                        Name
                      </th>

                      <th className="text-left p-4">
                        Subject
                      </th>

                      <th className="text-left p-4">
                        Email
                      </th>

                      <th className="text-left p-4">
                        Phone
                      </th>

                      <th className="text-left p-4">
                        Status
                      </th>
                    </tr>
                  </thead>

                  <tbody>

                    {paginatedTeachers.map(
                      (
                        teacher,
                      ) => (
                        <tr
                          key={
                            teacher.id
                          }
                          className="border-t"
                        >
                          <td className="p-4">
                            {
                              teacher.name
                            }
                          </td>

                          <td className="p-4">
                            {
                              teacher.subject
                            }
                          </td>

                          <td className="p-4">
                            {
                              teacher.email
                            }
                          </td>

                          <td className="p-4">
                            {
                              teacher.phone
                            }
                          </td>

                          <td className="p-4">

                            {teacher.status ===
                            "left"
                              ? "Left"
                              : "Active"}

                          </td>

                        </tr>
                      ),
                    )}

                  </tbody>

                </table>

              </div>
            )}

          </div>
        )}

      {/* ========================================= */}
      {/* CLASSROOM DETAILS */}
      {/* ========================================= */}

      {!loadingDetails &&
        selectedView ===
          "classrooms" && (
          <div className="bg-white mt-8 rounded-xl shadow overflow-hidden">

            <div className="p-6 border-b">

              <h2 className="text-xl font-semibold">
                Classroom Details
              </h2>

              <p className="text-gray-500 mt-2">
                Total Classrooms:{" "}
                {
                  sortedClassrooms.length
                }
              </p>

            </div>

            {sortedClassrooms.length ===
            0 ? (
              <div className="p-6 text-gray-500">
                No classrooms found.
              </div>
            ) : (
              <div className="overflow-x-auto">

                <table className="w-full">

                  <thead className="bg-gray-100">
                    <tr>
                      <th className="text-left p-4">
                        Class
                      </th>

                      <th className="text-left p-4">
                        Section
                      </th>

                      <th className="text-left p-4">
                        Capacity
                      </th>
                    </tr>
                  </thead>

                  <tbody>

                    {paginatedClassrooms.map(
                      (
                        classroom,
                      ) => (
                        <tr
                          key={
                            classroom.id
                          }
                          className="border-t"
                        >
                          <td className="p-4">
                            {
                              classroom.name
                            }
                          </td>

                          <td className="p-4">
                            {
                              classroom.section
                            }
                          </td>

                          <td className="p-4">
                            {
                              classroom.capacity
                            }
                          </td>
                        </tr>
                      ),
                    )}

                  </tbody>

                </table>

              </div>
            )}

          </div>
        )}

      {/* ========================================= */}
      {/* PAGINATION */}
      {/* ========================================= */}

      {!loadingDetails &&
        selectedView &&
        totalItems > 0 && (
          <div className="bg-white flex justify-between items-center p-5 border-t rounded-b-xl shadow">

            <p className="text-gray-600">
              Showing{" "}
              {startIndex + 1}
              {" - "}
              {Math.min(
                endIndex,
                totalItems,
              )}
              {" of "}
              {totalItems}
            </p>

            <div className="flex items-center gap-3">

              <button
                type="button"
                disabled={
                  currentPage === 1
                }
                onClick={() =>
                  setCurrentPage(
                    (
                      current,
                    ) =>
                      current - 1,
                  )
                }
                className="px-4 py-2 border rounded-lg disabled:opacity-40"
              >
                Previous
              </button>

              <span>
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
                    (
                      current,
                    ) =>
                      current + 1,
                  )
                }
                className="px-4 py-2 border rounded-lg disabled:opacity-40"
              >
                Next
              </button>

            </div>

          </div>
        )}

    </main>
  );
}

export default CollegeDashboard;