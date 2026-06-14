import { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api/api";

function RegisterPage() {
  const navigate = useNavigate();

  const [login, setLogin] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] =
    useState("");

  async function handleSubmit(e) {
    e.preventDefault();

    try {
      await api.post("/register", {
        login,
        email,
        password,
      });

      alert(
        "Регистрация выполнена успешно"
      );

      navigate("/login");
    } catch (err) {
      console.error(err);

      if (err.response?.data) {
        alert(err.response.data);
        return;
      }

      alert("Ошибка регистрации");
    }
  }

  return (
    <div className="container">
      <h1>Регистрация</h1>

      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Логин"
          value={login}
          onChange={(e) =>
            setLogin(e.target.value)
          }
        />

        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) =>
            setEmail(e.target.value)
          }
        />

        <input
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={(e) =>
            setPassword(e.target.value)
          }
        />

        <button type="submit">
          Зарегистрироваться
        </button>
      </form>
    </div>
  );
}

export default RegisterPage;