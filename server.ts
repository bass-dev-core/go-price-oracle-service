import express, { Request, Response } from 'express';

const app = express();
const port = process.env.PORT || 3000;

// Используем встроенный JSON-парсер
app.use(express.json());

app.post('/webhook', (req: Request, res: Response) => {
  const data = req.body;
  console.log('Получен webhook:', data);
  res.status(200).json({ message: 'Webhook получен успешно' });
});

app.listen(port, () => {
  console.log(`Сервер запущен и слушает порт ${port}`);
});
