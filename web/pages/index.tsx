import { GetServerSideProps } from 'next'
import { useState } from 'react'
import { useRouter } from 'next/router'
import axios from 'axios'

export default function AuthPage() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [signup, setSignup] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const handleSubmit = async (e: any) => {
    e.preventDefault()
    setError('')
    setSuccess('')
    try {
      if (signup) {
        await axios.post('/api/signup', { email, password })
        setSuccess('Signup successful! Please log in.')
        setSignup(false)
      } else {
        const res = await axios.post('/api/login', { email, password })
        if (res.data?.token) router.push('/dashboard')
      }
    } catch (err: any) {
      const msg = err.response?.data
      setError(typeof msg === 'string' ? msg : msg?.error || (signup ? 'Signup failed' : 'Login failed'))
    }
  }

  return (
    <main className="flex flex-col items-center justify-center min-h-screen p-4">
      <form onSubmit={handleSubmit} className="bg-white p-6 rounded shadow-md w-full max-w-md">
        <h1 className="text-2xl font-bold mb-4">{signup ? 'Sign Up' : 'Login'}</h1>
        <input type="email" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} className="w-full p-2 border rounded mb-3" />
        <input type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} className="w-full p-2 border rounded mb-4" />
        <button type="submit" className="w-full bg-black text-white p-2 rounded mb-2">{signup ? 'Sign Up' : 'Login'}</button>
        <button type="button" className="w-full text-blue-600 underline" onClick={() => { setSignup(!signup); setError(''); setSuccess(''); }}>
          {signup ? 'Already have an account? Login' : "Don't have an account? Sign Up"}
        </button>
        {error && <p className="text-red-600 mt-2">{error}</p>}
        {success && <p className="text-green-600 mt-2">{success}</p>}
      </form>
    </main>
  )
}

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  if (ctx.req.cookies['jwt']) {
    return { redirect: { destination: '/dashboard', permanent: false } }
  }
  return { props: {} }
}
