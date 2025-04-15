import { useEffect, useState } from 'react';
import { Subscription } from '../types/Subscription';

export default function Dashboard() {
  const [subs, setSubs] = useState<Subscription[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      const res = await fetch('/api/subscriptions/1'); // Replace '1' with actual user_id
      const data = await res.json();
      setSubs(data);
    };
    fetchData();
  }, []);

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">My Subscriptions</h1>
      <ul className="space-y-4">
        {subs.map(sub => (
          <li key={sub.id} className="p-4 border rounded shadow-sm">
            <h2 className="text-xl font-semibold">{sub.name}</h2>
            <p>₹{sub.cost} / {sub.billing_cycle}</p>
            <p>Category: {sub.category}</p>
            <p>Renewal: {sub.start_date}</p>
            <p className="text-sm italic text-gray-600">{sub.notes}</p>
          </li>
        ))}
      </ul>
    </div>
  );
}