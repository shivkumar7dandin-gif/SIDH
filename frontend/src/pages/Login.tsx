import { useState } from "react";
import { loginUser } from "../api/authApi";

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [showPassword, setShowPassword] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      setMessage("");

      const data = await loginUser(username, password);

      localStorage.setItem("token", data.token);
      localStorage.setItem("role", data.role);

      console.log("Login successful:", data);

      window.location.href = "/dashboard";
    } catch (error) {
      console.error("Login error:", error);
      setMessage("Invalid username or password");
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="bg-white p-8 rounded-xl shadow-md w-full max-w-md">
        <h1 className="text-3xl font-bold text-center text-gray-800">
          The Schools Attendance Book
        </h1>

        <p className="text-center text-gray-500 mt-2">
          Login to continue
        </p>

        <form onSubmit={handleLogin} className="mt-6">
          <div className="mb-4">
            <label className="block mb-2 font-medium">
              Username
            </label>

            <input
              type="text"
              value={username}
              onChange={(e) =>
                setUsername(e.target.value)
              }
              className="w-full border p-3 rounded-lg"
              placeholder="Enter username"
              required
            />
          </div>

          <div className="mb-4">
            <label className="block mb-2 font-medium">
              Password
            </label>

            <div className="relative">
              <input
                type={
                  showPassword
                    ? "text"
                    : "password"
                }
                value={password}
                onChange={(e) =>
                  setPassword(e.target.value)
                }
                className="w-full border p-3 pr-20 rounded-lg"
                placeholder="Enter password"
                required
              />

              <button
                type="button"
                onClick={() =>
                  setShowPassword(
                    !showPassword,
                  )
                }
                className="absolute right-3 top-1/2 -translate-y-1/2 text-blue-600 text-sm font-medium hover:text-blue-800"
              >
                {showPassword
                  ? "Hide"
                  : "Show"}
              </button>
            </div>
          </div>

          <button
            type="submit"
            className="w-full bg-blue-600 text-white p-3 rounded-lg hover:bg-blue-700"
          >
            Login
          </button>
        </form>

        {message && (
          <p className="text-center mt-4 text-red-600">
            {message}
          </p>
        )}
      </div>
    </div>
  );
}

export default Login;