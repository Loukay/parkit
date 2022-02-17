import express from 'express';

import morgan from 'morgan';

import { createServer } from 'http';

import { imagesRouter } from './routes/images.js';

import { createClient } from './clients/amqpClient.js';

import { normalizePort } from './server.js';

const app = express();

const port = normalizePort(process.env.PORT || '3001');

const server = createServer(app);

const channel = await createClient('amqp://localhost');

channel.sendToQueue('parking_images', Buffer.from('test message'));

app.use(morgan('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));

app.use('/images', imagesRouter);

app.set('port', port);

server.listen(port);
