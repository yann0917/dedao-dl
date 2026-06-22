import { createBrowserRouter, Navigate } from "react-router-dom"
import { ProtectedRoute } from "@/components/auth/ProtectedRoute"
import { AppShell } from "@/layouts/AppShell"
import { AudioDetailPage } from "@/pages/AudioDetailPage"
import { AudioGroupPage } from "@/pages/AudioGroupPage"
import { AudioArticleListPage } from "@/pages/AudioArticleListPage"
import { ArticleDetailPage } from "@/pages/ArticleDetailPage"
import { CategoryPage } from "@/pages/CategoryPage"
import { CourseArticleListPage } from "@/pages/CourseArticleListPage"
import { CourseDetailPage } from "@/pages/CourseDetailPage"
import { CoursePage } from "@/pages/CoursePage"
import { EbookCommentPage } from "@/pages/EbookCommentPage"
import { EbookDetailPage } from "@/pages/EbookDetailPage"
import { HomePage } from "@/pages/HomePage"
import { LoginPage } from "@/pages/LoginPage"
import { PurchasedAudioPage } from "@/pages/PurchasedAudioPage"
import { PurchasedCompassPage } from "@/pages/PurchasedCompassPage"
import { PurchasedEbookPage } from "@/pages/PurchasedEbookPage"
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
        element: <Navigate replace to="/purchased/courses" />,
      },
      {
        path: "purchased",
        element: <Navigate replace to="/purchased/courses" />,
      },
      {
        path: "purchased/courses",
        element: <CoursePage />,
      },
      {
        path: "purchased/ebooks",
        element: <PurchasedEbookPage />,
      },
      {
        path: "purchased/audios",
        element: <PurchasedAudioPage />,
      },
      {
        path: "purchased/compass",
        element: <PurchasedCompassPage />,
      },
      {
        path: "courses/:enid",
        element: <CourseDetailPage />,
      },
      {
        path: "courses/:enid/articles",
        element: <CourseArticleListPage />,
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
        path: "ebooks/:enid/comments",
        element: <EbookCommentPage />,
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
