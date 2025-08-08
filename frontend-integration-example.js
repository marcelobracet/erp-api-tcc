// 🚀 Frontend Integration Example - ERP API
// Copy this code to your frontend project

// ============================================================================
// CONFIGURAÇÃO DO AMBIENTE
// ============================================================================

const API_CONFIG = {
  development: {
    baseURL: "http://localhost:8080",
    timeout: 10000,
  },
  qa: {
    baseURL: "https://qa-api.seudominio.com",
    timeout: 15000,
  },
  production: {
    baseURL: "https://api.seudominio.com",
    timeout: 20000,
  },
};

const currentEnv = process.env.REACT_APP_ENVIRONMENT || "development";
const config = API_CONFIG[currentEnv];

// ============================================================================
// CLIENTE API COM AUTENTICAÇÃO AUTOMÁTICA
// ============================================================================

class ErpApiClient {
  constructor() {
    this.baseURL = config.baseURL;
    this.timeout = config.timeout;
  }

  // Método principal para fazer requisições
  async request(endpoint, options = {}) {
    const token = this.getAccessToken();

    const defaultHeaders = {
      "Content-Type": "application/json",
    };

    if (token) {
      defaultHeaders["Authorization"] = `Bearer ${token}`;
    }

    const requestConfig = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    };

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), this.timeout);

      const response = await fetch(`${this.baseURL}${endpoint}`, {
        ...requestConfig,
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      // Se token expirou, tentar renovar
      if (response.status === 401) {
        const refreshed = await this.refreshToken();
        if (refreshed) {
          // Reenviar requisição com novo token
          requestConfig.headers.Authorization = `Bearer ${this.getAccessToken()}`;
          return await fetch(`${this.baseURL}${endpoint}`, {
            ...requestConfig,
            signal: controller.signal,
          });
        } else {
          // Falha no refresh, fazer logout
          this.logout();
          throw new Error("Authentication failed");
        }
      }

      return response;
    } catch (error) {
      if (error.name === "AbortError") {
        throw new Error("Request timeout");
      }
      throw error;
    }
  }

  // Renovar token automaticamente
  async refreshToken() {
    try {
      const refreshToken = this.getRefreshToken();
      if (!refreshToken) {
        return false;
      }

      const response = await fetch(`${this.baseURL}/api/v1/auth/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (response.ok) {
        const data = await response.json();
        this.setTokens(data.access_token, data.refresh_token);
        return true;
      }

      return false;
    } catch (error) {
      console.error("Token refresh failed:", error);
      return false;
    }
  }

  // ============================================================================
  // MÉTODOS DE AUTENTICAÇÃO
  // ============================================================================

  async login(email, password) {
    try {
      const response = await this.request("/api/v1/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        const data = await response.json();
        this.setTokens(data.access_token, data.refresh_token);
        this.setUser(data.user);
        return data;
      } else {
        const errorData = await response.json();
        throw new Error(errorData.error || "Login failed");
      }
    } catch (error) {
      console.error("Login error:", error);
      throw error;
    }
  }

  async logout() {
    this.clearTokens();
    this.clearUser();
    // Redirecionar para login
    window.location.href = "/login";
  }

  // ============================================================================
  // MÉTODOS DE USUÁRIO
  // ============================================================================

  async registerUser(userData) {
    const response = await this.request("/api/v1/users/register", {
      method: "POST",
      body: JSON.stringify(userData),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || "Registration failed");
    }

    return await response.json();
  }

  async getUserProfile() {
    const response = await this.request("/api/v1/users/profile");

    if (!response.ok) {
      throw new Error("Failed to get user profile");
    }

    return await response.json();
  }

  async updateUser(userId, userData) {
    const response = await this.request(`/api/v1/users/${userId}`, {
      method: "PUT",
      body: JSON.stringify(userData),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || "Update failed");
    }

    return await response.json();
  }

  async deleteUser(userId) {
    const response = await this.request(`/api/v1/users/${userId}`, {
      method: "DELETE",
    });

    if (!response.ok) {
      throw new Error("Failed to delete user");
    }

    return true;
  }

  async listUsers(limit = 10, offset = 0) {
    const response = await this.request(
      `/api/v1/users?limit=${limit}&offset=${offset}`
    );

    if (!response.ok) {
      throw new Error("Failed to list users");
    }

    return await response.json();
  }

  async getUserCount() {
    const response = await this.request("/api/v1/users/count");

    if (!response.ok) {
      throw new Error("Failed to get user count");
    }

    return await response.json();
  }

  // ============================================================================
  // UTILITÁRIOS DE TOKEN E USUÁRIO
  // ============================================================================

  getAccessToken() {
    return localStorage.getItem("erp_access_token");
  }

  getRefreshToken() {
    return localStorage.getItem("erp_refresh_token");
  }

  setTokens(accessToken, refreshToken) {
    localStorage.setItem("erp_access_token", accessToken);
    localStorage.setItem("erp_refresh_token", refreshToken);
  }

  clearTokens() {
    localStorage.removeItem("erp_access_token");
    localStorage.removeItem("erp_refresh_token");
  }

  getUser() {
    const userStr = localStorage.getItem("erp_user");
    return userStr ? JSON.parse(userStr) : null;
  }

  setUser(user) {
    localStorage.setItem("erp_user", JSON.stringify(user));
  }

  clearUser() {
    localStorage.removeItem("erp_user");
  }

  isAuthenticated() {
    return !!this.getAccessToken();
  }

  hasRole(requiredRole) {
    const user = this.getUser();
    return user && user.role === requiredRole;
  }

  hasAnyRole(roles) {
    const user = this.getUser();
    return user && roles.includes(user.role);
  }
}

// ============================================================================
// EXEMPLOS DE USO
// ============================================================================

// Criar instância do cliente
const apiClient = new ErpApiClient();

// Exemplo 1: Login
async function handleLogin(email, password) {
  try {
    const result = await apiClient.login(email, password);
    console.log("Login successful:", result.user);

    // Redirecionar baseado no role
    if (result.user.role === "admin") {
      window.location.href = "/admin/dashboard";
    } else {
      window.location.href = "/dashboard";
    }
  } catch (error) {
    console.error("Login failed:", error.message);
    // Mostrar erro para o usuário
    showError(error.message);
  }
}

// Exemplo 2: Registrar usuário
async function handleRegister(userData) {
  try {
    const result = await apiClient.registerUser(userData);
    console.log("User registered:", result);

    // Fazer login automaticamente
    await handleLogin(userData.email, userData.password);
  } catch (error) {
    console.error("Registration failed:", error.message);
    showError(error.message);
  }
}

// Exemplo 3: Obter perfil do usuário
async function loadUserProfile() {
  try {
    const profile = await apiClient.getUserProfile();
    console.log("User profile:", profile);

    // Atualizar UI com dados do usuário
    updateUserInfo(profile);
  } catch (error) {
    console.error("Failed to load profile:", error.message);
  }
}

// Exemplo 4: Listar usuários (admin)
async function loadUsers(page = 1, limit = 10) {
  try {
    const offset = (page - 1) * limit;
    const result = await apiClient.listUsers(limit, offset);
    console.log("Users:", result);

    // Atualizar tabela de usuários
    updateUsersTable(result.users);
    updatePagination(result.total, page, limit);
  } catch (error) {
    console.error("Failed to load users:", error.message);
  }
}

// Exemplo 5: Atualizar usuário
async function handleUpdateUser(userId, userData) {
  try {
    const result = await apiClient.updateUser(userId, userData);
    console.log("User updated:", result);

    // Atualizar UI
    showSuccess("User updated successfully");
    loadUsers(); // Recarregar lista
  } catch (error) {
    console.error("Update failed:", error.message);
    showError(error.message);
  }
}

// Exemplo 6: Deletar usuário
async function handleDeleteUser(userId) {
  if (!confirm("Are you sure you want to delete this user?")) {
    return;
  }

  try {
    await apiClient.deleteUser(userId);
    console.log("User deleted successfully");

    // Atualizar UI
    showSuccess("User deleted successfully");
    loadUsers(); // Recarregar lista
  } catch (error) {
    console.error("Delete failed:", error.message);
    showError(error.message);
  }
}

// ============================================================================
// COMPONENTE REACT DE EXEMPLO
// ============================================================================

// Exemplo de componente React para login
function LoginForm() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      await handleLogin(email, password);
    } catch (error) {
      console.error("Login error:", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
        required
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
        required
      />
      <button type="submit" disabled={loading}>
        {loading ? "Logging in..." : "Login"}
      </button>
    </form>
  );
}

// Exemplo de componente React para dashboard
function Dashboard() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const profile = await apiClient.getUserProfile();
        setUser(profile);
      } catch (error) {
        console.error("Failed to load profile:", error);
      } finally {
        setLoading(false);
      }
    };

    if (apiClient.isAuthenticated()) {
      loadProfile();
    } else {
      window.location.href = "/login";
    }
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>Welcome, {user?.name}!</h1>
      <p>Role: {user?.role}</p>
      <p>Email: {user?.email}</p>

      {user?.role === "admin" && (
        <div>
          <h2>Admin Panel</h2>
          <button onClick={() => (window.location.href = "/admin/users")}>
            Manage Users
          </button>
        </div>
      )}

      <button onClick={() => apiClient.logout()}>Logout</button>
    </div>
  );
}

// ============================================================================
// UTILITÁRIOS DE UI
// ============================================================================

function showSuccess(message) {
  // Implementar notificação de sucesso
  console.log("Success:", message);
  // Exemplo: toast.success(message);
}

function showError(message) {
  // Implementar notificação de erro
  console.error("Error:", message);
  // Exemplo: toast.error(message);
}

function updateUserInfo(user) {
  // Atualizar elementos da UI com dados do usuário
  const nameElement = document.getElementById("user-name");
  const emailElement = document.getElementById("user-email");
  const roleElement = document.getElementById("user-role");

  if (nameElement) nameElement.textContent = user.name;
  if (emailElement) emailElement.textContent = user.email;
  if (roleElement) roleElement.textContent = user.role;
}

function updateUsersTable(users) {
  // Atualizar tabela de usuários
  const tableBody = document.getElementById("users-table-body");
  if (!tableBody) return;

  tableBody.innerHTML = users
    .map(
      (user) => `
    <tr>
      <td>${user.name}</td>
      <td>${user.email}</td>
      <td>${user.role}</td>
      <td>${user.is_active ? "Active" : "Inactive"}</td>
      <td>
        <button onclick="editUser('${user.id}')">Edit</button>
        <button onclick="deleteUser('${user.id}')">Delete</button>
      </td>
    </tr>
  `
    )
    .join("");
}

function updatePagination(total, currentPage, limit) {
  // Atualizar paginação
  const totalPages = Math.ceil(total / limit);
  const paginationElement = document.getElementById("pagination");

  if (!paginationElement) return;

  let paginationHTML = "";
  for (let i = 1; i <= totalPages; i++) {
    paginationHTML += `
      <button 
        onclick="loadUsers(${i})" 
        class="${i === currentPage ? "active" : ""}"
      >
        ${i}
      </button>
    `;
  }

  paginationElement.innerHTML = paginationHTML;
}

// ============================================================================
// EXPORTAR PARA USO EM OUTROS ARQUIVOS
// ============================================================================

export default apiClient;
export {
  ErpApiClient,
  handleLogin,
  handleRegister,
  loadUserProfile,
  loadUsers,
};
