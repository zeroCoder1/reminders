import type { NextApiRequest, NextApiResponse } from 'next'
import axios from 'axios'

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'GET') return res.status(405).json({ error: 'Method not allowed' })
  try {
    const backend = process.env.API_URL || 'https://reminders-production-2ada.up.railway.app'
    const token = req.cookies['jwt']
    const response = await axios.get(`${backend}/subscriptions/list`, {
      headers: { Authorization: token ? `Bearer ${token}` : '' }
    })
    res.status(200).json(response.data)
  } catch (err: any) {
    // Log error details for debugging
    console.error('Subscriptions List API error:', err?.response?.status, err?.response?.data, err?.message)
    const status = err.response?.status || 500
    const message = err.response?.data || 'Failed to fetch subscriptions'
    res.status(status).json({ error: typeof message === 'string' ? message : JSON.stringify(message) })
  }
}
