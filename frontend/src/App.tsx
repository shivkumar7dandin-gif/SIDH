import {
  BrowserRouter,
  Navigate,
  Route,
  Routes,
} from "react-router-dom";

import Login from "./pages/Login";
import CollegeDashboard from "./pages/CollegeDashboard";

import Students from "./pages/Students";
import AddStudent from "./pages/AddStudent";
import EditStudent from "./pages/EditStudent";

import Teachers from "./pages/Teachers";
import AddTeacher from "./pages/AddTeacher";
import EditTeacher from "./pages/EditTeacher";

import Classrooms from "./pages/Classrooms";

import AdminLayout from "./layouts/AdminLayout";

function App() {
  const token = localStorage.getItem("token");
  const role = localStorage.getItem("role");

  const isCollegeAdmin =
    Boolean(token) && role === "college_admin";

  return (
    <BrowserRouter>
      <Routes>
        {/* LOGIN */}
        <Route
          path="/login"
          element={
            isCollegeAdmin ? (
              <Navigate
                to="/dashboard"
                replace
              />
            ) : (
              <Login />
            )
          }
        />

        {/* PROTECTED ADMIN ROUTES */}
        <Route
          element={
            isCollegeAdmin ? (
              <AdminLayout />
            ) : (
              <Navigate
                to="/login"
                replace
              />
            )
          }
        >
          {/* DASHBOARD */}
          <Route
            path="/dashboard"
            element={<CollegeDashboard />}
          />

          {/* STUDENTS */}
          <Route
            path="/students"
            element={<Students />}
          />

          <Route
            path="/students/add"
            element={<AddStudent />}
          />

          <Route
            path="/students/:id/edit"
            element={<EditStudent />}
          />

          {/* TEACHERS */}
          <Route
            path="/teachers"
            element={<Teachers />}
          />

          <Route
            path="/teachers/add"
            element={<AddTeacher />}
          />

          <Route
            path="/teachers/:id/edit"
            element={<EditTeacher />}
          />

          {/* CLASSROOMS */}
          <Route
            path="/classrooms"
            element={<Classrooms />}
          />
        </Route>

        {/* UNKNOWN ROUTE */}
        <Route
          path="*"
          element={
            isCollegeAdmin ? (
              <Navigate
                to="/dashboard"
                replace
              />
            ) : (
              <Navigate
                to="/login"
                replace
              />
            )
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;