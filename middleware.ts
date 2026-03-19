import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// Auth is fully client-side (token lives in memory, refresh cookie on backend).
// We skip server-side cookie checks since the refresh_token cookie is set by
// the backend (different port), which browsers don't reliably share across ports.
// Client-side protection is handled by the AuthGuard in the root layout.
export function middleware(_request: NextRequest) {
  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|.*\\..*).*)', '/'],
}
