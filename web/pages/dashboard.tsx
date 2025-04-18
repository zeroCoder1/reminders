import { GetServerSideProps } from 'next'
import axios from 'axios'

export default function Dashboard({ subscriptions }: any) {
  return (
    <main className="p-6">
      <h1 className="text-2xl font-bold mb-4">Your Subscriptions</h1>
      <ul className="space-y-3">
        {subscriptions.map((s: any) => (
          <li key={s.id} className="border p-4 rounded shadow-sm">
            <p><strong>{s.name}</strong> - {s.type}</p>
            <p>{s.start_date} to {s.end_date || 'ongoing'} | {s.currency} {s.amount}</p>
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
    const res = await axios.get('api/subscriptions/list', {
      headers: { Authorization: `Bearer ${token}` }
    })
    return { props: { subscriptions: res.data } }
  } catch (e) {
    return { redirect: { destination: '/', permanent: false } }
  }
}
