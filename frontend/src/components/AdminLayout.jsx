import { Link, Outlet } from "react-router-dom";

function AdminLayout() {
  return (
    <div className="admin-layout">

      <h1>Админ-панель</h1>

      <nav className="admin-nav">

        <Link to="/admin/products">
          Товары
        </Link>
        
      </nav>

      <Outlet />

    </div>
  );
}

export default AdminLayout;