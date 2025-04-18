import { useState } from 'react'
import { useRouter } from 'next/router'
import axios from 'axios'

export default function Login() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')

  const handleLogin = async (e: any) => {
    e.preventDefault()
    try {
      const res = await axios.post('/api/login', { email, password })
      if (res.data?.success) router.push('/dashboard')
    } catch (err: any) {
      setError(err.response?.data || 'Login failed')
    }
  }

  return (
    <main className="flex flex-col items-center justify-center min-h-screen p-4">
      <form onSubmit={handleLogin} className="bg-white p-6 rounded shadow-md w-full max-w-md">
        <h1 className="text-2xl font-bold mb-4">Login</h1>
        <input type="email" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} className="w-full p-2 border rounded mb-3" />
        <input type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} className="w-full p-2 border rounded mb-4" />
        <button type="submit" className="w-full bg-black text-white p-2 rounded">Login</button>
        {error && <p className="text-red-600 mt-2">{error}</p>}
      </form>
    </main>
  )
}
