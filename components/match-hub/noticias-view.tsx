"use client"

import { useState, useEffect, useCallback } from "react"
import { ArrowLeft, Clock, Plus, Edit, Trash2, Eye, EyeOff } from "lucide-react"
import { apiClient } from "@/lib/api-client"
import { useAuth } from "@/lib/auth-context"
import { Button } from "@/components/ui/button"
import { AdBanner } from "@/components/match-hub/ad-banner"

interface Article {
  id: string
  author_id: string
  author_name: string
  title: string
  summary: string
  content: string
  image_url?: string
  category: string
  published: boolean
  published_at?: string
  created_at: string
}

type Category = "all" | "tendencias" | "tech" | "seguridad" | "negocios" | "general"

const CATEGORIES: { id: Category; label: string }[] = [
  { id: "all", label: "Todo" },
  { id: "tendencias", label: "Tendencias" },
  { id: "tech", label: "Tech" },
  { id: "seguridad", label: "Seguridad" },
  { id: "negocios", label: "Negocios" },
]

const CATEGORY_COLORS: Record<string, string> = {
  tendencias: "bg-orange-100 text-orange-700",
  tech: "bg-blue-100 text-blue-700",
  seguridad: "bg-red-100 text-red-700",
  negocios: "bg-green-100 text-green-700",
  general: "bg-gray-100 text-gray-700",
}

function timeAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime()
  const hrs = Math.floor(diff / 3600000)
  if (hrs < 1) return "hace menos de 1h"
  if (hrs < 24) return `hace ${hrs}h`
  return `hace ${Math.floor(hrs / 24)}d`
}

function AdminArticleForm({ article, onSave, onCancel }: {
  article?: Partial<Article>
  onSave: (data: Partial<Article>) => Promise<void>
  onCancel: () => void
}) {
  const [form, setForm] = useState({
    title: article?.title ?? "",
    summary: article?.summary ?? "",
    content: article?.content ?? "",
    image_url: article?.image_url ?? "",
    category: article?.category ?? "general",
    published: article?.published ?? false,
  })
  const [saving, setSaving] = useState(false)

  const handleSubmit = async () => {
    if (!form.title || !form.summary || !form.content) return
    setSaving(true)
    try {
      await onSave({
        ...form,
        image_url: form.image_url || undefined,
      })
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="bg-card border border-border rounded-xl p-4 mb-4 space-y-3">
      <h3 className="font-semibold text-card-foreground">{article?.id ? "Editar artículo" : "Nuevo artículo"}</h3>
      <input
        placeholder="Título*"
        value={form.title}
        onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
        className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
      />
      <input
        placeholder="Resumen* (máx 500 caracteres)"
        value={form.summary}
        onChange={(e) => setForm((f) => ({ ...f, summary: e.target.value.slice(0, 500) }))}
        className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
      />
      <textarea
        placeholder="Contenido completo*"
        value={form.content}
        onChange={(e) => setForm((f) => ({ ...f, content: e.target.value }))}
        rows={5}
        className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary resize-none"
      />
      <input
        placeholder="URL de imagen (opcional)"
        value={form.image_url}
        onChange={(e) => setForm((f) => ({ ...f, image_url: e.target.value }))}
        className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
      />
      <div className="flex gap-2">
        <select
          value={form.category}
          onChange={(e) => setForm((f) => ({ ...f, category: e.target.value }))}
          className="flex-1 border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
        >
          {["general", "tendencias", "tech", "seguridad", "negocios"].map((c) => (
            <option key={c} value={c}>{c}</option>
          ))}
        </select>
        <label className="flex items-center gap-2 text-sm text-card-foreground cursor-pointer">
          <input
            type="checkbox"
            checked={form.published}
            onChange={(e) => setForm((f) => ({ ...f, published: e.target.checked }))}
            className="rounded"
          />
          Publicar
        </label>
      </div>
      <div className="flex gap-2">
        <Button variant="outline" className="flex-1" onClick={onCancel}>Cancelar</Button>
        <Button
          className="flex-1 bg-primary hover:bg-primary/90"
          onClick={handleSubmit}
          disabled={saving || !form.title || !form.summary || !form.content}
        >
          {saving ? "Guardando..." : "Guardar"}
        </Button>
      </div>
    </div>
  )
}

interface AdData {
  id: string
  title: string
  description?: string
  image_url?: string
  cta_text: string
  cta_url: string
}

export function NoticiasView() {
  const { user } = useAuth()
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [category, setCategory] = useState<Category>("all")
  const [selected, setSelected] = useState<Article | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [editingArticle, setEditingArticle] = useState<Article | undefined>()
  const [activeAd, setActiveAd] = useState<AdData | null>(null)

  const isAdmin = user?.is_admin ?? false

  const loadArticles = useCallback(async () => {
    setLoading(true)
    try {
      const catParam = category !== "all" ? `&category=${category}` : ""
      const [data, ad] = await Promise.all([
        apiClient<Article[]>(`/api/v1/news?limit=20${catParam}`),
        apiClient<AdData | null>(`/api/v1/ads/active?badge=${user?.badge ?? "none"}`).catch(() => null),
      ])
      setArticles(data ?? [])
      setActiveAd(ad ?? null)
    } catch {
      //
    } finally {
      setLoading(false)
    }
  }, [category, user?.badge])

  useEffect(() => { loadArticles() }, [loadArticles])

  const handleCreate = async (data: Partial<Article>) => {
    await apiClient<Article>('/api/v1/news', { method: 'POST', body: data })
    setShowForm(false)
    loadArticles()
  }

  const handleUpdate = async (id: string, data: Partial<Article>) => {
    await apiClient<Article>(`/api/v1/news/${id}`, { method: 'PUT', body: data })
    setEditingArticle(undefined)
    loadArticles()
  }

  const handleDelete = async (id: string) => {
    if (!confirm("¿Eliminar este artículo?")) return
    await apiClient(`/api/v1/news/${id}`, { method: 'DELETE' })
    loadArticles()
  }

  const togglePublish = async (article: Article) => {
    await apiClient(`/api/v1/news/${article.id}`, {
      method: 'PUT',
      body: { published: !article.published },
    })
    loadArticles()
  }

  if (selected) {
    return (
      <div className="flex flex-col h-full">
        <div className="flex items-center gap-3 px-4 py-3 border-b border-border bg-card sticky top-0 z-10">
          <button onClick={() => setSelected(null)} className="text-muted-foreground hover:text-card-foreground">
            <ArrowLeft className="w-5 h-5" />
          </button>
          <h2 className="font-semibold text-sm text-card-foreground line-clamp-1 flex-1">{selected.title}</h2>
        </div>
        <div className="flex-1 overflow-y-auto pb-24">
          {selected.image_url && (
            <div className="w-full aspect-video bg-muted overflow-hidden">
              <img src={selected.image_url} alt={selected.title} className="w-full h-full object-cover" referrerPolicy="no-referrer" />
            </div>
          )}
          <div className="p-4">
            <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${CATEGORY_COLORS[selected.category] ?? "bg-gray-100 text-gray-700"}`}>
              {selected.category}
            </span>
            <h1 className="text-xl font-bold text-card-foreground mt-3 mb-2">{selected.title}</h1>
            <div className="flex items-center gap-2 text-xs text-muted-foreground mb-4">
              <Clock className="w-3 h-3" />
              <span>{timeAgo(selected.created_at)}</span>
              <span>·</span>
              <span>Por {selected.author_name}</span>
            </div>
            <p className="text-sm text-card-foreground leading-relaxed whitespace-pre-wrap">{selected.content}</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-border bg-card sticky top-0 z-10">
        <h1 className="text-lg font-bold text-card-foreground">Noticias</h1>
        {isAdmin && (
          <button
            onClick={() => { setEditingArticle(undefined); setShowForm(true) }}
            className="w-9 h-9 bg-primary rounded-full flex items-center justify-center text-primary-foreground shadow-sm"
          >
            <Plus className="w-4 h-4" />
          </button>
        )}
      </div>

      {/* Category pills */}
      <div className="flex gap-2 px-4 py-3 overflow-x-auto no-scrollbar border-b border-border/50 bg-card">
        {CATEGORIES.map((c) => (
          <button
            key={c.id}
            onClick={() => setCategory(c.id)}
            className={`px-4 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-colors ${
              category === c.id
                ? "bg-primary text-primary-foreground"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            }`}
          >
            {c.label}
          </button>
        ))}
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto pb-24 px-4 pt-3">
        {showForm && !editingArticle && (
          <AdminArticleForm
            onSave={handleCreate}
            onCancel={() => setShowForm(false)}
          />
        )}
        {editingArticle && (
          <AdminArticleForm
            article={editingArticle}
            onSave={(data) => handleUpdate(editingArticle.id, data)}
            onCancel={() => setEditingArticle(undefined)}
          />
        )}

        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin" />
          </div>
        ) : articles.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground text-sm">
              {isAdmin ? "No hay artículos. Crea el primero." : "No hay noticias disponibles."}
            </p>
          </div>
        ) : (
          articles.map((article, i) => (
            <div
              key={article.id}
              className={`bg-card border border-border rounded-xl overflow-hidden mb-3 cursor-pointer hover:shadow-sm transition-shadow ${
                !article.published ? "opacity-60" : ""
              }`}
            >
              {/* Featured hero for first article */}
              {i === 0 && article.image_url && (
                <div
                  className="w-full aspect-video bg-muted overflow-hidden"
                  onClick={() => setSelected(article)}
                >
                  <img src={article.image_url} alt={article.title} className="w-full h-full object-cover" referrerPolicy="no-referrer" />
                </div>
              )}
              <div className="p-4" onClick={() => setSelected(article)}>
                <div className="flex items-start gap-3">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className={`text-[10px] font-medium px-2 py-0.5 rounded-full ${CATEGORY_COLORS[article.category] ?? "bg-gray-100 text-gray-700"}`}>
                        {article.category}
                      </span>
                      {!article.published && isAdmin && (
                        <span className="text-[10px] font-medium px-2 py-0.5 rounded-full bg-yellow-100 text-yellow-700">Borrador</span>
                      )}
                    </div>
                    <h3 className="font-semibold text-sm text-card-foreground line-clamp-2 mb-1">{article.title}</h3>
                    <p className="text-xs text-muted-foreground line-clamp-2">{article.summary}</p>
                    <div className="flex items-center gap-2 mt-2 text-xs text-muted-foreground">
                      <Clock className="w-3 h-3" />
                      <span>{timeAgo(article.created_at)}</span>
                    </div>
                  </div>
                  {i > 0 && article.image_url && (
                    <div className="w-20 h-16 rounded-lg bg-muted overflow-hidden flex-shrink-0">
                      <img src={article.image_url} alt="" className="w-full h-full object-cover" referrerPolicy="no-referrer" />
                    </div>
                  )}
                </div>
              </div>

              {/* Admin actions */}
              {isAdmin && (
                <div className="flex gap-2 px-4 pb-3" onClick={(e) => e.stopPropagation()}>
                  <button
                    onClick={() => togglePublish(article)}
                    className="flex items-center gap-1 text-xs text-muted-foreground hover:text-card-foreground"
                  >
                    {article.published ? <EyeOff className="w-3.5 h-3.5" /> : <Eye className="w-3.5 h-3.5" />}
                    {article.published ? "Ocultar" : "Publicar"}
                  </button>
                  <button
                    onClick={() => { setEditingArticle(article); setShowForm(false) }}
                    className="flex items-center gap-1 text-xs text-muted-foreground hover:text-card-foreground"
                  >
                    <Edit className="w-3.5 h-3.5" />
                    Editar
                  </button>
                  <button
                    onClick={() => handleDelete(article.id)}
                    className="flex items-center gap-1 text-xs text-red-400 hover:text-red-600"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                    Eliminar
                  </button>
                </div>
              )}
            </div>
          ))
        )}

        {/* Ad banner at bottom of feed */}
        {activeAd && !loading && <AdBanner ad={activeAd} />}
      </div>
    </div>
  )
}
