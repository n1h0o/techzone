import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import api from "../api/api";

function AdminProductsPage() {
  const { user } = useAuth();

  const [products, setProducts] = useState([]);

  const [editingId, setEditingId] = useState(null);

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [price, setPrice] = useState("");
  const [stock, setStock] = useState("");
  const [imageURL, setImageURL] = useState("");

  useEffect(() => {
    loadProducts();
  }, []);

  async function loadProducts() {
    try {
      const res = await api.get("/admin/products");
      setProducts(res.data);
    } catch (err) {
      console.error(err);
      alert("Не удалось загрузить товары");
    }
  }

  function clearForm() {
    setEditingId(null);

    setName("");
    setDescription("");
    setPrice("");
    setStock("");
    setImageURL("");
  }

  function editProduct(product) {
    setEditingId(product.id);

    setName(product.name);
    setDescription(product.description);
    setPrice(product.price);
    setStock(product.stock);
    setImageURL(product.image_url || "");

    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  }

  async function changeStatus(id, isActive) {
    try {
      await api.patch(`/products/${id}/status`, {
        is_active: isActive,
      });

      alert(
        isActive
          ? "Товар восстановлен"
          : "Товар скрыт"
      );

      loadProducts();
    } catch (err) {
      console.error(err);
      alert("Ошибка");
    }
  }

  async function handleSubmit(e) {
    e.preventDefault();

    try {
      if (editingId) {
        await api.put(`/products/${editingId}`, {
          name,
          description,
          price: Number(price),
          stock: Number(stock),
          image_url: imageURL,
        });

        alert("Товар обновлён");
      } else {
        await api.post("/products", {
          name,
          description,
          price: Number(price),
          stock: Number(stock),
          image_url: imageURL,
        });

        alert("Товар создан");
      }

      clearForm();
      loadProducts();
    } catch (err) {
      console.error(err);
      alert("Ошибка");
    }
  }

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
  return (
    <div className="container">
      <h1>Управление товарами</h1>

      <form className="admin-form" onSubmit={handleSubmit}>
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Название"
        />

        <input
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Описание"
        />

        <input
          type="number"
          value={price}
          onChange={(e) => setPrice(e.target.value)}
          placeholder="Цена"
        />

        <input
          type="number"
          value={stock}
          onChange={(e) => setStock(e.target.value)}
          placeholder="Количество"
        />

        <input
          value={imageURL}
          onChange={(e) => setImageURL(e.target.value)}
          placeholder="URL изображения"
        />

        {imageURL && (
          <img
            src={imageURL}
            alt="preview"
            style={{
              width: "200px",
              height: "200px",
              objectFit: "cover",
              borderRadius: "12px",
              marginBottom: "15px",
            }}
          />
        )}

        <button
          className="btn-success"
          type="submit"
        >
          {editingId
            ? "Сохранить изменения"
            : "Добавить товар"}
        </button>

        {editingId && (
          <button
            className="btn-secondary"
            type="button"
            onClick={clearForm}
          >
            Отмена
          </button>
        )}
      </form>

      <hr />
            {products.length === 0 ? (
        <p>Товаров пока нет.</p>
      ) : (
        <div className="products-grid">
          {products.map((product) => (
            <div
              key={product.id}
              className={`product-card ${
                !product.is_active ? "product-hidden" : ""
              }`}
            >
              <div className="product-image">
                {product.image_url ? (
                  <img
                    src={product.image_url}
                    alt={product.name}
                  />
                ) : (
                  <div className="no-image">
                    📦
                  </div>
                )}
              </div>

              <h2>{product.name}</h2>

              <p className="description">
                {product.description}
              </p>

              <div className="product-info">
                <span className="price">
                  {product.price} ₽
                </span>

                <span className="stock">
                  Остаток: {product.stock}
                </span>
              </div>

              <p>
                Статус:{" "}
                {product.is_active ? (
                  <span className="status-active">
                    Активен
                  </span>
                ) : (
                  <span className="status-hidden">
                    Скрыт
                  </span>
                )}
              </p>

              <div className="admin-actions">
                <button
                  onClick={() =>
                    editProduct(product)
                  }
                >
                  Редактировать
                </button>

                {product.is_active ? (
                  <button
                    className="btn-danger"
                    onClick={() =>
                      changeStatus(
                        product.id,
                        false
                      )
                    }
                  >
                    Скрыть
                  </button>
                ) : (
                  <button
                    className="btn-success"
                    onClick={() =>
                      changeStatus(
                        product.id,
                        true
                      )
                    }
                  >
                    Восстановить
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default AdminProductsPage;