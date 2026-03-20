import type { AppRouterInstance } from "next/dist/shared/lib/app-router-context.shared-runtime"

export function goToProfile(id: string, router: AppRouterInstance): void {
  if (!id || !/^[0-9a-f-]{36}$/i.test(id)) return
  router.push(`/profile/${id}`)
}
