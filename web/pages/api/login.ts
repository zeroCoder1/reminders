import type { NextApiRequest, NextApiResponse } from 'next'
import axios from 'axios'

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'POST') return res.status(405).json({ error: 'Method not allowed' })
  try {
    const backend = process.env.API_URL || 'http://localhost:8080'
    const response = await axios.post(`${backend}/login`, req.body)
    if (response.data.token) {
      res.setHeader('Set-Cookie', `jwt=${response.data.token}; Path=/; HttpOnly`)
    }
    res.status(200).json(response.data)
  } catch (err: any) {
    const status = err.response?.status || 500
    const message = err.response?.data || 'Login failed'
    res.status(status).json({ error: typeof message === 'string' ? message : JSON.stringify(message) })
  }
}
