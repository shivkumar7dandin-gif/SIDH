import axios from "axios";

import { useEffect, useState } from "react";

import { useNavigate, useParams } from "react-router-dom";

import {
  createTeacherAssignment,
  deleteTeacherAssignment,
  getTeacherAssignments,
  getTeacherById,
  updateTeacherStatus,
} from "../api/teacherApi";

import type { Teacher, TeacherClassAssignment } from "../api/teacherApi";

import { getClassrooms } from "../api/classroomApi";

import type { Classroom } from "../api/classroomApi";

function TeacherProfile() {
  const navigate = useNavigate();

  const { id } = useParams();

  const [teacher, setTeacher] = useState<Teacher | null>(null);

  const [loading, setLoading] = useState(true);

  const [updatingStatus, setUpdatingStatus] = useState(false);

  const [error, setError] = useState("");

  // ------------------------------------------------
  // TEACHER CLASS ASSIGNMENT STATE
  // ------------------------------------------------

  const [assignments, setAssignments] = useState<TeacherClassAssignment[]>([]);

  const [classrooms, setClassrooms] = useState<Classroom[]>([]);

  const [selectedClassroomId, setSelectedClassroomId] = useState("");

  const [assignmentSubject, setAssignmentSubject] = useState("");

  const [loadingAssignments, setLoadingAssignments] = useState(false);

  const [savingAssignment, setSavingAssignment] = useState(false);

  // ------------------------------------------------
  // LOAD TEACHER
  // ------------------------------------------------

  const loadTeacher = async () => {
    if (!id) {
      setError("Teacher ID not found");
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      setError("");

      const data = await getTeacherById(id);

      setTeacher(data);
    } catch (err) {
      console.error("Get teacher profile error:", err);

      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Failed to load teacher profile");
      } else {
        setError("Failed to load teacher profile");
      }
    } finally {
      setLoading(false);
    }
  };

  // ------------------------------------------------
  // LOAD TEACHER ASSIGNMENTS
  // ------------------------------------------------

  const loadAssignments = async () => {
    if (!id) {
      return;
    }

    try {
      setLoadingAssignments(true);

      const data = await getTeacherAssignments(id);

      setAssignments(data);
    } catch (err) {
      console.error("Get teacher assignments error:", err);

      if (axios.isAxiosError(err)) {
        setError(
          err.response?.data?.error || "Failed to load teacher assignments",
        );
      } else {
        setError("Failed to load teacher assignments");
      }
    } finally {
      setLoadingAssignments(false);
    }
  };

  // ------------------------------------------------
  // LOAD CLASSROOMS
  // ------------------------------------------------

  const loadClassrooms = async () => {
    try {
      const data = await getClassrooms();

      setClassrooms(data);
    } catch (err) {
      console.error("Get classrooms error:", err);

      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Failed to load classrooms");
      } else {
        setError("Failed to load classrooms");
      }
    }
  };

  // ------------------------------------------------
  // INITIAL LOAD
  // ------------------------------------------------

  useEffect(() => {
    loadTeacher();
    loadAssignments();
    loadClassrooms();
  }, [id]);

  // ------------------------------------------------
  // UPDATE TEACHER STATUS
  // ------------------------------------------------

  const handleStatusChange = async () => {
    if (!teacher || !id) {
      return;
    }

    const currentStatus = teacher.status === "left" ? "left" : "active";

    const newStatus = currentStatus === "active" ? "left" : "active";

    const message =
      newStatus === "left"
        ? "Are you sure you want to mark this teacher as Left?"
        : "Are you sure you want to reactivate this teacher?";

    const confirmed = window.confirm(message);

    if (!confirmed) {
      return;
    }

    try {
      setUpdatingStatus(true);

      setError("");

      await updateTeacherStatus(id, newStatus);

      await loadTeacher();
    } catch (err) {
      console.error("Update teacher status error:", err);

      if (axios.isAxiosError(err)) {
        setError(
          err.response?.data?.error || "Failed to update teacher status",
        );
      } else {
        setError("Failed to update teacher status");
      }
    } finally {
      setUpdatingStatus(false);
    }
  };

  // ------------------------------------------------
  // CREATE CLASS ASSIGNMENT
  // ------------------------------------------------

  const handleCreateAssignment = async () => {
    if (!id) {
      return;
    }

    if (!selectedClassroomId) {
      setError("Please select a classroom");
      return;
    }

    if (!assignmentSubject.trim()) {
      setError("Please enter a subject");
      return;
    }

    try {
      setSavingAssignment(true);

      setError("");

      await createTeacherAssignment({
        teacher_id: id,
        classroom_id: selectedClassroomId,
        subject: assignmentSubject.trim(),
      });

      setSelectedClassroomId("");

      setAssignmentSubject("");

      await loadAssignments();
    } catch (err) {
      console.error("Create teacher assignment error:", err);

      if (axios.isAxiosError(err)) {
        setError(
          err.response?.data?.error || "Failed to create teacher assignment",
        );
      } else {
        setError("Failed to create teacher assignment");
      }
    } finally {
      setSavingAssignment(false);
    }
  };

  // ------------------------------------------------
  // DELETE CLASS ASSIGNMENT
  // ------------------------------------------------

  const handleDeleteAssignment = async (assignmentId: string) => {
    const confirmed = window.confirm(
      "Are you sure you want to remove this class assignment?",
    );

    if (!confirmed) {
      return;
    }

    try {
      setError("");

      await deleteTeacherAssignment(assignmentId);

      await loadAssignments();
    } catch (err) {
      console.error("Delete teacher assignment error:", err);

      if (axios.isAxiosError(err)) {
        setError(
          err.response?.data?.error || "Failed to remove teacher assignment",
        );
      } else {
        setError("Failed to remove teacher assignment");
      }
    }
  };

  // ------------------------------------------------
  // GET CLASSROOM DISPLAY NAME
  // ------------------------------------------------

  const getClassroomName = (classroomId: string) => {
    const classroom = classrooms.find((item) => item.id === classroomId);

    if (!classroom) {
      return "Unknown Classroom";
    }

    return `${classroom.name} - ${classroom.section}`;
  };

  // ------------------------------------------------
  // DISPLAY VALUE
  // ------------------------------------------------

  const displayValue = (value: string | number | null | undefined) => {
    if (value === undefined || value === null || value === "") {
      return "Not Added";
    }

    return value;
  };

  // ------------------------------------------------
  // FORMAT JOINING DATE
  // ------------------------------------------------

  const formatJoiningDate = (date?: string) => {
    if (!date) {
      return "Not Added";
    }

    const parsedDate = new Date(`${date}T00:00:00`);

    if (Number.isNaN(parsedDate.getTime())) {
      return date;
    }

    return parsedDate.toLocaleDateString("en-IN", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    });
  };

  // ------------------------------------------------
  // LOADING
  // ------------------------------------------------

  if (loading) {
    return <div className="p-6">Loading teacher profile...</div>;
  }

  // ------------------------------------------------
  // TEACHER NOT FOUND
  // ------------------------------------------------

  if (!teacher) {
    return (
      <div className="p-6">
        <div className="max-w-5xl mx-auto">
          <div className="bg-red-50 border border-red-200 text-red-700 p-4 rounded-lg">
            {error || "Teacher profile not found"}
          </div>

          <button
            onClick={() => navigate("/teachers")}
            className="mt-4 text-blue-600 hover:text-blue-800"
          >
            ← Back to Teachers
          </button>
        </div>
      </div>
    );
  }

  const status = teacher.status === "left" ? "left" : "active";

  const address = teacher.address;

  const addressParts = [
    address?.house_no,
    address?.street,
    address?.village,
    address?.city,
    address?.state,
  ].filter(Boolean);

  return (
    <div className="p-6">
      <div className="max-w-6xl mx-auto">
        {/* BACK */}

        <button
          onClick={() => navigate("/teachers")}
          className="text-blue-600 hover:text-blue-800 mb-5"
        >
          ← Back to Teachers
        </button>

        {/* ERROR */}

        {error && (
          <div className="mb-5 border border-red-300 bg-red-50 text-red-700 p-4 rounded-lg">
            {error}
          </div>
        )}

        {/* PROFILE HEADER */}

        <div className="bg-white border shadow-sm rounded-xl p-6 mb-6">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-5">
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 rounded-full bg-blue-100 text-blue-700 flex items-center justify-center text-2xl font-bold">
                {teacher.name?.charAt(0).toUpperCase() || "T"}
              </div>

              <div>
                <div className="flex flex-wrap items-center gap-3">
                  <h1 className="text-3xl font-semibold text-gray-900">
                    {teacher.name}
                  </h1>

                  <span
                    className={
                      status === "active"
                        ? "bg-green-100 text-green-700 px-3 py-1 rounded-full text-sm font-medium"
                        : "bg-red-100 text-red-700 px-3 py-1 rounded-full text-sm font-medium"
                    }
                  >
                    {status === "active" ? "Active" : "Left"}
                  </span>
                </div>

                <p className="text-gray-500 mt-1">
                  Employee ID:{" "}
                  <span className="font-medium text-gray-700">
                    {displayValue(teacher.employee_id)}
                  </span>
                </p>

                <p className="text-gray-500">{displayValue(teacher.subject)}</p>
              </div>
            </div>

            <div className="flex flex-wrap gap-3">
              <button
                onClick={() => navigate(`/teachers/${teacher.id}/edit`)}
                className="bg-blue-600 text-white px-5 py-2.5 rounded-lg hover:bg-blue-700"
              >
                Edit Profile
              </button>

              <button
                onClick={handleStatusChange}
                disabled={updatingStatus}
                className={
                  status === "active"
                    ? "bg-orange-600 text-white px-5 py-2.5 rounded-lg hover:bg-orange-700 disabled:opacity-50"
                    : "bg-green-600 text-white px-5 py-2.5 rounded-lg hover:bg-green-700 disabled:opacity-50"
                }
              >
                {updatingStatus
                  ? "Updating..."
                  : status === "active"
                    ? "Mark Left"
                    : "Reactivate"}
              </button>
            </div>
          </div>
        </div>

        {/* PROFILE INFORMATION */}

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* PERSONAL INFORMATION */}

          <section className="bg-white border shadow-sm rounded-xl p-6">
            <h2 className="text-lg font-semibold mb-5">Personal Information</h2>

            <div className="space-y-4">
              <ProfileRow
                label="Full Name"
                value={displayValue(teacher.name)}
              />

              <ProfileRow label="Age" value={displayValue(teacher.age)} />

              <ProfileRow label="Gender" value={displayValue(teacher.gender)} />

              <ProfileRow
                label="Username"
                value={displayValue(teacher.username)}
              />
            </div>
          </section>

          {/* PROFESSIONAL INFORMATION */}

          <section className="bg-white border shadow-sm rounded-xl p-6">
            <h2 className="text-lg font-semibold mb-5">
              Professional Information
            </h2>

            <div className="space-y-4">
              <ProfileRow
                label="Employee ID"
                value={displayValue(teacher.employee_id)}
              />

              <ProfileRow
                label="Subject"
                value={displayValue(teacher.subject)}
              />

              <ProfileRow
                label="Qualification"
                value={displayValue(teacher.qualification)}
              />

              <ProfileRow
                label="Experience"
                value={
                  teacher.experience_years !== undefined &&
                  teacher.experience_years !== null
                    ? `${teacher.experience_years} Years`
                    : "Not Added"
                }
              />

              <ProfileRow
                label="Joining Date"
                value={formatJoiningDate(teacher.joining_date)}
              />
            </div>
          </section>

          {/* CONTACT INFORMATION */}

          <section className="bg-white border shadow-sm rounded-xl p-6">
            <h2 className="text-lg font-semibold mb-5">Contact Information</h2>

            <div className="space-y-4">
              <ProfileRow label="Email" value={displayValue(teacher.email)} />

              <ProfileRow label="Phone" value={displayValue(teacher.phone)} />
            </div>
          </section>

          {/* ADDRESS */}

          <section className="bg-white border shadow-sm rounded-xl p-6">
            <h2 className="text-lg font-semibold mb-5">Address</h2>

            {addressParts.length > 0 ? (
              <div>
                <p className="text-gray-800 leading-7">
                  {addressParts.join(", ")}
                </p>

                {address?.pincode && (
                  <p className="text-gray-600 mt-1">
                    Pincode: {address.pincode}
                  </p>
                )}
              </div>
            ) : (
              <p className="text-gray-500">Address not added.</p>
            )}
          </section>
        </div>

        {/* ================================================= */}
        {/* CLASS ASSIGNMENTS */}
        {/* ================================================= */}

        <section className="bg-white border shadow-sm rounded-xl p-6 mt-6">
          <div className="mb-6">
            <h2 className="text-lg font-semibold text-gray-900">
              Class Assignments
            </h2>

            <p className="text-sm text-gray-500 mt-1">
              Assign classrooms and subjects handled by this teacher.
            </p>
          </div>

          {/* ASSIGNMENT FORM */}

          {status === "active" && (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
              {/* CLASSROOM */}

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Classroom
                </label>

                <select
                  value={selectedClassroomId}
                  onChange={(e) => setSelectedClassroomId(e.target.value)}
                  className="w-full border border-gray-300 rounded-lg px-3 py-2.5 bg-white"
                >
                  <option value="">Select Classroom</option>

                  {[...classrooms]
                    .sort((a, b) => {
                      const aNumber = parseInt(a.name);
                      const bNumber = parseInt(b.name);

                      return aNumber - bNumber;
                    })
                    .map((classroom) => (
                      <option key={classroom.id} value={classroom.id}>
                        {classroom.name} - {classroom.section}
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
                  value={assignmentSubject}
                  onChange={(e) => setAssignmentSubject(e.target.value)}
                  placeholder="Example: Science"
                  className="w-full border border-gray-300 rounded-lg px-3 py-2.5"
                />
              </div>

              {/* ASSIGN BUTTON */}

              <div className="flex items-end">
                <button
                  type="button"
                  onClick={handleCreateAssignment}
                  disabled={savingAssignment}
                  className="w-full bg-blue-600 text-white px-4 py-2.5 rounded-lg hover:bg-blue-700 disabled:opacity-50"
                >
                  {savingAssignment ? "Assigning..." : "Assign Class"}
                </button>
              </div>
            </div>
          )}

          {/* ASSIGNMENTS */}

          {loadingAssignments ? (
            <p className="text-gray-500">Loading class assignments...</p>
          ) : assignments.length === 0 ? (
            <div className="border border-dashed border-gray-300 rounded-lg p-6 text-center">
              <p className="text-gray-500">
                No classrooms assigned to this teacher.
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              {assignments.map((assignment) => (
                <div
                  key={assignment.id}
                  className="border border-gray-200 rounded-lg p-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4"
                >
                  <div>
                    <h3 className="font-semibold text-gray-900">
                      {getClassroomName(assignment.classroom_id)}
                    </h3>

                    <p className="text-sm text-gray-600 mt-1">
                      Subject:{" "}
                      <span className="font-medium text-gray-800">
                        {assignment.subject}
                      </span>
                    </p>
                  </div>

                  <button
                    type="button"
                    onClick={() => handleDeleteAssignment(assignment.id)}
                    className="text-red-600 hover:text-red-800 font-medium text-sm"
                  >
                    Remove
                  </button>
                </div>
              ))}
            </div>
          )}
        </section>
      </div>
    </div>
  );
}

// ------------------------------------------------
// PROFILE ROW
// ------------------------------------------------

type ProfileRowProps = {
  label: string;
  value: string | number;
};

function ProfileRow({ label, value }: ProfileRowProps) {
  return (
    <div className="flex justify-between gap-5 border-b border-gray-100 pb-3 last:border-b-0">
      <span className="text-gray-500">{label}</span>

      <span className="font-medium text-gray-900 text-right">{value}</span>
    </div>
  );
}

export default TeacherProfile;
