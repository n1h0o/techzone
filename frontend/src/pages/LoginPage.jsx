import { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api/api";
import { useAuth } from "../context/AuthContext";

function LoginPage() {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");

  const { setUser } = useAuth();
  const navigate = useNavigate();

  async function handleSubmit(e) {
    e.preventDefault();

    try {
      const res = await api.post("/login", {
        login,
        password,
      });

      localStorage.setItem(
        "token",
        res.data.token,
      );

      const me = await api.get("/me", {
        headers: {
          Authorization: `Bearer ${res.data.token}`,
        },
      });

      setUser(me.data);

      alert("Успешный вход");

      navigate("/");
    } catch (err) {
      console.error(err);
      alert("Ошибка авторизации");
    }
  }

  return (
    <div className="container">
      <h1>Вход</h1>

      <form onSubmit={handleSubmit}>
        <input
          value={login}
          onChange={(e) =>
            setLogin(e.target.value)
          }
          placeholder="Логин"
        />

        <input
          type="password"
          value={password}
          onChange={(e) =>
            setPassword(e.target.value)
          }
          placeholder="Пароль"
        />

        <button type="submit">
          Войти
        </button>
      </form>
    </div>
  );
}

export default LoginPage;