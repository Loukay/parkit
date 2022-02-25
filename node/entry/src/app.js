import express from 'express';
import morgan from 'morgan';
import helmet from 'helmet';
import 'dotenv/config';
import { createServer } from 'http';
import { imagesRouter } from './routes/images.js';
import { normalizePort } from './utils.js';
import bodyParser from 'body-parser';

const app = express();

const port = normalizePort(process.env.PORT || '3001');

const server = createServer(app);

app.use(helmet());
app.use(morgan('dev'));
app.use(bodyParser.raw({
    type: 'image/*',
    limit: '1mb'
}));

app.use('/images', imagesRouter);

app.set('port', port);

server.listen(port);
