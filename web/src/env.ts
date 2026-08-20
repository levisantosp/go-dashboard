import { z } from 'zod/mini'

const schema = z.object({
  VITE_API_URL: z.url()
})

export const env = schema.parse(import.meta.env)
