import { NavLink, Outlet, useNavigate } from "react-router-dom";

function AdminLayout() {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("role");

    navigate("/login");
    window.location.reload();
  };

  const menuClass = ({ isActive }: { isActive: boolean }) =>
    `block w-full p-3 rounded-lg ${
      isActive ? "bg-slate-700 text-white" : "text-gray-200 hover:bg-slate-800"
    }`;

  return (
    <div className="min-h-screen bg-gray-100 flex">
      {/* SIDEBAR */}
      <aside className="w-64 bg-slate-900 text-white min-h-screen p-6">
        <h1 className="text-xl font-bold mb-10">The Schools Attendance Book</h1>

        <nav className="space-y-3">
          {/* DASHBOARD */}
          <NavLink to="/dashboard" className={menuClass}>
            Dashboard
          </NavLink>

          {/* STUDENTS */}
          <NavLink to="/students" className={menuClass}>
            Students
          </NavLink>

          {/* TEACHERS */}
          <NavLink to="/teachers" className={menuClass}>
            Teachers
          </NavLink>

          {/* CLASSROOMS */}
          <NavLink to="/classrooms" className={menuClass}>
            Classrooms
          </NavLink>

          {/* ATTENDANCE */}
          <NavLink to="/attendance" className={menuClass}>
            Attendance
          </NavLink>

          {/* ASSESSMENTS */}
          <NavLink to="/assessments" className={menuClass}>
            Assessments
          </NavLink>

          {/* CLASS SCHEDULE */}
          <NavLink to="/timetable" className={menuClass}>
            Class Schedule
          </NavLink>
        </nav>
      </aside>

      {/* MAIN CONTENT */}
      <div className="flex-1">
        {/* HEADER */}
        <header className="bg-white shadow px-8 py-4 flex justify-between items-center">
          <h2 className="text-xl font-semibold">The Schools Attendance Book</h2>

          <button
            onClick={handleLogout}
            className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg"
          >
            Logout
          </button>
        </header>

        {/* PAGE CONTENT */}
        <main className="p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

export default AdminLayout;
