import express from 'express';
import cors from 'cors';
import subscriptionsRouter from './routes/subscriptions';
import authRouter from './routes/auth';
import dotenv from 'dotenv';
dotenv.config();

const app = express();
app.use(cors());
app.use(express.json());

app.use('/api/subscriptions', subscriptionsRouter);
app.use('/api/auth', authRouter);

const PORT = process.env.PORT || 5000;
app.listen(PORT, () => console.log(`Server running on port ${PORT}`));