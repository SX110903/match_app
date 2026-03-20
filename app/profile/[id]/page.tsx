"use client"
import { useParams, useRouter } from "next/navigation"
import { PublicProfileView } from "@/components/match-hub/public-profile-view"

export default function ProfilePage() {
  const { id } = useParams<{ id: string }>()
  const router = useRouter()
  return (
    <div className="max-w-[430px] mx-auto h-screen bg-background">
      <PublicProfileView userId={id} onClose={() => router.back()} />
    </div>
  )
}
