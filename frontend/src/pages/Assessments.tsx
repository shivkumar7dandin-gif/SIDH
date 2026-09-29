import { useEffect, useMemo, useState } from "react";

import { createAssessment } from "../api/assessmentApi";

import type { Assessment } from "../api/assessmentApi";

import { getClassrooms } from "../api/classroomApi";

import type { Classroom } from "../api/classroomApi";

import { getStudents } from "../api/studentApi";

import type { Student } from "../api/studentApi";

function Assessments() {
  const [classrooms, setClassrooms] = useState<Classroom[]>([]);

  const [students, setStudents] = useState<Student[]>([]);

  const [selectedClassroomId, setSelectedClassroomId] = useState("");

  const [loading, setLoading] = useState(true);

  const [saving, setSaving] = useState(false);

  const [message, setMessage] = useState("");

  const [error, setError] = useState("");

  const [form, setForm] = useState<Assessment>({
    student_id: "",
    classroom_id: "",
    academic_year: "2026-2027",
    assessment_date: "",
    assessment_no: 1,
    subject: "",
    total_marks: 100,
    obtained_marks: 0,
  });

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
        console.error("Failed to load assessment data:", err);

        setError("Failed to load classrooms or students");
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

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
  // CLASSROOM CHANGE
  // -----------------------------------------

  const handleClassroomChange = (classroomId: string) => {
    setSelectedClassroomId(classroomId);

    setMessage("");
    setError("");

    setForm((previous) => ({
      ...previous,

      classroom_id: classroomId,

      student_id: "",
    }));
  };

  // -----------------------------------------
  // FORM CHANGE
  // -----------------------------------------

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;

    setForm((previous) => ({
      ...previous,

      [name]:
        name === "assessment_no" ||
        name === "total_marks" ||
        name === "obtained_marks"
          ? Number(value)
          : value,
    }));

    setMessage("");
  };

  // -----------------------------------------
  // SUBMIT ASSESSMENT
  // -----------------------------------------

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    setMessage("");
    setError("");

    if (!selectedClassroomId) {
      setError("Please select classroom");
      return;
    }

    if (!form.student_id) {
      setError("Please select student");
      return;
    }

    if (!form.academic_year.trim()) {
      setError("Please enter academic year");
      return;
    }

    if (!form.assessment_date) {
      setError("Please select assessment date");
      return;
    }

    if (!form.subject.trim()) {
      setError("Please enter subject");
      return;
    }

    if (form.total_marks <= 0) {
      setError("Total marks must be greater than 0");
      return;
    }

    if (form.obtained_marks < 0) {
      setError("Obtained marks cannot be negative");
      return;
    }

    if (form.obtained_marks > form.total_marks) {
      setError("Obtained marks cannot be greater than total marks");
      return;
    }

    try {
      setSaving(true);

      const result = await createAssessment(form);

      setMessage(result.message || "Assessment created successfully!");

      setForm((previous) => ({
        ...previous,

        student_id: "",

        subject: "",

        obtained_marks: 0,
      }));

      setTimeout(() => {
        setMessage("");
      }, 4000);
    } catch (err) {
      console.error("Create assessment error:", err);

      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Something went wrong");
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
        <h1 className="text-3xl font-bold">Assessments</h1>

        <p className="text-gray-600 mt-2">
          Select a classroom and student, then enter assessment marks.
        </p>
      </div>

      {/* SUCCESS */}

      {message && (
        <div className="mb-5 flex items-center gap-3 rounded-lg border border-green-300 bg-green-100 p-4 text-green-800">
          <span className="text-xl">✓</span>

          <span className="font-medium">{message}</span>
        </div>
      )}

      {/* ERROR */}

      {error && (
        <div className="mb-5 rounded-lg border border-red-300 bg-red-100 p-4 text-red-700">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="bg-white rounded-xl shadow p-6">
        {/* CLASSROOM + STUDENT */}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div>
            <label className="block font-medium mb-2">Classroom</label>

            <select
              value={selectedClassroomId}
              onChange={(e) => handleClassroomChange(e.target.value)}
              className="w-full border rounded-lg p-3"
              required
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

          <div>
            <label className="block font-medium mb-2">Student</label>

            <select
              name="student_id"
              value={form.student_id}
              onChange={handleChange}
              disabled={!selectedClassroomId || loading}
              className="w-full border rounded-lg p-3 disabled:bg-gray-100"
              required
            >
              <option value="">Select Student</option>

              {filteredStudents.map((student) => (
                <option key={student.id} value={student.id}>
                  Roll {student.roll_number}
                  {" - "}
                  {student.name}
                </option>
              ))}
            </select>

            {selectedClassroomId && filteredStudents.length === 0 && (
              <p className="text-sm text-red-600 mt-2">
                No students found in this classroom.
              </p>
            )}
          </div>
        </div>

        {/* ACADEMIC YEAR + DATE */}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div>
            <label className="block font-medium mb-2">Academic Year</label>

            <input
              type="text"
              name="academic_year"
              placeholder="2026-2027"
              value={form.academic_year}
              onChange={handleChange}
              className="w-full border rounded-lg p-3"
              required
            />
          </div>

          <div>
            <label className="block font-medium mb-2">Assessment Date</label>

            <input
              type="date"
              name="assessment_date"
              value={form.assessment_date}
              onChange={handleChange}
              className="w-full border rounded-lg p-3"
              required
            />
          </div>
        </div>

        {/* ASSESSMENT NUMBER + SUBJECT */}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div>
            <label className="block font-medium mb-2">Assessment Number</label>

            <input
              type="number"
              name="assessment_no"
              min="1"
              value={form.assessment_no}
              onChange={handleChange}
              className="w-full border rounded-lg p-3"
              required
            />
          </div>

          <div>
            <label className="block font-medium mb-2">Subject</label>

            <input
              type="text"
              name="subject"
              value={form.subject}
              onChange={handleChange}
              placeholder="Example: Mathematics"
              className="w-full border rounded-lg p-3"
              required
            />
          </div>
        </div>

        {/* MARKS */}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <div>
            <label className="block font-medium mb-2">Total Marks</label>

            <input
              type="number"
              name="total_marks"
              min="1"
              value={form.total_marks}
              onChange={handleChange}
              className="w-full border rounded-lg p-3"
              required
            />
          </div>

          <div>
            <label className="block font-medium mb-2">Obtained Marks</label>

            <input
              type="number"
              name="obtained_marks"
              min="0"
              max={form.total_marks}
              value={form.obtained_marks}
              onChange={handleChange}
              className="w-full border rounded-lg p-3"
              required
            />
          </div>
        </div>

        {/* BUTTON */}

        <button
          type="submit"
          disabled={saving || loading}
          className="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-6 py-3 rounded-lg"
        >
          {saving ? "Saving..." : "Create Assessment"}
        </button>
      </form>
    </div>
  );
}

export default Assessments;
