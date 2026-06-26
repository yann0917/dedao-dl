import { createContext, useCallback, useContext, useEffect, useMemo, useState, type PropsWithChildren } from "react"
import { api, type AuthMeData, type UserInfo } from "@/api"

type AuthContextValue = {
  user: UserInfo | null
  loggedIn: boolean
  loading: boolean
  refreshAuth: () => Promise<AuthMeData>
  completeLogin: (user?: UserInfo | null) => Promise<void>
  switchAccount: (uidHazy: string) => Promise<AuthMeData>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: PropsWithChildren) {
  const [user, setUser] = useState<UserInfo | null>(null)
  const [loggedIn, setLoggedIn] = useState(false)
  const [loading, setLoading] = useState(true)

  const applyAuthState = useCallback((data: AuthMeData) => {
    setLoggedIn(data.loggedIn)
    setUser(data.user ?? null)
    return data
  }, [])

  const refreshAuth = useCallback(async () => {
    try {
      const data = await api.auth.me()
      return applyAuthState(data)
    } catch {
      const fallback = { loggedIn: false, user: undefined }
      return applyAuthState(fallback)
    } finally {
      setLoading(false)
    }
  }, [applyAuthState])

  useEffect(() => {
    void refreshAuth()
  }, [refreshAuth])

  const completeLogin = useCallback(
    async (nextUser?: UserInfo | null) => {
      if (nextUser) {
        applyAuthState({ loggedIn: true, user: nextUser })
        setLoading(false)
        return
      }

      await refreshAuth()
    },
    [applyAuthState, refreshAuth],
  )

  const switchAccount = useCallback(
    async (uidHazy: string) => {
      const data = await api.auth.switchAccount(uidHazy)
      applyAuthState(data)
      return data
    },
    [applyAuthState],
  )

  const logout = useCallback(async () => {
    await api.auth.logout()
    applyAuthState({ loggedIn: false, user: undefined })
  }, [applyAuthState])

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      loggedIn,
      loading,
      refreshAuth,
      completeLogin,
      switchAccount,
      logout,
    }),
    [completeLogin, loading, loggedIn, logout, refreshAuth, switchAccount, user],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error("useAuth 必须在 AuthProvider 内部使用")
  }

  return context
}
