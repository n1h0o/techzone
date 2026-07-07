import { createContext, useContext, useEffect, useState } from "react";
import api from "../api/api";

const CartContext = createContext(null);

export function CartProvider({ children }) {
  const [cartCount, setCartCount] = useState(0);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    refreshCart();
  }, []);

  async function refreshCart() {
    const token = localStorage.getItem("token");

    if (!token) {
      setCartCount(0);
      return;
    }

    try {
      setLoading(true);

      const res = await api.get("/cart");

      const items = res.data.items || [];

      const count = items.reduce(
        (sum, item) => sum + item.quantity,
        0
      );

      setCartCount(count);
    } catch (err) {
      console.error(err);
      setCartCount(0);
    } finally {
      setLoading(false);
    }
  }

  function clearCart() {
    setCartCount(0);
  }

  return (
    <CartContext.Provider
      value={{
        cartCount,
        loading,
        refreshCart,
        clearCart,
      }}
    >
      {children}
    </CartContext.Provider>
  );
}

export function useCart() {
  const context = useContext(CartContext);

  if (!context) {
    throw new Error(
      "useCart must be used inside CartProvider"
    );
  }

  return context;
}