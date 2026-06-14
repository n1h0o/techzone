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
    try {
      await api.post("/cart/items", {
        product_id: productId,
        quantity: 1,
      });

      alert("Товар добавлен в корзину");
    } catch (err) {
      console.error(err);
    }
  }

  return (
    <div className="container">
      <h1>Каталог товаров</h1>

      {products.length === 0 ? (
        <p>Товаров пока нет</p>
      ) : (
        products.map((product) => (
          <div
            key={product.id}
            className="card"
          >
            <h2>{product.name}</h2>

            <p>{product.description}</p>

            <p>
              <strong>Цена:</strong>{" "}
              {product.price} ₽
            </p>

            <p>
              <strong>Остаток:</strong>{" "}
              {product.stock}
            </p>

            <button
              onClick={() =>
                addToCart(product.id)
              }
            >
              Добавить в корзину
            </button>
          </div>
        ))
      )}
    </div>
  );
}

export default ProductsPage;