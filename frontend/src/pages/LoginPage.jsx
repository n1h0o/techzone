import { useState } from "react";
import api from "../api/api";

function LoginPage() {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");

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

      alert("Успешный вход");
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