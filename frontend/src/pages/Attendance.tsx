import { useEffect, useMemo, useState } from "react";

import { getClassrooms } from "../api/classroomApi";
import type { Classroom } from "../api/classroomApi";

import { getStudents } from "../api/studentApi";
import type { Student } from "../api/studentApi";

import {
  createBulkAttendance,
  getAttendance,
  getClassroomMonthlyAttendance,
} from "../api/attendanceApi";

import type {
  Attendance as AttendanceRecord,
  ClassroomMonthlyAttendanceSummary,
} from "../api/attendanceApi";

type AttendanceStatus = "Present" | "Absent";

type StudentAttendance = {
  student: Student;
  status: AttendanceStatus;
  alreadySaved: boolean;
};

// -----------------------------------------
// GET TODAY DATE
// -----------------------------------------

const getTodayDate = () => {
  const today = new Date();

  const year = today.getFullYear();

  const month = String(today.getMonth() + 1).padStart(2, "0");

  const day = String(today.getDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
};

// -----------------------------------------
// GET CURRENT MONTH
// Example: 2026-09
// -----------------------------------------

const getCurrentMonth = () => {
  const today = new Date();

  const year = today.getFullYear();

  const month = String(today.getMonth() + 1).padStart(2, "0");

  return `${year}-${month}`;
};

// -----------------------------------------
// MONTH NAME
// -----------------------------------------

const getMonthName = (year: number, month: number) => {
  return new Date(year, month - 1, 1).toLocaleString("en-IN", {
    month: "long",
    year: "numeric",
  });
};

function Attendance() {
  const [classrooms, setClassrooms] = useState<Classroom[]>([]);

  const [students, setStudents] = useState<Student[]>([]);

  const [savedAttendance, setSavedAttendance] = useState<AttendanceRecord[]>(
    [],
  );

  const [selectedClassroomId, setSelectedClassroomId] = useState("");

  const [academicYear, setAcademicYear] = useState("2026-2027");

  const today = getTodayDate();

  const [date, setDate] = useState(today);

  const [attendanceList, setAttendanceList] = useState<StudentAttendance[]>([]);

  const [loading, setLoading] = useState(true);

  const [loadingAttendance, setLoadingAttendance] = useState(false);

  const [saving, setSaving] = useState(false);

  const [message, setMessage] = useState("");

  const [error, setError] = useState("");

  // -----------------------------------------
  // ATTENDANCE HISTORY STATE
  // -----------------------------------------

  const [historyMonth, setHistoryMonth] = useState(getCurrentMonth());

  const [history, setHistory] =
    useState<ClassroomMonthlyAttendanceSummary | null>(null);

  const [historyLoading, setHistoryLoading] = useState(false);

  const [historyError, setHistoryError] = useState("");

  // -----------------------------------------
  // LOAD CLASSROOMS + STUDENTS
  // -----------------------------------------

  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        setError("");

        const [classroomData, studentData] = await Promise.all([
          getClassrooms(),
          getStudents(),
        ]);

        setClassrooms(classroomData);
        setStudents(studentData);
      } catch (err) {
        console.error("Failed to load data:", err);

        setError("Failed to load classrooms or students");
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  // -----------------------------------------
  // LOAD SAVED ATTENDANCE
  // -----------------------------------------

  useEffect(() => {
    const loadAttendance = async () => {
      if (!academicYear.trim()) {
        setSavedAttendance([]);
        return;
      }

      try {
        setLoadingAttendance(true);

        setError("");

        const data = await getAttendance(academicYear.trim());

        setSavedAttendance(data);
      } catch (err) {
        console.error("Failed to load attendance:", err);

        setSavedAttendance([]);

        if (err instanceof Error) {
          setError(err.message);
        }
      } finally {
        setLoadingAttendance(false);
      }
    };

    loadAttendance();
  }, [academicYear]);

  // -----------------------------------------
  // LOAD MONTHLY ATTENDANCE HISTORY
  // -----------------------------------------

  useEffect(() => {
    const loadHistory = async () => {
      if (!selectedClassroomId || !academicYear.trim() || !historyMonth) {
        setHistory(null);
        setHistoryError("");
        return;
      }

      try {
        setHistoryLoading(true);
        setHistoryError("");

        const data = await getClassroomMonthlyAttendance(
          selectedClassroomId,
          academicYear.trim(),
          historyMonth,
        );

        setHistory(data);
      } catch (err) {
        console.error("Failed to load attendance history:", err);

        setHistory(null);

        if (err instanceof Error) {
          setHistoryError(err.message);
        } else {
          setHistoryError("Failed to load attendance history");
        }
      } finally {
        setHistoryLoading(false);
      }
    };

    loadHistory();
  }, [selectedClassroomId, academicYear, historyMonth]);

  // -----------------------------------------
  // SORT CLASSROOMS
  // -----------------------------------------

  const sortedClassrooms = useMemo(() => {
    return [...classrooms].sort((a, b) => {
      const classA = parseInt(a.name) || 0;

      const classB = parseInt(b.name) || 0;

      if (classA !== classB) {
        return classA - classB;
      }

      return a.section.localeCompare(b.section);
    });
  }, [classrooms]);

  // -----------------------------------------
  // FILTER STUDENTS BY CLASSROOM
  // -----------------------------------------

  const filteredStudents = useMemo(() => {
    if (!selectedClassroomId) {
      return [];
    }

    return students
      .filter((student) => student.classroom_id === selectedClassroomId)
      .sort((a, b) => a.roll_number - b.roll_number);
  }, [students, selectedClassroomId]);

  // -----------------------------------------
  // SELECTED CLASSROOM
  // -----------------------------------------

  const selectedClassroom = useMemo(() => {
    return classrooms.find((classroom) => classroom.id === selectedClassroomId);
  }, [classrooms, selectedClassroomId]);

  // -----------------------------------------
  // PREPARE ATTENDANCE LIST
  // -----------------------------------------

  useEffect(() => {
    const list = filteredStudents.map((student) => {
      const existing = savedAttendance.find(
        (attendance) =>
          attendance.student_id === student.id &&
          attendance.classroom_id === selectedClassroomId &&
          attendance.academic_year === academicYear.trim() &&
          attendance.date === date,
      );

      return {
        student,

        status: existing?.attendance ?? "Present",

        alreadySaved: Boolean(existing),
      };
    });

    setAttendanceList(list);
  }, [
    filteredStudents,
    savedAttendance,
    selectedClassroomId,
    academicYear,
    date,
  ]);

  // -----------------------------------------
  // CHANGE SINGLE STUDENT ATTENDANCE
  // -----------------------------------------

  const changeAttendance = (studentId: string, status: AttendanceStatus) => {
    setAttendanceList((previous) =>
      previous.map((item) =>
        item.student.id === studentId
          ? {
              ...item,
              status,
            }
          : item,
      ),
    );

    setMessage("");
  };

  // -----------------------------------------
  // MARK ALL PRESENT
  // -----------------------------------------

  const markAllPresent = () => {
    setAttendanceList((previous) =>
      previous.map((item) => ({
        ...item,
        status: "Present",
      })),
    );

    setMessage("");
  };

  // -----------------------------------------
  // MARK ALL ABSENT
  // -----------------------------------------

  const markAllAbsent = () => {
    setAttendanceList((previous) =>
      previous.map((item) => ({
        ...item,
        status: "Absent",
      })),
    );

    setMessage("");
  };

  // -----------------------------------------
  // COUNTS
  // -----------------------------------------

  const presentCount = attendanceList.filter(
    (item) => item.status === "Present",
  ).length;

  const absentCount = attendanceList.filter(
    (item) => item.status === "Absent",
  ).length;

  const savedCount = attendanceList.filter((item) => item.alreadySaved).length;

  // -----------------------------------------
  // SAVE / UPDATE BULK ATTENDANCE
  // -----------------------------------------

  const handleSaveAttendance = async () => {
    setMessage("");
    setError("");

    if (!selectedClassroomId) {
      setError("Please select classroom");
      return;
    }

    if (!academicYear.trim()) {
      setError("Please enter academic year");
      return;
    }

    if (!date) {
      setError("Please select attendance date");
      return;
    }

    if (date !== today) {
      setError("Attendance can only be marked or updated for today's date");
      return;
    }

    if (attendanceList.length === 0) {
      setError("No students found in this classroom");
      return;
    }

    try {
      setSaving(true);

      const hadExisting = attendanceList.some((item) => item.alreadySaved);

      // ONE BULK API REQUEST
      await createBulkAttendance({
        classroom_id: selectedClassroomId,

        academic_year: academicYear.trim(),

        date,

        attendances: attendanceList.map((item) => ({
          student_id: item.student.id,

          attendance: item.status,
        })),
      });

      // RELOAD FROM MONGODB
      const freshAttendance = await getAttendance(academicYear.trim());

      setSavedAttendance(freshAttendance);

      // REFRESH HISTORY TOO
      if (historyMonth) {
        const freshHistory = await getClassroomMonthlyAttendance(
          selectedClassroomId,
          academicYear.trim(),
          historyMonth,
        );

        setHistory(freshHistory);
      }

      if (hadExisting) {
        setMessage("Attendance updated successfully!");
      } else {
        setMessage("Attendance saved successfully!");
      }

      setTimeout(() => {
        setMessage("");
      }, 4000);
    } catch (err) {
      console.error("Save attendance error:", err);

      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Failed to save attendance");
      }
    } finally {
      setSaving(false);
    }
  };

  // -----------------------------------------
  // UI
  // -----------------------------------------

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      {/* TITLE */}

      <div className="mb-8">
        <h1 className="text-3xl font-bold">Attendance</h1>

        <p className="text-gray-600 mt-2">
          Mark today's attendance and monitor historical attendance month-wise.
        </p>
      </div>

      {/* ERROR */}

      {error && (
        <div className="mb-5 bg-red-100 text-red-700 p-4 rounded-lg">
          {error}
        </div>
      )}

      {/* SUCCESS */}

      {message && (
        <div className="mb-5 flex items-center gap-3 rounded-lg border border-green-300 bg-green-100 p-4 text-green-800">
          <span className="text-xl">✓</span>

          <span className="font-medium">{message}</span>
        </div>
      )}

      {/* FILTERS */}

      <div className="bg-white rounded-xl shadow p-6 mb-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* CLASSROOM */}

          <div>
            <label className="block font-medium mb-2">Classroom</label>

            <select
              value={selectedClassroomId}
              onChange={(e) => setSelectedClassroomId(e.target.value)}
              className="w-full border rounded-lg p-3"
            >
              <option value="">Select Classroom</option>

              {sortedClassrooms.map((classroom) => (
                <option key={classroom.id} value={classroom.id}>
                  {classroom.name}

                  {classroom.section ? ` - Section ${classroom.section}` : ""}
                </option>
              ))}
            </select>
          </div>

          {/* ACADEMIC YEAR */}

          <div>
            <label className="block font-medium mb-2">Academic Year</label>

            <input
              type="text"
              value={academicYear}
              onChange={(e) => setAcademicYear(e.target.value)}
              placeholder="2026-2027"
              className="w-full border rounded-lg p-3"
            />
          </div>

          {/* DATE */}

          <div>
            <label className="block font-medium mb-2">Date</label>

            <input
              type="date"
              value={date}
              min={today}
              max={today}
              onChange={(e) => setDate(e.target.value)}
              className="w-full border rounded-lg p-3"
            />

            <p className="text-sm text-gray-500 mt-2">
              Attendance can only be marked or updated for today's date.
            </p>
          </div>
        </div>
      </div>

      {/* TODAY STUDENTS */}

      <div className="bg-white rounded-xl shadow overflow-hidden mb-10">
        <div className="p-6 border-b flex flex-wrap justify-between items-center gap-4">
          <div>
            <h2 className="text-xl font-semibold">Today's Attendance</h2>

            {selectedClassroomId && (
              <div className="text-gray-500 mt-1">
                <span>Total Students: {attendanceList.length}</span>

                {savedCount > 0 && (
                  <span className="ml-4 text-blue-600">
                    Already Marked: {savedCount}
                  </span>
                )}
              </div>
            )}
          </div>

          {attendanceList.length > 0 && (
            <div className="flex gap-3">
              <button
                type="button"
                onClick={markAllPresent}
                disabled={saving}
                className="bg-green-600 hover:bg-green-700 disabled:bg-green-400 text-white px-4 py-2 rounded-lg"
              >
                Mark All Present
              </button>

              <button
                type="button"
                onClick={markAllAbsent}
                disabled={saving}
                className="bg-red-600 hover:bg-red-700 disabled:bg-red-400 text-white px-4 py-2 rounded-lg"
              >
                Mark All Absent
              </button>
            </div>
          )}
        </div>

        {loading || loadingAttendance ? (
          <div className="p-8 text-center text-gray-500">
            Loading attendance...
          </div>
        ) : !selectedClassroomId ? (
          <div className="p-10 text-center text-gray-500">
            Select a classroom to view students.
          </div>
        ) : attendanceList.length === 0 ? (
          <div className="p-10 text-center text-gray-500">
            No students found in this classroom.
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="text-left px-6 py-4">Roll No</th>

                    <th className="text-left px-6 py-4">Student Name</th>

                    <th className="text-left px-6 py-4">Gender</th>

                    <th className="text-center px-6 py-4">Present</th>

                    <th className="text-center px-6 py-4">Absent</th>
                  </tr>
                </thead>

                <tbody>
                  {attendanceList.map((item) => (
                    <tr
                      key={item.student.id}
                      className="border-t hover:bg-gray-50"
                    >
                      <td className="px-6 py-4">{item.student.roll_number}</td>

                      <td className="px-6 py-4 font-medium">
                        {item.student.name}

                        {item.alreadySaved && (
                          <span className="ml-2 text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded">
                            Saved
                          </span>
                        )}
                      </td>

                      <td className="px-6 py-4">{item.student.gender}</td>

                      <td className="px-6 py-4 text-center">
                        <input
                          type="radio"
                          name={`attendance-${item.student.id}`}
                          checked={item.status === "Present"}
                          onChange={() =>
                            changeAttendance(item.student.id, "Present")
                          }
                          disabled={saving}
                          className="w-5 h-5 cursor-pointer"
                        />
                      </td>

                      <td className="px-6 py-4 text-center">
                        <input
                          type="radio"
                          name={`attendance-${item.student.id}`}
                          checked={item.status === "Absent"}
                          onChange={() =>
                            changeAttendance(item.student.id, "Absent")
                          }
                          disabled={saving}
                          className="w-5 h-5 cursor-pointer"
                        />
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* TODAY SUMMARY */}

            <div className="border-t p-6 flex flex-wrap justify-between items-center gap-5">
              <div className="flex flex-wrap gap-6 text-lg">
                <div>
                  Total: <strong>{attendanceList.length}</strong>
                </div>

                <div className="text-green-700">
                  Present: <strong>{presentCount}</strong>
                </div>

                <div className="text-red-700">
                  Absent: <strong>{absentCount}</strong>
                </div>
              </div>

              <button
                type="button"
                onClick={handleSaveAttendance}
                disabled={saving || loadingAttendance}
                className="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-6 py-3 rounded-lg"
              >
                {saving
                  ? "Saving..."
                  : savedCount > 0
                    ? "Update Attendance"
                    : "Save Attendance"}
              </button>
            </div>
          </>
        )}
      </div>

      {/* ===================================== */}
      {/* ATTENDANCE HISTORY */}
      {/* ===================================== */}

      <div className="bg-white rounded-xl shadow overflow-hidden">
        <div className="p-6 border-b">
          <div className="flex flex-wrap items-end justify-between gap-5">
            <div>
              <h2 className="text-2xl font-semibold">Attendance History</h2>

              <p className="text-gray-500 mt-1">
                View month-wise historical attendance for the selected
                classroom.
              </p>

              {selectedClassroom && (
                <p className="mt-2 font-medium text-blue-700">
                  {selectedClassroom.name}
                  {selectedClassroom.section
                    ? ` - Section ${selectedClassroom.section}`
                    : ""}
                </p>
              )}
            </div>

            <div className="w-full sm:w-64">
              <label className="block text-sm font-medium mb-2">
                Select Month
              </label>

              <input
                type="month"
                value={historyMonth}
                max={getCurrentMonth()}
                onChange={(e) => setHistoryMonth(e.target.value)}
                className="w-full border rounded-lg p-3"
              />
            </div>
          </div>
        </div>

        {historyError && (
          <div className="m-6 bg-red-100 text-red-700 p-4 rounded-lg">
            {historyError}
          </div>
        )}

        {!selectedClassroomId ? (
          <div className="p-10 text-center text-gray-500">
            Select a classroom above to view attendance history.
          </div>
        ) : historyLoading ? (
          <div className="p-10 text-center text-gray-500">
            Loading attendance history...
          </div>
        ) : !history ? (
          <div className="p-10 text-center text-gray-500">
            No attendance history available.
          </div>
        ) : (
          <>
            {/* MONTH TITLE */}

            <div className="px-6 pt-6">
              <h3 className="text-xl font-semibold">
                {getMonthName(history.year, history.month)}
              </h3>

              <p className="text-gray-500">
                Academic Year: {history.academic_year}
              </p>
            </div>

            {/* HISTORY SUMMARY CARDS */}

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4 p-6">
              <div className="border rounded-xl p-5">
                <p className="text-sm text-gray-500">Working Days</p>

                <p className="text-2xl font-bold mt-1">
                  {history.working_days}
                </p>
              </div>

              <div className="border rounded-xl p-5">
                <p className="text-sm text-gray-500">Marked Days</p>

                <p className="text-2xl font-bold mt-1 text-blue-700">
                  {history.attendance_marked_days}
                </p>
              </div>

              <div className="border rounded-xl p-5">
                <p className="text-sm text-gray-500">Total Students</p>

                <p className="text-2xl font-bold mt-1">
                  {history.total_students}
                </p>
              </div>

              <div className="border rounded-xl p-5">
                <p className="text-sm text-gray-500">Present Records</p>

                <p className="text-2xl font-bold mt-1 text-green-700">
                  {history.present}
                </p>
              </div>

              <div className="border rounded-xl p-5">
                <p className="text-sm text-gray-500">Attendance %</p>

                <p className="text-2xl font-bold mt-1">
                  {history.attendance_percentage.toFixed(2)}%
                </p>
              </div>
            </div>

            {/* IMPORTANT INFO */}

            <div className="mx-6 mb-6 rounded-lg bg-blue-50 border border-blue-200 p-4 text-sm text-blue-800">
              Attendance percentage is calculated using only marked attendance
              records. Unmarked days are not treated as absent.
            </div>

            {/* STUDENT HISTORY TABLE */}

            <div className="overflow-x-auto border-t">
              <table className="w-full">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="text-left px-6 py-4">Roll No</th>

                    <th className="text-left px-6 py-4">Student Name</th>

                    <th className="text-center px-6 py-4">Present</th>

                    <th className="text-center px-6 py-4">Absent</th>

                    <th className="text-center px-6 py-4">Marked Days</th>

                    <th className="text-center px-6 py-4">Attendance %</th>
                  </tr>
                </thead>

                <tbody>
                  {history.students.length === 0 ? (
                    <tr>
                      <td
                        colSpan={6}
                        className="p-10 text-center text-gray-500"
                      >
                        No attendance records found for this month.
                      </td>
                    </tr>
                  ) : (
                    history.students.map((student) => (
                      <tr
                        key={student.student_id}
                        className="border-t hover:bg-gray-50"
                      >
                        <td className="px-6 py-4">{student.roll_number}</td>

                        <td className="px-6 py-4 font-medium">
                          {student.name}
                        </td>

                        <td className="px-6 py-4 text-center text-green-700 font-medium">
                          {student.present}
                        </td>

                        <td className="px-6 py-4 text-center text-red-700 font-medium">
                          {student.absent}
                        </td>

                        <td className="px-6 py-4 text-center">
                          {student.attendance_marked_days}
                        </td>

                        <td className="px-6 py-4 text-center font-semibold">
                          {student.attendance_percentage.toFixed(2)}%
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

            {/* BOTTOM TOTAL */}

            <div className="border-t bg-gray-50 p-6 flex flex-wrap gap-8">
              <div>
                Present Records:{" "}
                <strong className="text-green-700">{history.present}</strong>
              </div>

              <div>
                Absent Records:{" "}
                <strong className="text-red-700">{history.absent}</strong>
              </div>

              <div>
                Marked Days: <strong>{history.attendance_marked_days}</strong>
              </div>

              <div>
                Overall Attendance:{" "}
                <strong>{history.attendance_percentage.toFixed(2)}%</strong>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default Attendance;
