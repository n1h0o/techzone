import { useEffect, useState } from "react";
import api from "../api/api";

function AdminPage() {
  const emptyForm = {
    name: "",
    description: "",
    price: "",
    stock: "",
  };

  const [products, setProducts] = useState([]);
  const [form, setForm] = useState(emptyForm);
  const [editingId, setEditingId] = useState(null);

  useEffect(() => {
    loadProducts();
  }, []);

  async function loadProducts() {
    try {
      const res = await api.get("/products");
      setProducts(res.data);
    } catch {
      alert("Не удалось загрузить товары");
    }
  }

  function handleChange(e) {
    setForm({
      ...form,
      [e.target.name]: e.target.value,
    });
  }

  async function handleSubmit(e) {
    e.preventDefault();

    try {
      if (editingId) {
        await api.put(`/products/${editingId}`, {
          name: form.name,
          description: form.description,
          price: Number(form.price),
          stock: Number(form.stock),
        });

        alert("Товар обновлен");
      } else {
        await api.post("/products", {
          name: form.name,
          description: form.description,
          price: Number(form.price),
          stock: Number(form.stock),
        });

        alert("Товар добавлен");
      }

      setForm(emptyForm);
      setEditingId(null);
      loadProducts();
    } catch (err) {
      console.error(err);
      alert("Ошибка");
    }
  }

  function editProduct(product) {
    setEditingId(product.id);

    setForm({
      name: product.name,
      description: product.description,
      price: product.price,
      stock: product.stock,
    });
  }

  async function deleteProduct(id) {
    if (!window.confirm("Удалить товар?")) {
      return;
    }

    try {
      await api.delete(`/products/${id}`);

      alert("Удалено");

      loadProducts();
    } catch (err) {
      console.error(err);
      alert("Ошибка удаления");
    }
  }

  return (
    <div className="container">
      <h1>Админ панель</h1>

      <form onSubmit={handleSubmit} className="admin-form">
        <input
          name="name"
          placeholder="Название"
          value={form.name}
          onChange={handleChange}
        />

        <input
          name="description"
          placeholder="Описание"
          value={form.description}
          onChange={handleChange}
        />

        <input
          name="price"
          type="number"
          placeholder="Цена"
          value={form.price}
          onChange={handleChange}
        />

        <input
          name="stock"
          type="number"
          placeholder="Количество"
          value={form.stock}
          onChange={handleChange}
        />

        <button type="submit">
          {editingId ? "Сохранить" : "Добавить"}
        </button>
      </form>

      <hr />

      {products.map((product) => (
        <div key={product.id} className="product-card">

          <h2>{product.name}</h2>

          <p>{product.description}</p>

          <p>{product.price} ₽</p>

          <p>Остаток: {product.stock}</p>

          <button
            onClick={() => editProduct(product)}
          >
            Редактировать
          </button>

          <button
            onClick={() => deleteProduct(product.id)}
          >
            Удалить
          </button>

        </div>
      ))}
    </div>
  );
}

export default AdminPage;