import {
  BrowserRouter,
  Routes,
  Route,
} from "react-router-dom";

import Navbar from "./components/Navbar";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import ProductsPage from "./pages/ProductsPage";
import CartPage from "./pages/CartPage";
import OrdersPage from "./pages/OrdersPage";

function App() {
  return (
    <BrowserRouter>
      <Navbar />

      <Routes>
        <Route
          path="/"
          element={<ProductsPage />}
        />

        <Route
          path="/products"
          element={<ProductsPage />}
        />

        <Route
          path="/login"
          element={<LoginPage />}
        />

        <Route
          path="/register"
          element={<RegisterPage />}
        />

        <Route
          path="/cart"
          element={<CartPage />}
        />

        <Route
          path="/orders"
          element={<OrdersPage />}
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;