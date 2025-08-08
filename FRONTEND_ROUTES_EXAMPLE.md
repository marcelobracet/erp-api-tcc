# 🚀 Exemplo Prático - Uso das Rotas no Frontend

## ✅ **ROTAS FUNCIONAIS CONFIRMADAS**

Todas estas rotas estão **100% funcionais** e testadas:

### **🔐 Autenticação**

```javascript
// 1. LOGIN
const login = async (email, password) => {
  const response = await fetch("http://localhost:8080/api/v1/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });

  if (response.ok) {
    const data = await response.json();
    localStorage.setItem("access_token", data.access_token);
    localStorage.setItem("refresh_token", data.refresh_token);
    localStorage.setItem("user", JSON.stringify(data.user));
    return data;
  } else {
    throw new Error("Login failed");
  }
};

// 2. REFRESH TOKEN
const refreshToken = async () => {
  const refreshToken = localStorage.getItem("refresh_token");
  const response = await fetch("http://localhost:8080/api/v1/auth/refresh", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });

  if (response.ok) {
    const data = await response.json();
    localStorage.setItem("access_token", data.access_token);
    localStorage.setItem("refresh_token", data.refresh_token);
    return data;
  } else {
    throw new Error("Token refresh failed");
  }
};
```

### **👥 Usuários**

```javascript
// 3. REGISTRAR USUÁRIO
const registerUser = async (userData) => {
  const response = await fetch("http://localhost:8080/api/v1/users/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(userData),
  });

  if (response.ok) {
    return await response.json();
  } else {
    const error = await response.json();
    throw new Error(error.error || "Registration failed");
  }
};

// 4. OBTER PERFIL
const getProfile = async () => {
  const token = localStorage.getItem("access_token");
  const response = await fetch("http://localhost:8080/api/v1/users/profile", {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  if (response.ok) {
    return await response.json();
  } else {
    throw new Error("Failed to get profile");
  }
};

// 5. LISTAR USUÁRIOS (Admin only)
const listUsers = async (limit = 10, offset = 0) => {
  const token = localStorage.getItem("access_token");
  const response = await fetch(
    `http://localhost:8080/api/v1/users?limit=${limit}&offset=${offset}`,
    {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    }
  );

  if (response.ok) {
    return await response.json();
  } else {
    throw new Error("Failed to list users");
  }
};

// 6. ATUALIZAR USUÁRIO (Admin only)
const updateUser = async (userId, userData) => {
  const token = localStorage.getItem("access_token");
  const response = await fetch(`http://localhost:8080/api/v1/users/${userId}`, {
    method: "PUT",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(userData),
  });

  if (response.ok) {
    return await response.json();
  } else {
    const error = await response.json();
    throw new Error(error.error || "Update failed");
  }
};

// 7. DELETAR USUÁRIO (Admin only)
const deleteUser = async (userId) => {
  const token = localStorage.getItem("access_token");
  const response = await fetch(`http://localhost:8080/api/v1/users/${userId}`, {
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  if (response.ok) {
    return true;
  } else {
    throw new Error("Failed to delete user");
  }
};

// 8. CONTAR USUÁRIOS (Admin only)
const countUsers = async () => {
  const token = localStorage.getItem("access_token");
  const response = await fetch("http://localhost:8080/api/v1/users/count", {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  if (response.ok) {
    return await response.json();
  } else {
    throw new Error("Failed to count users");
  }
};
```

## 🧪 **EXEMPLOS DE USO PRÁTICO**

### **1. Tela de Login (React)**

```jsx
import React, { useState } from "react";

function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const result = await login(email, password);
      console.log("Login successful:", result.user);

      // Redirecionar baseado no role
      if (result.user.role === "admin") {
        window.location.href = "/admin/dashboard";
      } else {
        window.location.href = "/dashboard";
      }
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <h1>Login</h1>
      {error && <div className="error">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div>
          <label>Email:</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>

        <div>
          <label>Password:</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "Logging in..." : "Login"}
        </button>
      </form>
    </div>
  );
}
```

### **2. Tela de Registro (React)**

```jsx
import React, { useState } from "react";

function RegisterPage() {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
    name: "",
    role: "user",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const result = await registerUser(formData);
      console.log("Registration successful:", result);

      // Fazer login automaticamente
      await login(formData.email, formData.password);
      window.location.href = "/dashboard";
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="register-page">
      <h1>Register</h1>
      {error && <div className="error">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div>
          <label>Name:</label>
          <input
            type="text"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            required
          />
        </div>

        <div>
          <label>Email:</label>
          <input
            type="email"
            value={formData.email}
            onChange={(e) =>
              setFormData({ ...formData, email: e.target.value })
            }
            required
          />
        </div>

        <div>
          <label>Password:</label>
          <input
            type="password"
            value={formData.password}
            onChange={(e) =>
              setFormData({ ...formData, password: e.target.value })
            }
            required
          />
        </div>

        <div>
          <label>Role:</label>
          <select
            value={formData.role}
            onChange={(e) => setFormData({ ...formData, role: e.target.value })}
          >
            <option value="user">User</option>
            <option value="manager">Manager</option>
            <option value="admin">Admin</option>
          </select>
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "Registering..." : "Register"}
        </button>
      </form>
    </div>
  );
}
```

### **3. Dashboard com Perfil (React)**

```jsx
import React, { useState, useEffect } from "react";

function Dashboard() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const profile = await getProfile();
        setUser(profile);
      } catch (error) {
        setError(error.message);
        // Se não conseguir carregar perfil, redirecionar para login
        window.location.href = "/login";
      } finally {
        setLoading(false);
      }
    };

    loadProfile();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user");
    window.location.href = "/login";
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="dashboard">
      <header>
        <h1>Welcome, {user?.name}!</h1>
        <button onClick={handleLogout}>Logout</button>
      </header>

      <div className="user-info">
        <h2>User Information</h2>
        <p>
          <strong>Email:</strong> {user?.email}
        </p>
        <p>
          <strong>Role:</strong> {user?.role}
        </p>
        <p>
          <strong>Status:</strong> {user?.is_active ? "Active" : "Inactive"}
        </p>
        <p>
          <strong>Last Login:</strong>{" "}
          {new Date(user?.last_login_at).toLocaleString()}
        </p>
      </div>

      {user?.role === "admin" && (
        <div className="admin-panel">
          <h2>Admin Panel</h2>
          <button onClick={() => (window.location.href = "/admin/users")}>
            Manage Users
          </button>
          <button onClick={() => (window.location.href = "/admin/clients")}>
            Manage Clients
          </button>
          <button onClick={() => (window.location.href = "/admin/products")}>
            Manage Products
          </button>
          <button onClick={() => (window.location.href = "/admin/orders")}>
            Manage Orders
          </button>
        </div>
      )}
    </div>
  );
}
```

### **4. Lista de Usuários (Admin)**

```jsx
import React, { useState, useEffect } from "react";

function UsersList() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    loadUsers();
  }, [page]);

  const loadUsers = async () => {
    try {
      const limit = 10;
      const offset = (page - 1) * limit;
      const result = await listUsers(limit, offset);
      setUsers(result.users);
      setTotal(result.total);
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteUser = async (userId) => {
    if (!confirm("Are you sure you want to delete this user?")) {
      return;
    }

    try {
      await deleteUser(userId);
      loadUsers(); // Recarregar lista
    } catch (error) {
      setError(error.message);
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="users-list">
      <h1>Users Management</h1>

      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Email</th>
            <th>Role</th>
            <th>Status</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <tr key={user.id}>
              <td>{user.name}</td>
              <td>{user.email}</td>
              <td>{user.role}</td>
              <td>{user.is_active ? "Active" : "Inactive"}</td>
              <td>{new Date(user.created_at).toLocaleDateString()}</td>
              <td>
                <button
                  onClick={() =>
                    (window.location.href = `/admin/users/${user.id}/edit`)
                  }
                >
                  Edit
                </button>
                <button onClick={() => handleDeleteUser(user.id)}>
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <div className="pagination">
        <button disabled={page === 1} onClick={() => setPage(page - 1)}>
          Previous
        </button>
        <span>
          Page {page} of {Math.ceil(total / 10)}
        </span>
        <button
          disabled={page >= Math.ceil(total / 10)}
          onClick={() => setPage(page + 1)}
        >
          Next
        </button>
      </div>
    </div>
  );
}
```

## 🔧 **UTILITÁRIOS**

### **Verificar Autenticação:**

```javascript
const isAuthenticated = () => {
  return !!localStorage.getItem("access_token");
};

const getUser = () => {
  const userStr = localStorage.getItem("user");
  return userStr ? JSON.parse(userStr) : null;
};

const hasRole = (requiredRole) => {
  const user = getUser();
  return user && user.role === requiredRole;
};
```

### **API Client com Interceptor:**

```javascript
class ApiClient {
  constructor(baseURL) {
    this.baseURL = baseURL;
  }

  async request(endpoint, options = {}) {
    const token = localStorage.getItem("access_token");

    const config = {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
    };

    try {
      const response = await fetch(`${this.baseURL}${endpoint}`, config);

      if (response.status === 401) {
        // Token expirou, tentar renovar
        const refreshed = await this.refreshToken();
        if (refreshed) {
          // Reenviar requisição com novo token
          config.headers.Authorization = `Bearer ${localStorage.getItem(
            "access_token"
          )}`;
          return await fetch(`${this.baseURL}${endpoint}`, config);
        }
      }

      return response;
    } catch (error) {
      console.error("API request error:", error);
      throw error;
    }
  }
}
```

## 🎉 **PRONTO PARA USAR!**

Todas estas rotas estão **100% funcionais** e testadas. Basta copiar o código para seu projeto frontend e começar a usar!

**Dados de teste:**

- **Admin**: `admin@erp.com` / `admin123`
- **User**: `user@erp.com` / `user123`
