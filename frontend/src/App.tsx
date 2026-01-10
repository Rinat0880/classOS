import { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useAuthStore } from './store/authStore';
import { ROUTES } from './constants';

// Pages (пока заглушки, создадим позже)
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Groups from './pages/Groups';
import Users from './pages/Users';
import Settings from './pages/Settings';

// Layout
import MainLayout from './components/layout/MainLayout';

// QueryClient для React Query
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000, // 5 минут
    },
  },
});

// Protected Route компонент
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  if (!isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} replace />;
  }

  return <>{children}</>;
};

function App() {
  const initialize = useAuthStore((state) => state.initialize);

  // Инициализация auth при загрузке
  useEffect(() => {
    initialize();
  }, [initialize]);

  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          {/* Публичный роут */}
          <Route path={ROUTES.LOGIN} element={<Login />} />

          {/* Защищенные роуты */}
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <MainLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Dashboard />} />
            <Route path={ROUTES.GROUPS} element={<Groups />} />
            <Route path={ROUTES.USERS} element={<Users />} />
            <Route path={ROUTES.SETTINGS} element={<Settings />} />
          </Route>

          {/* 404 - редирект на Dashboard */}
          <Route path="*" element={<Navigate to={ROUTES.DASHBOARD} replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
