import { useQuery } from '@tanstack/react-query'
import { createFileRoute, useRouter } from '@tanstack/react-router'
import { LogOut, Server } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Spinner } from '../components/ui/spinner'
import { getStatusCountQueryOptions } from '../gen/hooks'
import { auth } from '../lib/auth'

export const Route = createFileRoute('/')({
  component: App
})

function App() {
  const router = useRouter()
  const { isPending, error, data } = useQuery({
    ...getStatusCountQueryOptions(),
    staleTime: 60_000
  })

  async function handleSignOut() {
    await auth.signOut()
    await router.invalidate()
  }

  return (
    <div>
      <div className='flex items-center justify-between'>
        <h1 className='text-xl font-bold md:text-2xl'>Overview</h1>
        <Button variant='outline' onClick={handleSignOut}>
          <LogOut />
          Sign out
        </Button>
      </div>

      <div className='mt-10 flex flex-wrap justify-center gap-5'>
        <Card className='h-30 w-80'>
          <CardHeader className='flex justify-between'>
            <CardTitle>Status</CardTitle>
            <Server className='text-violet-400' />
          </CardHeader>

          <CardContent className='text-muted-foreground'>
            {isPending ? (
              <Spinner />
            ) : error ? (
              <span className='text-red-400'>
                An unexpected error has occurred
              </span>
            ) : (
              data.count
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
