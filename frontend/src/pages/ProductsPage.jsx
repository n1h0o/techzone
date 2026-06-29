import { useEffect, useState } from "react";
import api from "../api/api";

function ProductsPage() {
  const [products, setProducts] = useState([]);

  useEffect(() => {
    loadProducts();
  }, []);

  async function loadProducts() {
    try {
      const res = await api.get("/products");
      setProducts(res.data);
    } catch (err) {
      console.error(err);
    }
  }

  async function addToCart(productId) {
    const token = localStorage.getItem("token");

    if (!token) {
      alert(
        "Для добавления товара необходимо войти или зарегистрироваться."
      );
      return;
    }

    try {
      await api.post("/cart/items", {
        product_id: productId,
        quantity: 1,
      });

      alert("Товар добавлен в корзину");
    } catch (err) {
      console.error(err);
      alert("Не удалось добавить товар");
    }
  }

  return (
    <div className="container">
      <h1>Каталог товаров</h1>

      {products.length === 0 ? (
        <p>Товаров пока нет</p>
      ) : (
        <div className="products-grid">
          {products.map((product) => (
            <div
              key={product.id}
              className="product-card"
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
                  {Number(product.price).toLocaleString("ru-RU")} ₽
                </span>

                <span className="stock">
                  Остаток: {product.stock}
                </span>
              </div>

              <button
                className="buy-btn"
                onClick={() =>
                  addToCart(product.id)
                }
              >
                Добавить в корзину
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default ProductsPage;