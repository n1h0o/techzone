import {
  useEffect,
  useState,
  useCallback,
} from "react";
import api from "../api/api";
import { CartContext } from "./cart-context";

export function CartProvider({ children }) {
  const [cartCount, setCartCount] = useState(0);
  const [loading, setLoading] = useState(false);

  const refreshCart = useCallback(async () => {
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
  }, []);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      void refreshCart();
    }, 0);

    return () => window.clearTimeout(timeoutId);
  }, [refreshCart]);

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
