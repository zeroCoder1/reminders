import express from 'express';
import jwt from 'jsonwebtoken';
const router = express.Router();

router.post('/login', (req, res) => {
  const { email } = req.body;
  const token = jwt.sign({ email }, process.env.JWT_SECRET!, { expiresIn: '7d' });
  res.json({ token });
});

export default router;