import { useEffect, useMemo, useState, type FormEvent } from "react";

import axios from "axios";

import {
  createTimetable,
  deleteTimetable,
  getTimetables,
  updateTimetable,
  type Timetable as TimetableType,
} from "../api/timetableApi";

import { getTeachers, type Teacher } from "../api/teacherApi";

import { getClassrooms, type Classroom } from "../api/classroomApi";

const DAYS = [
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
  "sunday",
];

function Timetable() {
  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [classrooms, setClassrooms] = useState<Classroom[]>([]);
  const [timetables, setTimetables] = useState<TimetableType[]>([]);

  const [teacherID, setTeacherID] = useState("");
  const [classroomID, setClassroomID] = useState("");
  const [subject, setSubject] = useState("");
  const [dayOfWeek, setDayOfWeek] = useState("monday");
  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [academicYear, setAcademicYear] = useState("2026-2027");

  // Timetable currently being edited
  const [editingID, setEditingID] = useState<string | null>(null);

  const [loading, setLoading] = useState(false);
  const [deletingID, setDeletingID] = useState<string | null>(null);

  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  // ========================================
  // ACTIVE TEACHERS
  // ========================================

  const activeTeachers = useMemo(() => {
    return teachers.filter(
      (teacher) => teacher.status?.toLowerCase().trim() !== "left",
    );
  }, [teachers]);

  // ========================================
  // INITIAL LOAD
  // ========================================

  useEffect(() => {
    loadInitialData();
  }, []);

  const loadInitialData = async () => {
    try {
      setError("");

      const [teacherData, classroomData] = await Promise.all([
        getTeachers(),
        getClassrooms(),
      ]);

      setTeachers(Array.isArray(teacherData) ? teacherData : []);

      setClassrooms(Array.isArray(classroomData) ? classroomData : []);

      await loadTimetables("2026-2027");
    } catch (err) {
      console.error("Load timetable page error:", err);

      setError("Failed to load timetable data");
    }
  };

  // ========================================
  // LOAD TIMETABLES
  // ========================================

  const loadTimetables = async (year: string) => {
    if (!year.trim()) {
      setTimetables([]);
      return;
    }

    try {
      const data = await getTimetables(year.trim());

      setTimetables(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error("Load timetables error:", err);

      setTimetables([]);
    }
  };

  // ========================================
  // ACADEMIC YEAR CHANGE
  // ========================================

  const handleAcademicYearChange = (value: string) => {
    setAcademicYear(value);
    setMessage("");
    setError("");
  };

  // ========================================
  // TEACHER CHANGE
  // ========================================

  const handleTeacherChange = (selectedTeacherID: string) => {
    setTeacherID(selectedTeacherID);

    const selectedTeacher = teachers.find(
      (teacher) => teacher.id === selectedTeacherID,
    );

    if (selectedTeacher) {
      setSubject(selectedTeacher.subject || "");
    } else {
      setSubject("");
    }
  };

  // ========================================
  // RESET FORM
  // ========================================

  const resetForm = () => {
    setEditingID(null);

    setTeacherID("");
    setClassroomID("");
    setSubject("");
    setDayOfWeek("monday");
    setStartTime("");
    setEndTime("");
  };

  // ========================================
  // CREATE OR UPDATE
  // ========================================

  const handleScheduleClass = async (event: FormEvent) => {
    event.preventDefault();

    setMessage("");
    setError("");

    if (!teacherID) {
      setError("Please select a teacher");
      return;
    }

    if (!classroomID) {
      setError("Please select a classroom");
      return;
    }

    if (!subject.trim()) {
      setError("Please enter a subject");
      return;
    }

    if (!academicYear.trim()) {
      setError("Please enter academic year");
      return;
    }

    if (!dayOfWeek) {
      setError("Please select a day");
      return;
    }

    if (!startTime) {
      setError("Please select start time");
      return;
    }

    if (!endTime) {
      setError("Please select end time");
      return;
    }

    if (startTime >= endTime) {
      setError("End time must be after start time");
      return;
    }

    const payload = {
      teacher_id: teacherID,
      classroom_id: classroomID,
      subject: subject.trim(),
      day_of_week: dayOfWeek,
      start_time: startTime,
      end_time: endTime,
      academic_year: academicYear.trim(),
    };

    try {
      setLoading(true);

      // UPDATE MODE
      if (editingID) {
        await updateTimetable(editingID, payload);

        setMessage("Class schedule updated successfully");

        resetForm();
      } else {
        // CREATE MODE
        await createTimetable(payload);

        setMessage("Class scheduled successfully");

        // Keep teacher/year/day
        // for faster scheduling.
        setClassroomID("");
        setStartTime("");
        setEndTime("");
      }

      await loadTimetables(academicYear.trim());
    } catch (err) {
      console.error("Save timetable error:", err);

      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Failed to save class schedule");
      } else {
        setError("Failed to save class schedule");
      }
    } finally {
      setLoading(false);
    }
  };

  // ========================================
  // EDIT TIMETABLE
  // ========================================

  const handleEdit = (timetable: TimetableType) => {
    setMessage("");
    setError("");

    setEditingID(timetable.id);

    setTeacherID(timetable.teacher_id);

    setClassroomID(timetable.classroom_id);

    setSubject(timetable.subject);

    setDayOfWeek(timetable.day_of_week);

    setStartTime(timetable.start_time);

    setEndTime(timetable.end_time);

    setAcademicYear(timetable.academic_year);

    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  };

  // ========================================
  // CANCEL EDIT
  // ========================================

  const handleCancelEdit = () => {
    resetForm();

    setMessage("");
    setError("");
  };

  // ========================================
  // DELETE TIMETABLE
  // ========================================

  const handleDelete = async (timetable: TimetableType) => {
    const teacherName = getTeacherName(timetable.teacher_id);

    const confirmed = window.confirm(
      `Delete ${teacherName}'s ${timetable.subject} class on ${timetable.day_of_week} from ${timetable.start_time} to ${timetable.end_time}?`,
    );

    if (!confirmed) {
      return;
    }

    try {
      setDeletingID(timetable.id);

      setMessage("");
      setError("");

      await deleteTimetable(timetable.id);

      setMessage("Class schedule deleted successfully");

      if (editingID === timetable.id) {
        resetForm();
      }

      await loadTimetables(academicYear.trim());
    } catch (err) {
      console.error("Delete timetable error:", err);

      if (axios.isAxiosError(err)) {
        setError(
          err.response?.data?.error || "Failed to delete class schedule",
        );
      } else {
        setError("Failed to delete class schedule");
      }
    } finally {
      setDeletingID(null);
    }
  };

  // ========================================
  // TEACHER NAME
  // ========================================

  const getTeacherName = (teacherIDValue: string) => {
    const teacher = teachers.find((item) => item.id === teacherIDValue);

    return teacher?.name || "Teacher not found";
  };

  // ========================================
  // CLASSROOM NAME
  // ========================================

  const getClassroomName = (classroomIDValue: string) => {
    const classroom = classrooms.find((item) => item.id === classroomIDValue);

    if (!classroom) {
      return "Classroom not found";
    }

    return `${classroom.name} - Section ${classroom.section}`;
  };

  // ========================================
  // SORT TIMETABLE
  // ========================================

  const sortedTimetables = useMemo(() => {
    return [...timetables].sort((a, b) => {
      const dayA = DAYS.indexOf(a.day_of_week.toLowerCase());

      const dayB = DAYS.indexOf(b.day_of_week.toLowerCase());

      if (dayA !== dayB) {
        return dayA - dayB;
      }

      return a.start_time.localeCompare(b.start_time);
    });
  }, [timetables]);

  return (
    <main className="p-8">
      {/* PAGE TITLE */}

      <div>
        <h1 className="text-3xl font-bold">Class Schedule</h1>

        <p className="mt-2 text-gray-600">
          Assign teachers to classrooms and schedule their weekly classes.
        </p>
      </div>

      {/* SUCCESS MESSAGE */}

      {message && (
        <div className="mt-6 p-4 rounded-lg bg-green-100 text-green-700">
          {message}
        </div>
      )}

      {/* ERROR MESSAGE */}

      {error && (
        <div className="mt-6 p-4 rounded-lg bg-red-100 text-red-700">
          {error}
        </div>
      )}

      {/* ======================================== */}
      {/* CLASS SCHEDULE FORM */}
      {/* ======================================== */}

      <div className="bg-white rounded-xl shadow mt-8 p-6">
        <h2 className="text-xl font-semibold">
          {editingID ? "Edit Class Schedule" : "Schedule a Class"}
        </h2>

        {editingID && (
          <p className="text-sm text-blue-600 mt-2">
            You are editing an existing class schedule.
          </p>
        )}

        <form onSubmit={handleScheduleClass} className="mt-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* TEACHER */}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Teacher
              </label>

              <select
                value={teacherID}
                onChange={(event) => handleTeacherChange(event.target.value)}
                className="w-full border rounded-lg px-3 py-2"
              >
                <option value="">Select Teacher</option>

                {activeTeachers.map((teacher) => (
                  <option key={teacher.id} value={teacher.id}>
                    {teacher.name}
                    {teacher.subject ? ` - ${teacher.subject}` : ""}
                  </option>
                ))}
              </select>
            </div>

            {/* CLASSROOM */}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Classroom
              </label>

              <select
                value={classroomID}
                onChange={(event) => setClassroomID(event.target.value)}
                className="w-full border rounded-lg px-3 py-2"
              >
                <option value="">Select Classroom</option>

                {classrooms.map((classroom) => (
                  <option key={classroom.id} value={classroom.id}>
                    {classroom.name} - Section {classroom.section}
                  </option>
                ))}
              </select>
            </div>

            {/* SUBJECT */}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Subject
              </label>

              <input
                type="text"
                value={subject}
                onChange={(event) => setSubject(event.target.value)}
                placeholder="Mathematics"
                className="w-full border rounded-lg px-3 py-2"
              />
            </div>

            {/* ACADEMIC YEAR */}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Academic Year
              </label>

              <input
                type="text"
                value={academicYear}
                onChange={(event) =>
                  handleAcademicYearChange(event.target.value)
                }
                onBlur={() => loadTimetables(academicYear.trim())}
                placeholder="2026-2027"
                className="w-full border rounded-lg px-3 py-2"
              />
            </div>

            {/* DAY */}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Day
              </label>

              <select
                value={dayOfWeek}
                onChange={(event) => setDayOfWeek(event.target.value)}
                className="w-full border rounded-lg px-3 py-2"
              >
                {DAYS.map((day) => (
                  <option key={day} value={day}>
                    {day.charAt(0).toUpperCase() + day.slice(1)}
                  </option>
                ))}
              </select>
            </div>

            {/* START TIME */}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Start Time
              </label>

              <input
                type="time"
                value={startTime}
                onChange={(event) => setStartTime(event.target.value)}
                className="w-full border rounded-lg px-3 py-2"
              />
            </div>

            {/* END TIME */}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                End Time
              </label>

              <input
                type="time"
                value={endTime}
                onChange={(event) => setEndTime(event.target.value)}
                className="w-full border rounded-lg px-3 py-2"
              />
            </div>
          </div>

          {/* BUTTONS */}

          <div className="mt-6 flex gap-3">
            <button
              type="submit"
              disabled={loading}
              className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {loading
                ? editingID
                  ? "Updating..."
                  : "Scheduling..."
                : editingID
                  ? "Update Class"
                  : "Schedule Class"}
            </button>

            {editingID && (
              <button
                type="button"
                onClick={handleCancelEdit}
                disabled={loading}
                className="bg-gray-500 text-white px-6 py-2 rounded-lg hover:bg-gray-600 disabled:opacity-50"
              >
                Cancel
              </button>
            )}
          </div>
        </form>
      </div>

      {/* ======================================== */}
      {/* WEEKLY SCHEDULE */}
      {/* ======================================== */}

      <div className="bg-white rounded-xl shadow mt-8 overflow-hidden">
        <div className="p-6 border-b">
          <h2 className="text-xl font-semibold">Weekly Class Schedule</h2>

          <p className="text-gray-500 mt-1">Academic Year: {academicYear}</p>

          <p className="text-gray-500 mt-1">
            Total Scheduled Classes: {sortedTimetables.length}
          </p>
        </div>

        {sortedTimetables.length === 0 ? (
          <div className="p-6 text-gray-500">
            No classes scheduled for this academic year.
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-100">
                <tr>
                  <th className="text-left p-4">Teacher</th>

                  <th className="text-left p-4">Classroom</th>

                  <th className="text-left p-4">Subject</th>

                  <th className="text-left p-4">Day</th>

                  <th className="text-left p-4">Start Time</th>

                  <th className="text-left p-4">End Time</th>

                  <th className="text-left p-4">Actions</th>
                </tr>
              </thead>

              <tbody>
                {sortedTimetables.map((timetable) => (
                  <tr
                    key={timetable.id}
                    className={`border-t hover:bg-gray-50 ${
                      editingID === timetable.id ? "bg-blue-50" : ""
                    }`}
                  >
                    <td className="p-4 font-medium">
                      {getTeacherName(timetable.teacher_id)}
                    </td>

                    <td className="p-4">
                      {getClassroomName(timetable.classroom_id)}
                    </td>

                    <td className="p-4">{timetable.subject}</td>

                    <td className="p-4 capitalize">{timetable.day_of_week}</td>

                    <td className="p-4">{timetable.start_time}</td>

                    <td className="p-4">{timetable.end_time}</td>

                    <td className="p-4">
                      <div className="flex gap-2">
                        {/* EDIT */}

                        <button
                          type="button"
                          onClick={() => handleEdit(timetable)}
                          disabled={deletingID === timetable.id}
                          className="bg-yellow-500 text-white px-3 py-1 rounded hover:bg-yellow-600 disabled:opacity-50"
                        >
                          Edit
                        </button>

                        {/* DELETE */}

                        <button
                          type="button"
                          onClick={() => handleDelete(timetable)}
                          disabled={deletingID === timetable.id}
                          className="bg-red-600 text-white px-3 py-1 rounded hover:bg-red-700 disabled:opacity-50"
                        >
                          {deletingID === timetable.id
                            ? "Deleting..."
                            : "Delete"}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </main>
  );
}

export default Timetable;
