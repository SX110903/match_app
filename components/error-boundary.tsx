"use client"

import { Component, type ErrorInfo, type ReactNode } from "react"

interface Props {
  children: ReactNode
  fallback?: ReactNode
}

interface State {
  hasError: boolean
  error?: Error
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error("[ErrorBoundary]", error, info.componentStack)
  }

  render() {
    if (this.state.hasError) {
      return (
        this.props.fallback ?? (
          <div className="flex flex-col items-center justify-center h-full gap-4 p-8 text-center">
            <p className="text-destructive font-semibold text-base">Algo salió mal</p>
            <p className="text-sm text-muted-foreground">{this.state.error?.message}</p>
            <button
              className="text-sm text-primary underline"
              onClick={() => this.setState({ hasError: false })}
            >
              Reintentar
            </button>
          </div>
        )
      )
    }

    return this.props.children
  }
}
