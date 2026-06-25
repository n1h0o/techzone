import { useState } from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import api from "../api/api";

function AdminProductsPage() {
  const { user } = useAuth();

  const [name, setName] = useState("");
  const [description, setDescription] =
    useState("");
  const [price, setPrice] = useState("");
  const [stock, setStock] = useState("");

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (user.role !== "admin") {
    return (
      <div className="container">
        <h1>Доступ запрещён</h1>
      </div>
    );
  }

  async function handleSubmit(e) {
    e.preventDefault();

    try {
      await api.post("/products", {
        name,
        description,
        price: Number(price),
        stock: Number(stock),
      });

      alert("Товар успешно создан");

      setName("");
      setDescription("");
      setPrice("");
      setStock("");
    } catch (err) {
      console.error(err);
      alert("Ошибка создания товара");
    }
  }

  return (
    <div className="container">
      <h1>Добавление товара</h1>

      <form onSubmit={handleSubmit}>
        <input
          value={name}
          onChange={(e) =>
            setName(e.target.value)
          }
          placeholder="Название"
        />

        <input
          value={description}
          onChange={(e) =>
            setDescription(e.target.value)
          }
          placeholder="Описание"
        />

        <input
          type="number"
          value={price}
          onChange={(e) =>
            setPrice(e.target.value)
          }
          placeholder="Цена"
        />

        <input
          type="number"
          value={stock}
          onChange={(e) =>
            setStock(e.target.value)
          }
          placeholder="Количество"
        />

        <button type="submit">
          Добавить товар
        </button>
      </form>
    </div>
  );
}

export default AdminProductsPage;