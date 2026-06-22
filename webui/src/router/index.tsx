import { createBrowserRouter, Navigate } from "react-router-dom"
import { ProtectedRoute } from "@/components/auth/ProtectedRoute"
import { AppShell } from "@/layouts/AppShell"
import { AudioDetailPage } from "@/pages/AudioDetailPage"
import { AudioGroupPage } from "@/pages/AudioGroupPage"
import { AudioArticleListPage } from "@/pages/AudioArticleListPage"
import { ArticleDetailPage } from "@/pages/ArticleDetailPage"
import { CategoryPage } from "@/pages/CategoryPage"
import { CoursePage } from "@/pages/CoursePage"
import { EbookDetailPage } from "@/pages/EbookDetailPage"
import { HomePage } from "@/pages/HomePage"
import { LoginPage } from "@/pages/LoginPage"
import { UserCenterPage } from "@/pages/UserCenterPage"

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/",
    element: (
      <ProtectedRoute>
        <AppShell />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: "courses",
        element: <CoursePage />,
      },
      {
        path: "category",
        element: <CategoryPage />,
      },
      {
        path: "ebooks/:enid",
        element: <EbookDetailPage />,
      },
      {
        path: "audios/:enid",
        element: <AudioDetailPage />,
      },
      {
        path: "audio-groups/:enid",
        element: <AudioGroupPage />,
      },
      {
        path: "audio-groups/:enid/articles",
        element: <AudioArticleListPage />,
      },
      {
        path: "articles/:aType/:enid",
        element: <ArticleDetailPage />,
      },
      {
        path: "user",
        element: <UserCenterPage />,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate replace to="/" />,
  },
])
