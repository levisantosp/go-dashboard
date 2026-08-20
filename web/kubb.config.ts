import { pluginAxios } from '@kubb/plugin-axios'
import { pluginReactQuery } from '@kubb/plugin-react-query'
import { pluginTs } from '@kubb/plugin-ts'
import { pluginZod } from '@kubb/plugin-zod'
import { defineConfig } from 'kubb/config'
import { env } from './src/env'

export default defineConfig({
  input: 'http://localhost:3333/openapi.json',
  output: {
    path: './src/gen'
  },
  plugins: [
    pluginTs(),
    pluginZod(),
    pluginAxios({
      baseURL: 'http://localhost:3333'
    }),
    pluginReactQuery({
      client: 'axios'
    })
  ]
})
