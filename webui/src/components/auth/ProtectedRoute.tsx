import { type PropsWithChildren } from "react"
import { Loader2 } from "lucide-react"
import { Navigate } from "react-router-dom"
import { useAuth } from "@/providers/AuthProvider"

export function ProtectedRoute({ children }: PropsWithChildren) {
  const { loading, loggedIn } = useAuth()

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center text-slate-500">
        <Loader2 className="mr-2 h-5 w-5 animate-spin" />
        正在检查登录状态...
      </div>
    )
  }

  if (!loggedIn) {
    return <Navigate replace to="/login" />
  }

  return children
}
