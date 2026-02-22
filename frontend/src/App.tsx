import { useEffect } from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { Layout } from './components/Layout';
import { ProtectedRoute } from './components/ProtectedRoute';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { CatalogPage } from './pages/CatalogPage';
import { QuestionPage } from './pages/QuestionPage';
import { CreateQuestionPage } from './pages/CreateQuestionPage';
import { ReviewPage } from './pages/ReviewPage';
import { NotFoundPage } from './pages/NotFoundPage';
import { useAuthStore } from './stores/auth';

const router = createBrowserRouter([
  { path: '/login', element: <LoginPage /> },
  { path: '/register', element: <RegisterPage /> },
  {
    element: (
      <ProtectedRoute>
        <Layout />
      </ProtectedRoute>
    ),
    children: [
      { path: '/', element: <CatalogPage /> },
      { path: '/questions/new', element: <CreateQuestionPage /> },
      { path: '/questions/:id', element: <QuestionPage /> },
      { path: '/review', element: <ReviewPage /> },
    ],
  },
  { path: '*', element: <NotFoundPage /> },
]);

export default function App() {
  const init = useAuthStore((s) => s.init);

  useEffect(() => {
    init();
  }, [init]);

  return <RouterProvider router={router} />;
}
