import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import api from "../api/api";
import { useCart } from "../context/useCart";

function ProductsPage() {
  const [products, setProducts] = useState([]);

  const { refreshCart } = useCart();

  const loadProducts = useCallback(async () => {
    try {
      const res = await api.get("/products");
      setProducts(res.data);
    } catch (err) {
      console.error(err);
      toast.error("Не удалось загрузить товары");
    }
  }, []);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      void loadProducts();
    }, 0);

    return () => window.clearTimeout(timeoutId);
  }, [loadProducts]);

  async function addToCart(productId) {
    const token = localStorage.getItem("token");

    if (!token) {
      toast.warning(
        "Для добавления товара необходимо войти в аккаунт"
      );
      return;
    }

    try {
      await api.post("/cart/items", {
        product_id: productId,
        quantity: 1,
      });

      await refreshCart();

      toast.success("Товар добавлен в корзину 🛒");
    } catch (err) {
      console.error(err);

      if (err.response?.data) {
        toast.error(err.response.data);
      } else {
        toast.error("Не удалось добавить товар");
      }
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
                  {Number(product.price).toLocaleString(
                    "ru-RU"
                  )}{" "}
                  ₽
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
                🛒 Добавить в корзину
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default ProductsPage;
