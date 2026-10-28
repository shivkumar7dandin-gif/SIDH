import { useEffect, useState } from "react";
import axios from "axios";

import {
  createSubject,
  deleteSubject,
  getSubjects,
  type Subject,
} from "../../api/subjectApi";

const Subjects = () => {
  const [subjects, setSubjects] = useState<Subject[]>([]);

  const [name, setName] = useState("");
  const [code, setCode] = useState("");

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const loadSubjects = async () => {
    try {
      setLoading(true);
      setError("");

      const data = await getSubjects();

      const sortedSubjects = [...data].sort((a, b) =>
        a.name.localeCompare(b.name),
      );

      setSubjects(sortedSubjects);
    } catch (err) {
      console.error(err);
      setError("Failed to load subjects.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadSubjects();
  }, []);

  const handleCreate = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!name.trim()) {
      setError("Subject name is required.");
      return;
    }

    if (!code.trim()) {
      setError("Subject code is required.");
      return;
    }

    try {
      setSaving(true);
      setError("");
      setSuccess("");

      await createSubject({
        name: name.trim(),
        code: code.trim().toUpperCase(),
      });

      setName("");
      setCode("");

      setSuccess("Subject created successfully.");

      await loadSubjects();
    } catch (err) {
      console.error(err);

      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Failed to create subject.");
      } else {
        setError("Failed to create subject.");
      }
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (subject: Subject) => {
    const confirmed = window.confirm(
      `Are you sure you want to delete ${subject.name}?`,
    );

    if (!confirmed) {
      return;
    }

    try {
      setError("");
      setSuccess("");

      await deleteSubject(subject.id);

      setSuccess("Subject deleted successfully.");

      await loadSubjects();
    } catch (err) {
      console.error(err);

      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Failed to delete subject.");
      } else {
        setError("Failed to delete subject.");
      }
    }
  };

  return (
    <div
      style={{
        padding: "24px",
        maxWidth: "1100px",
        margin: "0 auto",
      }}
    >
      <div style={{ marginBottom: "24px" }}>
        <h1 style={{ marginBottom: "6px" }}>Subject Master</h1>

        <p style={{ margin: 0, color: "#666" }}>
          Manage subjects for your school.
        </p>
      </div>

      {error && (
        <div
          style={{
            padding: "12px",
            marginBottom: "16px",
            background: "#fee2e2",
            borderRadius: "6px",
          }}
        >
          {error}
        </div>
      )}

      {success && (
        <div
          style={{
            padding: "12px",
            marginBottom: "16px",
            background: "#dcfce7",
            borderRadius: "6px",
          }}
        >
          {success}
        </div>
      )}

      <div
        style={{
          border: "1px solid #ddd",
          borderRadius: "8px",
          padding: "20px",
          marginBottom: "28px",
        }}
      >
        <h2 style={{ marginTop: 0 }}>Add Subject</h2>

        <form onSubmit={handleCreate}>
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "minmax(200px, 1fr) minmax(150px, 250px)",
              gap: "16px",
              marginBottom: "16px",
            }}
          >
            <div>
              <label>Subject Name</label>

              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Example: Mathematics"
                style={{
                  display: "block",
                  width: "100%",
                  padding: "10px",
                  marginTop: "6px",
                  boxSizing: "border-box",
                }}
              />
            </div>

            <div>
              <label>Subject Code</label>

              <input
                type="text"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="Example: MATH"
                style={{
                  display: "block",
                  width: "100%",
                  padding: "10px",
                  marginTop: "6px",
                  boxSizing: "border-box",
                }}
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={saving}
            style={{
              padding: "10px 18px",
              cursor: saving ? "not-allowed" : "pointer",
            }}
          >
            {saving ? "Creating..." : "Add Subject"}
          </button>
        </form>
      </div>

      <div>
        <h2>Subjects</h2>

        {loading ? (
          <p>Loading subjects...</p>
        ) : subjects.length === 0 ? (
          <p>No subjects have been created yet.</p>
        ) : (
          <div
            style={{
              overflowX: "auto",
            }}
          >
            <table
              style={{
                width: "100%",
                borderCollapse: "collapse",
              }}
            >
              <thead>
                <tr>
                  <th style={tableHeaderStyle}>Subject</th>

                  <th style={tableHeaderStyle}>Code</th>

                  <th style={tableHeaderStyle}>Status</th>

                  <th style={tableHeaderStyle}>Action</th>
                </tr>
              </thead>

              <tbody>
                {subjects.map((subject) => (
                  <tr key={subject.id}>
                    <td style={tableCellStyle}>{subject.name}</td>

                    <td style={tableCellStyle}>{subject.code}</td>

                    <td style={tableCellStyle}>{subject.status}</td>

                    <td style={tableCellStyle}>
                      <button
                        type="button"
                        onClick={() => handleDelete(subject)}
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
    </div>
  );
};

const tableHeaderStyle: React.CSSProperties = {
  textAlign: "left",
  padding: "12px",
  borderBottom: "2px solid #ddd",
};

const tableCellStyle: React.CSSProperties = {
  padding: "12px",
  borderBottom: "1px solid #ddd",
};

export default Subjects;
