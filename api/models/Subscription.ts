export interface Subscription {
  id: string;
  user_id: string;
  name: string;
  category: string;
  cost: number;
  billing_cycle: 'Monthly' | 'Yearly';
  start_date: string;
  notes?: string;
}