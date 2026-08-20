import { Outlet, createRootRoute, redirect } from '@tanstack/react-router'
import { auth } from '../lib/auth'

const publicRoutes = new Set(['/signin'])

export const Route = createRootRoute({
  component: RootComponent,
  beforeLoad: async ({ location }) => {
    const isPublicRoute = publicRoutes.has(location.pathname)
    const session = await auth.fetchSession()

    if (!session && isPublicRoute) {
      return
    }

    if (!session?.isAdmin && !isPublicRoute) {
      await auth.signOut()

      throw redirect({
        to: '/signin',
        replace: true
      })
    }

    if (session?.isAdmin && isPublicRoute) {
      throw redirect({
        to: '/',
        replace: true
      })
    }
  }
})

function RootComponent() {
  return <Outlet />
}
