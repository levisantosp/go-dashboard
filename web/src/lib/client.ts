import { client } from '../gen/.kubb/client'

client.interceptors.request.use((req) => {
  req.withCredentials = true
  return req
})

// export default client
