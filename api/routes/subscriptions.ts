import express from 'express';
import pool from '../db';
const router = express.Router();

router.get('/:userId', async (req, res) => {
  const { userId } = req.params;
  const subs = await pool.query('SELECT * FROM subscriptions WHERE user_id = $1', [userId]);
  res.json(subs.rows);
});

router.post('/', async (req, res) => {
  const { user_id, name, cost, category, billing_cycle, start_date, notes } = req.body;
  const result = await pool.query(
    'INSERT INTO subscriptions (user_id, name, cost, category, billing_cycle, start_date, notes) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING *',
    [user_id, name, cost, category, billing_cycle, start_date, notes]
  );
  res.json(result.rows[0]);
});

router.delete('/:id', async (req, res) => {
  const { id } = req.params;
  await pool.query('DELETE FROM subscriptions WHERE id = $1', [id]);
  res.sendStatus(204);
});

export default router;