import { useCallback, useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { api, type UserCenterResponse } from "@/api"
import { useAuth } from "@/providers/AuthProvider"

export function useUserCenter() {
  const navigate = useNavigate()
  const { user, switchAccount: switchActiveAccount, logout: logoutAuth } = useAuth()
  const [data, setData] = useState<UserCenterResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [switchingUID, setSwitchingUID] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await api.user.center()
      setData(response)
    } catch (err) {
      setError(err instanceof Error ? err.message : "用户中心加载失败")
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load, user?.uid_hazy])

  const switchAccount = useCallback(async (uidHazy: string) => {
    setSwitchingUID(uidHazy)
    setError(null)
    try {
      await switchActiveAccount(uidHazy)
      await load()
    } catch (err) {
      setError(err instanceof Error ? err.message : "账号切换失败")
    } finally {
      setSwitchingUID(null)
    }
  }, [load, switchActiveAccount])

  const logout = useCallback(async () => {
    await logoutAuth()
    navigate("/login", { replace: true })
  }, [logoutAuth, navigate])

  return {
    data,
    loading,
    switchingUID,
    error,
    reload: load,
    switchAccount,
    logout,
  }
}
