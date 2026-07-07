import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Toaster } from "sonner";

import "./index.css";
import App from "./App.jsx";

import { AuthProvider } from "./context/AuthContext";
import { CartProvider } from "./context/CartContext";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <AuthProvider>
      <CartProvider>
        <App />

        <Toaster
          position="top-right"
          richColors
          closeButton
          duration={3000}
          expand={true}
        />
      </CartProvider>
    </AuthProvider>
  </StrictMode>
);