import { GetServerSideProps } from 'next'
import axios from 'axios'
import { useRouter } from 'next/router'

export default function Dashboard({ subscriptions }: any) {
  const router = useRouter()
  const handleLogout = async () => {
    document.cookie = 'jwt=; Path=/; Max-Age=0'
    router.push('/')
  }
  return (
    <main className="p-6 max-w-2xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Your Subscriptions</h1>
        <button onClick={handleLogout} className="text-red-600 underline">Logout</button>
      </div>
      <ul className="space-y-3">
        {subscriptions.length === 0 && <li>No subscriptions found.</li>}
        {subscriptions.map((s: any) => (
          <li key={s.id} className="border p-4 rounded shadow-sm">
            <div className="flex justify-between">
              <div>
                <p><strong>{s.name}</strong> <span className="text-gray-500">({s.type})</span></p>
                <p className="text-sm text-gray-600">{s.start_date} to {s.end_date || 'ongoing'}</p>
              </div>
              <div className="text-right">
                <p className="font-mono">{s.currency} {s.amount}</p>
              </div>
            </div>
          </li>
        ))}
      </ul>
    </main>
  )
}

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  const token = ctx.req.cookies['jwt']
  if (!token) return { redirect: { destination: '/', permanent: false } }
  try {
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || `http://${ctx.req.headers.host}`
    const res = await axios.get(`${baseUrl}/api/subscriptions/list`, {
      headers: { cookie: ctx.req.headers.cookie || '' }
    })
    return { props: { subscriptions: res.data } }
  } catch {
    return { redirect: { destination: '/', permanent: false } }
  }
}
