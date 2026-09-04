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
      <aside className="w-64 bg-slate-900 text-white min-h-screen p-6">
        <h1 className="text-xl font-bold mb-10">The Schools Attendance Book</h1>

        <nav className="space-y-3">
          <NavLink to="/dashboard" className={menuClass}>
            Dashboard
          </NavLink>

          <NavLink to="/students" className={menuClass}>
            Students
          </NavLink>

          <NavLink to="/teachers" className={menuClass}>
            Teachers
          </NavLink>

          <NavLink to="/classrooms" className={menuClass}>
            Classrooms
          </NavLink>

          <NavLink to="/attendance" className={menuClass}>
            Attendance
          </NavLink>

          <NavLink to="/assessments" className={menuClass}>
            Assessments
          </NavLink>
        </nav>
      </aside>

      <div className="flex-1">
        <header className="bg-white shadow px-8 py-4 flex justify-between items-center">
          <h2 className="text-xl font-semibold">The Schools Attendance Book</h2>

          <button
            onClick={handleLogout}
            className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg"
          >
            Logout
          </button>
        </header>

        <Outlet />
      </div>
    </div>
  );
}

export default AdminLayout;
