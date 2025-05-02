import type { NextApiRequest, NextApiResponse } from 'next'
import axios from 'axios'

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'GET') return res.status(405).end()
  try {
    const backend = process.env.API_URL || 'http://localhost:8080'
    const token = req.headers.authorization || req.cookies['jwt']
    const response = await axios.get(`${backend}/subscriptions/list`, {
      headers: { Authorization: token ? `Bearer ${token.replace(/^Bearer\s/, '')}` : '' }
    })
    res.status(200).json(response.data)
  } catch (err: any) {
    res.status(err.response?.status || 500).send(err.response?.data || 'Failed to fetch subscriptions')
  }
}
