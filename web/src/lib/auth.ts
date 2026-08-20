import { queryOptions, useQuery } from '@tanstack/react-query'
import { env } from '../env'
import { postAuthSignout } from '../gen/clients'
import {
  getAuthSessionQueryKey,
  getAuthSessionQueryOptions
} from '../gen/hooks'
import { queryClient } from '../main'

const query = queryOptions(getAuthSessionQueryOptions())

function createAuthClient() {
  async function fetchSession() {
    return await queryClient.fetchQuery(query)
  }

  function signIn() {
    window.location.href = `${env.VITE_API_URL}/auth/discord`
  }

  async function signOut() {
    await postAuthSignout()
    await queryClient.invalidateQueries({
      queryKey: getAuthSessionQueryKey()
    })
  }

  function useSession() {
    return useQuery({
      ...query,
      staleTime: 10 * 60_000
    })
  }

  return {
    fetchSession,
    signIn,
    signOut,
    useSession
  }
}

export const auth = createAuthClient()
